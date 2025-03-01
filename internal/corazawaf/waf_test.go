// Copyright 2022 Juan Pablo Tosso and the OWASP Coraza contributors
// SPDX-License-Identifier: Apache-2.0

package corazawaf

import (
	"io"
	"os"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	waf := NewWAF()
	waf.RequestBodyAccess = true
	waf.ResponseBodyAccess = true
	waf.RequestBodyLimit = 1044

	tx := waf.NewTransactionWithID("test")
	if !tx.RequestBodyAccess {
		t.Error("Request body access not enabled")
	}
	if !tx.ResponseBodyAccess {
		t.Error("Response body access not enabled")
	}
	if tx.RequestBodyLimit != 1044 {
		t.Error("Request body limit not set")
	}
	if tx.id != "test" {
		t.Error("ID not set")
	}
	tx = waf.NewTransactionWithID("")
	if tx.id == "" {
		t.Error("ID not set")
	}
	tx = waf.NewTransaction()
	if tx.id == "" {
		t.Error("ID not set")
	}
}

func TestSetDebugLogPath(t *testing.T) {
	waf := NewWAF()

	testCases := []struct {
		path   string
		writer io.Writer
	}{
		{
			path:   "/dev/stdout",
			writer: os.Stdout,
		},
		{
			path:   "/dev/stderr",
			writer: os.Stderr,
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.path, func(t *testing.T) {
			err := waf.SetDebugLogPath(tCase.path)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			l := waf.Logger.(*stdDebugLogger)
			if want, have := tCase.writer, l.Writer(); want != have {
				t.Error("unexpected logger writer")
			}
			_ = waf.SetDebugLogPath("")
		})
	}
}

func TestValidate(t *testing.T) {
	testCases := map[string]struct {
		customizer func(*WAF)
		expectErr  bool
	}{
		"default": {
			expectErr:  false,
			customizer: func(w *WAF) {},
		},
		"request body limit less than zero": {
			expectErr:  true,
			customizer: func(w *WAF) { w.RequestBodyLimit = -1 },
		},
		"request body limit greater than 1gb": {
			expectErr:  true,
			customizer: func(w *WAF) { w.RequestBodyLimit = _1gb + 1 },
		},
		"request body in memory limit less than zero": {
			expectErr:  true,
			customizer: func(w *WAF) { w.SetRequestBodyInMemoryLimit(-1) },
		},
		"request body limit less than request body in memory limit": {
			expectErr: true,
			customizer: func(w *WAF) {
				w.RequestBodyLimit = 10
				w.SetRequestBodyInMemoryLimit(11)
			}},
		"response body limit less than zero": {
			expectErr:  true,
			customizer: func(w *WAF) { w.ResponseBodyLimit = -1 },
		},
		"response body limit greater than 1gb": {
			expectErr:  true,
			customizer: func(w *WAF) { w.ResponseBodyLimit = _1gb + 1 },
		},
	}

	for name, tCase := range testCases {
		t.Run(name, func(t *testing.T) {
			waf := NewWAF()
			tCase.customizer(waf)
			err := waf.Validate()
			if tCase.expectErr {
				if err == nil {
					t.Fatalf("expected error: %s", err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}
			}

		})
	}
}
