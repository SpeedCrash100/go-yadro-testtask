package pkg

import "errors"

type State struct {
	table_count uint
	time_start  Time
	time_end    Time
	price       uint

	current_time Time

	client_set            map[string]struct{}
	clients_start_time    map[string]Time
	clients_current_table map[string]uint
	clients_profit        map[string]int

	tables_occupation []string

	queue Queue[string]

	events []Event
}

func MakeState() State {
	return State{
		client_set:            make(map[string]struct{}),
		clients_start_time:    make(map[string]Time),
		clients_current_table: make(map[string]uint),
		clients_profit:        make(map[string]int),
	}
}

// Are we know this client(It is in club)
func (s State) Known(client string) bool {
	_, ok := s.client_set[client]
	return ok
}

// Add client to known list
func (s *State) AddClient(client string) {
	s.client_set[client] = struct{}{}

	if _, ok := s.clients_profit[client]; !ok {
		s.clients_profit[client] = 0
	}

}

// Add client to known list
func (s *State) ClientLeave(client string) (uint, error) {

	if !s.Known(client) {
		return 0, errors.New("unknown client leaving")
	}

	if start_time, ok := s.clients_start_time[client]; ok {
		diff := int(s.current_time.Diff(start_time).HoursUp())
		s.clients_profit[client] += diff * int(s.price)
	}

	delete(s.client_set, client)
	return s.LeaveTable(client)
}

func (s State) TableBusy(number uint) bool {
	if s.table_count < number || s.table_count == 0 {
		panic("table number out of range")
	}

	return len(s.tables_occupation[number-1]) != 0
}

func (s *State) OccupyTable(number uint, client string) {
	if s.TableBusy(number) {
		panic("table is busy")
	}

	table_id := number - 1

	s.tables_occupation[table_id] = client
	s.clients_current_table[client] = table_id

	if _, ok := s.clients_start_time[client]; !ok {
		s.clients_start_time[client] = s.current_time
	}
}

func (s *State) LeaveTable(client string) (uint, error) {

	if table_id, ok := s.clients_current_table[client]; ok {
		delete(s.clients_current_table, client)
		s.tables_occupation[table_id] = ""
		return table_id + 1, nil
	}

	return 0, errors.New("client don't need to leave table")
}

func (s State) HaveEmptyTable() bool {
	for table_nmb := uint(1); table_nmb <= s.table_count; table_nmb++ {
		if !s.TableBusy(table_nmb) {
			return true
		}
	}
	return false
}
