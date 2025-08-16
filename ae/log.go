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
	"runtime/debug"
	"strings"
	"time"

	"github.com/akfaew/utils"
	"github.com/akfaew/utils/xctc"
)

var (
	trimprefix = ""
	projectID  = os.Getenv("GOOGLE_CLOUD_PROJECT")
	// serviceName holds the App Engine service name, if available.
	serviceName = os.Getenv("GAE_SERVICE")
)

type Log struct {
	ctx         context.Context
	httpRequest map[string]any
	data        map[string]any
	labels      map[string]string
}

type User interface {
	UserID() string
	UserEmail() string
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

	if ip := r.Header.Get("X-Appengine-User-Ip"); ip != "" {
		req["remoteIp"] = ip
	} else if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		if i := strings.Index(ip, ","); i >= 0 {
			ip = ip[:i]
		}
		req["remoteIp"] = strings.TrimSpace(ip)
	} else if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
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

func (log *Log) WithLabel(name, value string) *Log {
	l := *log
	if l.labels == nil {
		l.labels = map[string]string{}
	} else {
		labels := make(map[string]string, len(l.labels)+1)
		for k, v := range l.labels {
			labels[k] = v
		}
		l.labels = labels
	}
	l.labels[name] = value
	return &l
}

func (log *Log) WithDuration(d time.Duration) *Log {
	return log.WithField("duration_ms", d.Milliseconds())
}

func (log *Log) WithUser(u User) *Log {
	return log.WithLabel("user_id", u.UserID()).WithLabel("user_email", u.UserEmail())
}

func (log *Log) WithUserID(id string) *Log {
	return log.WithLabel("user_id", id)
}

// WithComponent sets the label "component" to the provided value.
func (log *Log) WithComponent(component string) *Log {
	return log.WithLabel("component", component)
}

// Set trimprefix to the path to the source code directory, so that we only log the filename and not the full path.
func init() {
	_, path, _, _ := runtime.Caller(1)

	trimprefix = filepath.Dir(path) + string(filepath.Separator)
}

func getlocation(skip int) sourceLocation {
	pc, file, line, _ := runtime.Caller(skip + 2)

	sl := sourceLocation{
		Function: runtime.FuncForPC(pc).Name(),
		File:     file,
		Line:     line,
	}

	return sl
}

type sourceLocation struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function,omitempty"`
}

type entry struct {
	Trace          string            `json:"logging.googleapis.com/trace,omitempty"`
	SpanID         string            `json:"logging.googleapis.com/spanId,omitempty"`
	SourceLocation sourceLocation    `json:"logging.googleapis.com/sourceLocation,omitempty"`
	Data           map[string]any    `json:"data"`
	Labels         map[string]string `json:"labels,omitempty"`
	Message        string            `json:"message,omitempty"`
	Severity       string            `json:"severity,omitempty"`
	HTTPRequest    map[string]any    `json:"httpRequest,omitempty"`
	StackTrace     string            `json:"stackTrace,omitempty"`
}

func (log *Log) write(severity, format string, a ...any) {
	e := entry{
		Message:  fmt.Sprintf(format, a...),
		Severity: severity,

		Data: map[string]any{},

		SourceLocation: getlocation(1),

		Labels:      map[string]string{},
		HTTPRequest: log.httpRequest,
	}

	// Default subsystem label to the App Engine service, if present.
	if serviceName != "" {
		e.Labels["subsystem"] = serviceName
	}

	for k, v := range log.data {
		e.Data[k] = v
	}

	for k, v := range log.labels {
		e.Labels[k] = v
	}

	if log.ctx != nil {
		if x := xctc.XCTC(log.ctx); x != "" {
			if trace, err := utils.ParseTraceParent(x); err == nil {
				e.Trace = fmt.Sprintf("projects/%s/traces/%s", projectID, trace.TraceID)
				e.SpanID = trace.ParentID
			}
		}
	}

	if severity == "ERROR" {
		e.StackTrace = strings.ReplaceAll(string(debug.Stack()), trimprefix, "")
	}

	if b, err := json.Marshal(e); err == nil {
		_, _ = os.Stdout.Write(append(b, '\n'))
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "{\"severity\":\"ERROR\",\"message\":\"json.Marshal(): %v\"}\n", err)
	}
}

func (log *Log) Debugf(format string, a ...any) {
	log.write("DEBUG", format, a...)
}

func (log *Log) Infof(format string, a ...any) {
	log.write("INFO", format, a...)
}

func (log *Log) Warningf(format string, a ...any) {
	log.write("WARNING", format, a...)
}

func (log *Log) Errorf(format string, a ...any) {
	log.write("ERROR", format, a...)
}

func (log *Log) Err(err error) {
	if err == nil {
		return
	}
	log.Errorf("%v", err)
}
