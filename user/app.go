package user

import (
	"context"
	"fmt"
	"net/http"
	"crypto/tls"
)

// Configuration holds application configuration
type Configuration struct {
	// Define your configuration fields here
}

// Application represents the main application struct
type Application struct {
	Configuration *Configuration
	Router        *http.ServeMux
}

// NewApplication creates an instance of the app and returns an app struct
func NewApplication(ctx context.Context, configuration *Configuration) (*Application, error) {
	router := http.NewServeMux()
	application := &Application{
		Configuration: configuration,
		Router:        router,
	}

	return application, nil
}

// HomeHandler is a simple handler function for the home route
func (application *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

// Run starts the application
func (application *Application) Run() {
	fmt.Println("Starting Application")

	// Register your handlers
	// application.Router.HandleFunc("/", application.HomeHandler)
	tlsConfig := &tls.Config{
		// You need to provide a valid certificate and private key
        	// You can obtain a certificate from a certificate authority (CA)
       	 	// or use a self-signed certificate for development/testing purposes
        	// CertFile: "path/to/cert.pem",
        	// KeyFile:  "path/to/key.pem",
	}
	// tlsConfig := &http.Transport{
	// 	TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	// }

	server := &http.Server{
		Addr:    ":8080", // Provide the desired port number
		Handler: application.Router,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Starting Server")

	err := server.ListenAndServeTLS("path/to/cert.pem", "path/to/key.pem")
	if err != nil {
		fmt.Println("Failed to start the server:", err)
	}
}