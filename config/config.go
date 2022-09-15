package config

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig, masterURL string
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "",
		"Paths to a kubeconfig. Only required if out-of-cluster.")

	flag.StringVar(&masterURL, "master", "",
		"The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

// GetConfig creates a *rest.Config for talking with a Kubernetes apiserver.
// - If --kubeconfig is set, will use the kubeconfig file at that location.
// - Otherwise, it assumes in-cluster
func GetConfig() (*rest.Config, error) {
	// Command-line flag
	if len(kubeconfig) > 0 {
		return clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	}

	// Environment variable
	if len(os.Getenv("KUBECONFIG")) > 0 {
		return clientcmd.BuildConfigFromFlags(masterURL, os.Getenv("KUBECONFIG"))
	}

	// Fallback to in-cluster
	if c, err := rest.InClusterConfig(); err == nil {
		return c, nil
	}

	// Fallback to default location
	if usr, err := user.Current(); err == nil {
		if c, err := clientcmd.BuildConfigFromFlags("", filepath.Join(usr.HomeDir, ".kube", "config")); err == nil {
			return c, nil
		}
	}

	return nil, fmt.Errorf("could not locate a kubeconfig")
}
