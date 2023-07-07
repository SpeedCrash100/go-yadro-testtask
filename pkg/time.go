package pkg

import (
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

func MakeTime(description string) (Time, error) {
	var t Time

	digits := strings.Split(description, ":")
	if len(digits) != 2 {
		return t, ErrInvalidTimeFormat
	}

	hours_str := digits[0]
	mins_str := digits[1]

	// If there are no leading zeros
	if len(hours_str) != 2 || len(mins_str) != 2 {
		return t, ErrInvalidTimeFormat
	}

	// Drop strings like "+4"
	if !unicode.IsDigit(rune(hours_str[0])) || !unicode.IsDigit(rune(mins_str[0])) {
		return t, ErrInvalidTimeFormat
	}

	hours, err := strconv.Atoi(hours_str)
	if err != nil {
		return t, err
	}

	if hours < 0 || 24 <= hours {
		return t, ErrTimeOutOfRange
	}

	mins, err := strconv.Atoi(mins_str)
	if err != nil {
		return t, err
	}

	if mins < 0 || 60 <= mins {
		return t, ErrTimeOutOfRange
	}

	t.Hour = uint8(hours)
	t.Minutes = uint8(mins)

	return t, nil
}

func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d", t.Hour, t.Minutes)
}

func (left Time) Less(right Time) bool {
	if left.Hour < right.Hour {
		return true
	} else if left.Hour == right.Hour {
		return left.Minutes < right.Minutes
	}

	return false
}

func (left Time) LessOrEquals(right Time) bool {
	return left.Less(right) || left == right
}

func (t Time) Between(start, end Time) bool {
	return start.LessOrEquals(t) && t.Less(end)
}

func (left Time) Add(right Time) Time {
	minutes := left.Minutes + right.Minutes
	add_to_hours := minutes / 60
	minutes = minutes % 60

	hours := left.Hour + right.Hour + add_to_hours

	return Time{hours, minutes}
}

func (left Time) Diff(right Time) Time {
	hours := left.Hour - right.Hour
	var minutes = left.Minutes - right.Minutes
	if left.Minutes < right.Minutes {
		hours--
		minutes = 60 - (right.Minutes - left.Minutes)
	}

	return Time{hours, minutes}
}

func (left Time) HoursUp() uint8 {
	hours := left.Hour
	if 0 < left.Minutes {
		hours++
	}

	return hours
}
