package main

import (
	"cart/configs"
	"cart/controllers"
	"cart/grpc"
	"cart/services"
	"flag"
	"fmt"
)

func main() {
	// Parse command line flags
	serverType := flag.String("server", "both", "Server type to run (grpc, http, or both)")
	flag.Parse()

	// Initialize environment configuration
	fmt.Println("Initializing Cart Service...")
	configs.InitEnv()

	// Initialize Redis client
	redis := services.GetRedisClient()
	defer redis.Close()

	// Start servers based on the server type
	switch *serverType {
	case "grpc":
		fmt.Println("Starting gRPC server only...")
		grpc.Init()
	case "http":
		fmt.Println("Starting HTTP server only...")
		controllers.InitHTTPServer()
	case "both":
		fmt.Println("Starting both gRPC and HTTP servers...")
		// Start the HTTP server in a goroutine
		go controllers.InitHTTPServer()
		// Start the gRPC server (this call is blocking)
		grpc.Init()
	default:
		fmt.Printf("Invalid server type: %s. Valid options are 'grpc', 'http', or 'both'\n", *serverType)
	}
}
