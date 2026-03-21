package kerberos

import (
	"encoding/json"
	"fmt"
	"time"
)

type AuthenticationService struct {
	ClientKeys map[string][]byte
	TGSKey     []byte
}

func (as *AuthenticationService) HandleASRequest(req ASRequest) (*ASResponse, error) {
	fmt.Printf("[AS] Handling authentication request from %s\n", req.ClientName)

	clientKey, ok := as.ClientKeys[req.ClientName]
	if !ok {
		return nil, fmt.Errorf("client not found: %s", req.ClientName)
	}

	tgsSessionKey := GenerateRandomKey()

	tgt := Ticket{
		ClientName:  req.ClientName,
		ServiceName: "TGS",
		SessionKey:  tgsSessionKey,
		Expiration:  time.Now().Add(8 * time.Hour),
	}

	tgtBytes, err := json.Marshal(tgt)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TGT: %v", err)
	}
	encryptedTGT, err := Encrypt(tgtBytes, as.TGSKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt TGT: %v", err)
	}

	sessionPart := SessionPart{
		SessionKey:  tgsSessionKey,
		ServiceName: "TGS",
	}

	sessionBytes, err := json.Marshal(sessionPart)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal session part: %v", err)
	}
	encryptedSession, err := Encrypt(sessionBytes, clientKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt session part: %v", err)
	}

	fmt.Printf("[AS] Successfully issued TGT for %s\n", req.ClientName)

	return &ASResponse{
		EncryptedTGT:         encryptedTGT,
		EncryptedSessionPart: encryptedSession,
	}, nil
}
