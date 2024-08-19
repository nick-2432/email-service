package main

import (
	"emailses/email"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	// Initialize the SES configuration
	config := email.SESConfig{
		Sender:          "nikhilanand2432@gmail.com",
		Receiver:        "nikhilanand2432@gmail.com",
		Subject:         "Sample Email Subject",
		TemplatePath:    "templates/email-template.html",
		Region:          "eu-north-1",
		AccessKeyId:     "--",
		SecretAccessKey: "------",
	}

	// Initialize the email service with the configuration
	emailService := email.NewEmailService(config)

	// Initialize the handler with the email service
	emailHandler := email.NewHandler(emailService)

	// Set up the router
	r := mux.NewRouter()
	emailHandler.Init(r)

	// Start the server
	fmt.Println("Server is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
