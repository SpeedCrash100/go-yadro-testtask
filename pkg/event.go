package pkg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

func NewInputEvent(description string, state State) (InputEvent, error) {
	pieces := strings.Split(description, " ")

	if len(pieces) < 3 {
		// 3 pieces minimum: time, id, client
		return nil, ErrInvalidEventFormat
	}

	id, err := strconv.Atoi(pieces[1])
	if err != nil {
		return nil, err
	}

	time, err := MakeTime(pieces[0])
	if err != nil {
		return nil, err
	}

	client := pieces[2]

	for _, ch := range client {
		if !((unicode.IsLetter(ch) && unicode.IsLower(ch)) || unicode.IsDigit(ch) || ch == '_') {
			return nil, ErrInvalidEventFormat
		}
	}

	remaining_pieces := pieces[3:]

	var event InputEvent

	switch id {
	case EVENT_ID_IN_CLIENT_ENTERED:
		event = NewClientEnteredInputEvent(time, client)
	case EVENT_ID_IN_CLIENT_TAKE_A_SEAT:
		event, err = NewClientTakeASeatInputEvent(time, client, remaining_pieces, state)
	case EVENT_ID_IN_CLIENT_CLIENT_WAITING:
		event = NewClientWaitingInputEvent(time, client)
	case EVENT_ID_IN_CLIENT_LEFT:
		event = NewClientLeftInputEvent(time, client)
	}

	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, ErrUnknownEventType
	}

	return event, nil
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
type InputEvent interface {
	Event

	// Change state according to event
	Translate(*State)
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

// Base event that must hold client
type ClientAssociatedEvent struct {
	BaseEvent
	client string
}

func MakeClientAssociatedEvent(id int, time Time, client string) ClientAssociatedEvent {
	return ClientAssociatedEvent{BaseEvent{time, id}, client}
}

func (e ClientAssociatedEvent) String() string {
	return e.BaseEvent.String() + " " + e.client
}

// Client Entered INPUT event
type ClientEnteredInputEvent struct {
	ClientAssociatedEvent
}

func NewClientEnteredInputEvent(time Time, client string) InputEvent {
	return &ClientEnteredInputEvent{MakeClientAssociatedEvent(EVENT_ID_IN_CLIENT_ENTERED, time, client)}
}

func (e *ClientEnteredInputEvent) Translate(s *State) {

	if s.Known(e.client) {
		error_event := NewErrorOutputEvent(e, MSG_CLIENT_HAS_ALREADY_IN_CLUB)
		s.events = append(s.events, error_event)
		return
	}

	if !e.Time().Between(s.time_start, s.time_end) {
		error_event := NewErrorOutputEvent(e, MSG_CLIENT_HAS_ARRIVED_NOT_IN_TIME)
		s.events = append(s.events, error_event)
		return
	}

	s.AddClient(e.client)
}

// Client take a seat
type ClientTakeASeatInputEvent struct {
	ClientAssociatedEvent
	table_nmb uint
}

func NewClientTakeASeatInputEvent(time Time, client string, parts []string, state State) (InputEvent, error) {
	if len(parts) != 1 {
		return nil, ErrInvalidEventFormat
	}

	table_nmb, err := strconv.ParseUint(parts[0], 10, 0)
	if err != nil {
		return nil, err
	}

	if table_nmb == 0 || state.table_count < uint(table_nmb) {
		return nil, ErrInvalidEventFormat
	}

	return &ClientTakeASeatInputEvent{MakeClientAssociatedEvent(EVENT_ID_IN_CLIENT_TAKE_A_SEAT, time, client), uint(table_nmb)}, nil
}

func (e *ClientTakeASeatInputEvent) Translate(s *State) {
	if s.TableBusy(e.table_nmb) {
		error_event := NewErrorOutputEvent(e, MSG_PLACE_IS_BUSY)
		s.events = append(s.events, error_event)
		return
	}

	if !s.Known(e.client) {
		error_event := NewErrorOutputEvent(e, MSG_CLIENT_UNKNOWN)
		s.events = append(s.events, error_event)
		return
	}

	s.OccupyTable(e.table_nmb, e.client)
}

func (e *ClientTakeASeatInputEvent) String() string {
	return e.ClientAssociatedEvent.String() + " " + fmt.Sprintf("%d", e.table_nmb)
}

type ClientWaitingInputEvent struct {
	ClientAssociatedEvent
}

func NewClientWaitingInputEvent(time Time, client string) InputEvent {
	return &ClientWaitingInputEvent{MakeClientAssociatedEvent(EVENT_ID_IN_CLIENT_CLIENT_WAITING, time, client)}
}

func (e *ClientWaitingInputEvent) Translate(s *State) {

	if s.HaveEmptyTable() {
		error_event := NewErrorOutputEvent(e, MSG_WAITING_WHILE_HAVE_FREE_SPACE)
		s.events = append(s.events, error_event)
		return
	}

	if s.queue.IsFull() {
		s.ClientLeave(e.client)
		event := NewClientLeftOutputEvent(s.current_time, e.client)
		s.events = append(s.events, event)
		return
	}

	s.queue.Push(e.client)

}

type ClientLeftInputEvent struct {
	ClientAssociatedEvent
}

func NewClientLeftInputEvent(time Time, client string) InputEvent {
	return &ClientLeftInputEvent{MakeClientAssociatedEvent(EVENT_ID_IN_CLIENT_LEFT, time, client)}
}

func (e *ClientLeftInputEvent) Translate(s *State) {
	if !s.Known(e.client) {
		error_event := NewErrorOutputEvent(e, MSG_CLIENT_UNKNOWN)
		s.events = append(s.events, error_event)
		return
	}

	freed_table, err := s.ClientLeave(e.client)
	if err != nil {
		return
	}

	if !s.queue.IsEmpty() {
		client_to_place, _ := s.queue.Pop()
		s.OccupyTable(freed_table, client_to_place)
		occupy_event := NewClientTakenSeatOutputEvent(s.current_time, client_to_place, int(freed_table))
		s.events = append(s.events, occupy_event)
	}

}

type ClientLeftOutputEvent struct {
	ClientAssociatedEvent
}

func NewClientLeftOutputEvent(time Time, client string) Event {
	return &ClientLeftOutputEvent{MakeClientAssociatedEvent(EVENT_ID_OUT_CLIENT_LEFT, time, client)}
}

type ClientTakenSeatOutputEvent struct {
	ClientAssociatedEvent
	table_nmb uint
}

func NewClientTakenSeatOutputEvent(time Time, client string, table_nmb int) Event {
	return &ClientTakenSeatOutputEvent{MakeClientAssociatedEvent(EVENT_ID_OUT_CLIENT_LEFT, time, client), uint(table_nmb)}
}

func (e *ClientTakenSeatOutputEvent) String() string {
	return e.ClientAssociatedEvent.String() + " " + fmt.Sprintf("%d", e.table_nmb)
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
