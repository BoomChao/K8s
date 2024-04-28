package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ReconcileNode struct {
	Client    client.Client
	Clientset kubernetes.Interface
}

func (rp *ReconcileNode) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: request.Name,
		},
	}

	err := rp.Client.Get(ctx, request.NamespacedName, node)
	if err != nil {
		panic(err)
	}

	fmt.Println(node.Name)

	node.Annotations["test"] = "test"
	if _, err := rp.Clientset.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{}); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
