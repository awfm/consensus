package integration

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/alvalor/consensus"
	"github.com/alvalor/consensus/cache"
	"github.com/alvalor/consensus/message"
	"github.com/alvalor/consensus/mocks"
	"github.com/alvalor/consensus/model"
	"github.com/alvalor/consensus/test/fixture"
)

type Participant struct {

	// participant configuration
	log            zerolog.Logger
	selfID         model.Hash
	participantIDs []model.Hash
	genesisID      model.Hash
	round          uint64
	stop           []Condition
	ignore         []error

	// state data
	vertexDB   map[model.Hash]*model.Vertex
	proposalDB map[model.Hash]*message.Proposal
	voteCache  map[model.Hash](map[model.Hash]*message.Vote)
	proposalQ  chan *message.Proposal
	voteQ      chan *message.Vote

	// real dependencies
	pcache *cache.ProposalCache
	vcache *cache.VoteCache

	// dependency mocks
	net    *mocks.Network
	state  *mocks.State
	build  *mocks.Builder
	sign   *mocks.Signer
	verify *mocks.Verifier

	// processor instance
	pro *consensus.Processor
}

func NewParticipant(t require.TestingT, options ...Option) *Participant {

	// configure time function for zerolog
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }

	// initialize the default values for the participant
	selfID := fixture.Hash(t)
	p := Participant{

		log:            zerolog.New(ioutil.Discard),
		selfID:         selfID,
		participantIDs: []model.Hash{selfID},
		genesisID:      model.ZeroHash,
		round:          0,
		stop:           []Condition{AfterRound(10, errFinished)},
		ignore:         []error{consensus.ObsoleteProposal{}, consensus.ObsoleteVote{}},

		vertexDB:   make(map[model.Hash]*model.Vertex),
		proposalDB: make(map[model.Hash]*message.Proposal),
		voteCache:  make(map[model.Hash](map[model.Hash]*message.Vote)),
		proposalQ:  make(chan *message.Proposal, 1024),
		voteQ:      make(chan *message.Vote, 1024),

		pcache: cache.NewProposalCache(),
		vcache: cache.NewVoteCache(),

		net:    &mocks.Network{},
		state:  &mocks.State{},
		build:  &mocks.Builder{},
		sign:   &mocks.Signer{},
		verify: &mocks.Verifier{},
	}

	// apply the options
	for _, option := range options {
		option(&p)
	}

	// program round-robin state behaviour
	p.state.On("Round").Return(
		func() uint64 {
			return p.round
		},
		nil,
	)
	p.state.On("Set", mock.Anything).Return(
		func(height uint64) error {
			if height <= p.round {
				return fmt.Errorf("must set higher round (set: %d, round: %d)", height, p.round)
			}
			p.round = height
			return nil
		},
	)
	p.state.On("Leader", mock.Anything).Return(
		func(height uint64) model.Hash {
			src := rand.NewSource(int64(height))
			r := rand.New(src)
			index := r.Intn(len(p.participantIDs))
			leader := p.participantIDs[index]
			return leader
		},
		nil,
	)
	p.state.On("Threshold", mock.Anything).Return(
		func() uint {
			threshold := uint(len(p.participantIDs) * 2 / 3)
			return threshold
		},
		nil,
	)

	// program loopback network mock behaviour
	p.net.On("Broadcast", mock.Anything).Return(
		func(proposal *message.Proposal) error {
			p.proposalQ <- proposal
			vertexID := proposal.Vertex.ID()
			p.log.Debug().
				Str("message", "proposal").
				Uint64("height", proposal.Height).
				Hex("vertex", vertexID[:]).
				Hex("parent", proposal.Parent.VertexID[:]).
				Hex("arc", proposal.ArcID[:]).
				Hex("proposer", proposal.SignerID[:]).
				Msg("proposal looped")
			return nil
		},
	)
	p.net.On("Transmit", mock.Anything, mock.Anything).Return(
		func(vote *message.Vote, recipientID model.Hash) error {
			if recipientID != p.selfID {
				return fmt.Errorf("invalid recipient (%x)", recipientID)
			}
			p.voteQ <- vote
			p.log.Debug().
				Str("message", "vote").
				Uint64("height", vote.Height).
				Hex("vertex", vote.VertexID[:]).
				Hex("voter", vote.SignerID[:]).
				Msg("vote looped")
			return nil
		},
	)

	// program random builder behaviour
	p.build.On("Arc").Return(
		func() model.Hash {
			return fixture.Hash(t)
		},
		nil,
	)

	// program no-signature signer behaviour
	p.sign.On("Self").Return(
		func() model.Hash {
			return p.selfID
		},
		nil,
	)
	p.sign.On("Proposal", mock.Anything).Return(
		func(vertex *model.Vertex) *message.Proposal {
			vertex.SignerID = p.selfID
			proposal := message.Proposal{
				Vertex:    vertex,
				Signature: nil,
			}
			return &proposal
		},
		nil,
	)
	p.sign.On("Vote", mock.Anything).Return(
		func(vertex *model.Vertex) *message.Vote {
			vote := message.Vote{
				Height:    vertex.Height,
				VertexID:  vertex.ID(),
				SignerID:  p.selfID,
				Signature: nil,
			}
			return &vote
		},
		nil,
	)

	// program always-valid verifier behaviour
	p.verify.On("Proposal", mock.Anything).Return(nil)
	p.verify.On("Vote", mock.Anything).Return(nil)

	// inject dependencies into processor
	p.pro = consensus.NewProcessor(p.state, p.net, p.build, p.sign, p.verify, p.pcache, p.vcache)

	return &p
}

func (p *Participant) Run() error {

	// create the ignore function to easily check errors
	ignore := Ignore(Combine(p.ignore...))

	for {

		// check stop condition first
		for _, condition := range p.stop {
			err := condition(p)
			if err != nil {
				return fmt.Errorf("encountered stop condition: %w", err)
			}
		}

		// set up logging entry and error
		var err error

		// wait for a message to process
		select {
		case proposal := <-p.proposalQ:
			vertexID := proposal.Vertex.ID()
			p.log.Debug().
				Uint64("height", proposal.Height).
				Hex("vertex", vertexID[:]).
				Hex("parent", proposal.Parent.VertexID[:]).
				Hex("arc", proposal.ArcID[:]).
				Hex("proposer", proposal.SignerID[:]).
				Msg("proposal received")
			err = p.pro.OnProposal(proposal)
		case vote := <-p.voteQ:
			p.log.Debug().
				Uint64("height", vote.Height).
				Hex("vertex", vote.VertexID[:]).
				Hex("voter", vote.SignerID[:]).
				Msg("vote received")
			err = p.pro.OnVote(vote)
		case <-time.After(100 * time.Millisecond):
			continue
		}

		// check if we should log a message
		if err == nil || ignore(err) {
			continue
		}

		return fmt.Errorf("encountered critical error: %w", err)
	}
}
