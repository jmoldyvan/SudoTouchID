package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const sudoLocalFilePath = "/etc/pam.d/sudo_local"

func enabledLine() []string {
	return []string{"auth", "sufficient", "pam_tid.so"}
}

func disabledLine() []string {
	return []string{"#auth", "sufficient", "pam_tid.so"}
}

func main() {

	enable := flag.Bool("enable", false, "enable touch id for sudo")
	disable := flag.Bool("disable", false, "disable touch id for sudo")
	flag.Parse()

	if *enable && *disable || !*enable && !*disable {
		log.Fatalf("must provide either -enable or -disable")
	}

	file, newFile, err := openOrCreateFile(sudoLocalFilePath)
	if err != nil {
		log.Fatalf("failed to open or create file %s", sudoLocalFilePath)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("failed to close file: %v", err)
		}
	}()

	if *enable {
		enableTouchID(file, newFile)
	}

	if *disable {
		disableTouchID(file, newFile)
	}
}

func enableTouchID(file *os.File, newFile bool) {
	fmt.Println("enable touch id for sudo")

	// Write enabled line if it's a new file
	if newFile {
		err := writeLineToFile(file, enabledLine())
		if err != nil {
			log.Fatalf("Failed to write enabled line to file: %v", err)
		}
		return
	}

	// Replace disabled line with enabled line in an existing file
	err := findAndReplaceLineInFile(file, disabledLine(), enabledLine())
	if err != nil {
		log.Fatalf("Failed to write enabled line to file: %v", err)
	}
}

// Same as enabled, just swapping the parameters
func disableTouchID(file *os.File, newFile bool) {
	fmt.Println("disable touch id for sudo")
	if newFile {
		err := writeLineToFile(file, disabledLine())
		if err != nil {
			log.Fatalf("failed to write disabled line to file: %v", err)
		}
		return
	}

	err := findAndReplaceLineInFile(file, enabledLine(), disabledLine())
	if err != nil {
		log.Fatalf("failed to write disabled line to file: %v", err)
	}

}

func writeLineToFile(file *os.File, line []string) error {
	// Truncate and seek to the beginning
	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	_, err := writer.WriteString(strings.Join(line, " ") + "\n")
	if err != nil {
		return err
	}
	return writer.Flush()
}

func findAndReplaceLineInFile(file *os.File, findLine []string, replaceLine []string) error {
	// Read the existing lines
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Reset file for writing
	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	writer := bufio.NewWriter(file)

	// Flag to check if replacement was done
	replaced := false
	for _, line := range lines {
		currentLine := strings.Fields(line)
		if equalSlice(currentLine, findLine) {
			_, err := writer.WriteString(strings.Join(replaceLine, " ") + "\n")
			if err != nil {
				return err
			}
			replaced = true
		} else {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				return err
			}
		}
	}

	// Append the line if it was not replaced
	if !replaced {
		_, err := writer.WriteString(strings.Join(replaceLine, " ") + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func equalSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func openOrCreateFile(filename string) (*os.File, bool, error) {
	fmt.Printf("%s \n", filename)
	fileExists := fileExists(filename)
	fmt.Printf("%t", fileExists)
	var file *os.File
	var newfile bool
	var err error
	if fileExists {
		file, err = openFile(filename)
		newfile = false
	} else {
		file, err = createFile(filename)
		newfile = true
	}
	return file, newfile, err
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

// nolint:gosec
func openFile(filename string) (*os.File, error) {
	// #nosec G304: filename is trusted
	return os.OpenFile(filename, os.O_RDWR, 0o644)
}

// nolint:gosec
func createFile(filename string) (*os.File, error) {
	return os.Create(filename)
}
