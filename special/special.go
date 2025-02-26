package special

import (
	"context"

	"github.com/mt1976/frantic-aegis/app/web/security"
)

func GetUserCode(ctx context.Context) string {
	return security.Current_UserCode(ctx)
}

func GetUserKey(ctx context.Context) string {
	return security.Current_UserKey(ctx)
}
