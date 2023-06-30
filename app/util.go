package main

import "os"

func createFileIfDNE(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err := os.Create(path)
		if err != nil {
			return err
		}
	}

	return nil
}
