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
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	appgenv1 "github.com/dgoodwin/argocd-appgenerator/api/v1"
)

const (
	argoCDSecretTypeLabel   = "argocd.argoproj.io/secret-type"
	argoCDSecretTypeCluster = "cluster"
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

func (r *ApplicationGeneratorReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("applicationgenerator", req.NamespacedName)
	log := r.Log.WithValues("appgenerator", req.NamespacedName)

	ag := &appgenv1.ApplicationGenerator{}
	err := r.Client.Get(context.Background(), req.NamespacedName, ag)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("not found", "ApplicationGenerator", req.NamespacedName)
		}
		return ctrl.Result{}, err
	}
	log.Info("reconciling")

	// your logic here

	return ctrl.Result{}, nil
}

func (r *ApplicationGeneratorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appgenv1.ApplicationGenerator{}).
		Watches(
			&source.Kind{Type: &corev1.Secret{}},
			&clusterSecretEventHandler{
				Client: mgr.GetClient(),
				Log:    ctrl.Log.WithName("eventhandler").WithName("clustersecret"),
			}).
		Complete(r)
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
