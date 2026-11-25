// Command inari-ext-argocd is the Inari ArgoCD reference backend extension.
// It is launched and supervised by the control plane's Extension Host as a
// hashicorp go-plugin subprocess. Runtime dependencies use only public SDKs:
// INARI_AGENT_GATEWAY_URL, INARI_EXTENSION_CLIENT_ID,
// INARI_EXTENSION_CLIENT_SECRET, and INARI_EXTENSION_TOKEN_URL configure the
// typed Agent Gateway client. INARI_MANAGED_PROJECTS and INARI_ACTION_TIMEOUT
// retain their ArgoCD-specific behavior.
package main

import (
	"context"
	"log"

	"github.com/7K-Hiroba/argocd/extension/internal/argocd"
	"github.com/7K-Hiroba/argocd/extension/internal/extplugin"
	"github.com/7K-Hiroba/argocd/extension/internal/runtimeconfig"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalf("inari-ext-argocd: %v", err)
	}
}

func run(ctx context.Context) error {
	gateway, err := runtimeconfig.GatewayFromEnv(ctx)
	if err != nil {
		return err
	}
	plugin, err := extplugin.Build(extplugin.Config{
		Gateway: gateway,
		Scoping: argocd.Scoping{ManagedProjects: runtimeconfig.ManagedProjects()},
		Timeout: runtimeconfig.ActionTimeout(),
	})
	if err != nil {
		return err
	}
	// A nil gateway is intentional: handshake and capabilities still work,
	// while health and every action fail closed without machine configuration.
	return plugin.Serve(ctx)
}
