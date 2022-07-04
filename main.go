package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
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

	webhookSlack()

	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", filepath.Join("/Users/stazdx", ".kube", "config"), "(optional) absolute path to the kubeconfig file")

	flag.Parse()

	configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfig}
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: "microk8s"}

	kconf, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()
	if err != nil {
		fmt.Println(nil, err)
	}

	ctx := context.Background()

	// config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	// if err != nil {
	// 	panic(err)
	// }

	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	panic(err)
	// }

	clientset2, err := kubernetes.NewForConfig(kconf)
	if err != nil {
		panic(err)
	}

	list2, err := clientset2.AppsV1().Deployments(apiv1.NamespaceDefault).List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println(list2.Items[0])

	// list, err := clientset.AppsV1().Deployments(apiv1.NamespaceDefault).List(ctx, metav1.ListOptions{})
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(list.Items)

}

func webhookSlack() {
	attachment := slack.Attachment{
		Color:         "good",
		Fallback:      "You successfully posted by Incoming Webhook URL!",
		AuthorName:    "Staz Dx",
		AuthorSubname: "github.com",
		AuthorLink:    "https://github.com/stazdx",
		AuthorIcon:    "https://avatars2.githubusercontent.com/u/1691541",
		Text:          "<!channel> All text in Slack uses the same system of escaping: chat messages, direct messages, file comments, etc. :smile:\nSee <https://api.slack.com/docs/message-formatting#linking_to_channels_and_users>",
		Footer:        "slack api",
		FooterIcon:    "https://platform.slack-edge.com/img/default_application_icon.png",
		Ts:            json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	msg := slack.WebhookMessage{
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.PostWebhook("SLACK_WEBHOOK_URL", &msg)
	if err != nil {
		fmt.Println(err)
	}
}
