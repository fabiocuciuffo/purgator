package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var DIRECTORIES_TO_DELETE = []string{
	"node_modules",
	"vendor",
}

var NO_DELETE bool = true

var TOTAL_SIZE float64 = 0

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fsys := os.DirFS(dir)
	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if !d.IsDir() {
			return nil
		}

		if !contains(DIRECTORIES_TO_DELETE, d.Name()) {
			return nil
		}

		TOTAL_SIZE, err = calculateTotalSize(path)
		if err != nil {
			log.Fatal(err)
		}

		err = os.RemoveAll(path)
		if err != nil {
			log.Fatal(err)
		}

		err = os.RemoveAll(path)

		if err != nil {
			log.Fatal(err)
		}

		NO_DELETE = false

		fmt.Printf("\033[36m%s a été supprimé\033[0m\n", path)

		return fs.SkipDir
	})

	if NO_DELETE {
		fmt.Println("\033[33mAucun répertoire supprimé.\033[0m")
	} else {
		fmt.Printf("\033[32mTaille totale: %.2f Mo\033[0m\n\n", TOTAL_SIZE/(1024*1024))
	}
}

func contains[T comparable](slice []T, val T) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func calculateTotalSize(path string) (float64, error) {
	var totalSize float64

	err := filepath.Walk(path, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		totalSize += float64(info.Size())
		return nil
	})

	if err != nil {
		return 0, err
	}
	return totalSize, nil
}
