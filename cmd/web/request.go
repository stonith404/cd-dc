package web

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func NewRequestContext(serviceName string, r *http.Request) context.Context {
	requestId := generateRequestID()
	return context.WithValue(context.Background(), "RequestContext", RequestContext{
		ServiceName: serviceName,
		RequestID:   requestId,
		Logger:      log.New(log.Writer(), fmt.Sprintf("[%s-%s] ", serviceName, requestId), log.LstdFlags),
	})
}

// Random 5 character string
func generateRequestID() string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 5)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

type RequestContext struct {
	ServiceName string
	RequestID   string
	Logger      *log.Logger
}
