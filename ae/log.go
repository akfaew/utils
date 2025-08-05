package ae

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/akfaew/utils/ae/stackdriver-gae-logrus-plugin"
)

var trimprefix = ""

type Log struct {
	*logrus.Entry
}

func NewLog(ctx context.Context) *Log {
	return &Log{
		logrus.WithContext(ctx),
	}
}

// Set trimprefix to the path to the source code directory, so that we only log the filename and not the full path.
func init() {
	_, path, _, _ := runtime.Caller(1)

	trimprefix = filepath.Dir(path) + string(filepath.Separator)

	formatter := stackdriver.GAEStandardFormatter(
		stackdriver.WithProjectID(os.Getenv("GOOGLE_CLOUD_PROJECT")),
	)
	logrus.SetFormatter(formatter)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel) // Log the debug severity or above.
}

func logctx(skip int) (file string, line int) {
	_, file, line, _ = runtime.Caller(skip + 1)
	file = strings.TrimPrefix(file, trimprefix)

	return
}

func (log *Log) Debugf(format string, a ...any) {
	log.Entry.Debugf(format, a...)
}

func (log *Log) Debugfd(format string, a ...any) {
	file, line := logctx(1)

	log.Entry.Debugf("%s:%d %s", file, line, fmt.Sprintf(format, a...))
}

func (log *Log) DebugJSON(v any) {
	if b, err := json.MarshalIndent(v, "", "\t"); err != nil {
		log.Entry.Debugf("json.MarshalIndent(): err=%v", err)
	} else {
		log.Entry.Debugf("%s", string(b))
	}
}

func (log *Log) Infof(format string, a ...any) {
	log.Entry.Infof(format, a...)
}

func (log *Log) Infofd(format string, a ...any) {
	file, line := logctx(1)

	log.Entry.Infof("%s:%d %s", file, line, fmt.Sprintf(format, a...))
}

func (log *Log) Warningf(format string, a ...any) {
	log.Entry.Warningf(format, a...)
}

func (log *Log) Warningfd(format string, a ...any) {
	file, line := logctx(1)

	log.Entry.Warningf("%s:%d %s", file, line, fmt.Sprintf(format, a...))
}

func (log *Log) Errorf(format string, a ...any) {
	log.Entry.Errorf(format, a...)
}

func (log *Log) Errorfd(format string, a ...any) {
	file, line := logctx(1)

	log.Entry.Errorf("%s:%d %s", file, line, fmt.Sprintf(format, a...))
}
