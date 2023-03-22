package main

import (
	"context"
	"fmt"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// create a Kubernetes client using in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// define a label selector to find completed and evicted pods
	labelSelector := metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      "status.phase",
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{"Succeeded", "Failed", "Completed", "Error"},
			},
		},
	}

	// define a pod deletion policy
	policy := metav1.DeletePropagationForeground

	// create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// get the list of pods to delete
	podList, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// delete each pod in the list
	for _, pod := range podList.Items {
		err := clientset.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{
			PropagationPolicy: &policy,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Deleted pod %s/%s\n", pod.Namespace, pod.Name)
	}
}

---
func main() {
	// create a Kubernetes client
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
			panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
			panic(err.Error())
	}

	// delete pods with status "Completed", "Evicted", and "Error"
	deletePodsWithStatus(clientset, "Completed")
	deletePodsWithStatus(clientset, "Evicted")
	deletePodsWithStatus(clientset, "Error")
}

func deletePodsWithStatus(clientset kubernetes.Interface, status string) {
	// create a pod client
	podClient := clientset.CoreV1().Pods(v1.NamespaceAll)

	// list all pods with the given status
	pods, err := podClient.List(context.TODO(), v1.ListOptions{
			FieldSelector: fmt.Sprintf("status.phase=%s", status),
	})
	if err != nil {
			panic(errors.Wrapf(err, "failed to list pods with status %s", status).Error())
	}

	// delete each pod
	for _, pod := range pods.Items {
			err := podClient.Delete(context.TODO(), pod.Name, v1.DeleteOptions{})
			if errors.IsNotFound(err) {
					continue
			}
			if err != nil {
					panic(errors.Wrapf(err, "failed to delete pod %s", pod.Name).Error())
			}
			fmt.Printf("Deleted pod %s with status %s\n", pod.Name, status)
	}
}
