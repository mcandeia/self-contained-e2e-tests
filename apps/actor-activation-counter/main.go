/*
Copyright 2022 The Dapr Authors
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
	"log"
	"net/http"
	"sync/atomic"

	"github.com/dapr/go-sdk/actor"
	dapr "github.com/dapr/go-sdk/client"
	daprd "github.com/dapr/go-sdk/service/http"
)

func testActorFactory() actor.Server {
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	return &TestActor{
		daprClient: client,
	}
}

var counter = atomic.Int64{}

type TestActor struct {
	actor.ServerImplBase
	daprClient dapr.Client
}

func (t *TestActor) Type() string {
	return "fake-actor-type"
}

// user defined functions
func (t *TestActor) Inc(ctx context.Context, req any) (any, error) {
	counter.Add(1)
	return "ok", nil
}

func (t *TestActor) Get(ctx context.Context, req any) (int, error) {
	return int(counter.Load()), nil
}

func main() {
	s := daprd.NewService(":8080")
	s.RegisterActorImplFactory(testActorFactory)
	log.Println("started")
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error listenning: %v", err)
	}
}
