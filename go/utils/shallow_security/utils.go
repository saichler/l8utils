package main

import (
	"bytes"
	"os"
)

func SeekResource(path string, filename string) string {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return ""
	}
	if fileInfo.Name() == filename {
		return path
	}
	if fileInfo.IsDir() {
		files, err := os.ReadDir(path)
		if err != nil {
			return ""
		}
		for _, file := range files {
			found := SeekResource(pathOf(path, file), filename)
			if found != "" {
				return found
			}
		}
	}
	return ""
}

func pathOf(path string, file os.DirEntry) string {
	buff := bytes.Buffer{}
	buff.WriteString(path)
	buff.WriteString("/")
	buff.WriteString(file.Name())
	return buff.String()
}
