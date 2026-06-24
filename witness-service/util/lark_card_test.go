// Copyright (c) 2026 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package util

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeLarkCardCallback_ParsesAction(t *testing.T) {
	cases := []struct {
		name       string
		action     string
		wantAction string
	}{
		{"approve button", "approve", "approve"},
		{"reject button", "reject", "reject"},
		{"missing action", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actionField := ""
			if tc.action != "" {
				actionField = `,"action":"` + tc.action + `"`
			}
			body := []byte(`{"event":{"operator":{"open_id":"ou_test"},"action":{"value":{"transferID":"0xdead","cashier":"c1","nonce":"n1"` + actionField + `},"form_value":{}}}}`)
			cb, err := DecodeLarkCardCallback(body)
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if cb.Action != tc.wantAction {
				t.Fatalf("Action: got %q, want %q", cb.Action, tc.wantAction)
			}
			if cb.TransferID != "0xdead" {
				t.Fatalf("TransferID: got %q", cb.TransferID)
			}
			if cb.OpenID != "ou_test" {
				t.Fatalf("OpenID: got %q", cb.OpenID)
			}
			if cb.Nonce != "n1" {
				t.Fatalf("Nonce: got %q", cb.Nonce)
			}
		})
	}
}

func TestDecodeLarkCardCallback_ParsesMuteHeight(t *testing.T) {
	cases := []struct {
		name   string
		height string // raw JSON value (quoted string or bare number)
	}{
		{"height as string", `"28401933"`},
		{"height as number", `28401933`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body := []byte(`{"event":{"operator":{"open_id":"ou_x"},"action":{"value":{"cashier":"c1","height":` +
				tc.height + `,"nonce":"n2","action":"mute"}}}}`)
			cb, err := DecodeLarkCardCallback(body)
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if cb.Action != "mute" {
				t.Fatalf("Action: got %q, want mute", cb.Action)
			}
			if cb.Height != 28401933 {
				t.Fatalf("Height: got %d, want 28401933", cb.Height)
			}
			if cb.Cashier != "c1" {
				t.Fatalf("Cashier: got %q", cb.Cashier)
			}
		})
	}
}

func TestSendLarkStaleWarningCard_EmptyWebhook(t *testing.T) {
	if err := SendLarkStaleWarningCard("", LarkStaleWarning{Height: 1}); err == nil {
		t.Fatal("expected error for empty webhook")
	}
}

func TestSendLarkStaleWarningCard_PostsMuteButton(t *testing.T) {
	var got []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	if err := SendLarkStaleWarningCard(srv.URL, LarkStaleWarning{
		Cashier: "0xcash", Height: 42, Nonce: "n",
	}); err != nil {
		t.Fatalf("send: %v", err)
	}
	s := string(got)
	for _, want := range []string{`"msg_type":"interactive"`, `"action":"mute"`, `"height":"42"`, "Mute"} {
		if !strings.Contains(s, want) {
			t.Fatalf("card payload missing %q:\n%s", want, s)
		}
	}
}

func TestSendLarkStaleWarningCard_MuteButtonRoutesByGuardKey(t *testing.T) {
	// When the stale height belongs to a previous cashier, the body shows the
	// owning cashier but the Mute button must route back to the posting guard.
	var got []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	if err := SendLarkStaleWarningCard(srv.URL, LarkStaleWarning{
		Cashier: "0xPREVIOUS", GuardKey: "0xCURRENT", Height: 9, Nonce: "n",
	}); err != nil {
		t.Fatalf("send: %v", err)
	}
	s := string(got)
	// Mute button routes to the guard key.
	if !strings.Contains(s, `"cashier":"0xCURRENT"`) {
		t.Fatalf("mute button should route by GuardKey 0xCURRENT:\n%s", s)
	}
	// The owning cashier is shown in the card body.
	if !strings.Contains(s, "0xPREVIOUS") {
		t.Fatalf("card body should display owning cashier 0xPREVIOUS:\n%s", s)
	}
}
