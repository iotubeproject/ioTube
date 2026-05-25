// Copyright (c) 2026 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package util

import (
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
