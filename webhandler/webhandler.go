package webhandler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/akfaew/utils/ae"
	"github.com/gorilla/mux"
)

var ErrorContextFunc func(message, meditation string) interface{} = nil

type WebHandler func(http.ResponseWriter, *http.Request) *WebError

type WebTemplate struct {
	t *template.Template
}

type Templates map[string]*WebTemplate

type WebError struct {
	Code    int    // HTTP response code
	Error   error  // for the logs
	Message string // for the user
}

func WebErrorf(code int, err error, format string, v ...interface{}) *WebError {
	return &WebError{
		Code:    code,
		Error:   err,
		Message: fmt.Sprintf(format, v...),
	}
}

func WebErrorInternal(err error) *WebError {
	return &WebError{
		Code:    http.StatusInternalServerError,
		Error:   err,
		Message: "Internal server error",
	}
}

func WebErrorInternalf(err error, format string, v ...interface{}) *WebError {
	return &WebError{
		Code:    http.StatusInternalServerError,
		Error:   err,
		Message: fmt.Sprintf(format, v...),
	}
}

// ParseTemplate parses the template nesting it in the base template
func ParseTemplate(filenames ...string) *WebTemplate {
	if len(filenames) == 0 {
		return nil
	}

	paths := []string{}
	for _, f := range filenames {
		paths = append(paths, "templates/"+f)
	}
	tmpl := template.New("").Funcs(template.FuncMap{
		"format_html": FormatHTML,
		"contains":    strings.Contains,
		"has_prefix":  strings.HasPrefix,
		"has_suffix":  strings.HasSuffix,
	})
	tmpl = template.Must(tmpl.ParseFiles(paths...))

	return &WebTemplate{tmpl.Lookup(filenames[0])}
}

// ServeHTTP renders an error template and logs the failure if a problem occurs.
func (fn WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *WebError
		ctx := r.Context()
		lg := ae.NewLog(ctx)
		m := meditation()

		// Log for the site admin. No logging occurs if no error is passed.
		if e.Error != nil || len(e.Message) > 0 {
			switch e.Code / 100 {
			case 4:
				lg.Infof("Handler error %d. err=\"%v\", msg=\"%s\", meditation=%s",
					e.Code, e.Error, e.Message, m)
			default:
				lg.Errorf("Handler error %d. err=\"%v\", msg=\"%s\", meditation=%s",
					e.Code, e.Error, e.Message, m)
			}
		}

		// Error for the user
		tmpl, err := template.ParseFiles("templates/error.html")
		if err != nil { // if TemplateError does not exist
			http.Error(w, e.Message, e.Code)
			return
		}

		buf := new(bytes.Buffer)
		var vars interface{}
		if ErrorContextFunc != nil {
			vars = ErrorContextFunc(e.Message, m)
		}
		if err := tmpl.Execute(buf, vars); err != nil {
			lg.Errorfd("err=%v", err)
			http.Error(w, "Internal error executing template", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(e.Code)
		if _, err := buf.WriteTo(w); err != nil {
			lg.Errorfd("err=%v", err)
			http.Error(w, "Internal error executing template", http.StatusInternalServerError)
			return
		}
	}
}

func WebHandle(r *mux.Router, method string, path string, handler func(w http.ResponseWriter, r *http.Request) *WebError) {
	r.Methods(method).Path(path).Handler(WebHandler(handler))
}

// Executor uses webContext() to obtain a list of variables to pass to the underlying html template.
func (tmpl *WebTemplate) Executor(webContext func(http.ResponseWriter, *http.Request, *WebTemplate) (interface{}, *WebError)) func(http.ResponseWriter, *http.Request) *WebError {
	return func(w http.ResponseWriter, r *http.Request) *WebError {
		// Get the web context
		vars, weberr := webContext(w, r, tmpl)
		if weberr != nil {
			return weberr
		}

		// Execute the template
		buf := new(bytes.Buffer)
		if err := tmpl.t.Execute(buf, vars); err != nil {
			return WebErrorf(http.StatusInternalServerError, err, "Internal error executing template")
		}

		// Return the result
		if _, err := buf.WriteTo(w); err != nil {
			return WebErrorInternal(err)
		}

		return nil
	}
}
