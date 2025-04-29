package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		return
	}
	fmt.Println("Current Working Directory is:", path)
	fmt.Println("Checking for all files in the directory")

	files := listFiles(path)

	dirMap, err := loadMappings("mappings.json")
	if err != nil {
		log.Fatalf("Error loading mappings: %v\n", err)
		return
	}

	// Iterate over all files
	for _, file := range files {
		if file.Name() == "file_organiser.go" || file.Name() == "mappings.json" {
			continue
		}

		if !file.IsDir() {
			ext := strings.ReplaceAll(strings.ToLower(filepath.Ext(file.Name())), " ", "")
			dirName, exists := dirMap[ext]
			var dirPath string
			if exists {
				dirPath = filepath.Join(path, dirName)
			} else {
				dirPath = filepath.Join(path, "others")
			}

			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
					log.Printf("Error creating directory %s: %v\n", dirPath, err)
					continue
				}
			}

			newPath := filepath.Join(dirPath, file.Name())
			MoveFile(path, newPath, file)
		}
	}
}

func loadMappings(filename string) (map[string]string, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read mappings file: %v", err)
	}

	var dirMap map[string]string
	if err := json.Unmarshal(file, &dirMap); err != nil {
		return nil, fmt.Errorf("failed to parse mappings file: %v", err)
	}

	return dirMap, nil
}

func listFiles(path string) []os.DirEntry {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("Error reading directory %s: %v\n", path, err)
	}
	return files
}

func MoveFile(source, destination string, file os.DirEntry) {
	originalPath := filepath.Join(source, file.Name())
	newPath := destination
	fmt.Printf("Moving file %s to %s\n", originalPath, newPath)

	if err := os.Rename(originalPath, newPath); err != nil {
		log.Printf("Error moving file %s to %s: %v\n", originalPath, newPath, err)
	}
}
