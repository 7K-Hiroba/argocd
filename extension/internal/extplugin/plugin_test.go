package extplugin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	agentv1 "github.com/7K-Inari/inari-api/gen/go/inari/agent/v1"
	pluginv1 "github.com/7K-Inari/inari-api/gen/go/inari/plugin/v1"
	pluginsdk "github.com/7K-Inari/inari-plugin-sdk"
	"github.com/7K-Inari/inari-plugin-sdk/testkit"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	"github.com/7K-Hiroba/argocd/extension/internal/argocd"
	"github.com/7K-Inari/inari-plugin-sdk/host/agentgateway"
)

// fakeGateway records public Agent Gateway calls and returns a scripted ack or error.
type fakeGateway struct {
	clusterID string
	command   *agentv1.InvokeAction
	ack       *agentv1.CommandAck
	err       error
}

var _ agentgateway.Gateway = (*fakeGateway)(nil)

func (f *fakeGateway) InvokeAction(_ context.Context, clusterID string, command *agentv1.InvokeAction) (*agentv1.CommandAck, error) {
	f.clusterID = clusterID
	f.command = command
	return f.ack, f.err
}

var testAuth = pluginsdk.AuthContext{PrincipalID: "user-1", TenantID: "tenant-acme"}
var scoping = argocd.Scoping{ManagedProjects: []string{"inari"}}

func build(t *testing.T, gw agentgateway.Gateway) *testkit.Client {
	t.Helper()
	p, err := Build(Config{Gateway: gw, Scoping: scoping})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	return testkit.Run(t, p)
}

// withUserSession mimics the control plane's extension proxy injecting the
// per-user SSO session credential (oidc-sso-session) as connection metadata.
func withUserSession(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx,
		pluginsdk.HeaderAuthMethod, pluginsdk.AuthMethodOIDCSSOSession.String(),
		pluginsdk.HeaderDownstreamAuthorization, "Bearer "+token,
	)
}

func invoke(t *testing.T, c *testkit.Client, action string, payload any) (*pluginsdk.Response, error) {
	t.Helper()
	return invokeWithSession(t, c, action, payload, "user-token-1")
}

func invokeWithSession(t *testing.T, c *testkit.Client, action string, payload any, token string) (*pluginsdk.Response, error) {
	t.Helper()
	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if token != "" {
		ctx = withUserSession(ctx, token)
	}
	return c.Invoke(ctx, action, testAuth, b)
}

// TestDownstreamTokenSchemeStripping proves the gateway hop receives the
// raw token regardless of how the injected credential was schemed/cased,
// and that a scheme-only value (no credential material) fails closed.
func TestDownstreamTokenSchemeStripping(t *testing.T) {
	payload := map[string]any{
		"clusterId": "cluster-1",
		"app":       map[string]any{"name": "a", "namespace": "argocd", "project": "inari"},
	}
	for _, tc := range []struct {
		name       string
		injected   string
		want       string
		wantReauth bool
	}{
		{"bearer", "Bearer tok-1", "tok-1", false},
		{"bearer lowercase", "bearer tok-2", "tok-2", false},
		{"bearer uppercase", "BEARER tok-3", "tok-3", false},
		{"raw token, no scheme", "tok-4", "tok-4", false},
		{"double space after scheme", "Bearer  tok-5", "tok-5", false},
		{"scheme only, no token", "Bearer", "", true},
		{"scheme only lowercase", "bearer", "", true},
		{"scheme + trailing space", "Bearer ", "", true},
		{"whitespace only", " ", "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gw := &fakeGateway{ack: &agentv1.CommandAck{CommandId: "cmd-1", Result: agentv1.CommandResult_COMMAND_RESULT_APPLIED}}
			c := build(t, gw)
			b, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			ctx = metadata.AppendToOutgoingContext(ctx,
				pluginsdk.HeaderAuthMethod, pluginsdk.AuthMethodOIDCSSOSession.String(),
				pluginsdk.HeaderDownstreamAuthorization, tc.injected,
			)
			_, err = c.Invoke(ctx, argocd.ActionSync, testAuth, b)
			if tc.wantReauth {
				var pe *pluginsdk.Error
				if !errors.As(err, &pe) || pe.Code != pluginsdk.CodeUnauthenticated {
					t.Fatalf("injected %q: expected fail-closed reauth, got err=%v", tc.injected, err)
				}
				if gw.command != nil {
					t.Fatal("gateway must not be called without credential material")
				}
				return
			}
			if err != nil {
				t.Fatalf("injected %q: %v", tc.injected, err)
			}
			if gw.command == nil {
				t.Fatalf("injected %q: gateway was not called", tc.injected)
			}
		})
	}
}

