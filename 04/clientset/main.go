package main

import (
	"context"
	"flag"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"time"
)

func main() {

	kubeconfig := filepath.Join(homedir.HomeDir(), "/workspace/clientgo/.kube", "config")
	flag.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "path to kubeconfig file")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		_ = fmt.Errorf("err:%s", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	coreV1 := clientset.CoreV1()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pod, err := coreV1.Pods("kube-system").Get(ctx, "kube-proxy-8f8j8", v1.GetOptions{})
	if err != nil {
		log.Fatalf("Error getting Pod: %s", err)
	} else {
		log.Printf("Pod Name: %s", pod.Name)
		log.Printf("Pod Name: %s", pod.Namespace)
	}
}
