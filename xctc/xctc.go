package xctc

import (
	"context"
)

type xctcKeyType string

const XctcKey xctcKeyType = "xctc"

// XCTC returns the XCloudTraceContent value from the context.
func XCTC(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	xctc, ok := ctx.Value(XctcKey).(string)
	if !ok {
		return ""
	}
	return xctc
}
