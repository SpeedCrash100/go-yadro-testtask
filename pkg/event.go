package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

// Event ids
const (
	EVENT_ID_UNKNOWN = 0

	EVENT_ID_IN_CLIENT_ENTERED        = 1
	EVENT_ID_IN_CLIENT_TAKE_A_SEAT    = 2
	EVENT_ID_IN_CLIENT_CLIENT_WAITING = 3
	EVENT_ID_IN_CLIENT_LEFT           = 4

	EVENT_ID_OUT_CLIENT_LEFT        = 11
	EVENT_ID_OUT_CLIENT_TAKE_A_SEAT = 12
	EVENT_ID_OUT_ERROR              = 13
)

// Error messages
const (
	MSG_CLIENT_HAS_ALREADY_IN_CLUB     = "YouShallNotPass"
	MSG_CLIENT_HAS_ARRIVED_NOT_IN_TIME = "NotOpenYet"
	MSG_PLACE_IS_BUSY                  = "PlaceIsBusy"
	MSG_CLIENT_UNKNOWN                 = "ClientUnknown"
	MSG_WAITING_WHILE_HAVE_FREE_SPACE  = "ICanWaitNoLonger!"
)

var (
	ErrInvalidEventFormat = errors.New("invalid event format")
	ErrUnknownEventType   = errors.New("invalid event type")
)

func ReadEvent(s *bufio.Scanner) (EventInput, error) {
	var base BaseEvent

	if err := base.Read(s); err != nil {
		return nil, err
	}

	id := base.Id()

	switch id {
	case EVENT_ID_IN_CLIENT_ENTERED:
		client_event := MakeClientAssociatedEvent(base)
		if err := client_event.Read(s); err != nil {
			return nil, err
		}
		return NewClientEnteredEvent(client_event), nil

	}
	return nil, ErrUnknownEventType
}

// Base Event interface
type Event interface {
	// We need to write events
	fmt.Stringer

	// Get time of event
	Time() Time

	// Get Event ID
	Id() int
}

// Input events can change state and
type EventInput interface {
	Event

	// Change state according to event
	Translate(*State)

	// Read event additional information from scanner
	Read(*bufio.Scanner) error
}

// Base event for all events
type BaseEvent struct {
	time Time
	id   int
}

func (e BaseEvent) String() string {
	return e.time.String() + " " + strconv.Itoa(e.id)
}

func (e BaseEvent) Time() Time {
	return e.time
}

func (e BaseEvent) Id() int {
	return e.id
}

func (e *BaseEvent) Read(s *bufio.Scanner) error {
	if err := e.time.ReadFrom(s); err != nil {
		return err
	}

	if !s.Scan() {
		if s.Err() == nil {
			// EOF
			return ErrInvalidEventFormat
		}
		return s.Err()
	}

	id_str := s.Text()
	id, err := strconv.Atoi(id_str)
	if err != nil {
		return err
	}

	e.id = id

	return nil
}

// Base event that must hold client
type ClientAssociatedEvent struct {
	BaseEvent
	client string
}

func MakeClientAssociatedEvent(parent BaseEvent) ClientAssociatedEvent {
	return ClientAssociatedEvent{parent, ""}
}

func (e ClientAssociatedEvent) String() string {
	return e.BaseEvent.String() + " " + e.client
}

func (e *ClientAssociatedEvent) Read(s *bufio.Scanner) error {

	if !s.Scan() {
		if s.Err() == nil {
			// EOF
			return ErrInvalidEventFormat
		}
		return s.Err()
	}

	client := s.Text()
	client_chars := []rune(client)

	for i := 0; i < len(client_chars); i++ {
		ch := client_chars[i]
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
			continue
		}
		return ErrInvalidEventFormat
	}

	e.client = client

	return nil
}

// Client Entered INPUT event
type ClientEnteredInputEvent struct {
	ClientAssociatedEvent
}

func NewClientEnteredEvent(parent ClientAssociatedEvent) EventInput {
	out := ClientEnteredInputEvent{parent}
	return &out
}

func (e *ClientEnteredInputEvent) Translate(s *State) {
	if s.Known(e.client) {
		error_event := NewErrorOutputEvent(e, MSG_CLIENT_HAS_ALREADY_IN_CLUB)
		fmt.Fprintln(s.writer, error_event)
		return
	}

	if !e.Time().Between(s.time_start, s.time_end) {
		error_event := NewErrorOutputEvent(e, MSG_CLIENT_HAS_ARRIVED_NOT_IN_TIME)
		fmt.Fprintln(s.writer, error_event)
		return
	}

	s.AddClient(e.client)
}

// Error event
type ErrorOutputEvent struct {
	BaseEvent
	message string
}

func NewErrorOutputEvent(parent Event, msg string) Event {
	return ErrorOutputEvent{
		BaseEvent{
			time: parent.Time(),
			id:   EVENT_ID_OUT_ERROR,
		},
		msg,
	}
}

func (e ErrorOutputEvent) String() string {
	return e.BaseEvent.String() + " " + e.message
}
