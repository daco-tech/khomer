package main

import (
	"khomer/config"
	"net/http"

	"github.com/common-nighthawk/go-figure"
	"github.com/labstack/echo/v4"

	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// ServiceController keeps a registry of the Services that exist in the cluster
type ServiceController struct {
	informerFactory informers.SharedInformerFactory
	serviceInformer coreinformers.ServiceInformer

	services []string
}

func (c *ServiceController) Run(stopper chan struct{}) error {
	// start all the shared informers that have been created by the factory so far
	c.informerFactory.Start(stopper)

	// wait for initial synchronization of the local cache
	if !cache.WaitForCacheSync(stopper, c.serviceInformer.Informer().HasSynced) {
		return fmt.Errorf("timed out waiting for caches to sync")
	}

	return nil
}

func (c *ServiceController) onAdd(obj interface{}) {
	service := obj.(*v1.Service)
	fmt.Printf("Service added: %s\n", service.Name)

	c.services = append(c.services, service.Name)
}

func (c *ServiceController) onUpdate(oldObj interface{}, newObj interface{}) {
	oldservice := oldObj.(*v1.Service)
	newservice := newObj.(*v1.Service)
	fmt.Printf("Service updated: %s => %s\n", oldservice.Name, newservice.Name)

	for i, v := range c.services {
		if v == oldservice.Name {
			c.services = append(c.services[:i], c.services[i+1:]...)
			break
		}
	}

	c.services = append(c.services, newservice.Name)
}

func (c *ServiceController) onDelete(obj interface{}) {
	service := obj.(*v1.Service)
	fmt.Printf("Service deleted: %s\n", service.Name)

	for i, v := range c.services {
		if v == service.Name {
			c.services = append(c.services[:i], c.services[i+1:]...)
			break
		}
	}
}

// an informer internally consists of a watcher, a lister and an in-memory cache.
func createServiceController(informerFactory informers.SharedInformerFactory) *ServiceController {
	serviceInformer := informerFactory.Core().V1().Services()

	c := &ServiceController{
		informerFactory: informerFactory,
		serviceInformer: serviceInformer,
	}

	serviceInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.onAdd,
			UpdateFunc: c.onUpdate,
			DeleteFunc: c.onDelete,
		},
	)

	return c
}

func main() {
	// init
	myFigure := figure.NewFigure("kHomer", "", true)
	myFigure.Print()

	// creates the in-cluster config
	config, err := config.GetConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// create an Informer to monitor Services
	factory := informers.NewSharedInformerFactory(clientset, 0)
	controller := createServiceController(factory)

	stop := make(chan struct{})
	defer close(stop)
	if err := controller.Run(stop); err != nil {
		panic(err.Error())
	}

	// webserver
	e := echo.New()
	e.HideBanner = true
	e.GET("/", func(c echo.Context) error {
		for {
			fmt.Printf("There are %d services in the cluster\n", len(controller.services))

			// get pods in all the namespaces by omitting namespace
			// Or specify namespace to get pods in particular namespace
			pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

			// Examples for error handling:
			// - Use helper functions e.g. errors.IsNotFound()
			// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
			_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "example-xxxxx", metav1.GetOptions{})
			if errors.IsNotFound(err) {
				fmt.Printf("Pod example-xxxxx not found in default namespace\n")
			} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
				fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
			} else if err != nil {
				panic(err.Error())
			} else {
				fmt.Printf("Found example-xxxxx pod in default namespace\n")
			}

			// time.Sleep(10 * time.Second)

			strMsg := fmt.Sprintf("There are %d pods in the cluster\nThere are %d services in the cluster\n", len(pods.Items), len(controller.services))
			return c.String(http.StatusOK, strMsg)
		}
	})
	e.Logger.Fatal(e.Start(":1323"))

}
