package main

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

func main() {
	fmt.Println("Hello world!")

	k8sConnect()

}

func k8sConnect() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatal(err)
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return cs
}
