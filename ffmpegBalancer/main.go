// takes a filename (and count of flines in it, cause i'm lazy)
// converts each line in it with ffmpeg, balanced onto the number of cores in the system
package main

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/cheggaaa/pb"
	"github.com/cryptix/go/logging"
)

var log = logging.Logger("convert")

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: convert <files.lst> <count>")
	}

	inputf, err := os.Open(os.Args[1])
	logging.CheckFatal(err)
	defer inputf.Close()

	var wg sync.WaitGroup
	jobs := make(chan string)
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(jobs, &wg)
	}

	wmaScanner := bufio.NewScanner(inputf)

	total, err := strconv.Atoi(os.Args[2])
	logging.CheckFatal(err)

	bar := pb.StartNew(total)
	for wmaScanner.Scan() {
		fname := wmaScanner.Text()
		if !strings.Contains(fname, "/wav/") {
			jobs <- fname
			bar.Increment()
		}
	}

	logging.CheckFatal(wmaScanner.Err())
	close(jobs)

	wg.Wait()
	bar.Finish()
}

func worker(jobs <-chan string, wg *sync.WaitGroup) {
	for job := range jobs {

		dir, file := filepath.Split(job)
		os.Mkdir(filepath.Join(dir, "wav"), 0700)
		// log.Info("converting", file)
		cmd := exec.Command("ffmpeg", "-y", "-i", file, "-acodec", "aac", "-strict", "-2", "-b:a", "192k", strings.Replace(file, ".wav", ".m4a", 1))
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Error(string(out))
			log.Fatal(err)
		}

		logging.CheckFatal(os.Rename(job, filepath.Join(dir, "wav", file)))
	}
	wg.Done()
}
