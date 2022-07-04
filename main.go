package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	// v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	fmt.Println("--Testing--")

	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", filepath.Join("/Users/stazdx", ".kube", "config"), "(optional) absolute path to the kubeconfig file")

	flag.Parse()

	configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfig}
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: "eks_saleor-backendv2-qa"}

	kconf, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()
	if err != nil {
		fmt.Println(nil, err)
	}

	ctx := context.Background()

	fmt.Println(ctx)

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	clientset, err := kubernetes.NewForConfig(config)

	list, err := clientset.AppsV1().Deployments(apiv1.NamespaceDefault).List(ctx, metav1.ListOptions{})

	fmt.Println(kconf, list)

}
