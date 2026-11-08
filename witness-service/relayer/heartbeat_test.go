package relayer

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestHeartbeatVersionField(t *testing.T) {
	for in, want := range map[string]string{
		"v7.1.0":                  "v7.1.0",
		"v7.1.0-3-gabc1234-dirty": "v7.1.0-3-gabc1234-dirty",
		"NoBuildInfo":             "NoBuildInfo",
		"0123456789abcdef0123456789abcdef01234567": "0123456789abcdef0123456789abcdef01234567",
		"":                         "",
		"v7.1.0'; DROP TABLE x;--": "",
		"v7 1":                     "",
		"v7\n":                     "",
		strings.Repeat("a", 64):    strings.Repeat("a", 64),
		strings.Repeat("a", 65):    "",
		"v7.1.0é":                  "",
	} {
		if got := heartbeatVersionField(in); got != want {
			t.Errorf("heartbeatVersionField(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIsHeartbeatWitness(t *testing.T) {
	addr := common.HexToAddress("0x5160").Bytes()
	for _, tc := range []struct {
		name           string
		addr           []byte
		member, loaded bool
		want           bool
	}{
		{"active witness", addr, true, true, true},
		{"not a witness", addr, false, true, false},
		{"witness set unavailable", addr, true, false, false},
		{"ed25519-sized address", make([]byte, 32), true, true, false},
		{"empty", nil, true, true, false},
	} {
		v := &submissionTestValidator{member: tc.member, loaded: tc.loaded}
		if got := isHeartbeatWitness(v, tc.addr); got != tc.want {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}
