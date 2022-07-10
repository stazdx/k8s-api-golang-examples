package main

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

func main() {
	fmt.Println("Hello world!")

	cs := k8sConnect()

	podWatcher(cs)

}

func podWatcher(cs *kubernetes.Clientset) {

	ctx := getContext()

	pods, _ := cs.CoreV1().Pods(apiv1.NamespaceDefault).List(ctx, metav1.ListOptions{})

	if pods.Items != nil {
		for _, pod := range pods.Items {
			fmt.Println(pod.Name, " -> ", pod.Status)
		}
	}
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

func getContext() context.Context {
	ctx := context.Background()

	return ctx
}
