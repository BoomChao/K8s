package main

import (
	"context"
	"flag"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func main() {
	ctx := context.Background()

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "")
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	mgr, err := manager.New(config, manager.Options{})
	if err != nil {
		panic(err)
	}

	/*
		c, err := controller.New("pod-controller", mgr, controller.Options{
			Reconciler: &ReconcilePod{Client: mgr.GetClient()},
		})
		if err != nil {
			panic(err)
		}

		err = c.Watch(source.Kind(mgr.GetCache(), &corev1.Pod{}), &handler.EnqueueRequestForObject{})
		if err != nil {
			panic(err)
		}

	*/

	c, err := controller.New("node-controller", mgr, controller.Options{
		Reconciler: &ReconcileNode{
			Client:    mgr.GetClient(),
			Clientset: clientset,
		},
	})
	if err != nil {
		panic(err)
	}

	err = c.Watch(source.Kind(mgr.GetCache(), &corev1.Node{}), &handler.EnqueueRequestForObject{})
	if err != nil {
		panic(err)
	}

	if err := mgr.Start(ctx); err != nil {
		panic(err)
	}
}
