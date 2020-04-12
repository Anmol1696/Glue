package main

import (
	"context"
)

var (
	RequestIdCtxKey = &contextKey{"RequestId"}
	IdentityCtxKey  = &contextKey{"IdentityId"}
)

const (
	prog        = "glue"
	version     = "v0"
	description = "is a test app for initiating creation"
	envPrefix   = "GLUE_"
)

// copied and modified from net/http/http.go
// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

func (k *contextKey) String() string { return prog + " context value " + k.name }

// Extracts request id from context
func getRequestId(ctx context.Context) string {
	id, ok := ctx.Value(RequestIdCtxKey).(string)
	if !(ok) {
		return "no-request-id"
	}
	return id
}
