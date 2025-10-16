package main

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// Create in-cluster config (works inside a pod)
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error creating in-cluster config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}

	namespace := "default"
	serviceName := "my-service"

	// Get the Service
	service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Error fetching service: %v", err)
	}

	// Extract the selector
	selector := service.Spec.Selector
	if len(selector) == 0 {
		log.Fatalf("Service %s has no selectors", serviceName)
	}

	// Convert map to label selector string
	labelSelector := metav1.FormatLabelSelector(&metav1.LabelSelector{MatchLabels: selector})

	// Get pods matching the selector
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Fatalf("Error listing pods: %v", err)
	}

	fmt.Printf("Pods for service %s:\n", serviceName)
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}
}
