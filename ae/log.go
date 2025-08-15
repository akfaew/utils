package ae

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/akfaew/utils"
	"github.com/akfaew/utils/xctc"
)

var (
	trimprefix = ""
	projectID  = os.Getenv("GOOGLE_CLOUD_PROJECT")
)

type Log struct {
	ctx         context.Context
	httpRequest map[string]any
}

func NewLog(ctx context.Context) *Log {
	return &Log{ctx: ctx}
}

// NewLogFromRequest extracts HTTP request metadata into a top-level httpRequest
// field recognized by Cloud Logging. App Engine's nginx logs populate
// protoPayload, but logs written via stdout can't set protoPayload directly.
// httpRequest provides similar functionality in Logs Explorer.
func NewLogFromRequest(r *http.Request) *Log {
	log := &Log{ctx: r.Context()}
	req := map[string]any{
		"requestMethod": r.Method,
		"requestUrl":    r.URL.String(),
		"userAgent":     r.UserAgent(),
	}
	if ref := r.Referer(); ref != "" {
		req["referer"] = ref
	}
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		req["remoteIp"] = ip
	} else if r.RemoteAddr != "" {
		req["remoteIp"] = r.RemoteAddr
	}
	log.httpRequest = req
	return log
}

// Set trimprefix to the path to the source code directory, so that we only log the filename and not the full path.
func init() {
	_, path, _, _ := runtime.Caller(1)

	trimprefix = filepath.Dir(path) + string(filepath.Separator)
}

func logctx(skip int) (file string, line int) {
	_, file, line, _ = runtime.Caller(skip + 1)
	file = strings.TrimPrefix(file, trimprefix)

	return
}

type entry struct {
	Trace       string         `json:"logging.googleapis.com/trace,omitempty"`
	SpanID      string         `json:"logging.googleapis.com/spanId,omitempty"`
	Data        map[string]any `json:"data"`
	Message     string         `json:"message,omitempty"`
	Severity    string         `json:"severity,omitempty"`
	HTTPRequest map[string]any `json:"httpRequest,omitempty"`
}

func (log *Log) write(severity, msg string) {
	e := entry{
		Data:        map[string]any{},
		Message:     msg,
		Severity:    severity,
		HTTPRequest: log.httpRequest,
	}

	if log.ctx != nil {
		if x := xctc.XCTC(log.ctx); x != "" {
			if trace, err := utils.ParseTraceParent(x); err == nil {
				e.Trace = fmt.Sprintf("projects/%s/traces/%s", projectID, trace.TraceID)
				e.SpanID = trace.ParentID
			}
		}
	}

	if b, err := json.Marshal(e); err == nil {
		_, _ = os.Stdout.Write(append(b, '\n'))
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "{\"severity\":\"ERROR\",\"message\":\"json.Marshal(): %v\"}\n", err)
	}
}

func (log *Log) Debugf(format string, a ...any) {
	log.write("DEBUG", fmt.Sprintf(format, a...))
}

func (log *Log) Debugfd(format string, a ...any) {
	file, line := logctx(1)

	log.write("DEBUG", fmt.Sprintf("%s:%d %s", file, line, fmt.Sprintf(format, a...)))
}

func (log *Log) DebugJSON(v any) {
	if b, err := json.MarshalIndent(v, "", "\t"); err != nil {
		log.write("DEBUG", fmt.Sprintf("json.MarshalIndent(): err=%v", err))
	} else {
		log.write("DEBUG", string(b))
	}
}

func (log *Log) Infof(format string, a ...any) {
	log.write("INFO", fmt.Sprintf(format, a...))
}

func (log *Log) Infofd(format string, a ...any) {
	file, line := logctx(1)

	log.write("INFO", fmt.Sprintf("%s:%d %s", file, line, fmt.Sprintf(format, a...)))
}

func (log *Log) Warningf(format string, a ...any) {
	log.write("WARNING", fmt.Sprintf(format, a...))
}

func (log *Log) Warningfd(format string, a ...any) {
	file, line := logctx(1)

	log.write("WARNING", fmt.Sprintf("%s:%d %s", file, line, fmt.Sprintf(format, a...)))
}

func (log *Log) Errorf(format string, a ...any) {
	log.write("ERROR", fmt.Sprintf(format, a...))
}

func (log *Log) Errorfd(format string, a ...any) {
	file, line := logctx(1)

	log.write("ERROR", fmt.Sprintf("%s:%d %s", file, line, fmt.Sprintf(format, a...)))
}
