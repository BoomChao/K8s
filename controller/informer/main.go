package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	factory := informers.NewSharedInformerFactory(clientset, time.Second*10)

	podInformer := factory.Core().V1().Pods().Informer()

	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {},
		// 由于设定了resync的时间为10s,所以这个调用会10s执行一次;即informer会10s将全量的数据塞入到updateFunc的处理逻辑中
		UpdateFunc: func(oldObj, newObj interface{}) {
			pod := newObj.(*corev1.Pod)
			podName := pod.GetName()
			fmt.Println(time.Now(), " Pod updated -> ", podName)
		},
		DeleteFunc: func(obj interface{}) {},
	})

	stopCh := make(chan struct{})
	defer close(stopCh)

	factory.Start(stopCh)
	// 等待第一次list全量获取数据
	if !cache.WaitForCacheSync(stopCh, podInformer.HasSynced) {
		return
	}

	// 这里直接从上面的informer的cache里面拿了
	lists := podInformer.GetIndexer().List()
	for id := range lists {
		pod := lists[id].(*corev1.Pod)
		fmt.Println(pod.GetName())
	}

	<-stopCh
}
