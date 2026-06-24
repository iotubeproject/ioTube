// Copyright (c) 2026 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// LarkApprovalRequest captures the fields rendered to the admin in the Lark
// interactive card. All fields are pre-formatted strings.
type LarkApprovalRequest struct {
	TransferID string
	Cashier    string
	Token      string
	Recipient  string
	Amount     string  // raw token amount in base units
	AmountUSD  string  // human-readable USD equivalent, e.g. "$61,234.56"; empty if unknown
	TxHash     string
	TxHashURL  string  // full explorer URL; if non-empty TxHash is rendered as a hyperlink
	Nonce      string
}

// LarkStaleWarning captures the fields rendered to the admin when a witness
// cannot fetch a stale block's source data (e.g. the source chain/API pruned
// that history). The card carries a single Mute button.
type LarkStaleWarning struct {
	Cashier string
	Height  uint64
	Nonce   string
}

// LarkCallback is the decoded payload from a Lark `card.action.trigger` event
// targeting our approval handler.
type LarkCallback struct {
	OpenID     string
	TransferID string
	Cashier    string
	Token      string
	Tidx       uint64
	// Height is the stale block height carried by a "mute" card.
	Height    uint64
	Nonce     string
	Timestamp int64
	// Action is "approve", "reject", or "mute" — set by the button the admin
	// clicked.
	Action string
}

// SendLarkApprovalCard posts an interactive card to a Lark bot incoming
// webhook. The card carries the transfer details and Approve / Reject buttons.
// The first admin to click either button locks the decision.
func SendLarkApprovalCard(webhook string, req LarkApprovalRequest) error {
	if webhook == "" {
		return fmt.Errorf("lark webhook is empty")
	}
	approveValue := map[string]string{
		"transferID": req.TransferID,
		"cashier":    req.Cashier,
		"token":      req.Token,
		"nonce":      req.Nonce,
		"action":     "approve",
	}
	rejectValue := map[string]string{
		"transferID": req.TransferID,
		"cashier":    req.Cashier,
		"token":      req.Token,
		"nonce":      req.Nonce,
		"action":     "reject",
	}
	card := map[string]interface{}{
		"config": map[string]bool{"wide_screen_mode": true},
		"header": map[string]interface{}{
			"template": "red",
			"title": map[string]string{
				"tag":     "plain_text",
				"content": "Cross-chain transfer approval required",
			},
		},
		"elements": []interface{}{
			map[string]interface{}{
				"tag": "div",
				"text": map[string]string{
					"tag":     "lark_md",
					"content": buildCardBody(req),
				},
			},
			map[string]interface{}{
				"tag": "action",
				"actions": []interface{}{
					map[string]interface{}{
						"tag":  "button",
						"type": "primary",
						"text": map[string]string{
							"tag":     "plain_text",
							"content": "Approve",
						},
						"value": approveValue,
					},
					map[string]interface{}{
						"tag":  "button",
						"type": "danger",
						"text": map[string]string{
							"tag":     "plain_text",
							"content": "Reject",
						},
						"value": rejectValue,
					},
				},
			},
		},
	}
	body, err := json.Marshal(map[string]interface{}{
		"msg_type": "interactive",
		"card":     card,
	})
	if err != nil {
		return err
	}
	resp, err := http.Post(webhook, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("lark card post failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}
	return nil
}

// SendLarkStaleWarningCard posts an interactive card warning that a witness
// cannot fetch a stale block's source data. The single Mute button records the
// (cashier, height) so the witness stops retrying/alerting on that height.
func SendLarkStaleWarningCard(webhook string, req LarkStaleWarning) error {
	if webhook == "" {
		return fmt.Errorf("lark webhook is empty")
	}
	heightStr := strconv.FormatUint(req.Height, 10)
	muteValue := map[string]string{
		"cashier": req.Cashier,
		"height":  heightStr,
		"nonce":   req.Nonce,
		"action":  "mute",
	}
	body := fmt.Sprintf(
		"**Cashier:** %s\n**Block height:** %s\n\n"+
			"The relayer asked this witness to re-pull block **%s**, but the source chain/API no "+
			"longer has that block's data, so this transfer **cannot be signed automatically**.\n\n"+
			"**Action required:** find, sign, and submit it manually (see runbook), then click "+
			"**Mute** to stop these alerts. Clicking **Mute** without signing abandons the transfer "+
			"on this witness.",
		req.Cashier, heightStr, heightStr,
	)
	card := map[string]interface{}{
		"config": map[string]bool{"wide_screen_mode": true},
		"header": map[string]interface{}{
			"template": "orange",
			"title": map[string]string{
				"tag":     "plain_text",
				"content": "⚠️ Witness can't fetch transfer — source data deleted",
			},
		},
		"elements": []interface{}{
			map[string]interface{}{
				"tag": "div",
				"text": map[string]string{
					"tag":     "lark_md",
					"content": body,
				},
			},
			map[string]interface{}{
				"tag": "action",
				"actions": []interface{}{
					map[string]interface{}{
						"tag":  "button",
						"type": "danger",
						"text": map[string]string{
							"tag":     "plain_text",
							"content": "Mute",
						},
						"value": muteValue,
					},
				},
			},
		},
	}
	payload, err := json.Marshal(map[string]interface{}{
		"msg_type": "interactive",
		"card":     card,
	})
	if err != nil {
		return err
	}
	resp, err := http.Post(webhook, "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("lark card post failed: status=%d body=%s", resp.StatusCode, string(respBody))
	}
	return nil
}

// buildCardBody assembles the Lark Markdown body shown to the admin. The TxHash
// is rendered as a hyperlink when TxHashURL is set; AmountUSD is appended to
// the Amount line when present.
func buildCardBody(req LarkApprovalRequest) string {
	amountLine := "**Amount:** " + req.Amount
	if req.AmountUSD != "" {
		amountLine += " (" + req.AmountUSD + ")"
	}
	txHashLine := "**TxHash:** "
	if req.TxHashURL != "" {
		txHashLine += fmt.Sprintf("[%s](%s)", req.TxHash, req.TxHashURL)
	} else {
		txHashLine += req.TxHash
	}
	return fmt.Sprintf(
		"**Cashier:** %s\n**Token:** %s\n**Recipient:** %s\n%s\n%s",
		req.Cashier, req.Token, req.Recipient, amountLine, txHashLine,
	)
}

// VerifyLarkCallbackSignature implements Lark's HMAC-SHA256 signing scheme for
// webhook event subscriptions: signature = base64(HMAC-SHA256(secret,
// timestamp + "\n" + nonce + "\n" + body)). Reject on any mismatch — never
// short-circuit comparison.
func VerifyLarkCallbackSignature(secret, timestamp, nonce string, body []byte, gotB64 string) bool {
	if secret == "" || gotB64 == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("\n"))
	mac.Write([]byte(nonce))
	mac.Write([]byte("\n"))
	mac.Write(body)
	want := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(want), []byte(gotB64))
}

