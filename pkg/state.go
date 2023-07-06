package pkg

import "bufio"

type State struct {
	table_count uint
	time_start  Time
	time_end    Time
	price       uint

	clients_start_time map[string]Time

	writer *bufio.Writer
}
