package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	// v1 "k8s.io/api/apps/v1"
	"github.com/slack-go/slack"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	fmt.Println("--Testing--")

	// send slack notification
	// webhookSlack()

	var kubeconfig *string
	usr, _ := os.UserHomeDir()

	// kubeconfig flag
	kubeconfig = flag.String("kubeconfig", filepath.Join(usr, ".kube", "config"), "(optional) absolute path to the kubeconfig file")

	flag.Parse()

	configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfig}

	// setting custom context
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: "microk8s"}

	kconf, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()
	if err != nil {
		fmt.Println(nil, err)
	}

	ctx := context.Background()

	// new k8s connection
	clientset, err := kubernetes.NewForConfig(kconf)
	if err != nil {
		panic(err)
	}

	// get deployments in default namespace
	deployments, err := clientset.AppsV1().Deployments(apiv1.NamespaceDefault).List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// print results - Deployments
	fmt.Println("\n--------- DEPLOYMENTS --------- \n\n", deployments.Items[0])

	// list all pods in default namespace
	pods, _ := clientset.CoreV1().Pods(apiv1.NamespaceDefault).List(ctx, metav1.ListOptions{})
	fmt.Println("\n--------- PODS ---------\n\n", pods.Items)

	if pods.Items != nil {
		for _, pod := range pods.Items {
			fmt.Println(pod.Name, " -> ", pod.Status)
		}
	}

	// list all pods in default namespace with label selector
	custom_pod, _ := clientset.CoreV1().Pods(apiv1.NamespaceDefault).List(ctx, metav1.ListOptions{LabelSelector: "app=demo"})
	fmt.Println("\n--------- CUSTOM POD ---------\n\n", custom_pod.Items)

	PodName := ""

	if custom_pod.Items != nil {
		for _, pod := range custom_pod.Items {
			// if the Pod doesn't have ready status
			if pod.Status.ContainerStatuses[0].Ready != true {
				// if the status is waiting or else (terminating)
				if pod.Status.ContainerStatuses[0].State.Waiting != nil {
					// send Slack notification with error
					webhookSlack(pod.Name, "Waiting", pod.Status.ContainerStatuses[0].State.Waiting.Reason,
						pod.Status.ContainerStatuses[0].State.Waiting.Message)
				} else {
					// send Slack notification with error
					webhookSlack(pod.Name, "Terminated", pod.Status.ContainerStatuses[0].State.Terminated.Reason,
						pod.Status.ContainerStatuses[0].State.Terminated.Message)
				}
			}
			PodName = pod.Name
			fmt.Println(pod.Name, " -> Ready:", pod.Status.ContainerStatuses[0].Ready,
				pod.Status.ContainerStatuses[0].State.Running)
		}
	}

	// get custom pod in default namespace - Name: test-5f6778868d-grcn7
	pod, err := clientset.CoreV1().Pods(apiv1.NamespaceDefault).Get(ctx, PodName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	// print pod status
	fmt.Println("\n--------- POD STATUS --------- \n\n", pod.Status)

}

func webhookSlack(rs string, status string, reason string, message string) {
	attachment := slack.Attachment{
		Color:         "#3AA3E3",
		Fallback:      "Kubernetes cluster has changes!",
		AuthorName:    "Staz Dx",
		AuthorSubname: "github.com",
		AuthorLink:    "https://github.com/stazdx",
		AuthorIcon:    "https://avatars2.githubusercontent.com/u/1691541",
		Text:          "<!channel> Resource Details:\n :notebook: name: `" + rs + "` \n :eyes: Status: `" + status + "` \n :bangbang: Reason: `" + reason + "` \n :warning: Message: `" + message + "`",
		Footer:        "slack api",
		FooterIcon:    "https://platform.slack-edge.com/img/default_application_icon.png",
		Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook("https://hooks.slack.com/services/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", &msg)
	if err != nil {
		fmt.Println(err)
	}
}