// DecodeLarkCardCallback parses a Lark `card.action.trigger` HTTP body into
// our LarkCallback type. The event JSON shape we care about:
//
//	{
//	  "event": {
//	    "operator": { "open_id": "ou_..." },
//	    "action":   { "value": { "transferID": "...", "cashier": "...", "token": "...", "tidx": 12, "nonce": "...", "action": "approve" } }
//	  }
//	}
//
// Extra fields are ignored.
func DecodeLarkCardCallback(body []byte) (LarkCallback, error) {
	var raw struct {
		Timestamp string `json:"ts"`
		Event     struct {
			Operator struct {
				OpenID string `json:"open_id"`
			} `json:"operator"`
			Action struct {
				Value map[string]interface{} `json:"value"`
			} `json:"action"`
		} `json:"event"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return LarkCallback{}, err
	}
	out := LarkCallback{OpenID: raw.Event.Operator.OpenID}
	if v, ok := raw.Event.Action.Value["transferID"].(string); ok {
		out.TransferID = v
	}
	if v, ok := raw.Event.Action.Value["cashier"].(string); ok {
		out.Cashier = v
	}
	if v, ok := raw.Event.Action.Value["token"].(string); ok {
		out.Token = v
	}
	if v, ok := raw.Event.Action.Value["nonce"].(string); ok {
		out.Nonce = v
	}
	if v, ok := raw.Event.Action.Value["action"].(string); ok {
		out.Action = v
	}
	if v, ok := raw.Event.Action.Value["tidx"]; ok {
		switch n := v.(type) {
		case float64:
			out.Tidx = uint64(n)
		case string:
			parsed, err := strconv.ParseUint(n, 10, 64)
			if err == nil {
				out.Tidx = parsed
			}
		}
	}
	if v, ok := raw.Event.Action.Value["height"]; ok {
		switch n := v.(type) {
		case float64:
			out.Height = uint64(n)
		case string:
			parsed, err := strconv.ParseUint(n, 10, 64)
			if err == nil {
				out.Height = parsed
			}
		}
	}
	if raw.Timestamp != "" {
		if ts, err := strconv.ParseInt(raw.Timestamp, 10, 64); err == nil {
			out.Timestamp = ts
		}
	}
	return out, nil
}

// FreshTimestamp returns true iff `ts` (unix seconds, possibly with millis
// stripped) is within ±maxSkew of `now`.
func FreshTimestamp(ts int64, now time.Time, maxSkew time.Duration) bool {
	if ts == 0 {
		return false
	}
	delta := now.Unix() - ts
	if delta < 0 {
		delta = -delta
	}
	return time.Duration(delta)*time.Second <= maxSkew
}
