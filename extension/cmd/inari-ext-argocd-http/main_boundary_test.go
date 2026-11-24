package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	agentv1 "github.com/7K-Inari/inari-api/gen/go/inari/agent/v1"
)

func TestActionsHandlerRejectsOversizedBodyWithoutInvokingGateway(t *testing.T) {
	gw := &fakeGateway{ack: &agentv1.CommandAck{CommandId: "cmd-large", Result: agentv1.CommandResult_COMMAND_RESULT_APPLIED}}
	h := buildHandler(t, gw)
	// Keep a complete valid value inside the old 1 MiB truncation boundary and
	// hide a second value beyond it. The entire request must be rejected.
	body := syncBody() + strings.Repeat(" ", (1<<20)-len(syncBody())) + `{}`
	r := httptest.NewRequest(http.MethodPost, "/actions/argocd.sync", strings.NewReader(body))
	r.Header.Set("X-Inari-User", "user-1")
	r.Header.Set("X-Inari-Org", "tenant-acme")
	r.Header.Set("X-Inari-Auth-Method", "oidc-sso-session")
	r.Header.Set("X-Inari-Downstream-Authorization", "Bearer sso-token-123")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
	if gw.command != nil {
		t.Fatal("gateway must not be called for an oversized action body")
	}
}

func TestActionsHandlerAcceptsBodyAtSizeLimit(t *testing.T) {
	gw := &fakeGateway{ack: &agentv1.CommandAck{CommandId: "cmd-limit", Result: agentv1.CommandResult_COMMAND_RESULT_APPLIED}}
	h := buildHandler(t, gw)
	body := syncBody() + strings.Repeat(" ", (1<<20)-len(syncBody()))
	r := httptest.NewRequest(http.MethodPost, "/actions/argocd.sync", strings.NewReader(body))
	r.Header.Set("X-Inari-User", "user-1")
	r.Header.Set("X-Inari-Org", "tenant-acme")
	r.Header.Set("X-Inari-Auth-Method", "oidc-sso-session")
	r.Header.Set("X-Inari-Downstream-Authorization", "Bearer sso-token-123")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
	if gw.command == nil {
		t.Fatal("gateway was not called for a valid body at the size limit")
	}
}

func TestActionsHandlerRequiresHostIdentity(t *testing.T) {
	for _, tc := range []struct {
		name string
		user string
		org  string
	}{
		{name: "missing user", org: "tenant-acme"},
		{name: "missing org", user: "user-1"},
		{name: "missing both"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gw := &fakeGateway{ack: &agentv1.CommandAck{CommandId: "cmd-auth", Result: agentv1.CommandResult_COMMAND_RESULT_APPLIED}}
			h := buildHandler(t, gw)
			r := httptest.NewRequest(http.MethodPost, "/actions/argocd.sync", strings.NewReader(syncBody()))
			r.Header.Set("X-Inari-User", tc.user)
			r.Header.Set("X-Inari-Org", tc.org)
			r.Header.Set("X-Inari-Auth-Method", "oidc-sso-session")
			r.Header.Set("X-Inari-Downstream-Authorization", "Bearer sso-token-123")
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("status %d body %s", w.Code, w.Body.String())
			}
			if gw.command != nil {
				t.Fatal("gateway must not be called without a complete host identity")
			}
		})
	}
}
