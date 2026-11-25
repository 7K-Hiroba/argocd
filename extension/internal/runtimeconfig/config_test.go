package runtimeconfig

import (
	"context"
	"testing"
	"time"

	"github.com/7K-Hiroba/argocd/extension/internal/argocd"
)

func TestGatewayFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantNil bool
		wantErr bool
	}{
		{name: "not configured fails closed at invocation", env: map[string]string{}, wantNil: true},
		{name: "url requires extension credentials", env: map[string]string{
			"INARI_AGENT_GATEWAY_URL": "https://inari-server.example.com",
		}, wantNil: true, wantErr: true},
		{name: "invalid url rejected", env: map[string]string{
			"INARI_AGENT_GATEWAY_URL":       "localhost:8080",
			"INARI_EXTENSION_CLIENT_ID":     "ext-inari-ext-argocd",
			"INARI_EXTENSION_CLIENT_SECRET": "secret",
			"INARI_EXTENSION_TOKEN_URL":     "https://keycloak.example.com/token",
		}, wantNil: true, wantErr: true},
		{name: "complete public client config", env: map[string]string{
			"INARI_AGENT_GATEWAY_URL":       "https://inari-server.example.com",
			"INARI_EXTENSION_CLIENT_ID":     "ext-inari-ext-argocd",
			"INARI_EXTENSION_CLIENT_SECRET": "secret",
			"INARI_EXTENSION_TOKEN_URL":     "https://keycloak.example.com/token",
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, key := range []string{"INARI_AGENT_GATEWAY_URL", "INARI_EXTENSION_CLIENT_ID", "INARI_EXTENSION_CLIENT_SECRET", "INARI_EXTENSION_TOKEN_URL"} {
				t.Setenv(key, tc.env[key])
			}
			gateway, err := GatewayFromEnv(context.Background())
			if (err != nil) != tc.wantErr {
				t.Fatalf("GatewayFromEnv error = %v, wantErr %v", err, tc.wantErr)
			}
			if (gateway == nil) != tc.wantNil {
				t.Fatalf("GatewayFromEnv gateway = %T, wantNil %v", gateway, tc.wantNil)
			}
		})
	}
}

func TestManagedProjects(t *testing.T) {
	t.Setenv("INARI_MANAGED_PROJECTS", "")
	if got := ManagedProjects(); len(got) != 1 || got[0] != "inari" {
		t.Fatalf("default projects = %v", got)
	}
	t.Setenv("INARI_MANAGED_PROJECTS", " inari, team-a,,team-b ")
	got := ManagedProjects()
	if len(got) != 3 || got[0] != "inari" || got[1] != "team-a" || got[2] != "team-b" {
		t.Fatalf("projects = %v", got)
	}
}

func TestActionTimeout(t *testing.T) {
	t.Setenv("INARI_ACTION_TIMEOUT", "")
	if got := ActionTimeout(); got != argocd.DefaultTimeout {
		t.Fatalf("default timeout = %v, want %v", got, argocd.DefaultTimeout)
	}
	t.Setenv("INARI_ACTION_TIMEOUT", "12s")
	if got := ActionTimeout(); got != 12*time.Second {
		t.Fatalf("timeout = %v", got)
	}
	t.Setenv("INARI_ACTION_TIMEOUT", "invalid")
	if got := ActionTimeout(); got != argocd.DefaultTimeout {
		t.Fatalf("invalid timeout fallback = %v", got)
	}
}
