package main

import (
	"fmt"
	bolingcavalryv1 "k8s_customize_controller/pkg/apis/bolingcavalry/v1"
	clientset "k8s_customize_controller/pkg/client/clientset/versioned"
	informers "k8s_customize_controller/pkg/client/informers/externalversions/bolingcavalry/v1"
	listers "k8s_customize_controller/pkg/client/listers/bolingcavalry/v1"
	"time"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

	runtime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

const controllerAgentName = "student-controller"

const (
	SuccessSynced          = "synced"
	MessageResourcesSynced = "Student synced successfully"
)

type Controller struct {
	//k8s clientset
	kubeclientset kubernetes.Interface

	//clientset for our own API group
	studentclient clientset.Interface

	listers.StudentLister
	studentSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface

	recoder record.EventRecorder
}

func NewController(kubeclientset kubernetes.Interface, studentclient clientset.Interface, studentInformer informers.StudentInformer) *Controller {

	runtime.Must(scheme.AddToScheme(scheme.Scheme))
	glog.V(4).Info("Creating event broadcaster") //指定日志的级别
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})

	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset: kubeclientset,
		studentclient: studentclient,
		StudentLister: studentInformer.Lister(),
		studentSynced: studentInformer.Informer().HasSynced,
		workqueue:     workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Students"),
		recoder:       recorder,
	}

	glog.Info("Setting up event handlers")

	//set up an event handler for when Student resoucres change
	studentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueStudent,
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldStudent := oldObj.(*bolingcavalryv1.Student)
			newStudent := newObj.(*bolingcavalryv1.Student)
			if oldStudent.ResourceVersion == newStudent.ResourceVersion {
				//版本一致，就表示没有实际更新的操作，立即返回
				return
			}
			controller.enqueueStudent(newObj)
		},
		DeleteFunc: controller.enqueueStudentForDelete,
	})

	return controller
}

//在此处开始controller的业务
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	glog.Info("Start controller, Start a cache data synchronization")
	if ok := cache.WaitForCacheSync(stopCh, c.studentSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("worker start")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("worker already started")
	<-stopCh
	glog.Info("worker already ended")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {

	}
}

//取数据处理
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}

	//We wrap this block in a func so we can defer c.workqueue.Done
	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}

		//在syncHandler中处理业务
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s':%s", key, err.Error())
		}

		c.workqueue.Forget(obj)
		glog.Info(fmt.Sprintf("Successfully synced '%s'", key))
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

//处理
func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	//从缓存中取对象
	student, err := c.StudentLister.Students(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			glog.Info(fmt.Sprintf("The Student object is deleted, please perform the actual deletion here: %s/%s ...", namespace, name))
			return nil
		}
		runtime.HandleError(fmt.Errorf("failed to list sutdent by: %s/%s", namespace, name))
		return err
	}

	glog.Info(fmt.Sprintf("Here is the desired status of the student %#v", student))
	glog.Info("实际状态是从业务层面得到的，此处应该取实际状态与期望状态做对比，并根据差异做出相应(新增或者删除)")

	c.recoder.Event(student, corev1.EventTypeNormal, SuccessSynced, MessageResourcesSynced)

	return nil
}

//数据先入缓存，再入队列
func (c *Controller) enqueueStudent(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}

	//将key放入到队列中
	c.workqueue.AddRateLimited(key)
}

//删除操作
func (c *Controller) enqueueStudentForDelete(obj interface{}) {
	var key string
	var err error
	key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}
