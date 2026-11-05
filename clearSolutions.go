package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	err := filepath.Walk(".", stripSolutions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Done! Stripped solution sections from all Go files.")
}

func stripSolutions(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// Skip directories and non-.go files
	if info.IsDir() ||
		!strings.HasSuffix(path, ".go") ||
		strings.HasSuffix(path, "clearSolutions.go") {
		return nil
	}

	// Read the file
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", path, err)
	}

	// Process the content
	lines := strings.Split(string(content), "\n")
	var result []string
	i := 0

	for i < len(lines) {
		// Look for the start marker
		if strings.Contains(lines[i], "// !->") {
			// Add the marker line
			result = append(result, lines[i])

			// add a blank line after the start marker
			indent := lines[i][:strings.Index(lines[i], "// !->")]
			result = append(result, indent)

			// Skip lines until we find the end marker
			i++
			for i < len(lines) && !strings.Contains(lines[i], "// <-!") {
				i++
			}
			// Add the end marker line
			if i < len(lines) {
				result = append(result, lines[i])
			}
			i++
		} else {
			result = append(result, lines[i])
			i++
		}
	}

	// Write back to file
	output := strings.Join(result, "\n")
	err = os.WriteFile(path, []byte(output), 0644)
	if err != nil {
		return fmt.Errorf("error writing %s: %w", path, err)
	}

	fmt.Printf("Processed: %s\n", path)
	return nil
}
