package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	deploymentName := flag.String("deployment", "", "deployment name")
	imageName := flag.String("image", "", "new image name")
	appName := flag.String("app", "app", "application name")

	flag.Parse()

	if *deploymentName == "" || *imageName == "" {
		panic("Not specify deployment name or image name")
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deployment, err := clientset.AppsV1().Deployments("default").Get(context.Background(), *deploymentName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	if errors.IsNotFound(err) {
		fmt.Println("Deployment not found")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting deployment %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found deployemnt\n")
		name := deployment.GetName()
		fmt.Println("name ->", name)

		containers := deployment.Spec.Template.Spec.Containers

		found := false
		for i := range containers {
			c := containers
			if c[i].Name == *appName {
				found = true
				fmt.Println("Old image ->", c[i].Image)
				fmt.Println("New image ->", *imageName)
				c[i].Image = *imageName
			}
		}

		if !found {
			fmt.Println("Wrong")
			os.Exit(0)
		}

		_, err := clientset.AppsV1().Deployments("default").Update(context.Background(), deployment, metav1.UpdateOptions{})
		if err != nil {
			panic(err.Error())
		}
	}

}
