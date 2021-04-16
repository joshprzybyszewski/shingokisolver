package compete

import (
	"log"
	"os"
	"sync"
)

var (
	flushLogsChan = make(chan struct{}, 0)
	logsWritten   sync.WaitGroup
)

func initLogs() {
	// TODO rm temp/*
}

func flushLogs() {
	close(flushLogsChan)
	logsWritten.Wait()
	flushLogsChan = make(chan struct{}, 0)
}

func writeToFile(
	filename string,
	bytes []byte,
) {
	logsWritten.Add(1)

	go func() {
		<-flushLogsChan
		defer logsWritten.Done()

		f, err := os.OpenFile(filename,
			os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("error opening file: %v\n", err)
			return
		}
		defer f.Close()

		if _, err := f.Write(bytes); err != nil {
			log.Printf("error WriteString file: %v\n", err)
			return
		}
	}()
}
