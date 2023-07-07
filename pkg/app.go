package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var (
	ErrEOF                 = errors.New("unexpected end of file")
	ErrInvalidOrderOfEvent = errors.New("invalid order of events")
)

type App struct {
	state State

	input  *bufio.Scanner
	output io.Writer
}

func NewApp(input io.Reader, output io.Writer) App {
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)
	return App{MakeState(), scanner, output}
}

func (app *App) Process() error {
	if err := app.readClubInfo(); err != nil {
		return err
	}

	var prev_time Time

	for app.input.Scan() {
		str := app.input.Text()
		if len(str) == 0 {
			continue
		}

		event, err := NewInputEvent(str)
		if err != nil {
			fmt.Fprintln(app.output, str)
			return err
		}

		// If event AFTER close then we need to generate client left event now
		if app.state.time_end.Less(event.Time()) {
			app.state.OnClubClose()
		}

		// Invalid order of events
		if !prev_time.LessOrEquals(event.Time()) {
			fmt.Fprintln(app.output, event)
			return ErrInvalidOrderOfEvent
		}

		prev_time = event.Time()
		app.state.current_time = event.Time()
		app.state.events = append(app.state.events, event)

		event.Translate(&app.state)

	}

	// Means it wasn't run in loop because there were not event after close
	if app.state.current_time.LessOrEquals(app.state.time_end) {
		app.state.OnClubClose()
	}

	fmt.Fprintln(app.output, app.state.time_start)

	for _, e := range app.state.events {
		fmt.Fprintln(app.output, e)
	}

	fmt.Fprintln(app.output, app.state.time_end)

	for i := 0; i < int(app.state.table_count); i++ {
		fmt.Fprintf(app.output, "%d %d %v\n", i+1, app.state.tables_profit[i], app.state.tables_usage[i])
	}

	return nil
}

func (app *App) readClubInfo() error {
	// First line - Get tables count
	if !app.input.Scan() {
		if app.input.Err() == nil {
			return ErrEOF
		}
		return app.input.Err()
	}

	tables_str := app.input.Text()
	tables_count, err := strconv.ParseUint(tables_str, 10, 0)
	if err != nil {
		fmt.Fprintln(app.output, tables_str)
		return err
	}

	app.state.InitTables(uint(tables_count))

	// Second line - Get open times
	if !app.input.Scan() {
		if app.input.Err() == nil {
			return ErrEOF
		}
		return app.input.Err()
	}

	times_str := app.input.Text()
	times_strs := strings.Split(times_str, " ")
	if len(times_strs) != 2 {
		fmt.Fprintln(app.output, times_str)
		return ErrInvalidTimeFormat
	}

	start_time, err := MakeTime(times_strs[0])
	if err != nil {
		fmt.Fprintln(app.output, tables_str)
		return err
	}
	end_time, err := MakeTime(times_strs[1])
	if err != nil {
		fmt.Fprintln(app.output, tables_str)
		return err
	}

	app.state.time_start = start_time
	app.state.time_end = end_time
	if !(start_time.Less(end_time)) {
		fmt.Fprintln(app.output, tables_str)
		return ErrInvalidTimeFormat
	}

	// Third line - Get tables count
	if !app.input.Scan() {
		if app.input.Err() == nil {
			return ErrEOF
		}
		return app.input.Err()
	}

	price_str := app.input.Text()
	price, err := strconv.ParseUint(price_str, 10, 0)
	if err != nil {
		fmt.Fprintln(app.output, price_str)
		return err
	}
	app.state.price = uint(price)

	return nil

}
