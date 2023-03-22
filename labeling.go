package main

import (
    "context"
    "fmt"
    "os"
    "time"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/labels"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)

func main() {
    // Create a Kubernetes clientset
    config, err := rest.InClusterConfig()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to get Kubernetes config: %v\n", err)
        os.Exit(1)
    }
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create Kubernetes clientset: %v\n", err)
        os.Exit(1)
    }

    // Create a context with a timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Get a list of all namespaces except kube-system
    nsList, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{
        LabelSelector: labels.SelectorFromSet(map[string]string{"name": "kube-system"}).String(),
    })
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to list namespaces: %v\n", err)
        os.Exit(1)
    }

    // Label each namespace with test=true
    for _, ns := range nsList.Items {
        if ns.Name == "kube-system" {
            continue
        }
        ns.Labels["test"] = "true"
        _, err := clientset.CoreV1().Namespaces().Update(ctx, &ns, metav1.UpdateOptions{})
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed to update namespace %s: %v\n", ns.Name, err)
            os.Exit(1)
        }
        fmt.Printf("Labeled namespace %s with test=true\n", ns.Name)
    }
}
