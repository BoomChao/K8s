package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ReconcilePod struct {
	Client client.Client
}

func (rp *ReconcilePod) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      request.Name,
			Namespace: request.Namespace,
		},
	}

	err := rp.Client.Get(ctx, request.NamespacedName, pod)
	if err != nil {
		panic(err)
	}

	fmt.Println(pod.Name)

	return reconcile.Result{}, nil
}
