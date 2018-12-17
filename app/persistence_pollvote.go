package app

import (
	"database/sql"
	"log"

	"gopkg.in/src-d/go-kallax.v1"
)

//PollVoteHandler ...
//go:generate moq -out pollvotehandler_moq.go . PollVoteHandler
type PollVoteHandler interface {
	PollAlreadyVotedByUser(pollID kallax.ULID, userID kallax.ULID) (bool, error)
	VotesFor(pollID kallax.ULID, option string) int64
	SaveVote(vote PollVote) PollVote
}

//IPollVoteStore ...
//go:generate moq -out ipollvotestore_moq.go . IPollVoteStore
type IPollVoteStore interface {
	Save(record *PollVote) (updated bool, err error)
	// FindOne(q *PollVoteQuery) (*PollVote, error)
	Count(q *PollVoteQuery) (int64, error)
}

//PollVoteHandlerImpl ...
type PollVoteHandlerImpl struct {
	Store IPollVoteStore
}

//NewPollVoteHandler ...
func NewPollVoteHandler(db *sql.DB) *PollVoteHandlerImpl {
	return &PollVoteHandlerImpl{
		Store: NewPollVoteStore(db),
	}
}

//PollAlreadyVotedByUser ...
func (h PollVoteHandlerImpl) PollAlreadyVotedByUser(pollID, userID kallax.ULID) (bool, error) {
	query := NewPollVoteQuery().
		FindByPollID(pollID).
		FindByUserID(userID)

	count, err := h.Store.Count(query)

	return count > 0, err
}

//SaveVote ...
func (h PollVoteHandlerImpl) SaveVote(vote PollVote) PollVote {
	log.Println("Registering vote", vote)

	h.Store.Save(&vote)

	return vote
}

//VotesFor ...
func (h PollVoteHandlerImpl) VotesFor(pollID kallax.ULID, option string) int64 {
	query := NewPollVoteQuery().FindByPollID(pollID).FindByChosenOption(option)

	votesOption, err := h.Store.Count(query)
	if err != nil {
		return 0
	}

	return votesOption
}
