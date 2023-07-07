package pkg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

type testCase struct {
	name string
	in   *os.File
	out  *os.File
}

// Compares 2 stream line by line
// Returns error if they differs or cannot be read
func compareReaders(left, right io.Reader) error {
	scanner_left := bufio.NewScanner(left)
	scanner_right := bufio.NewScanner(right)

	line_count := 1
	for scanner_left.Scan() && scanner_right.Scan() {
		left_line := scanner_left.Text()
		right_line := scanner_right.Text()

		if left_line != right_line {
			return fmt.Errorf("file differs on line: %d", line_count)
		}

		line_count++
	}

	if scanner_left.Err() != nil {
		return scanner_left.Err()
	}

	if scanner_right.Err() != nil {
		return scanner_right.Err()
	}

	// Means either left or right reached EOF
	if scanner_left.Err() == nil && scanner_right.Err() == nil {
		// Means that one of file can be read futher
		if scanner_left.Scan() || scanner_right.Scan() {
			return fmt.Errorf("file differs on the final line: %d", line_count)
		}
	}

	return nil
}

func TestApp(t *testing.T) {
	// Read test_cases
	input_files_dir := "../test_cases/input"
	output_files_dir := "../test_cases/output"

	// _, filename, _, _ := runtime.Caller(0)
	// t.Logf("%s", filename)

	input_files, err := os.ReadDir(input_files_dir)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	test_cases := []testCase{}

	for _, file := range input_files {
		if file.IsDir() {
			continue
		}

		in, err := os.Open(input_files_dir + "/" + file.Name())
		if err != nil {
			continue
		}
		out, err := os.Open(output_files_dir + "/" + file.Name())
		if err != nil {
			continue
		}

		test_cases = append(test_cases, testCase{file.Name(), in, out})
	}

	if len(test_cases) == 0 {
		t.Fatalf("Testcases not found")
	}

	for _, tc := range test_cases {
		t.Run(tc.name, func(t *testing.T) {
			real_output := bytes.NewBufferString("")
			app := NewApp(tc.in, real_output)

			err := app.Process() // Ignore errors here
			t.Logf("app process error: %v", err)

			output := real_output.String()
			output_stream := strings.NewReader(output)

			if err := compareReaders(output_stream, tc.out); err != nil {
				t.Errorf("compare error: %v", err)
			}

		})
	}

}
