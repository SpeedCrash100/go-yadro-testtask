package pkg

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestReadTime(t *testing.T) {
	test_cases := []struct {
		in    string // Input string
		fail  bool   // Should fail
		hours uint8  // Expected hours after read (if fail=false)
		mins  uint8  // Expected minutes after read (if fail=false)
	}{
		{"", true, 0, 0},
		{":", true, 0, 0},
		{"00:01", false, 0, 1},
		{"+0:00", true, 0, 0},
		{"00:+0", true, 0, 0},
		{"00:", true, 0, 0},
		{":00", true, 0, 0},
		{"12:00", false, 12, 0},
		{"24:00", true, 0, 0},
		{"00:70", true, 0, 0},
		{"-1:30", true, 0, 0},
		{"12:30 70:30", false, 12, 30},
	}

	for _, tc := range test_cases {
		t.Run("ReadTime: "+tc.in, func(t *testing.T) {
			s := bufio.NewScanner(strings.NewReader(tc.in))
			s.Split(bufio.ScanWords)

			var time Time
			err := time.ReadFrom(s)

			if tc.fail {
				if err == nil {
					t.Errorf("Expected failure of reading for '%s'", tc.in)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected failure of reading for '%s': %s", tc.in, err)
			}

			if tc.hours != time.Hour {
				t.Errorf("Read invalid hours: %d != %d", tc.hours, time.Hour)
			}

			if tc.mins != time.Minutes {
				t.Errorf("Read invalid minutes: %d != %d", tc.mins, time.Minutes)
			}

		})
	}
}

func TestTimeToString(t *testing.T) {
	test_cases := []struct {
		hours uint8  // Hours
		mins  uint8  // Mins
		str   string // Expected string generated
	}{
		{12, 0, "12:00"},
		{12, 30, "12:30"},
		{0, 0, "00:00"},
	}

	for _, tc := range test_cases {
		t.Run("Time To String: "+tc.str, func(t *testing.T) {
			s := fmt.Sprintf("%v", Time{tc.hours, tc.mins})

			if s != tc.str {
				t.Errorf("Invalid time '%s' convertion. Result is '%s'", tc.str, s)
			}

		})
	}
}
