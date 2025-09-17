package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"bibliotecaup.com/pkg/discovery/consul"
	discovery "bibliotecaup.com/pkg/registry"
	"bibliotecaup.com/reservation/internal/controller/reservation"
	httpHandler "bibliotecaup.com/reservation/internal/handler/http"
	"bibliotecaup.com/reservation/internal/repository/memory"
)

const serviceName = "reservation"

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()
	log.Printf("Starting reservation service on port %d", port)
	registry, err := consul.NewRegistry(os.Getenv("CONSUL_HTTP_ADDR"))
	if err != nil {
		log.Fatalf("Error creating Consul registry: %v", err)
	}
	// registry, err := consul.NewRegistry("localhost:8500")
	// if err != nil {
	// 	log.Printf("Error creating Consul registry: %v", err)
	// 	panic(err)
	// }
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("reservation:%d", port)); err != nil {
		log.Printf("Error registering service in Consul: %v", err)
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)
	repo := memory.New()
	ctrl := reservation.New(repo)
	h := httpHandler.New(ctrl)
	http.Handle("/reservation", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Printf("Error starting HTTP server: %v", err)
		panic(err)
	}
}
