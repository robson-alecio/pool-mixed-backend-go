package app

import (
	"database/sql"
	"log"

	kallax "gopkg.in/src-d/go-kallax.v1"
)

//PollHandler ...
//go:generate moq -out pollhandler_moq.go . PollHandler
type PollHandler interface {
	SavePoll(poll Poll) Poll
	FindPollByID(ID kallax.ULID) (*Poll, error)
}

//IPollStore ...
//go:generate moq -out ipollstore_moq.go . IPollStore
type IPollStore interface {
	Save(record *Poll) (updated bool, err error)
	FindOne(q *PollQuery) (*Poll, error)
}

//PollHandlerImpl ...
type PollHandlerImpl struct {
	Store         IPollStore
	OptionHandler PollOptionHandler
}

//NewPollHandler ...
func NewPollHandler(db *sql.DB, optionHandler PollOptionHandler) *PollHandlerImpl {
	return &PollHandlerImpl{
		Store:         NewPollStore(db),
		OptionHandler: optionHandler,
	}
}

//PollOptionHandler ...
//go:generate moq -out polloptionhandler_moq.go . PollOptionHandler
type PollOptionHandler interface {
	SavePollOption(poll PollOption) PollOption
	DeletePollOption(id kallax.ULID) error
	FindPollOptions(id kallax.ULID) ([]*PollOption, error)
	ExistsOption(pollID kallax.ULID, candidate string) (bool, error)
}

//IPollOptionStore ...
//go:generate moq -out ipolloptionstore_moq.go . IPollOptionStore
type IPollOptionStore interface {
	Save(record *PollOption) (updated bool, err error)
	Delete(record *PollOption) error
	FindOne(q *PollOptionQuery) (*PollOption, error)
	FindAll(q *PollOptionQuery) ([]*PollOption, error)
	Count(q *PollOptionQuery) (int64, error)
}

//PollOptionHandlerImpl ...
type PollOptionHandlerImpl struct {
	Store IPollOptionStore
}

//NewPollOptionHandler ...
func NewPollOptionHandler(db *sql.DB) *PollOptionHandlerImpl {
	return &PollOptionHandlerImpl{
		Store: NewPollOptionStore(db),
	}
}

//SavePoll ...
func (h PollHandlerImpl) SavePoll(poll Poll) Poll {
	log.Println("Saving Poll", poll)

	h.Store.Save(&poll)
	return poll
}

//FindPollByID ...
func (h PollHandlerImpl) FindPollByID(ID kallax.ULID) (*Poll, error) {
	query := NewPollQuery().FindByID(ID)
	poll, err := h.Store.FindOne(query)

	if err != nil {
		return poll, err
	}

	options, errOption := h.OptionHandler.FindPollOptions(ID)
	if errOption != nil {
		return poll, errOption
	}

	poll.Options = options

	return poll, nil
}

// SavePollOption ...
func (h PollOptionHandlerImpl) SavePollOption(pollOption PollOption) PollOption {
	log.Println("Adding Poll Option", pollOption)

	h.Store.Save(&pollOption)
	return pollOption
}

// DeletePollOption ...
func (h PollOptionHandlerImpl) DeletePollOption(id kallax.ULID) error {
	log.Println("Removing Poll Option", id)

	query := NewPollOptionQuery().FindByID(id)

	opt, err := h.Store.FindOne(query)

	if err != nil {
		return err
	}

	return h.Store.Delete(opt)
}

// FindPollOptions ...
func (h PollOptionHandlerImpl) FindPollOptions(id kallax.ULID) ([]*PollOption, error) {
	query := NewPollOptionQuery().FindByOwner(id)
	return h.Store.FindAll(query)
}

//ExistsOption ...
func (h PollOptionHandlerImpl) ExistsOption(pollID kallax.ULID, candidate string) (bool, error) {
	query := NewPollOptionQuery().
		FindByOwner(pollID).
		FindByContent(candidate)

	count, err := h.Store.Count(query)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
