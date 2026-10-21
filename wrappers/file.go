package wrappers

import (
	"bufio"
	"os"
)

func ReadFile(path string, visit func(string) error) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := visit(scanner.Text()); err != nil {
			return err
		}
	}

	return scanner.Err()
}
