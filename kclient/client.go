package kclient

import (
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Kclient struct {
	namespace string
	Clientset *kubernetes.Clientset
}

func InitKclient(namespace string) *Kclient {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("This APP is needed active in Kubernetes Cluster\n", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}
	return &Kclient{namespace: namespace, Clientset: clientset}
}

func (k *Kclient) EditPodAnnotationsIP(name, ip string) error {
	pod, err := k.Clientset.CoreV1().Pods(k.namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Can't get %s pod ERROR:%v", name, err)
	}
	pod.Annotations["Ip"] = ip
	_, err = k.Clientset.CoreV1().Pods(k.namespace).Update(pod)
	if err != nil {
		fmt.Printf("Can't get %s pod ERROR:%v", name, err)
	}
	return err
}
