package main

import (
	"context"
	"flag"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
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

	dyClient, err := dynamic.NewForConfig(config)
	if err != nil {
		_ = fmt.Errorf("err:%s", err)
	}
	resource := &schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := dyClient.Resource(*resource).Namespace("kube-system").Get(ctx, "kube-proxy-8f8j8", v1.GetOptions{})

	if err != nil {
		log.Fatalf("Error getting Pod: %s", err)
	}

	// 将 res.Object 赋值给 unstructured.Unstructured 的 Object 字段
	pod := unstructured.Unstructured{Object: res.Object}

	// 打印 Pod 名称和命名空间
	log.Printf("Pod Name: %s", pod.GetName())
	log.Printf("Pod Namespace: %s", pod.GetNamespace())

}
