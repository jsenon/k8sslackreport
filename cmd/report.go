package cmd

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
	var podsothers int
	var podsrunning int
	var podssuccess int
	var pvcbound int
	var pvcothers int

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
	// List Pods interface
	pods, err := client.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	// List Nodes interface
	nodes, err := client.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	// List Namespaces interface
	namespaces, err := client.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	// List Persistant Volume Claims interface
	pvc, err := client.CoreV1().PersistentVolumeClaims("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, n := range pods.Items {
		switch state := n.Status.Phase; state {
		case "Running":
			podsrunning = podsrunning + 1
		case "Succeeded":
			podssuccess = podssuccess + 1
		default:
			podsothers = podsothers + 1
		}
	}

	for _, n := range pvc.Items {
		switch state := n.Status.Phase; state {
		case "Bound":
			pvcbound = pvcbound + 1
		default:
			pvcothers = pvcothers + 1
		}
	}
	fmt.Printf("There are %d nodes in the cluster\n", len(nodes.Items))
	fmt.Printf("There are %d namespaces in the cluster\n", len(namespaces.Items))
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	fmt.Printf("There are %d running pods in the cluster\n", podsrunning)
	fmt.Printf("There are %d completed jobs in the cluster\n", podssuccess)
	fmt.Printf("There are %d failed status in the cluster\n", podsothers)
	fmt.Printf("There are %d pvc in the cluster\n", len(pvc.Items))
	fmt.Printf("There are %d pvc bound in the cluster\n", pvcbound)
	fmt.Printf("There are %d pvc not bound in the cluster\n", pvcothers)
	msg := "There are " + conv(len(nodes.Items)) + " nodes in the cluster \n" +
		"There are " + conv(len(namespaces.Items)) + " namespaces in the cluster\n" +
		"There are " + conv(len(pods.Items)) + " pods in the cluster\n" + "      " + conv(podsrunning) + " Runing, " + conv(podssuccess) + " Completed, " + conv(podsothers) + " failed\n" +
		"There are " + conv(len(pvc.Items)) + " pvc in the cluster \n" + "      " + conv(pvcbound) + " Bound, " + conv(pvcothers) + " not Bound"
	publish(msg)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func publish(msg string) {
	url := os.Getenv("SLACK_URL")

	values := map[string]string{"text": msg}
	b, err := json.Marshal(values)
	if err != nil {
		fmt.Println(err)
	}
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	httpclient := &http.Client{Transport: tr}
	rs, err := httpclient.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}
	defer rs.Body.Close() // nolint: errcheck
}

func conv(n int) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}
