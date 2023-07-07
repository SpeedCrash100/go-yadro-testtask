package pkg

import "errors"

type State struct {
	table_count uint
	time_start  Time
	time_end    Time
	price       uint

	current_time Time

	client_set            map[string]struct{}
	clients_current_table map[string]uint

	tables_occupation []string
	tables_profit     []uint
	tables_start_time []Time
	tables_usage      []Time

	queue Queue[string]

	events []Event
}

func MakeState() State {
	return State{
		client_set:            make(map[string]struct{}),
		clients_current_table: make(map[string]uint),
	}
}

func (s *State) InitTables(size uint) {
	s.table_count = size
	s.tables_occupation = make([]string, size)
	s.tables_profit = make([]uint, size)
	s.tables_start_time = make([]Time, size)
	s.tables_usage = make([]Time, size)

	s.queue = NewQueue[string](int(size))
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

// Add client to known list
func (s *State) ClientLeave(client string) (uint, error) {

	if !s.Known(client) {
		return 0, errors.New("unknown client leaving")
	}

	delete(s.client_set, client)
	return s.LeaveTable(client)
}

func (s State) Clients() []string {
	out := make([]string, 0)

	for k := range s.client_set {
		out = append(out, k)
	}

	return out
}

func (s *State) OnClubClose() {
	s.current_time = s.time_end

	for _, cl := range s.Clients() {
		s.LeaveTable(cl)
		event := NewClientLeftOutputEvent(s.time_end, cl)
		s.events = append(s.events, event)
	}
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
	s.tables_start_time[table_id] = s.current_time

	s.clients_current_table[client] = table_id

}

func (s *State) LeaveTable(client string) (uint, error) {

	if table_id, ok := s.clients_current_table[client]; ok {
		delete(s.clients_current_table, client)
		s.tables_occupation[table_id] = ""

		usage := s.current_time.Diff(s.tables_start_time[table_id])
		profit := uint(usage.HoursUp()) * s.price

		s.tables_profit[table_id] += profit
		s.tables_usage[table_id] = s.tables_usage[table_id].Add(usage)

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
