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
	"k8s.io/apimachinery/pkg/runtime"
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
	r.Log.V(1).Info("DING!")

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

func (h *clusterSecretEventHandler) Update(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	//h.Log.Info("processing cluster secret update", e)
}

func (h *clusterSecretEventHandler) Delete(e event.DeleteEvent, q workqueue.RateLimitingInterface) {
	//h.Log.Info("processing cluster secret delete", e)
}

func (h *clusterSecretEventHandler) Generic(e event.GenericEvent, q workqueue.RateLimitingInterface) {
	//h.Log.Info("processing cluster secret generic", e)
}

// Create implements handler.EventHandler
func (h *clusterSecretEventHandler) Create(e event.CreateEvent, q workqueue.RateLimitingInterface) {
	// Check for label, lookup all ApplicationGenerators that might match secret, queue them all
	if e.Meta.GetLabels()[argoCDSecretTypeLabel] != argoCDSecretTypeCluster {
		return
	}

	h.Log.Info("processing cluster secret create", "namespace", e.Meta.GetNamespace(), "name", e.Meta.GetName())
	h.Log.V(5).Info("listing all ApplicationGenerators")

	//appGensList := appgenv1.ApplicationGeneratorList{}
	//err := h.Client.List(context.Background(), appGensList)

	//h.EnqueueRequestForOwner.Create(e, q)
}
