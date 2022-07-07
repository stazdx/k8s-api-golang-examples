package main

import (
	"fmt"

	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

func main() {
	fmt.Println("Hello!")

}

func k8sConnect() {
	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatal(err)
	}

	fmt.Println(config)
}
