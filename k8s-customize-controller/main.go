package main

import (
	"flag"
	clientset "k8s_customize_controller/pkg/client/clientset/versioned"
	informers "k8s_customize_controller/pkg/client/informers/externalversions"
	"path/filepath"
	"time"

	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// var (
// 	masterURL  string
// 	kubeconfig string
// )

func main() {
	flag.Parse() //这个必须添加，因为glob要求使用前必须先进行flag.Parse

	//处理信号量
	// stopCh := singals.SetupSignalHandler()	//这个好像没什么用

	kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "")

	//处理入参
	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	studentClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	studentInformerFactory := informers.NewSharedInformerFactory(studentClient, time.Second*30)

	//得到controller
	controller := NewController(kubeClient, studentClient, studentInformerFactory.Bolingcavalry().V1().Students())

	stopCh := make(chan struct{})

	//启动informer
	go studentInformerFactory.Start(stopCh)

	//controller开始处理消息
	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controllers: %s", err.Error())
	}

}

// func init() {
// 	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if ...")
// 	flag.StringVar(&masterURL, "master", "", "The address of kubernetes API server...")
// }
