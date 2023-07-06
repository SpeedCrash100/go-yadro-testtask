package pkg

import (
	"io"
)

type State struct {
	table_count uint
	time_start  Time
	time_end    Time
	price       uint

	current_time Time

	client_set         map[string]struct{}
	clients_start_time map[string]Time

	writer io.Writer
}

// Are we know this client(It is in club)
func (s State) Known(client string) bool {
	_, ok := s.client_set[client]
	return ok
}

// Add client to known list
func (s *State) AddClient(client string) {
	s.client_set[client] = struct{}{}
}
