package extplugin

import (
	"context"

	"github.com/7K-Inari/inari-plugin-sdk/host/agentgateway"
)

// health reports NOT_SERVING when the Agent Gateway is not configured, so the
// host can surface a degraded extension instead of silent failures.
type health struct {
	gw agentgateway.Gateway
}

func (h health) Health(_ context.Context) (bool, string) {
	if h.gw == nil {
		return false, "agent gateway is not configured; actions fail closed"
	}
	return true, "agent gateway configured"
}
