package security

import (
	"context"

	"github.com/mt1976/frantic-core/contextHandler"
)

// PageWranlger - Get the SessionID and SessionKey from the context and provide to the page being built by the source application
// Returns the SessionKey Identifier and SessionKey and an error if one occurs (reserved for future use)
func PageWranlger(ctx context.Context) (string, string, error) {
	return SESSION_KEY, contextHandler.GetSessionID(ctx), nil
}
