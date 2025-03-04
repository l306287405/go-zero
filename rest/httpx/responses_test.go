package httpx

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/l306287405/go-zero/core/logx"
)

type message struct {
	Name string `json:"name"`
}

func init() {
	logx.Disable()
}

func TestError(t *testing.T) {
	const (
		body        = "foo"
		wrappedBody = `"foo"`
	)

	tests := []struct {
		name          string
		input         string
		errorHandler  func(error) (int, interface{})
		expectHasBody bool
		expectBody    string
		expectCode    int
	}{
		{
			name:          "default error handler",
			input:         body,
			expectHasBody: true,
			expectBody:    body,
			expectCode:    http.StatusBadRequest,
		},
		{
			name:  "customized error handler return string",
			input: body,
			errorHandler: func(err error) (int, interface{}) {
				return http.StatusForbidden, err.Error()
			},
			expectHasBody: true,
			expectBody:    wrappedBody,
			expectCode:    http.StatusForbidden,
		},
		{
			name:  "customized error handler return error",
			input: body,
			errorHandler: func(err error) (int, interface{}) {
				return http.StatusForbidden, err
			},
			expectHasBody: true,
			expectBody:    body,
			expectCode:    http.StatusForbidden,
		},
		{
			name:  "customized error handler return nil",
			input: body,
			errorHandler: func(err error) (int, interface{}) {
				return http.StatusForbidden, nil
			},
			expectHasBody: false,
			expectBody:    "",
			expectCode:    http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := tracedResponseWriter{
				headers: make(map[string][]string),
			}
			if test.errorHandler != nil {
				lock.RLock()
				prev := errorHandler
				lock.RUnlock()
				SetErrorHandler(test.errorHandler)
				defer func() {
					lock.Lock()
					errorHandler = prev
					lock.Unlock()
				}()
			}
			Error(&w, errors.New(test.input))
			assert.Equal(t, test.expectCode, w.code)
			assert.Equal(t, test.expectHasBody, w.hasBody)
			assert.Equal(t, test.expectBody, strings.TrimSpace(w.builder.String()))
		})
	}
}

func TestErrorWithHandler(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	Error(&w, errors.New("foo"), func(w http.ResponseWriter, err error) {
		http.Error(w, err.Error(), 499)
	})
	assert.Equal(t, 499, w.code)
	assert.True(t, w.hasBody)
	assert.Equal(t, "foo", strings.TrimSpace(w.builder.String()))
}

func TestOk(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	Ok(&w)
	assert.Equal(t, http.StatusOK, w.code)
}

func TestOkJson(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	msg := message{Name: "anyone"}
	OkJson(&w, msg)
	assert.Equal(t, http.StatusOK, w.code)
	assert.Equal(t, "{\"name\":\"anyone\"}", w.builder.String())
}

func TestWriteJsonTimeout(t *testing.T) {
	// only log it and ignore
	w := tracedResponseWriter{
		headers: make(map[string][]string),
		timeout: true,
	}
	msg := message{Name: "anyone"}
	WriteJson(&w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

func TestWriteJsonLessWritten(t *testing.T) {
	w := tracedResponseWriter{
		headers:     make(map[string][]string),
		lessWritten: true,
	}
	msg := message{Name: "anyone"}
	WriteJson(&w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

type tracedResponseWriter struct {
	headers     map[string][]string
	builder     strings.Builder
	hasBody     bool
	code        int
	lessWritten bool
	timeout     bool
}

func (w *tracedResponseWriter) Header() http.Header {
	return w.headers
}

func (w *tracedResponseWriter) Write(bytes []byte) (n int, err error) {
	if w.timeout {
		return 0, http.ErrHandlerTimeout
	}

	n, err = w.builder.Write(bytes)
	if w.lessWritten {
		n -= 1
	}
	w.hasBody = true

	return
}

func (w *tracedResponseWriter) WriteHeader(code int) {
	w.code = code
}