func TestCapabilities(t *testing.T) {
	c := build(t, &fakeGateway{})
	ctx := context.Background()
	info, err := c.Info(ctx)
	if err != nil {
		t.Fatalf("Info: %v", err)
	}
	if info.GetName() != "inari-ext-argocd" {
		t.Fatalf("name: %v", info.GetName())
	}
	actions, err := c.Capabilities(ctx)
	if err != nil {
		t.Fatalf("Capabilities: %v", err)
	}
	want := []string{argocd.ActionRefresh, argocd.ActionResourceAction, argocd.ActionRollback, argocd.ActionSync}
	if len(actions) != len(want) {
		t.Fatalf("actions: %v", actions)
	}
	for i, a := range actions {
		if a.GetName() != want[i] {
			t.Fatalf("action %d: got %q want %q", i, a.GetName(), want[i])
		}
		if len(a.GetInputSchema()) == 0 || len(a.GetOutputSchema()) == 0 {
			t.Fatalf("action %q missing schemas", a.GetName())
		}
	}
}

func TestDeclaresOIDCSSOSessionAuthMethod(t *testing.T) {
	p, err := Build(Config{Gateway: &fakeGateway{}, Scoping: scoping})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	info, err := pluginsdk.NewGRPCService(p).GetInfo(context.Background(), &pluginv1.GetInfoRequest{})
	if err != nil {
		t.Fatalf("GetInfo: %v", err)
	}
	methods := info.GetAuthMethods()
	if len(methods) != 1 {
		t.Fatalf("expected exactly one declared auth method, got %v", methods)
	}
	m := methods[0]
	if m.GetType() != pluginv1.AuthMethod_TYPE_OIDC_SSO_SESSION {
		t.Fatalf("auth method type: %v", m.GetType())
	}
	if m.GetAudience() != "argocd" {
		t.Fatalf("audience (session provider key): %q", m.GetAudience())
	}
	if !m.GetIsDefault() {
		t.Fatal("oidc-sso-session must be the default auth method")
	}
}

func TestFailsClosedWithoutUserSession(t *testing.T) {
	gw := &fakeGateway{ack: &agentv1.CommandAck{CommandId: "cmd-1", Result: agentv1.CommandResult_COMMAND_RESULT_APPLIED}}
	c := build(t, gw)
	_, err := invokeWithSession(t, c, argocd.ActionSync, map[string]any{
		"clusterId": "cluster-1",
		"app":       map[string]any{"name": "a", "namespace": "argocd", "project": "inari"},
	}, "")
	var pe *pluginsdk.Error
	if !errors.As(err, &pe) || pe.Code != pluginsdk.CodeUnauthenticated {
		t.Fatalf("expected CodeUnauthenticated without a user session, got %v", err)
	}
	if gw.command != nil {
		t.Fatal("gateway must not be called without a user session")
	}
}

func TestDoesNotPlaceUserTokenInCommand(t *testing.T) {
	gw := &fakeGateway{ack: &agentv1.CommandAck{CommandId: "cmd-1", Result: agentv1.CommandResult_COMMAND_RESULT_APPLIED}}
	c := build(t, gw)
	const token = "sso-session-token-xyz"
	_, err := invokeWithSession(t, c, argocd.ActionSync, map[string]any{
		"clusterId": "cluster-1",
		"app":       map[string]any{"name": "web-shop", "namespace": "argocd", "project": "inari"},
	}, token)
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	raw, err := proto.Marshal(gw.command)
	if err != nil {
		t.Fatalf("marshal command: %v", err)
	}
	if bytes.Contains(raw, []byte(token)) {
		t.Fatal("user token leaked into the command payload")
	}
	if gw.command.GetUserCredentialRef() != "" {
		t.Fatalf("extension must not self-assert user_credential_ref, got %q", gw.command.GetUserCredentialRef())
	}
}

func TestExpiredDownstreamSessionSignalsReauth(t *testing.T) {
	gw := &fakeGateway{ack: &agentv1.CommandAck{
		CommandId: "cmd-2",
		Result:    agentv1.CommandResult_COMMAND_RESULT_FAILED,
		Message:   argocd.ReauthSignalPrefix + ": argocd session expired (401)",
	}}
	c := build(t, gw)
	_, err := invoke(t, c, argocd.ActionSync, map[string]any{
		"clusterId": "cluster-1",
		"app":       map[string]any{"name": "a", "namespace": "argocd", "project": "inari"},
	})
	var pe *pluginsdk.Error
	if !errors.As(err, &pe) || pe.Code != pluginsdk.CodeUnauthenticated {
		t.Fatalf("expected CodeUnauthenticated for expired downstream session, got %v", err)
	}
	if pe.Details["code"] != "reauth_required" || pe.Details["provider"] != "argocd" {
		t.Fatalf("expected typed re-auth details, got %v", pe.Details)
	}
}

