/*
Copyright 2022.

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
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	addonsv1 "cluster-api-addon-flux/api/v1alpha1"
)

// SetupWithManager sets up the controller with the Manager.
func (r *FluxAddonReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := ctrl.LoggerFrom(ctx)

	c, err := ctrl.NewControllerManagedBy(mgr).
		For(&addonsv1.FluxAddon{}).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log, r.WatchFilterValue)).
		Build(r)
	if err != nil {
		return fmt.Errorf("creating controller: %w", err)
	}

	if err := c.Watch(
		&source.Kind{Type: &clusterv1.Cluster{}},
		handler.EnqueueRequestsFromMapFunc(r.ClusterToFluxAddonMapper),
		predicates.ResourceNotPausedAndHasFilterLabel(log, r.WatchFilterValue),
	); err != nil {
		return fmt.Errorf("watching capi cluster: %w", err)
	}

	return nil
}

// FluxAddonReconciler reconciles a FluxAddon object
type FluxAddonReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	WatchFilterValue string
}

//+kubebuilder:rbac:groups=addons.cluster.x-k8s.io,resources=fluxaddons,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=addons.cluster.x-k8s.io,resources=fluxaddons/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=addons.cluster.x-k8s.io,resources=fluxaddons/finalizers,verbs=update
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters,verbs=list;watch
//+kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=kubeadmcontrolplanes,verbs=list;get;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FluxAddon object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *FluxAddonReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := log.FromContext(ctx)

	addon := &addonsv1.FluxAddon{}
	if err := r.Client.Get(ctx, req.NamespacedName, addon); err != nil {
		if apierrors.IsNotFound(err) {
			log.V(2).Info("FluxAddon resource not found, skipping")
			return ctrl.Result{}, nil
		}
	}

	patchHelper, err := patch.NewHelper(addon, r.Client)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("creating addon patch helper: %w", err)
	}

	defer func() {
		if err := patchAddon(ctx, patchHelper, addon); err != nil {
			reterr = err
			log.Error(err, "failed to patch FluxAdd")
			return
		}
		log.V(2).Info("Patched FluxAddon")
	}()

	selector := addon.Spec.ClusterSelector
	clusterList, err := r.listClustersWithLabels(ctx, addon.Namespace, selector)
	if err != nil {
		//TODO:conditions

		return ctrl.Result{}, err
	}
	addon.SetMatchingClusters(clusterList.Items)

	log.V(2).Info("Finding FluxAddonInstance for FluxAdd")
	label := map[string]string{
		addonsv1.FluxAddonLabelName: addon.Name,
	}

	return ctrl.Result{}, nil
}

func (r *FluxAddonReconciler) ClusterToFluxAddonMapper(o client.Object) []ctrl.Request {
	cluster, ok := o.(*clusterv1.Cluster)
	if !ok {
		//TODO: log error
		return nil
	}

	fluxAddons := &addonsv1.FluxAddonList{}

	listOpts := []client.ListOption{
		client.MatchingLabels{
			clusterv1.ClusterLabelName: cluster.Name,
		},
	}

	if err := r.Client.List(context.TODO(), fluxAddons, listOpts...); err != nil {
		//TODO: log error
		return nil
	}

	results := []ctrl.Request{}
	for _, addon := range fluxAddons.Items {
		results = append(results, ctrl.Request{
			NamespacedName: client.ObjectKey{
				Namespace: addon.GetNamespace(),
				Name:      addon.Labels[addonsv1.FluxAddonLabelName],
			},
		})
	}

	return results
}

func (r *FluxAddonReconciler) listClustersWithLabels(ctx context.Context, namespace string, selector metav1.LabelSelector) (*clusterv1.ClusterList, error) {
	clusterList := &clusterv1.ClusterList{}
	if err := r.Client.List(ctx, clusterList, client.InNamespace(namespace), client.MatchingLabels(selector.MatchLabels)); err != nil {
		return nil, err
	}

	return clusterList, nil
}

func (r *FluxAddonReconciler) listClustersWithLabels(ctx context.Context, namespace string

func patchAddon(ctx context.Context, patchHelper *patch.Helper, addon *addonsv1.FluxAddon) error {
	conditions.SetSummary(addon,
		conditions.WithConditions(
		//TODO: add conditionms
		),
	)

	return patchHelper.Patch(
		ctx,
		addon,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			clusterv1.ReadyCondition,
			//TODO: addon conditions
		}},
		patch.WithStatusObservedGeneration{},
	)
}
