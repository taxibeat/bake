package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/taxibeat/bake/docker/component/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/Shopify/sarama.v1"
)

func main() {
	redisAddr := os.Getenv("REDIS")
	redisClient := redis.NewClient(redisAddr)
	_, err := redisClient.Set(context.Background(), "testservice", "foo", time.Second).Result()
	if err != nil {
		log.Fatal(err)
	}

	mongoAddr := os.Getenv("MONGO")

	opts := options.Client()
	rs := "rs0"
	opts.ReplicaSet = &rs
	opts.ApplyURI("mongodb://" + mongoAddr)
	mongoClient, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatal(err)
	}

	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	kafkaAddr := os.Getenv("KAFKA")
	kafkaClient, err := sarama.NewClient([]string{kafkaAddr}, nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = kafkaClient.Topics()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", Health)
	port := os.Getenv("PORT")
	fmt.Println("Running on port:", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// Health is a simple health endpoint.
func Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}