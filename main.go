package main

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	kubeconfig, err := ctrl.GetConfig()
	if err != nil {
		fmt.Printf("failed to get kubeconfig: %v\n", err)
		os.Exit(1)
	}

	kubeclient, err := client.New(kubeconfig, client.Options{})
	if err != nil {
		fmt.Printf("failed to create client: %v\n", err)
		os.Exit(1)
	}

	oldSA := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tester",
			Namespace: "default",
		},
		ImagePullSecrets: []corev1.LocalObjectReference{
			{
				Name: "key1",
			},
		},
	}

	if err := kubeclient.Create(context.TODO(), oldSA); err != nil {
		fmt.Printf("failed to create sa: %v\n", err)
		os.Exit(1)
	}

	existingResource := &corev1.ServiceAccount{}
	key := types.NamespacedName{
		Name:      "tester",
		Namespace: "default",
	}
	if err := kubeclient.Get(context.TODO(), key, existingResource); err != nil {
		fmt.Printf("failed to fetch serviceaccount: %v\n", err)
		os.Exit(1)
	}

	existingResource.ImagePullSecrets = []corev1.LocalObjectReference{
		{
			Name: "key1",
		},
		{
			Name: "key2",
		},
	}

	if err := kubeclient.Update(context.TODO(), existingResource); err != nil {
		fmt.Printf("failed to update serviceaccount: %v\n", err)
		os.Exit(1)
	}
}
