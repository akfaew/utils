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
	data        map[string]any
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

func (log *Log) WithField(name string, value any) *Log {
	l := *log
	if l.data == nil {
		l.data = map[string]any{}
	} else {
		data := make(map[string]any, len(l.data)+1)
		for k, v := range l.data {
			data[k] = v
		}
		l.data = data
	}
	l.data[name] = value
	return &l
}

// Set trimprefix to the path to the source code directory, so that we only log the filename and not the full path.
func init() {
	_, path, _, _ := runtime.Caller(1)

	trimprefix = filepath.Dir(path) + string(filepath.Separator)
}

func logctx(skip int) (file string, line int, function string) {
	pc, file, line, _ := runtime.Caller(skip + 1)
	file = strings.TrimPrefix(file, trimprefix)
	if fn := runtime.FuncForPC(pc); fn != nil {
		function = fn.Name()
	}

	return
}

type sourceLocation struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function,omitempty"`
}

type entry struct {
	Trace          string          `json:"logging.googleapis.com/trace,omitempty"`
	SpanID         string          `json:"logging.googleapis.com/spanId,omitempty"`
	SourceLocation *sourceLocation `json:"logging.googleapis.com/sourceLocation,omitempty"`
	Data           map[string]any  `json:"data"`
	Message        string          `json:"message,omitempty"`
	Severity       string          `json:"severity,omitempty"`
	HTTPRequest    map[string]any  `json:"httpRequest,omitempty"`
}

func (log *Log) write(severity, msg string, sl *sourceLocation) {
	e := entry{
		Data:           map[string]any{},
		Message:        msg,
		Severity:       severity,
		HTTPRequest:    log.httpRequest,
		SourceLocation: sl,
	}

	for k, v := range log.data {
		e.Data[k] = v
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
	log.write("DEBUG", fmt.Sprintf(format, a...), nil)
}

func (log *Log) Debugfd(format string, a ...any) {
	file, line, function := logctx(1)
	sl := &sourceLocation{File: file, Line: line, Function: function}

	log.write("DEBUG", fmt.Sprintf("%s:%d %s", file, line, fmt.Sprintf(format, a...)), sl)
}

func (log *Log) DebugJSON(v any) {
	if b, err := json.MarshalIndent(v, "", "\t"); err != nil {
		log.write("DEBUG", fmt.Sprintf("json.MarshalIndent(): err=%v", err), nil)
	} else {
		log.write("DEBUG", string(b), nil)
	}
}

func (log *Log) Infof(format string, a ...any) {
	log.write("INFO", fmt.Sprintf(format, a...), nil)
}

func (log *Log) Infofd(format string, a ...any) {
	file, line, function := logctx(1)
	sl := &sourceLocation{File: file, Line: line, Function: function}

	log.write("INFO", fmt.Sprintf("%s:%d %s", file, line, fmt.Sprintf(format, a...)), sl)
}

func (log *Log) Warningf(format string, a ...any) {
	log.write("WARNING", fmt.Sprintf(format, a...), nil)
}

func (log *Log) Warningfd(format string, a ...any) {
	file, line, function := logctx(1)
	sl := &sourceLocation{File: file, Line: line, Function: function}

	log.write("WARNING", fmt.Sprintf("%s:%d %s", file, line, fmt.Sprintf(format, a...)), sl)
}

func (log *Log) Errorf(format string, a ...any) {
	log.write("ERROR", fmt.Sprintf(format, a...), nil)
}

func (log *Log) Errorfd(format string, a ...any) {
	file, line, function := logctx(1)
	sl := &sourceLocation{File: file, Line: line, Function: function}

	log.write("ERROR", fmt.Sprintf("%s:%d %s", file, line, fmt.Sprintf(format, a...)), sl)
}
