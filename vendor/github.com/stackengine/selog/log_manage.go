//
package selog

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

var ()

type LogWriter struct {
	lck        sync.Mutex
	log_max    int
	line_max   int
	line_cur   int
	out        io.Writer
	log_dir    string
	log_prefix string
}

func (lw *LogWriter) houseKeeping() {
	var fnames []string

	lw.lck.Lock()
	defer lw.lck.Unlock()

	files, _ := ioutil.ReadDir(lw.log_dir)
	for _, f := range files {
		if strings.HasPrefix(f.Name(), lw.log_prefix) {
			fnames = append(fnames, f.Name())
		}
	}
	fmt.Printf("** Names: %#v\n", fnames[:len(fnames)-lw.log_max])

}

func NewLogWriter(lineMax int, logMax int, logdir string, logprefix string) (*LogWriter, error) {
	// ensure log directory
	if err := os.MkdirAll(logdir, 0755); err != nil {
		return nil, err
	}
	lw := &LogWriter{
		line_max:   lineMax,
		log_max:    logMax,
		log_dir:    logdir,
		log_prefix: logprefix}
	lw.houseKeeping()
	return lw, nil
}

func (lw *LogWriter) Write(b []byte) (int, error) {
	lw.lck.Lock()
	defer lw.lck.Unlock()

	if lw.line_cur >= lw.line_max {
		lw.Cycle()
	}

	lw.out.Write(b)
	lw.line_cur += 1
	return 0, nil
}

func (lw *LogWriter) SetMaxLines(max int) error {
	lw.lck.Lock()
	defer lw.lck.Unlock()
	lw.line_max = max
	return nil
}

func (lw *LogWriter) cycle() error {

	return nil
}

func (lw *LogWriter) Cycle() error {
	lw.lck.Lock()
	defer lw.lck.Unlock()
	lw.line_cur = 0
	return lw.cycle()
}
