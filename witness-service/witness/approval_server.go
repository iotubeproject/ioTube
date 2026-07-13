// Copyright (c) 2026 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/iotexproject/ioTube/witness-service/util"
)

// ApprovalServer exposes a single POST /lark/callback endpoint that authorizes
// previously-paused transfers. Multiple cashiers share one server; the
// embedded `cashier` field in the Lark card's action.value routes to the right
// ApprovalGuard.
type ApprovalServer struct {
	listenAddr    string
	signingSecret string
	guards        map[string]*ApprovalGuard

	mu          sync.Mutex
	failuresPer map[string]int // open_id → consecutive failed action attempts
}

// NewApprovalServer returns a server bound to the given listen address (e.g.
// ":9082"). signingSecret is the Lark workspace's webhook verification
// secret; pass "" to disable signature checking (insecure — only for local
// testing). guards keys MUST match the `cashier` value embedded in each card.
func NewApprovalServer(listenAddr, signingSecret string, guards map[string]*ApprovalGuard) *ApprovalServer {
	return &ApprovalServer{
		listenAddr:    listenAddr,
		signingSecret: signingSecret,
		guards:        guards,
		failuresPer:   make(map[string]int),
	}
}

// Start binds the listen address and launches the HTTP server in a goroutine.
// It returns an error if the port cannot be bound so the caller can fail fast.
func (s *ApprovalServer) Start() error {
	if s == nil || s.listenAddr == "" {
		return nil
	}
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/lark/callback", s.handleCallback)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok")
	})
	srv := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		log.Printf("approval server listening on %s\n", s.listenAddr)
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("approval server stopped: %v\n", err)
		}
	}()
	return nil
}

func (s *ApprovalServer) handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		http.Error(w, "read body", http.StatusBadRequest)
		return
	}

	timestamp := r.Header.Get("X-Lark-Request-Timestamp")
	nonce := r.Header.Get("X-Lark-Request-Nonce")
	gotSig := r.Header.Get("X-Lark-Signature")

	// Lark URL verification handshake: POST {"type":"url_verification","challenge":"...","token":"..."}.
	// Reply with the challenge before any signature check so the workspace can
	// register the URL.
	var probe struct {
		Type      string `json:"type"`
		Challenge string `json:"challenge"`
	}
	if json.Unmarshal(body, &probe) == nil && probe.Type == "url_verification" {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"challenge": probe.Challenge})
		return
	}

	if s.signingSecret != "" {
		if !util.VerifyLarkCallbackSignature(s.signingSecret, timestamp, nonce, body, gotSig) {
			log.Printf("approval server: signature mismatch from %s\n", r.RemoteAddr)
			http.Error(w, "signature mismatch", http.StatusUnauthorized)
			return
		}
	}

	// Decode + replay protection
	cb, err := util.DecodeLarkCardCallback(body)
	if err != nil {
		http.Error(w, "bad payload", http.StatusBadRequest)
		return
	}
	if cb.Timestamp != 0 && !util.FreshTimestamp(cb.Timestamp, time.Now(), 5*time.Minute) {
		http.Error(w, "stale request", http.StatusBadRequest)
		return
	}
	guard, ok := s.guards[cb.Cashier]
	if !ok {
		http.Error(w, "unknown cashier", http.StatusNotFound)
		return
	}
	if guard.SeenNonce(nonce) {
		http.Error(w, "duplicate", http.StatusConflict)
		return
	}

	var acted bool
	var actErr error
	switch cb.Action {
	case "approve":
		acted, actErr = guard.Approve(cb)
	case "reject":
		acted, actErr = guard.Reject(cb)
	case "mute":
		acted, actErr = guard.Mute(cb)
	default:
		http.Error(w, "unknown action", http.StatusBadRequest)
		return
	}
	if actErr != nil {
		s.recordFailure(cb.OpenID)
		respondToast(w, false, fmt.Sprintf("action failed: %v", actErr))
		return
	}
	s.clearFailures(cb.OpenID)
	switch cb.Action {
	case "mute":
		if !acted {
			respondToast(w, true, "height already muted")
		} else {
			respondToast(w, true, "height muted")
		}
	case "reject":
		if !acted {
			respondToast(w, false, "transfer already decided (handled by another admin?)")
		} else {
			respondToast(w, true, "transfer rejected")
		}
	default:
		if !acted {
			// Row was not in approval state — already handled by another admin.
			respondToast(w, false, "transfer already decided (handled by another admin?)")
		} else {
			respondToast(w, true, "transfer approved")
		}
	}
}

func (s *ApprovalServer) recordFailure(openID string) {
	if openID == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failuresPer[openID]++
	if s.failuresPer[openID] >= 3 {
		util.Alert(fmt.Sprintf("[witness:approval] open_id=%s has %d consecutive failed approvals",
			openID, s.failuresPer[openID]))
	}
}

func (s *ApprovalServer) clearFailures(openID string) {
	if openID == "" {
		return
	}
	s.mu.Lock()
	delete(s.failuresPer, openID)
	s.mu.Unlock()
}

func respondToast(w http.ResponseWriter, ok bool, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	level := "info"
	if !ok {
		level = "error"
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"toast": map[string]interface{}{
			"type":    level,
			"content": msg,
		},
	})
}
