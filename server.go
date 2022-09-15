package main

import (
	"khomer/config"
	"net/http"

	"github.com/common-nighthawk/go-figure"
	"github.com/labstack/echo/v4"

	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

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

	// webserver
	e := echo.New()
	e.HideBanner = true
	e.GET("/", func(c echo.Context) error {
		for {
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

			time.Sleep(10 * time.Second)

			strMsg := fmt.Sprintf("There are %d pods in the cluster\n", len(pods.Items))
			return c.String(http.StatusOK, strMsg)
		}
		return c.String(http.StatusOK, "No connection with the kubernetes cluster.")
	})
	e.Logger.Fatal(e.Start(":1323"))

}
