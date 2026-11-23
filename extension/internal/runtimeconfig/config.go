// Package runtimeconfig builds extension runtime dependencies from the process environment.
package runtimeconfig

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/7K-Hiroba/argocd/extension/internal/argocd"
	"github.com/7K-Inari/inari-plugin-sdk/host/agentgateway"
	"golang.org/x/oauth2/clientcredentials"
)

// GatewayFromEnv builds the public Agent Gateway client. An unset endpoint
// leaves the plugin available for handshake/capability discovery while every
// action fails closed. Partial machine-identity configuration is an error.
func GatewayFromEnv(ctx context.Context) (agentgateway.Gateway, error) {
	baseURL := strings.TrimSpace(os.Getenv("INARI_AGENT_GATEWAY_URL"))
	if baseURL == "" {
		return nil, nil
	}
	clientID := strings.TrimSpace(os.Getenv("INARI_EXTENSION_CLIENT_ID"))
	clientSecret := strings.TrimSpace(os.Getenv("INARI_EXTENSION_CLIENT_SECRET"))
	tokenURL := strings.TrimSpace(os.Getenv("INARI_EXTENSION_TOKEN_URL"))
	if clientID == "" || clientSecret == "" || tokenURL == "" {
		return nil, fmt.Errorf("INARI_EXTENSION_CLIENT_ID, INARI_EXTENSION_CLIENT_SECRET, and INARI_EXTENSION_TOKEN_URL are required when INARI_AGENT_GATEWAY_URL is set")
	}
	tokens := (&clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}).TokenSource(ctx)
	gateway, err := agentgateway.NewClient(agentgateway.Config{BaseURL: baseURL, TokenSource: tokens})
	if err != nil {
		return nil, fmt.Errorf("configure agent gateway: %w", err)
	}
	return gateway, nil
}

// ManagedProjects returns the configured ArgoCD project allowlist.
func ManagedProjects() []string {
	raw := os.Getenv("INARI_MANAGED_PROJECTS")
	if raw == "" {
		return []string{"inari"}
	}
	var projects []string
	for _, project := range strings.Split(raw, ",") {
		if project = strings.TrimSpace(project); project != "" {
			projects = append(projects, project)
		}
	}
	return projects
}

// ActionTimeout returns the configured action timeout.
func ActionTimeout() time.Duration {
	if raw := os.Getenv("INARI_ACTION_TIMEOUT"); raw != "" {
		if timeout, err := time.ParseDuration(raw); err == nil && timeout > 0 {
			return timeout
		}
	}
	return argocd.DefaultTimeout
}
