package compete

import (
	"log"
	"os"
	"sync"
)

const (
	fileDir = `./temp/`
)

var (
	flushLogsChan = make(chan struct{})
	logsWritten   sync.WaitGroup

	curSubdir string
)

func initLogs(subDir string) error {
	curSubdir = fileDir + subDir
	return os.Mkdir(curSubdir, os.ModePerm)
}

func flushLogs() {
	close(flushLogsChan)
	logsWritten.Wait()
	flushLogsChan = make(chan struct{})
	curSubdir = ``
}

func writeToFile(
	filename string,
	bytes []byte,
) {
	logsWritten.Add(1)
	dir := curSubdir

	go func() {
		<-flushLogsChan
		defer logsWritten.Done()

		f, err := os.OpenFile(dir+filename,
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
