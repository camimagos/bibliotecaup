package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"bibliotecaup.com/cubicle/internal/controller/cubicle"
	metadatagateway "bibliotecaup.com/cubicle/internal/gateway/metadata/http"
	reservationgateway "bibliotecaup.com/cubicle/internal/gateway/reservation/http"
	httphandler "bibliotecaup.com/cubicle/internal/handler/http"
	"bibliotecaup.com/pkg/discovery/consul"
	discovery "bibliotecaup.com/pkg/registry"
)

const serviceName = "cubicle"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "Port to run the HTTP server on")
	flag.Parse()
	log.Printf("Starting Cubicle service on port %d...", port)

	registry, err := consul.NewRegistry(os.Getenv("CONSUL_HTTP_ADDR"))
	if err != nil {
		log.Fatalf("Error creating Consul registry: %v", err)
	}
	// registry, err := consul.NewRegistry("localhost:8500")
	// if err != nil {
	// 	panic(err)
	// }
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
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
	metadataGateway := metadatagateway.New(registry)
	reservationgateway := reservationgateway.New(registry)
	ctrl := cubicle.New(reservationgateway, metadataGateway)
	h := httphandler.New(ctrl)
	http.Handle("/cubicle", http.HandlerFunc(h.GetCubicleDetails))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
