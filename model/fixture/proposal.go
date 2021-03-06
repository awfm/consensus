package fixture

import (
	"testing"

	"github.com/awfm/consensus/model/base"
	"github.com/awfm/consensus/model/message"
)

func Proposal(t testing.TB, options ...func(*message.Proposal)) *message.Proposal {
	proposal := message.Proposal{
		Candidate: Vertex(t),
		Quorum:    Quorum(t),
		Signature: Sig(t),
	}
	for _, option := range options {
		option(&proposal)
	}
	return &proposal
}

func WithCandidate(candidate *base.Vertex) func(*message.Proposal) {
	return func(proposal *message.Proposal) {
		proposal.Candidate = candidate
	}
}
