package kerberos

import (
	"encoding/json"
	"fmt"
	"time"
)

type TicketGrantingService struct {
	TGSKey      []byte
	ServiceKeys map[string][]byte
}

func (tgs *TicketGrantingService) HandleTGSRequest(req TGSRequest) (*TGSResponse, error) {
	fmt.Printf("[TGS] Handling ticket request for service: %s\n", req.ServiceName)

	tgtBytes, err := Decrypt(req.EncryptedTGT, tgs.TGSKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt TGT: %v", err)
	}

	var tgt Ticket
	if err := json.Unmarshal(tgtBytes, &tgt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TGT: %v", err)
	}

	authenticatorBytes, err := Decrypt(req.EncryptedAuthenticator, tgt.SessionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt authenticator: %v", err)
	}

	var authenticator Authenticator
	if err := json.Unmarshal(authenticatorBytes, &authenticator); err != nil {
		return nil, fmt.Errorf("failed to unmarshal authenticator: %v", err)
	}

	if authenticator.ClientName != tgt.ClientName {
		return nil, fmt.Errorf("client name mismatch")
	}

	if time.Now().After(tgt.Expiration) {
		return nil, fmt.Errorf("TGT expired")
	}

	serviceKey, ok := tgs.ServiceKeys[req.ServiceName]
	if !ok {
		return nil, fmt.Errorf("unknown service: %s", req.ServiceName)
	}

	serviceSessionKey := GenerateRandomKey()

	serviceTicket := Ticket{
		ClientName:  tgt.ClientName,
		ServiceName: req.ServiceName,
		SessionKey:  serviceSessionKey,
		Expiration:  time.Now().Add(1 * time.Hour),
	}

	serviceTicketBytes, err := json.Marshal(serviceTicket)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal service ticket: %v", err)
	}
	encryptedServiceTicket, err := Encrypt(serviceTicketBytes, serviceKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt service ticket: %v", err)
	}

	sessionPart := SessionPart{
		SessionKey:  serviceSessionKey,
		ServiceName: req.ServiceName,
	}

	sessionBytes, err := json.Marshal(sessionPart)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal session part: %v", err)
	}
	encryptedSession, err := Encrypt(sessionBytes, tgt.SessionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt session part: %v", err)
	}

	fmt.Printf("[TGS] Successfully issued service ticket for %s\n", req.ServiceName)

	return &TGSResponse{
		EncryptedServiceTicket: encryptedServiceTicket,
		EncryptedSessionPart:   encryptedSession,
	}, nil
}
