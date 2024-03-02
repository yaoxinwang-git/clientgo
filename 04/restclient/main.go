package main

import (
	"context"
	"flag"
	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"time"
)

func main() {
	//config
	//masterUrl 为空 会从kubeConfig中获取
	// 解析 kubeconfig 文件路径
	kubeconfig := filepath.Join(homedir.HomeDir(), "/workspace/clientgo/.kube", "config")
	flag.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "path to kubeconfig file")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		_ = fmt.Errorf("err:%s", err)
	}
	//config.Host = "http://192.168.0.161:6443"
	config.GroupVersion = &v1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs
	config.APIPath = "/api"
	//config.CAFile = "/Users/xinwang_yao/workspace/clientgo/ca.crt"

	//client
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		_ = fmt.Errorf("err:%s", err)
	}
	//get data
	pod := v1.Pod{}
	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = restClient.Get().Namespace("kube-system").Resource("pods").Name("kube-proxy-8f8j8").Do(ctx).Into(&pod)
	if err != nil {
		log.Fatalf("Error getting Pod: %s", err)
	} else {
		log.Printf("Pod Name: %s", pod.Name)
		log.Printf("Pod Name: %s", pod.Namespace)
	}
}
