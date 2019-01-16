package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/appengine/log"
)

func LogDebugf(ctx context.Context, format string, a ...interface{}) {
	log.Debugf(ctx, format, a...)
}

func LogDebugfd(ctx context.Context, format string, a ...interface{}) {
	file, line := logctx(1)

	log.Debugf(ctx, "%s:%d %s", file, line, fmt.Sprintf(format, a...))
}

func LogDebugJSON(ctx context.Context, v interface{}) {
	if b, err := json.MarshalIndent(v, "", "\t"); err != nil {
		log.Debugf(ctx, "json.Marshal(): err=%v", err)
	} else {
		log.Debugf(ctx, string(b))
	}
}

func LogInfof(ctx context.Context, format string, a ...interface{}) {
	log.Infof(ctx, format, a...)
}

func LogInfofd(ctx context.Context, format string, a ...interface{}) {
	file, line := logctx(1)

	log.Infof(ctx, "%s:%d %s", file, line, fmt.Sprintf(format, a...))
}

func LogWarningf(ctx context.Context, format string, a ...interface{}) {
	log.Warningf(ctx, format, a...)
}

func LogWarningfd(ctx context.Context, format string, a ...interface{}) {
	file, line := logctx(1)

	log.Warningf(ctx, "%s:%d %s", file, line, fmt.Sprintf(format, a...))
}

func LogErrorf(ctx context.Context, format string, a ...interface{}) {
	log.Errorf(ctx, format, a...)
}

func LogErrorfd(ctx context.Context, format string, a ...interface{}) {
	file, line := logctx(1)

	log.Errorf(ctx, "%s:%d %s", file, line, fmt.Sprintf(format, a...))
}

func LogCriticalf(ctx context.Context, format string, a ...interface{}) {
	log.Criticalf(ctx, format, a...)
}

func LogCriticalfd(ctx context.Context, format string, a ...interface{}) {
	file, line := logctx(1)

	log.Criticalf(ctx, "%s:%d %s", file, line, fmt.Sprintf(format, a...))
}
