package models

import "github.com/clerk/clerk-sdk-go/v2"

type ClerkUser struct {
	Id        string `json:"id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Username  string `json:"username"`
	EMail     string `json:"email"`
}

type ClerkWebhook struct {
	Type string    `json:"type"`
	Data ClerkUser `json:"data"`
}

type ClerkErrorResponse struct {
	Errors       []clerk.Error `json:"errors"`
	ClerkTraceId string        `json:"clerk_trace_id"`
}
