package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidTimeFormat = errors.New("invalid time in input files")
	ErrTimeOutOfRange    = errors.New("time format valid but values are out of range")
)

type Time struct {
	Hour    uint8
	Minutes uint8
}

func (t *Time) ReadFrom(reader *bufio.Scanner) error {
	if !reader.Scan() {
		if reader.Err() == nil {
			// EOF
			return errors.New("empty string")
		}
		return reader.Err()
	}

	text := reader.Text()
	digits := strings.Split(text, ":")
	if len(digits) != 2 {
		return ErrInvalidTimeFormat
	}

	hours_str := digits[0]
	mins_str := digits[1]

	// If there are no leading zeros
	if len(hours_str) != 2 || len(mins_str) != 2 {
		return ErrInvalidTimeFormat
	}

	// Drop strings like "+4"
	if !unicode.IsDigit(rune(hours_str[0])) || !unicode.IsDigit(rune(mins_str[0])) {
		return ErrInvalidTimeFormat
	}

	hours, err := strconv.Atoi(hours_str)
	if err != nil {
		return err
	}

	if hours < 0 || 24 <= hours {
		return ErrTimeOutOfRange
	}

	mins, err := strconv.Atoi(mins_str)
	if err != nil {
		return err
	}

	if mins < 0 || 60 <= mins {
		return ErrTimeOutOfRange
	}

	t.Hour = uint8(hours)
	t.Minutes = uint8(mins)

	return nil
}

func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d", t.Hour, t.Minutes)
}
