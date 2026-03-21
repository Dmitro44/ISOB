package kerberos

import (
	"encoding/json"
	"fmt"
	"time"
)

type Client struct {
	Name     string
	Password string
}

func (c *Client) RunFullFlow(as *AuthenticationService, tgs *TicketGrantingService, appService *ApplicationService) {
	fmt.Printf("\n[Client] Starting Kerberos authentication for %s\n", c.Name)
	fmt.Printf("[Client] Target service: %s\n\n", appService.Name)

	clientKey := DeriveKey(c.Password)

	fmt.Println("--- PHASE 1: AS EXCHANGE ---")
	asRequest := ASRequest{
		ClientName:  c.Name,
		ServiceName: "TGS",
	}

	asResponse, err := as.HandleASRequest(asRequest)
	if err != nil {
		fmt.Printf("[Client] AS Exchange failed: %v\n", err)
		return
	}

	sessionPartBytes, err := Decrypt(asResponse.EncryptedSessionPart, clientKey)
	if err != nil {
		fmt.Printf("[Client] Failed to decrypt AS response: %v\n", err)
		return
	}

	var tgsSessionPart SessionPart
	if err := json.Unmarshal(sessionPartBytes, &tgsSessionPart); err != nil {
		fmt.Printf("[Client] Failed to parse session part: %v\n", err)
		return
	}

	fmt.Printf("[Client] Successfully obtained TGT\n\n")

	fmt.Println("--- PHASE 2: TGS EXCHANGE ---")
	tgsAuthenticator := Authenticator{
		ClientName: c.Name,
		Timestamp:  time.Now(),
	}

	tgsAuthenticatorBytes, err := json.Marshal(tgsAuthenticator)
	if err != nil {
		fmt.Printf("[Client] Failed to marshal authenticator: %v\n", err)
		return
	}
	encryptedTGSAuthenticator, err := Encrypt(tgsAuthenticatorBytes, tgsSessionPart.SessionKey)
	if err != nil {
		fmt.Printf("[Client] Failed to encrypt authenticator: %v\n", err)
		return
	}

	tgsRequest := TGSRequest{
		EncryptedTGT:           asResponse.EncryptedTGT,
		EncryptedAuthenticator: encryptedTGSAuthenticator,
		ServiceName:            appService.Name,
	}

	tgsResponse, err := tgs.HandleTGSRequest(tgsRequest)
	if err != nil {
		fmt.Printf("[Client] TGS Exchange failed: %v\n", err)
		return
	}

	serviceSessionPartBytes, err := Decrypt(tgsResponse.EncryptedSessionPart, tgsSessionPart.SessionKey)
	if err != nil {
		fmt.Printf("[Client] Failed to decrypt TGS response: %v\n", err)
		return
	}

	var serviceSessionPart SessionPart
	if err := json.Unmarshal(serviceSessionPartBytes, &serviceSessionPart); err != nil {
		fmt.Printf("[Client] Failed to parse service session part: %v\n", err)
		return
	}

	fmt.Printf("[Client] Successfully obtained service ticket for %s\n\n", appService.Name)

	fmt.Println("--- PHASE 3: AP EXCHANGE ---")
	serviceAuthenticator := Authenticator{
		ClientName: c.Name,
		Timestamp:  time.Now(),
	}

	serviceAuthenticatorBytes, err := json.Marshal(serviceAuthenticator)
	if err != nil {
		fmt.Printf("[Client] Failed to marshal service authenticator: %v\n", err)
		return
	}
	encryptedServiceAuthenticator, err := Encrypt(serviceAuthenticatorBytes, serviceSessionPart.SessionKey)
	if err != nil {
		fmt.Printf("[Client] Failed to encrypt service authenticator: %v\n", err)
		return
	}

	apRequest := APRequest{
		EncryptedServiceTicket: tgsResponse.EncryptedServiceTicket,
		EncryptedAuthenticator: encryptedServiceAuthenticator,
	}

	response, err := appService.HandleRequest(apRequest)
	if err != nil {
		fmt.Printf("[Client] Service access failed: %v\n", err)
		return
	}

	fmt.Printf("[Client] Successfully authenticated to %s\n", appService.Name)
	fmt.Printf("[Client] Received response: %s\n\n", response)
}
