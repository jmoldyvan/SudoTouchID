package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnabledLine(t *testing.T) {
	expected := []string{"auth", "sufficient", "pam_tid.so"}
	if got := enabledLine(); !equalSlice(got, expected) {
		t.Errorf("enabledLine() = %v, want %v", got, expected)
	}
}

func TestDisabledLine(t *testing.T) {
	expected := []string{"#auth", "sufficient", "pam_tid.so"}
	if got := disabledLine(); !equalSlice(got, expected) {
		t.Errorf("disabledLine() = %v, want %v", got, expected)
	}
}

func TestEnableTouchID(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Write some initial content to the file
	initialContent := []string{"#auth", "sufficient", "pam_tid.so"}
	err = writeLineToFile(file, initialContent)
	if err != nil {
		t.Fatalf("Failed to write initial content to file: %v", err)
	}

	// Call the function being tested
	enableTouchID(file, false)

	// Verify the updated content in the file
	expectedContent := []string{"auth", "sufficient", "pam_tid.so"}
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	got := strings.Fields(scanner.Text())
	if !equalSlice(got, expectedContent) {
		t.Errorf("enableTouchID() did not update the file correctly, got: %v, want: %v", got, expectedContent)
	}
}

func TestEnableTouchID_NewFile(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Call the function being tested
	enableTouchID(file, true)

	// Verify the content in the new file
	expectedContent := []string{"auth", "sufficient", "pam_tid.so"}
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	got := strings.Fields(scanner.Text())
	if !equalSlice(got, expectedContent) {
		t.Errorf("enableTouchID() did not write to the new file correctly, got: %v, want: %v", got, expectedContent)
	}
}

func TestDisableTouchID(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Write some initial content to the file
	initialContent := []string{"auth", "sufficient", "pam_tid.so"}
	err = writeLineToFile(file, initialContent)
	if err != nil {
		t.Fatalf("Failed to write initial content to file: %v", err)
	}

	// Call the function being tested
	disableTouchID(file, false)

	// Verify the updated content in the file
	expectedContent := []string{"#auth", "sufficient", "pam_tid.so"}
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	got := strings.Fields(scanner.Text())
	if !equalSlice(got, expectedContent) {
		t.Errorf("disableTouchID() did not update the file correctly, got: %v, want: %v", got, expectedContent)
	}
}

func TestDisableTouchID_NewFile(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Call the function being tested
	disableTouchID(file, true)

	// Verify the content in the new file
	expectedContent := []string{"#auth", "sufficient", "pam_tid.so"}
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	got := strings.Fields(scanner.Text())
	if !equalSlice(got, expectedContent) {
		t.Errorf("disableTouchID() did not write to the new file correctly, got: %v, want: %v", got, expectedContent)
	}
}

func TestWriteLineToFile(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	line := []string{"test", "line"}

	err = writeLineToFile(file, line)
	if err != nil {
		t.Fatalf("Failed to write line to file: %v", err)
	}

	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	got := strings.Fields(scanner.Text())
	if !equalSlice(got, line) {
		t.Errorf("writeLineToFile() did not write the line correctly, got: %v, want: %v", got, line)
	}
}

func TestWriteLineToFile_Error(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	file.Close() // Close the file to simulate a write error.

	err = writeLineToFile(file, []string{"test", "line"})
	if err == nil {
		t.Errorf("Expected error when writing to closed file, got none")
	}
}

func TestFindAndReplaceLineInFile(t *testing.T) {
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Write some initial content to the file
	initialContent := []string{"auth", "sufficient", "pam_tid.so"}
	err = writeLineToFile(file, initialContent)
	if err != nil {
		t.Fatalf("Failed to write initial content to file: %v", err)
	}

	// Call the function being tested
	findLine := []string{"auth", "sufficient", "pam_tid.so"}
	replaceLine := []string{"auth", "required", "pam_tid.so"}
	err = findAndReplaceLineInFile(file, findLine, replaceLine)
	if err != nil {
		t.Fatalf("Failed to replace line in file: %v", err)
	}

	// Verify the updated content in the file
	expectedContent := []string{"auth", "required", "pam_tid.so"}
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	got := strings.Fields(scanner.Text())
	if !equalSlice(got, expectedContent) {
		t.Errorf(
			"findAndReplaceLineInFile() did not update the file correctly, got: %v, want: %v",
			got,
			expectedContent,
		)
	}
}

func TestEqualSlice(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		b    []string
		want bool
	}{
		{
			name: "Equal slices",
			a:    []string{"auth", "sufficient", "pam_tid.so"},
			b:    []string{"auth", "sufficient", "pam_tid.so"},
			want: true,
		},
		{
			name: "Different lengths",
			a:    []string{"auth", "sufficient", "pam_tid.so"},
			b:    []string{"auth", "sufficient"},
			want: false,
		},
		{
			name: "Different elements",
			a:    []string{"auth", "sufficient", "pam_tid.so"},
			b:    []string{"auth", "required", "pam_tid.so"},
			want: false,
		},
		{
			name: "Empty slices",
			a:    []string{},
			b:    []string{},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := equalSlice(tt.a, tt.b); got != tt.want {
				t.Errorf("equalSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenOrCreateFile_Error(t *testing.T) {
	_, _, err := openOrCreateFile("/invalid/path/to/file")
	if err == nil {
		t.Errorf("Expected error, got none")
	}
}

func TestCreateFile(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "testfile")

	file, err := createFile(filename)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// Check if the file exists
	if _, err := os.Stat(filename); err != nil {
		t.Errorf("File does not exist: %v", err)
	}
}
