// Command inari-e2e-gateway is the colocated e2e counterpart of the control
// plane's Agent Gateway + inari-agent pair, for the kind e2e. It serves the
// generated public Connect contract and runs one agent session against a
// reachable ArgoCD API.
package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"github.com/7K-Hiroba/argocd/extension/e2e/harness"
	"github.com/7K-Hiroba/argocd/extension/internal/agentstub"
	"github.com/7K-Inari/inari-api/gen/go/inari/extensions/v1/extensionsv1connect"
)

func main() {
	listen := envOr("LISTEN_ADDR", "127.0.0.1:9090")
	clusterID := envOr("CLUSTER_ID", "e2e-kind")

	listener, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatalf("listen %s: %v", listen, err)
	}
	gateway := harness.NewServer()
	rpcServer := connect.NewServer()
	extensionsv1connect.RegisterAgentGatewayHandler(rpcServer, gateway)
	mux := http.NewServeMux()
	connecthttp.Mount(mux, rpcServer)
	// Test-only OAuth token endpoint: the plugin uses the same standard
	// client-credentials flow as production, with no static-token bypass.
	mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"e2e-extension-token","token_type":"Bearer","expires_in":300}`))
	})
	server := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}

	sess, disconnect := gateway.RegisterAgent(clusterID)
	defer disconnect()
	argo := &agentstub.ArgoCDClient{
		BaseURL:         envOr("ARGOCD_BASE_URL", "http://127.0.0.1:8080"),
		BearerToken:     os.Getenv("ARGOCD_TOKEN"),
		ManagedProjects: projects(),
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go harness.RunAgent(ctx, sess, argo.Execute)
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	log.Printf("e2e gateway on %s, agent session for cluster %q against %s (managed projects: %v)",
		listen, clusterID, argo.BaseURL, argo.ManagedProjects)
	if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("serve: %v", err)
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func projects() []string {
	return strings.Split(envOr("INARI_MANAGED_PROJECTS", "inari"), ",")
}
