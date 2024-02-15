package dataparse

import (
	"io"
	"net/http"
	"os"
)

func DownloadPhoto(url, name string) error {

	fileExists := func(filename string) bool {
		_, err := os.Stat(filename)
		return !os.IsNotExist(err)
	}

	filename := "img/" + name + ".png"

	if !fileExists(filename) {

		response, err := http.Get(url)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			return err
		}
	}

	return nil
}
