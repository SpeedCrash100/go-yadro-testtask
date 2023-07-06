package pkg

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

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
			s := bufio.NewScanner(strings.NewReader(tc.in))
			s.Split(bufio.ScanWords)
			err := e.Read(s)
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
		{"12:32", true},
		{"12:32 ", true},
		{"12:32 1 client_1", false},
		{"12:32 1 client()_1", true},
		{"12:62 1", true},
	}

	for _, tc := range test_cases {
		t.Run(tc.in, func(t *testing.T) {
			var e ClientAssociatedEvent
			s := bufio.NewScanner(strings.NewReader(tc.in))
			s.Split(bufio.ScanWords)
			err := e.Read(s)
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
