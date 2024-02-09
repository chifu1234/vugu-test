package main

import (
	"context"
	"fmt"
	"github.com/vugu/vugu"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Root struct {
	List       corev1.PodList `vugu:"data"`
	IsLoading  bool           `vugu:"data"`
	Namespace  string         `vugu:"data"`
	Namespaces []string       `vugu:"data"`
	client     client.Client
}

func (c *Root) Init(vugu.InitCtx) {
	fmt.Println("init")
	var err error
	cfg := &rest.Config{
		Host: "localhost:8080",
	}
	// Create a new client to interact with Kubernetes objects.
	c.client, err = client.New(cfg, client.Options{})
	if err != nil {
		panic(fmt.Errorf("error creating client: %v", err))
	}

}

func (c *Root) UpdateNamesapce(event vugu.DOMEvent) {
	c.Namespace = event.PropString("target", "value")
	c.UpdateData(event)

}

// getData
func (c *Root) UpdateData(event vugu.DOMEvent) {
	ctx := context.Background()
	ee := event.EventEnv()
	go func() {
		// Create an empty Pod object to store the result.
		podList := &corev1.PodList{}

		if err := c.client.List(ctx, podList, client.InNamespace(c.Namespace)); err != nil {
			fmt.Printf("failed to get Pod: %v", err)
		}
		ee.Lock()
		// reset
		c.List = corev1.PodList{}
		c.List = *podList

		ee.UnlockRender()
	}()
}
