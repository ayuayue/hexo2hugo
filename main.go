package main

import (
	"fmt"
	"hexo2hugo/file"
	"os"
)

func main() {
	// Check if a directory path is provided as a command-line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory-path>")
		os.Exit(1)
	}
	dirPath := os.Args[1]
	filePathNames := file.GetAllMDFileName(dirPath)
	file.ReadAllMDFile(filePathNames)
}
