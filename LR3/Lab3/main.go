package main

import (
	"kerberos-go/kerberos"
)

func main() {
	clientName := "dmitry"
	clientPassword := "password123"
	serviceName := "fileserver"

	clientKey := kerberos.DeriveKey(clientPassword)
	tgsKey := kerberos.GenerateRandomKey()
	serviceKey := kerberos.GenerateRandomKey()

	authenticationService := &kerberos.AuthenticationService{
		ClientKeys: map[string][]byte{
			clientName: clientKey,
		},
		TGSKey: tgsKey,
	}

	ticketGrantingService := &kerberos.TicketGrantingService{
		TGSKey: tgsKey,
		ServiceKeys: map[string][]byte{
			serviceName: serviceKey,
		},
	}

	applicationService := &kerberos.ApplicationService{
		Name: serviceName,
		Key:  serviceKey,
	}

	client := &kerberos.Client{
		Name:     clientName,
		Password: clientPassword,
	}

	client.RunFullFlow(authenticationService, ticketGrantingService, applicationService)
}
