package main

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {

	kubeconfig := filepath.Join(homedir.HomeDir(), "/workspace/clientgo/.kube", "config")
	flag.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "path to kubeconfig file")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		_ = fmt.Errorf("err:%s", err)
	}

	// 创建 Discovery Client
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 确认 Pods 资源是可用的（这一步实际上对于获取 Pods 不是必需的，但它展示了 Discovery Client 的用法）
	_, apiResourceLists, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		panic(err.Error())
	}
	foundPodResource := false
	for _, apiResourceList := range apiResourceLists {
		if apiResourceList.GroupVersion == "v1" {
			for _, apiResource := range apiResourceList.APIResources {
				if apiResource.Kind == "Pod" {
					foundPodResource = true
					break
				}
			}
		}
		if foundPodResource {
			break
		}
	}

	if !foundPodResource {
		fmt.Println("Pod resource not found in the cluster")
		return
	}

	// 创建 Core Client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 使用 Core Client 获取 Pods 并打印名称
	pods, err := clientset.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range pods.Items {
		fmt.Printf("Pod Name: %s\n", pod.GetName())
	}

}
