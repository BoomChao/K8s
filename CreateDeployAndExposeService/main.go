package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/utils/pointer"
)

const (
	NAMESPACE       = "test-clientset"
	DEPLOYMENT_NAME = "client-test-deployment"
	SERVICE_NAME    = "client-test-service"
)

//在特定NameSpace内创建一个Deployment之后再创建一个Service暴露服务
func main() {
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "")
	}

	operate := flag.String("operate", "create", "operate type : create or clean")

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	fmt.Printf("operation is %v\n", *operate)

	if *operate == "clean" {
		clean(clientset)
	} else {
		createNamespace(clientset)
		createDeployment(clientset)
		createService(clientset)
	}
}

//清理本次实验的所有资源
func clean(clientset *kubernetes.Clientset) {
	//删除Service
	if err := clientset.CoreV1().Services(NAMESPACE).Delete(context.TODO(), SERVICE_NAME, metav1.DeleteOptions{}); err != nil {
		panic(err)
	}

	//删除Deployment
	if err := clientset.AppsV1().Deployments(NAMESPACE).Delete(context.TODO(), DEPLOYMENT_NAME, metav1.DeleteOptions{}); err != nil {
		panic(err)
	}

	//删除namespace
	if err := clientset.CoreV1().Namespaces().Delete(context.TODO(), NAMESPACE, metav1.DeleteOptions{}); err != nil {
		panic(err)
	}
}

//新建Namespace
func createNamespace(clientset *kubernetes.Clientset) {
	namespaceClient := clientset.CoreV1().Namespaces()

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: NAMESPACE,
		},
	}

	result, err := namespaceClient.Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Create namespace %s \n", result.GetName())
}

//新建一个Deployment
func createDeployment(clientset *kubernetes.Clientset) {
	deploymentClient := clientset.AppsV1().Deployments(NAMESPACE)

	//实例化Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: DEPLOYMENT_NAME,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "tomcat",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "tomcat",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "tomcat",
							Image:           "tomcat:8.0.18-jre8", //这是所需要的镜像
							ImagePullPolicy: "IfNotPresent",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolSCTP,
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := deploymentClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Create deployment %s \n", result.GetName())
}

//新建Service,暴露上面的Deployment
func createService(clientset *kubernetes.Clientset) {
	//得到Service的客户端
	serviceClient := clientset.CoreV1().Services(NAMESPACE)

	//实例化一个数据结构
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: SERVICE_NAME,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     8080, //将Pod的8080端口映射到node的30080端口
					NodePort: 30080,
				},
			},
			Selector: map[string]string{
				"app": "tomcat",
			},
		},
	}

	result, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Create service %s\n", result.GetName())
}