func TestSyncRoundTrip(t *testing.T) {
	gw := &fakeGateway{ack: &agentv1.CommandAck{
		CommandId: "cmd-1",
		Result:    agentv1.CommandResult_COMMAND_RESULT_APPLIED,
		Message:   "synced",
	}}
	c := build(t, gw)
	resp, err := invoke(t, c, argocd.ActionSync, map[string]any{
		"clusterId": "cluster-1",
		"app":       map[string]any{"name": "web-shop", "namespace": "argocd", "project": "inari"},
		"prune":     true,
	})
	if err != nil {
		t.Fatalf("Invoke: %v", err)
	}
	var res Result
	if err := json.Unmarshal(resp.Result, &res); err != nil {
		t.Fatalf("result: %v", err)
	}
	if res.Outcome != "applied" || res.Message != "synced" {
		t.Fatalf("result: %+v", res)
	}
	if gw.clusterID != "cluster-1" {
		t.Fatalf("cluster routing = %q", gw.clusterID)
	}
	if gw.command.GetAction() != "sync" { // agent-side verb is bare

		t.Fatalf("command action: %v", gw.command.GetAction())
	}
	if gw.command.GetParameters().AsMap()["prune"] != true {
		t.Fatalf("params not forwarded: %v", gw.command.GetParameters().AsMap())
	}
}

func TestFailsClosedWithoutGateway(t *testing.T) {
	c := build(t, nil)
	_, err := invoke(t, c, argocd.ActionSync, map[string]any{
		"clusterId": "cluster-1",
		"app":       map[string]any{"name": "a", "namespace": "argocd", "project": "inari"},
	})
	var pe *pluginsdk.Error
	if !errors.As(err, &pe) || pe.Code != pluginsdk.CodeUnavailable {
		t.Fatalf("expected CodeUnavailable, got %v", err)
	}
}

func TestFailsClosedWhenGatewayErrors(t *testing.T) {
	c := build(t, &fakeGateway{err: errors.New("agent disconnected")})
	_, err := invoke(t, c, argocd.ActionRefresh, map[string]any{
		"clusterId": "cluster-1",
		"app":       map[string]any{"name": "a", "namespace": "argocd", "project": "inari"},
	})
	var pe *pluginsdk.Error
	if !errors.As(err, &pe) || pe.Code != pluginsdk.CodeUnavailable {
		t.Fatalf("expected CodeUnavailable, got %v", err)
	}
}

func TestRejectsUnmanagedProject(t *testing.T) {
	gw := &fakeGateway{}
	c := build(t, gw)
	_, err := invoke(t, c, argocd.ActionSync, map[string]any{
		"clusterId": "cluster-1",
		"app":       map[string]any{"name": "a", "namespace": "argocd", "project": "default"},
	})
	var pe *pluginsdk.Error
	if !errors.As(err, &pe) || pe.Code != pluginsdk.CodeInvalidArgument {
		t.Fatalf("expected CodeInvalidArgument, got %v", err)
	}
	if gw.command != nil {
		t.Fatal("gateway must not be called for unmanaged projects")
	}
}

func TestRejectsUnknownFields(t *testing.T) {
	c := build(t, &fakeGateway{})
	_, err := invoke(t, c, argocd.ActionSync, map[string]any{
		"clusterId": "cluster-1",
		"app":       map[string]any{"name": "a", "namespace": "argocd", "project": "inari"},
		"tenantId":  "tenant-evil",
	})
	var pe *pluginsdk.Error
	if !errors.As(err, &pe) || pe.Code != pluginsdk.CodeInvalidArgument {
		t.Fatalf("expected CodeInvalidArgument, got %v", err)
	}
}

func TestAgentFailureSurfaces(t *testing.T) {
	c := build(t, &fakeGateway{ack: &agentv1.CommandAck{
		CommandId: "cmd-9",
		Result:    agentv1.CommandResult_COMMAND_RESULT_FAILED,
		Message:   "application not found",
	}})
	_, err := invoke(t, c, argocd.ActionRollback, map[string]any{
		"clusterId":  "cluster-1",
		"app":        map[string]any{"name": "a", "namespace": "argocd", "project": "inari"},
		"revisionId": 2,
	})
	var pe *pluginsdk.Error
	if !errors.As(err, &pe) || pe.Code != pluginsdk.CodeInternal {
		t.Fatalf("expected CodeInternal, got %v", err)
	}
}

func TestHealth(t *testing.T) {
	ok, _, err := build(t, &fakeGateway{}).Health(context.Background())
	if err != nil || !ok {
		t.Fatalf("with gateway: serving=%v err=%v", ok, err)
	}
	ok, msg, err := build(t, nil).Health(context.Background())
	if err != nil {
		t.Fatalf("Health: %v", err)
	}
	if ok || msg == "" {
		t.Fatalf("without gateway should not serve: %v %q", ok, msg)
	}
}
