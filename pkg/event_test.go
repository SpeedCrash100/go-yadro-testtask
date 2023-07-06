package pkg

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func NewStringScanner(str string) *bufio.Scanner {
	s := bufio.NewScanner(strings.NewReader(str))
	s.Split(bufio.ScanWords)
	return s
}

func TestErrorOutputEventFormat(t *testing.T) {
	parent := BaseEvent{
		time: Time{
			Hour:    14,
			Minutes: 45,
		},
		id: EVENT_ID_IN_CLIENT_ENTERED,
	}

	error_event := NewErrorOutputEvent(parent, MSG_CLIENT_UNKNOWN)

	expected_string := "14:45 13 " + MSG_CLIENT_UNKNOWN + "\n"
	s := fmt.Sprintln(error_event)
	if s != expected_string {
		t.Errorf("%s != %s", expected_string, s)
	}
}

func TestBaseEventReader(t *testing.T) {
	test_cases := []struct {
		in   string
		fail bool
	}{
		{"", true},
		{"12:32", true},
		{"12:32 ", true},
		{"12:32 1", false},
		{"12:62 1", true},
	}

	for _, tc := range test_cases {
		t.Run(tc.in, func(t *testing.T) {
			var e BaseEvent
			err := e.Read(NewStringScanner(tc.in))
			if tc.fail {
				if err == nil {
					t.Errorf("Expected failure '%s'", tc.in)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected failure '%s': '%s'", tc.in, err)
			}
		})
	}

}

func TestClientEventReader(t *testing.T) {
	test_cases := []struct {
		in   string
		fail bool
	}{
		{"", true},
		{"client_1", false},
		{"client()_1", true},
		{"12:62 1", true},
	}

	for _, tc := range test_cases {
		t.Run(tc.in, func(t *testing.T) {
			var e ClientAssociatedEvent
			err := e.Read(NewStringScanner(tc.in))
			if tc.fail {
				if err == nil {
					t.Errorf("Expected failure '%s'", tc.in)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected failure '%s': '%s'", tc.in, err)
			}

		})
	}
}

func TestClientEnteredInputEvent(t *testing.T) {
	var state State
	state.time_start = Time{9, 0}
	state.time_end = Time{20, 0}

	buffer := bytes.NewBufferString("")
	state.writer = buffer

	id_str := strconv.Itoa(EVENT_ID_IN_CLIENT_ENTERED)

	init_string := "08:32 " + id_str + " client1"

	event, err := ReadEvent(NewStringScanner(init_string))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if _, ok := event.(*ClientEnteredInputEvent); !ok {
		t.Fatalf("Invalid event type")
	}

	if event.Time().Hour != 8 {
		t.Errorf("Invalid time")
	}

	event.Translate(&state)

	if state.Known("client1") {
		t.Errorf("client1 must not be known")
	}

	expected_str := "08:32 13 " + MSG_CLIENT_HAS_ARRIVED_NOT_IN_TIME + "\n"
	if expected_str != buffer.String() {
		t.Errorf("'%s' != '%s'", expected_str, buffer.String())
	}

}
