package kerberos

import (
	"encoding/json"
	"fmt"
	"time"
)

type ApplicationService struct {
	Name string
	Key  []byte
}

func (s *ApplicationService) HandleRequest(req APRequest) (string, error) {
	fmt.Printf("[%s] Received access request\n", s.Name)

	serviceTicketBytes, err := Decrypt(req.EncryptedServiceTicket, s.Key)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt service ticket: %v", err)
	}

	var serviceTicket Ticket
	if err := json.Unmarshal(serviceTicketBytes, &serviceTicket); err != nil {
		return "", fmt.Errorf("failed to unmarshal service ticket: %v", err)
	}

	authenticatorBytes, err := Decrypt(req.EncryptedAuthenticator, serviceTicket.SessionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt authenticator: %v", err)
	}

	var authenticator Authenticator
	if err := json.Unmarshal(authenticatorBytes, &authenticator); err != nil {
		return "", fmt.Errorf("failed to unmarshal authenticator: %v", err)
	}

	if authenticator.ClientName != serviceTicket.ClientName {
		return "", fmt.Errorf("client identity mismatch")
	}

	if time.Now().After(serviceTicket.Expiration) {
		return "", fmt.Errorf("service ticket expired")
	}

	fmt.Printf("[%s] Client %s authenticated successfully\n", s.Name, authenticator.ClientName)

	return fmt.Sprintf("Hello %s! Here is your protected data from %s.", authenticator.ClientName, s.Name), nil
}
