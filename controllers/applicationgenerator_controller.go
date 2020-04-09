/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/client-go/util/workqueue"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	argoapi "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"

	appgenv1 "github.com/dgoodwin/argocd-appgenerator/api/v1"
)

const (
	argoCDSecretTypeLabel   = "argocd.argoproj.io/secret-type"
	argoCDSecretTypeCluster = "cluster"
	appGeneratorLabel       = "appgenerator.rm-rf.ca/from-generator"
)

// ApplicationGeneratorReconciler reconciles a ApplicationGenerator object
type ApplicationGeneratorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=appgenerator.rm-rf.ca,resources=applicationgenerators,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=appgenerator.rm-rf.ca,resources=applicationgenerators/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups=argoproj.io,resources=applications,verbs=get;list;watch;create;update;patch;delete

func (r *ApplicationGeneratorReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("applicationgenerator", req.NamespacedName)
	log := r.Log.WithValues("appgenerator", req.NamespacedName)

	// TODO: ArgoCD only acts on Applications in it's namespace

	ag := &appgenv1.ApplicationGenerator{}
	err := r.Client.Get(context.Background(), req.NamespacedName, ag)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("not found", "ApplicationGenerator", req.NamespacedName)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	log.Info("reconciling")

	// List all Clusters:
	clusterSecretList := &corev1.SecretList{}
	secretLabels := map[string]string{
		argoCDSecretTypeLabel: argoCDSecretTypeCluster,
	}
	for k, v := range ag.Spec.ClusterSelector.MatchLabels {
		secretLabels[k] = v
	}
	if err := r.Client.List(context.Background(), clusterSecretList, client.MatchingLabels(secretLabels)); err != nil {
		return ctrl.Result{}, err
	}
	log.Info("clusters matching labels", "count", len(clusterSecretList.Items))

	// Track the list of applications that should exist for this generator, used to
	// cleanup orphans after we're done.
	expectedApps := map[string]bool{}

	for _, clusterSecret := range clusterSecretList.Items {
		appName := GetName(ag.Name, clusterSecret.Name, validation.DNS1123SubdomainMaxLength)
		expectedApps[appName] = true
		app := &argoapi.Application{
			ObjectMeta: metav1.ObjectMeta{
				// Should Applications always go to same namespaces as the AppGenerator? If not, we cannot use OwnerReferences.
				// TODO: no, everything should be in argocd namespace
				Namespace: ag.Namespace,
				Name:      appName,
				Labels: map[string]string{
					appGeneratorLabel: ag.Name,
				},
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: ag.TypeMeta.APIVersion,
						Kind:       ag.TypeMeta.Kind,
						Name:       ag.Name,
						UID:        ag.ObjectMeta.UID,
					},
				},
				// TODO: owner ref to the ApplicationGenerator for free cleanup?
			},
			Spec: ag.Spec.ApplicationSpec,
		}
		// TODO: Check if application already exists and update if so.
		existingApp := &argoapi.Application{}
		err = r.Client.Get(context.Background(), types.NamespacedName{Namespace: app.Namespace, Name: app.Name}, existingApp)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Info("app does not yet exist, creating...", "name", app.Name)
				if err := r.Client.Create(context.Background(), app); err != nil {
					return ctrl.Result{}, err
				}
				log.Info("created app", "name", app.Name)
			} else {
				return ctrl.Result{}, err
			}
		} else {
			log.Info("app already exists", "name", app.Name)
			origApp := existingApp.DeepCopy()

			// Set expected ObjectMeta:
			existingApp.ObjectMeta.Labels[appGeneratorLabel] = ag.Name
			existingApp.OwnerReferences = []metav1.OwnerReference{
				{
					APIVersion: ag.TypeMeta.APIVersion,
					Kind:       ag.TypeMeta.Kind,
					Name:       ag.Name,
					UID:        ag.ObjectMeta.UID,
				},
			}

			existingApp.Spec = ag.Spec.ApplicationSpec

			if !reflect.DeepEqual(existingApp, origApp) {
				log.Info("app has changed, updating...")
				return ctrl.Result{}, r.Client.Update(context.Background(), existingApp)
			} else {
				log.Info("app is in expected state, no update required")
			}
		}

	}

	// Cleanup any orphaned Applications that should no longer exist:
	log.Info("expected apps", "apps", expectedApps)
	appList := &argoapi.ApplicationList{}
	labels := map[string]string{
		appGeneratorLabel: ag.Name,
	}
	if err := r.Client.List(context.Background(), appList, client.MatchingLabels(labels)); err != nil {
		return ctrl.Result{}, err
	}
	log.Info("applications matching labels", "count", len(appList.Items))
	for _, app := range appList.Items {
		if !expectedApps[app.Name] {
			log.Info("found orphaned app to delete", "app", app.Name)
			if err := r.Client.Delete(context.Background(), &app); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

func (r *ApplicationGeneratorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appgenv1.ApplicationGenerator{}).
		// TODO: watch Applications?
		Watches(
			&source.Kind{Type: &corev1.Secret{}},
			&clusterSecretEventHandler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("eventhandler").WithName("clustersecret"),
			}).
		Complete(r)
	// TODO: watch Applications and respond on changes if we own them.
}

var _ handler.EventHandler = &clusterSecretEventHandler{}

type clusterSecretEventHandler struct {
	//handler.EnqueueRequestForOwner
	Log    logr.Logger
	Client client.Client
}

func (h *clusterSecretEventHandler) Create(e event.CreateEvent, q workqueue.RateLimitingInterface) {
	h.queueRelatedAppGenerators(q, e.Meta)
}

func (h *clusterSecretEventHandler) Update(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	h.queueRelatedAppGenerators(q, e.MetaNew)
}

func (h *clusterSecretEventHandler) Delete(e event.DeleteEvent, q workqueue.RateLimitingInterface) {
	h.queueRelatedAppGenerators(q, e.Meta)
}

func (h *clusterSecretEventHandler) Generic(e event.GenericEvent, q workqueue.RateLimitingInterface) {
	h.queueRelatedAppGenerators(q, e.Meta)
}

func (h *clusterSecretEventHandler) queueRelatedAppGenerators(q workqueue.RateLimitingInterface, meta metav1.Object) {
	// Check for label, lookup all ApplicationGenerators that might match secret, queue them all
	if meta.GetLabels()[argoCDSecretTypeLabel] != argoCDSecretTypeCluster {
		return
	}

	h.Log.Info("processing cluster secret event", "namespace", meta.GetNamespace(), "name", meta.GetName())

	appGensList := &appgenv1.ApplicationGeneratorList{}
	err := h.Client.List(context.Background(), appGensList)
	if err != nil {
		h.Log.Error(err, "unable to list ApplicationGenerators")
		return
	}
	h.Log.Info("listed ApplicationGenerators", "count", len(appGensList.Items))
	for _, ag := range appGensList.Items {
		// TODO: only queue the AppGenerator if the labels match this cluster
		req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: ag.Namespace, Name: ag.Name}}
		q.Add(req)
	}
}
