package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	//读取config配置文件
	kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	//实例化clientset对象
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("err %s, creating clientset\n", err.Error())
	}

	/*
		//ToDo:创建Pod之间先检查这个要创建的Pod是否存在,存在则删除Pod
		err = clientset.CoreV1().Pods(corev1.NamespaceDefault).Delete(context.Background(), "mynginx", metav1.DeleteOptions{})
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 5)
	*/

	//Pod的信息
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mynginx", //这是要创建的Pod的名字
			Labels: map[string]string{
				"app": "web",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx",
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}

	res, err := clientset.CoreV1().Pods(corev1.NamespaceDefault).Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	//输出创建的Pod的名字
	fmt.Println(res.Name)
}
