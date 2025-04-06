package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sync"
	"time"
)

var TIME_START = time.Now().UnixMilli()

var DIRECTORIES_TO_DELETE = []string{
	"node_modules",
	"vendor",
}

var NO_DELETE bool = true

var TOTAL_SIZE float64 = 0

var TIME_ELPASED float64 = 0

var m = sync.Mutex{}

var wg = sync.WaitGroup{}

func main() {
	printMem := flag.Bool("print-mem", false, "bool print memory")
	flag.Parse()
	var dir string
	args := os.Args[1:]
	var tmpArgs []string
	for i, v := range args {
		if string(v[0]) == "-" {
			tmpArgs = append(args[:i], args[i+1:]...)
		}
	}
	args = tmpArgs
	if len(args) > 1 {
		log.Fatal("One argument expected only.")
		os.Exit(1)
	}
	if len(args) == 1 {
		dir = args[0]
	} else {
		dir, _ = os.Getwd()
	}
	exec(dir)

	wg.Wait()

	printExecLogs()

	if *printMem {
		printMemUsage()
	}
}

func exec(basePath string) {
	fsys := os.DirFS(basePath)

	var pathsToProcess []string

	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		if !d.IsDir() {
			return nil
		}

		if !slices.Contains(DIRECTORIES_TO_DELETE, d.Name()) {
			return nil
		}

		fullPath := path
		if string(basePath[len(basePath)-1]) == "/" {
			fullPath = basePath + path
		} else {
			fullPath = basePath + "/" + path
		}

		pathsToProcess = append(pathsToProcess, fullPath)

		return fs.SkipDir
	})

	for _, path := range pathsToProcess {
		wg.Add(1)
		go processDirectory(path)
	}
}

func processDirectory(path string) {
	defer wg.Done()

	currentSize, err := calculateTotalSize(path)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	m.Lock()
	TOTAL_SIZE += currentSize
	NO_DELETE = false
	m.Unlock()

	err = os.RemoveAll(path)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("\033[36m%s a été supprimé\033[0m\n", path)
}

func printExecLogs() {
	if NO_DELETE {
		fmt.Println("\033[33mAucun répertoire supprimé.\033[0m")
	} else {
		TIME_ELPASED = float64(time.Now().UnixMilli()-TIME_START) / 1000
		fmt.Printf("\033[32mEspace libéré %.2fMo en %.3fs\033[0m\n", TOTAL_SIZE/(1024*1024), TIME_ELPASED)
	}
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

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB\n", bToMb(m.Alloc))
	fmt.Printf("TotalAlloc = %v MiB\n", bToMb(m.TotalAlloc))
	fmt.Printf("Sys = %v MiB\n", bToMb(m.Sys))
	fmt.Printf("NumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
