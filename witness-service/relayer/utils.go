package relayer

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcwallet/wallet/txauthor"
	"github.com/btcsuite/btcwallet/wallet/txsizes"
	"github.com/btcsuite/btcwallet/wtxmgr"
)

// byAmount defines the methods needed to satisify sort.Interface to
// sort credits by their output amount.
type byAmount []wtxmgr.Credit

func (s byAmount) Len() int           { return len(s) }
func (s byAmount) Less(i, j int) bool { return s[i].Amount < s[j].Amount }
func (s byAmount) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func makeChangeScriptSource(addr *btcutil.AddressTaproot) *txauthor.ChangeSource {
	newChangeScript := func() ([]byte, error) {
		return txscript.PayToAddrScript(addr)
	}
	return &txauthor.ChangeSource{
		ScriptSize: txsizes.P2TRPkScriptSize,
		NewScript:  newChangeScript,
	}
}
