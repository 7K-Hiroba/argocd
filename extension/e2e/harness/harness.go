// Package harness provides a test-only Agent Gateway and agent session. It
// uses the generated public Connect contract; no transport schema is mirrored.
package harness

import (
	"context"
	"fmt"
	"sync"
	"time"

	"connectrpc.com/connect/v2"
	agentv1 "github.com/7K-Inari/inari-api/gen/go/inari/agent/v1"
	extensionsv1 "github.com/7K-Inari/inari-api/gen/go/inari/extensions/v1"
	"github.com/7K-Inari/inari-api/gen/go/inari/extensions/v1/extensionsv1connect"
)

// Session is one connected agent: commands in, acks out.
type Session struct {
	Commands chan *agentv1.InvokeAction
	acks     map[string]chan *agentv1.CommandAck
	mu       sync.Mutex
}

func newSession() *Session {
	return &Session{Commands: make(chan *agentv1.InvokeAction, 16), acks: map[string]chan *agentv1.CommandAck{}}
}

// Ack delivers the agent's acknowledgement for a command.
func (s *Session) Ack(ack *agentv1.CommandAck) {
	s.mu.Lock()
	ch, ok := s.acks[ack.GetCommandId()]
	s.mu.Unlock()
	if ok {
		ch <- ack
	}
}

// Server implements the generated Agent Gateway handler contract.
type Server struct {
	extensionsv1connect.UnimplementedAgentGatewayHandler
	mu       sync.Mutex
	sessions map[string]*Session
	// AgentTimeout bounds the wait for an agent's ack.
	AgentTimeout time.Duration
}

// NewServer builds a gateway server.
func NewServer() *Server {
	return &Server{AgentTimeout: 30 * time.Second, sessions: map[string]*Session{}}
}

// RegisterAgent connects an agent session for clusterID; the returned cancel
// function disconnects it (subsequent invokes fail closed with UNAVAILABLE).
func (s *Server) RegisterAgent(clusterID string) (*Session, func()) {
	sess := newSession()
	s.mu.Lock()
	s.sessions[clusterID] = sess
	s.mu.Unlock()
	return sess, func() {
		s.mu.Lock()
		delete(s.sessions, clusterID)
		s.mu.Unlock()
	}
}

// InvokeAction routes one public command to a connected agent session.
func (s *Server) InvokeAction(ctx context.Context, cmd *agentv1.InvokeAction) (*agentv1.CommandAck, error) {
	info, ok := connect.CallInfoForServerContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal, "missing Connect call info")
	}
	header := info.RequestHeader()
	cluster := header.Get(extensionsv1.MetadataCluster)
	tenant := header.Get(extensionsv1.MetadataTenant)
	if cluster == "" || tenant == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, "tenant and cluster metadata are required")
	}
	if header.Get("Authorization") == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, "extension identity is required")
	}

	s.mu.Lock()
	sess, ok := s.sessions[cluster]
	s.mu.Unlock()
	if !ok {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Sprintf("no connected agent for cluster %q", cluster))
	}

	ch := make(chan *agentv1.CommandAck, 1)
	sess.mu.Lock()
	if cmd.GetCommandId() == "" {
		sess.mu.Unlock()
		return nil, connect.NewError(connect.CodeInvalidArgument, "commandId is required for ack correlation")
	}
	if _, dup := sess.acks[cmd.GetCommandId()]; dup {
		sess.mu.Unlock()
		return nil, connect.NewError(connect.CodeAborted, fmt.Sprintf("duplicate commandId %q already in flight", cmd.GetCommandId()))
	}
	// Raw user credentials arrive only on this in-memory hop. The test
	// control plane mints a vault reference before the command reaches the
	// agent, exactly as the production contract requires.
	if credential := header.Get(extensionsv1.MetadataUserCredential); credential != "" {
		cmd.UserCredentialRef = "vaultref/" + cmd.GetCommandId()
	}
	sess.acks[cmd.GetCommandId()] = ch
	sess.mu.Unlock()
	defer func() {
		sess.mu.Lock()
		delete(sess.acks, cmd.GetCommandId())
		sess.mu.Unlock()
	}()

	select {
	case sess.Commands <- cmd:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	select {
	case ack := <-ch:
		return ack, nil
	case <-time.After(s.AgentTimeout):
		return nil, connect.NewError(connect.CodeDeadlineExceeded, fmt.Sprintf("agent did not ack command %q", cmd.GetCommandId()))
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// RunAgent drains the session's commands, executes each against ArgoCD, and
// acks. executor is typically (&agentstub.ArgoCDClient{...}).Execute.
func RunAgent(ctx context.Context, sess *Session, executor func(context.Context, *agentv1.InvokeAction) *agentv1.CommandAck) {
	for {
		select {
		case <-ctx.Done():
			return
		case cmd := <-sess.Commands:
			sess.Ack(executor(ctx, cmd))
		}
	}
}
