package integration

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/alvalor/consensus/message"
	"github.com/alvalor/consensus/mocks"
	"github.com/alvalor/consensus/model"
)

func Network(participants []*Participant) {

	// create a list & a registry for participants
	participantIDs := make([]model.Hash, 0, len(participants))
	registry := make(map[model.Hash]*Participant)
	for _, p := range participants {
		participantIDs = append(participantIDs, p.selfID)
		registry[p.selfID] = p
	}

	// update the set of participants for each participant and update the
	// network mock to properly connect them
	for i := range participants {
		sender := participants[i]

		sender.participantIDs = participantIDs
		*sender.net = mocks.Network{}
		sender.net.On("Broadcast", mock.Anything).Return(
			func(proposal *message.Proposal) error {
				for _, receiver := range participants {
					receiver.proposalQ <- proposal
				}
				return nil
			},
		)
		sender.net.On("Transmit", mock.Anything, mock.Anything).Return(
			func(vote *message.Vote, recipientID model.Hash) error {
				receiver, exists := registry[recipientID]
				if !exists {
					return fmt.Errorf("invalid recipient (%x)", recipientID)
				}
				receiver.voteQ <- vote
				return nil
			},
		)
	}
}
