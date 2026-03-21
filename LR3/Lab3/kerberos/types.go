package kerberos

import "time"

type Ticket struct {
	ClientName  string
	ServiceName string
	SessionKey  []byte
	Expiration  time.Time
}

type Authenticator struct {
	ClientName string
	Timestamp  time.Time
}

type ASRequest struct {
	ClientName  string
	ServiceName string
}

type ASResponse struct {
	EncryptedTGT         []byte
	EncryptedSessionPart []byte
}

type TGSRequest struct {
	EncryptedTGT           []byte
	EncryptedAuthenticator []byte
	ServiceName            string
}

type TGSResponse struct {
	EncryptedServiceTicket []byte
	EncryptedSessionPart   []byte
}

type APRequest struct {
	EncryptedServiceTicket []byte
	EncryptedAuthenticator []byte
}

type SessionPart struct {
	SessionKey  []byte
	ServiceName string
}
