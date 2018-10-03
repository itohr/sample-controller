/*
Copyright 2018 NTT corp..

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

func main() {
	mrg, err := builder.SimpleController().
		ForType(&appsv1.ReplicaSet{}).
		Owns(&corev1.Pod{}).
		Build(&ReplicaSetController{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting the Cmd.")

	log.Fatal(mrg.Start(signals.SetupSignalHandler()))
}

type ReplicaSetController struct {
	client.Client
}

func (a *ReplicaSetController) InjectClient(c client.Client) error {
	a.Client = c
	return nil
}

func (a *ReplicaSetController) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	rs := &appsv1.ReplicaSet{}
	err := a.Get(context.TODO(), req.NamespacedName, rs)
	if err != nil {
		return reconcile.Result{}, err
	}

	pods := &corev1.PodList{}
	err = a.List(context.TODO(), client.InNamespace(req.Namespace).MatchingLabels(rs.Spec.Template.Labels), pods)
	if err != nil {
		return reconcile.Result{}, err
	}

	rs.Labels["selector-pod-count"] = fmt.Sprintf("%v", len(pods.Items))
	err = a.Update(context.TODO(), rs)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

//func main() {
//	// Get a config to talk to the apiserver
//	cfg, err := config.GetConfig()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Create a new Cmd to provide shared dependencies and start components
//	mgr, err := manager.New(cfg, manager.Options{})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	log.Printf("Registering Components.")
//
//	// Setup Scheme for all resources
//	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
//		log.Fatal(err)
//	}
//
//	// Setup all Controllers
//	if err := controller.AddToManager(mgr); err != nil {
//		log.Fatal(err)
//	}
//
//	log.Printf("Starting the Cmd.")
//
//	// Start the Cmd
//	log.Fatal(mgr.Start(signals.SetupSignalHandler()))
//}
