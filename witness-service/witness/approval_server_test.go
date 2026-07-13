// Copyright (c) 2026 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestApprovalServer_MuteCallback drives a Lark "mute" card callback end-to-end
// through the HTTP handler and asserts the height is muted on the routed guard.
func TestApprovalServer_MuteCallback(t *testing.T) {
	rec := newFakeRecorder(map[string]int64{})
	g := NewApprovalGuard("c", time.Hour, nil, nil, nil, newFakePrices(nil), rec, "", "")
	g.larkAlerter = func(string) {}
	// signingSecret "" disables signature verification for the test.
	s := NewApprovalServer(":0", "", map[string]*ApprovalGuard{"c": g})

	body := `{"event":{"operator":{"open_id":"ou_a"},"action":{"value":{"cashier":"c","height":"77","nonce":"cn","action":"mute"}}}}`
	req := httptest.NewRequest(http.MethodPost, "/lark/callback", strings.NewReader(body))
	w := httptest.NewRecorder()

	s.handleCallback(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if !g.IsHeightMuted(77) {
		t.Fatal("height 77 should be muted after the callback")
	}
	if !strings.Contains(w.Body.String(), "muted") {
		t.Fatalf("expected a 'muted' toast, got %s", w.Body.String())
	}

	// A second identical callback is idempotent — still muted, no error toast.
	req2 := httptest.NewRequest(http.MethodPost, "/lark/callback", strings.NewReader(body))
	w2 := httptest.NewRecorder()
	s.handleCallback(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("second callback status=%d body=%s", w2.Code, w2.Body.String())
	}
	if !g.IsHeightMuted(77) {
		t.Fatal("height 77 should remain muted")
	}
}
