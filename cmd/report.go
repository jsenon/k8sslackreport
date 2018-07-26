package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var api string

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Launch Report",
	Long: `Launch report
	       on slack channel
           `,
	Run: func(cmd *cobra.Command, args []string) {
		if api == "internal" {
			Report("internal")
		}
		Report("external")
	},
}

func init() {
	reportCmd.PersistentFlags().StringVar(&api, "api", "", "api type: internal or external")
	rootCmd.AddCommand(reportCmd)
}

func Report(api string) {
	var kubeconfig *string
	var client *kubernetes.Clientset
	// var podsStore cache.Store
	// var podStorekube cache.Store
	// var nodesStore cache.Store
	// var eventStore cache.Store
	// var eventallStore cache.Store

	// ctx := context.Background()

	fmt.Println("You have selected api: ", api)
	// Internal k8s api
	if api == "internal" {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		client, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	}
	// External k8s api based on .kube/config
	if api == "external" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(homeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		flag.Parse()
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}

		// create the clientset
		client, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

	}
	pods, err := client.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
