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
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bsm/redislock"
	"github.com/dapr/go-sdk/actor"
	dapr "github.com/dapr/go-sdk/client"
	daprd "github.com/dapr/go-sdk/service/http"
	redis "github.com/go-redis/redis/v9"
)

func testActorFactory(redisClient *redis.Client) func() actor.Server {
	return func() actor.Server {
		client, err := dapr.NewClient()
		if err != nil {
			panic(err)
		}
		return &TestActor{
			daprClient: client,
			locker:     redislock.New(redisClient),
		}
	}
}

type TestActor struct {
	actor.ServerImplBase
	daprClient dapr.Client
	locker     *redislock.Client
}

func (t *TestActor) Type() string {
	return "fake-actor-type"
}

// user defined functions
func (t *TestActor) Lock(ctx context.Context, req any) (any, error) {
	lockTimeout := time.Second
	// Try to obtain lock.
	lock, err := t.locker.Obtain(ctx, t.ID(), lockTimeout, nil)
	if err == redislock.ErrNotObtained {
		return nil, errors.New("resource was locked!")
	}

	if err == nil {
		if releaseErr := lock.Release(ctx); releaseErr != nil {
			time.Sleep(lockTimeout) // sleep to make sure that the lock will be automatically released
		}
	}
	return "succeed", nil
}

func main() {
	// Connect to redis.
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    os.Getenv("REDIS_SVC_URL"),
	})
	defer client.Close()

	s := daprd.NewService(":8080")
	s.RegisterActorImplFactory(testActorFactory(client))
	log.Println("started")
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error listenning: %v", err)
	}
}
