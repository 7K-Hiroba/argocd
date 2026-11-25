// Command inari-ext-argocd-http serves the ArgoCD extension in dial mode
// (plan §5.8): the pluginv1 contract over ConnectRPC (HTTP) plus the
// POST /actions/<name> surface the control-plane extension proxy forwards
// (/api/extensions/<name>/actions/<action> → /actions/<action>, with the
// authenticated identity in X-Inari-User / X-Inari-Org).
//
// The sibling command (cmd/inari-ext-argocd) speaks the go-plugin subprocess
// protocol for host-supervised exec mode; this binary runs the extension as a
// standalone network service so the Extension Host's dial mode (registered
// endpoint + VerifyHandshake + reverse proxy) can reach it.
//
// Environment mirrors the subprocess command, plus:
//
//	INARI_HTTP_ADDR   listen address (default :8080)
package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	pluginv1 "github.com/7K-Inari/inari-api/gen/go/inari/plugin/v1"
	"github.com/7K-Inari/inari-api/gen/go/inari/plugin/v1/pluginv1connect"
	pluginsdk "github.com/7K-Inari/inari-plugin-sdk"
	"google.golang.org/grpc/metadata"

	"github.com/7K-Hiroba/argocd/extension/internal/argocd"
	"github.com/7K-Hiroba/argocd/extension/internal/extplugin"
	"github.com/7K-Hiroba/argocd/extension/internal/runtimeconfig"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("inari-ext-argocd-http: %v", err)
	}
}

func run() error {
	gateway, err := runtimeconfig.GatewayFromEnv(context.Background())
	if err != nil {
		return err
	}
	cfg := extplugin.Config{
		Gateway: gateway,
		Scoping: argocd.Scoping{ManagedProjects: runtimeconfig.ManagedProjects()},
		Timeout: runtimeconfig.ActionTimeout(),
	}

	p, err := extplugin.Build(cfg)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	// Contract surface (Verify handshake, health, capabilities).
	rpcServer := connect.NewServer()
	pluginv1connect.RegisterPluginContractServiceHandler(rpcServer, connectAdapter{s: pluginsdk.NewGRPCService(p)})
	connecthttp.Mount(mux, rpcServer)
	// Proxy surface: POST /actions/<action>.
	mux.HandleFunc("/actions/", actionsHandler(p))

	addr := os.Getenv("INARI_HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("inari-ext-argocd-http: dial mode listening on %s (contract at /%s/)", addr, pluginv1connect.PluginContractServiceName)
	server := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	return server.ListenAndServe()
}

// connectAdapter bridges the SDK's gRPC-shaped service to the ConnectRPC
// handler interface (same messages, plain request/response in connect v2).
type connectAdapter struct {
	s pluginv1.PluginContractServiceServer
}

func (a connectAdapter) GetInfo(ctx context.Context, req *pluginv1.GetInfoRequest) (*pluginv1.GetInfoResponse, error) {
	r, err := a.s.GetInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (a connectAdapter) GetCapabilities(ctx context.Context, req *pluginv1.GetCapabilitiesRequest) (*pluginv1.GetCapabilitiesResponse, error) {
	r, err := a.s.GetCapabilities(ctx, req)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (a connectAdapter) Invoke(ctx context.Context, req *pluginv1.InvokeRequest) (*pluginv1.InvokeResponse, error) {
	r, err := a.s.Invoke(ctx, req)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (a connectAdapter) HealthCheck(ctx context.Context, req *pluginv1.HealthCheckRequest) (*pluginv1.HealthCheckResponse, error) {
	r, err := a.s.HealthCheck(ctx, req)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Host-injected headers carrying the resolved per-user downstream
// credential (oidc-sso-session). Set only by the control plane's extension
// proxy (strip-then-inject); mirrored here into the contract's connection
// metadata. Values are credentials: never log them.
const (
	headerAuthMethod              = "X-Inari-Auth-Method"
	headerDownstreamAuthorization = "X-Inari-Downstream-Authorization"
	maxActionPayloadBytes         = 1 << 20
)

// actionsHandler maps the proxy-forwarded POST /actions/<action> onto the
// contract's Invoke, lifting identity from the injected headers (never from
// the payload — §5.8).
func actionsHandler(p *pluginsdk.Plugin) http.HandlerFunc {
	s := pluginsdk.NewGRPCService(p)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"detail":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		action := strings.TrimPrefix(r.URL.Path, "/actions/")
		if action == "" || strings.Contains(action, "/") {
			http.Error(w, `{"detail":"action name required"}`, http.StatusBadRequest)
			return
		}
		payload, err := io.ReadAll(io.LimitReader(r.Body, maxActionPayloadBytes+1))
		if err != nil {
			http.Error(w, `{"detail":"read body"}`, http.StatusBadRequest)
			return
		}
		if len(payload) > maxActionPayloadBytes {
			http.Error(w, `{"detail":"request body too large"}`, http.StatusRequestEntityTooLarge)
			return
		}
		if len(payload) == 0 {
			payload = []byte("{}")
		}
		ctx := r.Context()
		// Forward the host-injected downstream credential as connection
		// metadata (the SDK enforces the declared auth method fail-closed).
		if v := r.Header.Get(headerAuthMethod); v != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, pluginsdk.HeaderAuthMethod, v)
		}
		if v := r.Header.Get(headerDownstreamAuthorization); v != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, pluginsdk.HeaderDownstreamAuthorization, v)
		}
		// The contract service is invoked in-process here: bridge the
		// outgoing metadata to incoming so the SDK sees the connection
		// metadata exactly as it would over the wire.
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			ctx = metadata.NewIncomingContext(ctx, md)
		}
		if strings.TrimSpace(r.Header.Get(headerDownstreamAuthorization)) == "" {
			pluginsdk.WriteHTTPError(w, pluginsdk.ReauthError(argocd.SSOProvider))
			return
		}
		resp, err := s.Invoke(ctx, &pluginv1.InvokeRequest{
			Action: action,
			AuthContext: &pluginv1.AuthContext{
				PrincipalId: r.Header.Get("X-Inari-User"),
				TenantId:    r.Header.Get("X-Inari-Org"),
			},
			Payload:   payload,
			RequestId: r.Header.Get("X-Request-Id"),
		})
		if err != nil {
			http.Error(w, `{"detail":"invoke failed"}`, http.StatusBadGateway)
			return
		}
		if pe := resp.GetError(); pe != nil {
			pluginsdk.WriteHTTPError(w, pluginsdk.ErrorFromProto(pe))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(resp.GetResult())
	}
}
