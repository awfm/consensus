package fixture

import (
	"testing"

	"github.com/alvalor/consensus/model/message"
)

func Quorum(t testing.TB) *message.Quorum {
	quorum := message.Quorum{
		SignerIDs: Hashes(t, 3),
		Signature: Sig(t),
	}
	return &quorum
}
