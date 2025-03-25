package security

import (
	"context"

	"github.com/mt1976/frantic-core/messageHelpers"
)

// UserMessageFromContext builds a UserMessage from the current context
func UserMessageFromContext(ctx context.Context) messageHelpers.UserMessage {
	resp := messageHelpers.UserMessage{}
	resp.Key = Current_UserKey(ctx)
	resp.Code = Current_UserCode(ctx)
	resp.Locale = Current_UserLocale(ctx)
	resp.Theme = Current_SessionTheme(ctx)
	resp.Timezone = Current_SessionTimezone(ctx)
	resp.Source = ""
	resp.Spare1 = ""
	resp.Spare2 = ""
	return resp

}
