package consensus

import (
	"github.com/alvalor/consensus/message"
	"github.com/alvalor/consensus/model"
)

type Buffer interface {
	Tally(vote *message.Vote) error
	Votes(blockID model.Hash) ([]*message.Vote, error)
	Clear(blockID model.Hash) error
}
