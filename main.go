package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

const (
	torrentsDir   = "torrents"
	logFile       = "failed_downloads.log"
	baseURLFormat = "https://archive.org/download/%s/%s_archive.torrent"
)

var (
	csvFile    string
	numWorkers int
)

func init() {
	flag.StringVar(&csvFile, "csv", "search.csv", "path to the CSV file")
	flag.IntVar(&numWorkers, "workers", 10, "number of concurrent workers")
	flag.Parse()
}

func main() {
	if err := os.MkdirAll(torrentsDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create torrents directory: %v\n", err)
	}

	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v\n", err)
	}
	defer file.Close()

	failedLogFile, err := os.Create(logFile)
	if err != nil {
		log.Fatalf("Failed to create log file: %v\n", err)
	}
	defer failedLogFile.Close()
	failedLogWriter := bufio.NewWriter(failedLogFile)
	defer failedLogWriter.Flush()

	logger := log.New(os.Stdout, "", log.LstdFlags)
	failedLogger := log.New(failedLogWriter, "", log.LstdFlags)

	reader := csv.NewReader(file)

	headers, err := reader.Read()
	if err != nil {
		log.Fatalf("Failed to read CSV header: %v\n", err)
	}

	var identifierIndex int
	for i, header := range headers {
		if header == "identifier" {
			identifierIndex = i
			break
		}
	}

	identifiers := make(chan string)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(identifiers, logger, failedLogger, &wg)
	}

	go func() {
		for {
			record, err := reader.Read()
			if err == io.EOF {
				close(identifiers)
				break
			}
			if err != nil {
				logger.Printf("Failed to read CSV record: %v\n", err)
				continue
			}
			if len(record) == 0 {
				continue
			}

			identifiers <- record[identifierIndex]
		}
	}()

	wg.Wait()
}

func worker(identifiers chan string, logger, failedLogger *log.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	for identifier := range identifiers {
		url := fmt.Sprintf(baseURLFormat, identifier, identifier)
		filename := filepath.Join(torrentsDir, identifier+"_archive.torrent")

		if fileExists(filename) {
			logger.Printf("File %s already exists, skipping download.\n", filename)
			continue
		}

		logger.Printf("Downloading %s\n", filename)
		if err := downloadFile(url, filename); err != nil {
			failedLogger.Printf("Failed to download %s: %v\n", url, err)
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func downloadFile(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.Remove(filename)
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}
