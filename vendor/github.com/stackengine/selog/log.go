// Package selog is a logging package for stackengine apps
package selog

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
)

// Level is a debugging level from Trace (0) to Fatal (5)
type Level uint

const (
	// Fatal is the Fatal log level
	Fatal = iota
	// Err is the Error log level
	Err = iota
	// Warn is a Warn log level
	Warn = iota
	// Info is an information log level
	Info = iota
	// Debug is a debug log level
	Debug = iota
	// Trace is a trace log level and always prints
	Trace = iota
)

var logLevelNames = []string{
	"Fatal ",
	"Err   ",
	"Warn  ",
	"Info  ",
	"Debug ",
	"Trace ",
}

func (f Level) String() string {
	if f > Debug {
		return "Unknown"
	}
	return logLevelNames[f]
}

// Flag is an additive logging option.
type Flag uint8

const (
	NoAddToLogRing Flag = 1 << iota
	NoPrint
	NoPrefix
	NoTs
	NoFileNo
)

const (
	NoMeta = NoPrefix | NoTs | NoFileNo
)

var flagNames = map[Flag]string{
	NoAddToLogRing: "NoAddToLogRing",
	NoPrint:        "NoPrint",
	NoPrefix:       "NoPrefix",
	NoTs:           "NoTs",
	NoFileNo:       "NoFileNo",
}

func (f Flag) String() string {
	str := flagNames[f]
	if len(str) < 1 {
		return fmt.Sprintf("Unknown SeLogFlag: %d", f)
	}
	return str
}

var (
	logBuffer *Ring
	levels              = make(map[string]LogSetting)
	logs                = make(map[string]*Log)
	output    io.Writer = os.Stderr
	logsLock            = sync.RWMutex{}
)

func SetOutput(w io.Writer) {
	logsLock.Lock()
	output = w
	for _, logger := range logs {
		logger.out = output
	}
	logsLock.Unlock()
}

// Entry is a log entry.
type Entry struct {
	Ts      string `json:"ts"`
	Level   string `json:"level"`
	Prefix  string `json:"prefix"`
	FileNo  string `json:"file_no" codec:"file_no"`
	Message string `json:"message"`
}

func GetSettings() map[string]LogSetting {
	vals := make(map[string]LogSetting)
	for name, logger := range logs {
		vals[name] = logger.level
	}
	return vals
}

// LogSetting is a log level and an enabled indicator.
type LogSetting struct {
	Level   Level
	Enabled bool
}

// Log represents a logger type
type Log struct {
	flags  Flag
	level  LogSetting
	prefix string
	out    io.Writer
	buf    []byte
	buf_mu sync.Mutex
}

// RegisterCliFlags will register "disable-logging" and "debug" flags to the
// given app. When the app function is executed, ProcessCliFlags should be
// called to set the logging options in this package.
func RegisterCliFlags(app *cli.App) {
	loggers := strings.Join(Loggers(), ",")
	flags := []cli.Flag{
		cli.BoolFlag{
			Name:  "disable-logging",
			Usage: "disable logging",
		},
		cli.StringFlag{
			Name:   "debug",
			Usage:  fmt.Sprintf("enable logs at the debug level with a comma seperated list subset of: %v", loggers),
			EnvVar: "DEBUG",
		},
	}
	app.Flags = append(app.Flags, flags...)
}

// ProcessCliFlags parses the global flags added during RegisterCliFlags.
func ProcessCliFlags(ctx *cli.Context) {
	if ctx.GlobalBool("disable-logging") {
		Disable("all")
	} else {
		Enable("all")
	}

	debug := ctx.GlobalString("debug")
	for _, name := range strings.Split(debug, ",") {
		if len(name) <= 0 {
			continue
		}
		Enable(name)
		SetLevel(name, Debug)
	}
}

// SetLevel will the log level on the named logged. If "all" is passed, it will
// the given log level on all the log level. SetLevel returns true if the level
// was set successfully.
func SetLevel(name string, level Level) bool {
	logsLock.RLock()
	defer logsLock.RUnlock()
	if name == "all" || name == "*" {
		for _, logger := range logs {
			logger.setLevel(level)
		}
		return true
	}
	l := logs[name]
	if l == nil {
		lvl := levels[name]
		lvl.Level = level
		levels[name] = lvl
		return false
	}
	l.setLevel(level)
	return true
}

// SetLevels will set the log level for the named loggers. If "all" is passed,
// it will set the given log level on all the log levels.
func SetLevels(level Level, names ...string) {
	for _, name := range names {
		nm := strings.TrimSpace(name)
		if nm == "all" || nm == "*" {
			SetLevel(nm, level)
			return
		}
	}

	for _, name := range names {
		nm := strings.TrimSpace(name)
		if nm == "" {
			continue
		}

		SetLevel(nm, level)
	}
}

// Disable the logger under the name. "all" applies
// to all registered loggers.
func Disable(name string) bool {
	logsLock.RLock()
	defer logsLock.RUnlock()
	if name == "all" || name == "*" {
		for _, l := range logs {
			l.Disable()
		}
		return true
	}
	l := logs[name]
	if l == nil {
		lvl := levels[name]
		lvl.Enabled = false
		levels[name] = lvl
		return false
	}
	l.Disable()
	return true
}

// Enable turns on the logger under the name. "all" applies
// to all registered loggers.
func Enable(name string) bool {
	logsLock.RLock()
	defer logsLock.RUnlock()
	if name == "all" || name == "*" {
		for _, l := range logs {
			l.Enable()
		}
		return true
	}
	l := logs[name]
	if l == nil {
		lvl := levels[name]
		lvl.Enabled = false
		levels[name] = lvl
		return false
	}
	l.Enable()
	return true
}

// Register creates a new named logger and registers it. If the named logger already
// exists, Register() panics. "all" is a reserved name.
func Register(name string, flags Flag) *Log {
	logsLock.Lock()
	defer logsLock.Unlock()
	if logger, found := logs[name]; found {
		return logger
	}
	if name == "all" || name == "*" {
		panic("logger 'all' and '*' are reserved logger names")
	}
	log := &Log{}
	log.flags = flags
	log.out = output

	log.prefix = name
	logs[name] = log
	if lvl, ok := levels["all"]; ok {
		log.level = lvl
		log.Enable()
		return log
	}
	if lvl, ok := levels[name]; ok {
		log.level = lvl
	} else {
		log.Enable()
	}
	return log
}

// Loggers returns a list of registered logger names
func Loggers() []string {
	logsLock.RLock()
	defer logsLock.RUnlock()
	keys := make([]string, len(logs))
	i := 0
	for key := range logs {
		keys[i] = key
		i++
	}
	return keys
}

func (l *Log) getLevel() Level {
	return l.level.Level
}

func (l *Log) setLevel(v Level) {
	if v > Trace {
		v = Trace
	}
	l.level.Level = v
	l.level.Enabled = true
}

// ServeHTTP implements the negroni middleware interface. This allows
// an selog.Log instance to act as a middleware logger.
func (l *Log) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	l.Printf("Started %s %s", r.Method, r.URL.Path)

	next(rw, r)

	res := rw.(negroni.ResponseWriter)
	l.Printf("Completed %v %s in %v", res.Status(), http.StatusText(res.Status()), time.Since(start))
}

// Enable enables the logger.
func (l *Log) Enable() {
	l.level.Enabled = true
}

// Disable disables the logger.
func (l *Log) Disable() {
	l.level.Enabled = false
}

func (l *Log) writeLog(s string, level Level) error {
	now := time.Now()
	var (
		file string
		line int
		err  error
	)

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	year, month, day := now.Date()
	hour, min, sec := now.Clock()

	ent := &Entry{
		Prefix: l.prefix,
		Ts: fmt.Sprintf("%4d/%02d/%02d %02d:%02d:%02d",
			year, month, day, hour, min, sec),
		Level: logLevelNames[level],
	}

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	ent.FileNo = fmt.Sprintf("%s:%d", file, line)

	if len(s) > 0 && s[len(s)-1] != '\n' {
		ent.Message = fmt.Sprintf("%s\n", s)
	} else {
		ent.Message = s
	}
	sPrefix := fmt.Sprintf("[ %-18.18s ] [ %-6.6s ]",
		ent.Prefix, ent.Level)

	// TODO: employ a buffer recycling pattern to avoid
	// synchronizing on this buffer.
	l.buf_mu.Lock()
	defer l.buf_mu.Unlock()
	l.buf = l.buf[:0]
	if l.flags&NoPrefix == 0 {
		l.buf = append(l.buf, sPrefix...)
		l.buf = append(l.buf, ' ')
	}
	if l.flags&NoTs == 0 {
		l.buf = append(l.buf, ent.Ts...)
		l.buf = append(l.buf, ' ')
	}
	if l.flags&NoFileNo == 0 {
		l.buf = append(l.buf, ent.FileNo...)
		l.buf = append(l.buf, ' ')
	}
	l.buf = append(l.buf, ent.Message...)
	err = nil

	// this is good for testing and not cluttering up
	// stderr
	if l.flags&NoPrint == 0 {
		_, err = l.out.Write(l.buf)
	}
	if logBuffer != nil && l.flags&NoAddToLogRing == 0 {
		logBuffer.Add(ent)
	}
	return err
}

// Printf writes the log message at debug level or greater.
func (l *Log) Printf(format string, v ...interface{}) {
	if l.level.Enabled && l.level.Level >= Debug {
		l.writeLog(fmt.Sprintf(format, v...), Debug)
	}
}

// InfoPrintf writes the log message at info level or greater.
func (l *Log) InfoPrintf(format string, v ...interface{}) {
	if l.level.Enabled && l.level.Level >= Info {
		l.writeLog(fmt.Sprintf(format, v...), Info)
	}
}

// ErrPrintf writes the log message at any level.
func (l *Log) ErrPrintf(format string, v ...interface{}) {
	l.writeLog(fmt.Sprintf(format, v...), Err)
}

// WarnPrintf writes the log message at warn level or greater.
func (l *Log) WarnPrintf(format string, v ...interface{}) {
	if l.level.Enabled && l.level.Level >= Warn {
		l.writeLog(fmt.Sprintf(format, v...), Warn)
	}
}

// WarnPrintln writes the log message at warn level or greater.
func (l *Log) WarnPrintln(v ...interface{}) {
	if l.level.Enabled && l.level.Level >= Warn {
		l.writeLog(fmt.Sprintln(v...), Warn)
	}
}

// WarnPrint writes the log message at warn level or greater.
func (l *Log) WarnPrint(v ...interface{}) {
	if l.level.Enabled && l.level.Level >= Warn {
		l.writeLog(fmt.Sprint(v...), Warn)
	}
}

// TracePrintf writes the log message at Trace level or greater.
func (l *Log) TracePrintf(format string, v ...interface{}) {
	if l.level.Enabled && l.level.Level >= Trace {
		l.writeLog(fmt.Sprintf(format, v...), Trace)
	}
}

// TracePrintln writes the log message at Trace level or greater.
func (l *Log) TracePrintln(v ...interface{}) {
	if l.level.Enabled && l.level.Level >= Trace {
		l.writeLog(fmt.Sprintln(v...), Trace)
	}
}

// TracePrint writes the log message at Trace level or greater.
func (l *Log) TracePrint(v ...interface{}) {
	if l.level.Enabled && l.level.Level >= Trace {
		l.writeLog(fmt.Sprint(v...), Trace)
	}
}

// Print writes the log message at debug level or greater.
func (l *Log) Print(v ...interface{}) {
	if l.level.Enabled && l.level.Level >= Debug {
		l.writeLog(fmt.Sprint(v...), Debug)
	}
}

// InfoPrint writes the log message at info level or greater.
func (l *Log) InfoPrint(v ...interface{}) {
	if l.level.Enabled && l.level.Level >= Info {
		l.writeLog(fmt.Sprint(v...), Info)
	}
}

// ErrPrint writes the log message at any level.
func (l *Log) ErrPrint(v ...interface{}) {
	l.writeLog(fmt.Sprint(v...), Err)
}

// Println writes the log message at debug level or greater.
func (l *Log) Println(v ...interface{}) {
	if l.level.Enabled && l.level.Level >= Debug {
		l.writeLog(fmt.Sprintln(v...), Debug)
	}
}

// InfoPrintln writes the log message at info level or greater.
func (l *Log) InfoPrintln(v ...interface{}) {
	if l.level.Enabled && l.level.Level >= Info {
		l.writeLog(fmt.Sprintln(v...), Info)
	}
}

// ErrPrintln writes the log message at any level.
func (l *Log) ErrPrintln(v ...interface{}) {
	l.writeLog(fmt.Sprintln(v...), Err)
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func (l *Log) Fatal(v ...interface{}) {
	l.writeLog(fmt.Sprint(v...), Fatal)
	os.Exit(42)
}

// Fatalf is equivalent to Print() followed by a call to os.Exit(1).
func (l *Log) Fatalf(formatStr string, v ...interface{}) {
	l.writeLog(fmt.Sprintf(formatStr, v...), Fatal)
	os.Exit(42)
}

// SetGlobalRingBuffer sets the ring buffer for the selog package.
func SetGlobalRingBuffer(r *Ring) {
	logBuffer = r
}

// ring buffer for logging
type Ring struct {
	quitCh   chan chan error
	logInCh  chan *Entry
	logOutCh chan *Entry

	rdrs map[*Reader]*Reader
	lock *sync.RWMutex
}

// ConfigLogRing initializes logring buffering.
func ConfigLogRing(inSize, outSize int) *Ring {
	newRing := initLogRing(inSize, outSize)
	if newRing != nil {
		go newRing.run()
	}
	return newRing
}

func initLogRing(inSize, outSize int) *Ring {
	LogBuf := &Ring{
		logInCh:  make(chan *Entry, inSize),
		logOutCh: make(chan *Entry, outSize),
		quitCh:   make(chan chan error),
		rdrs:     make(map[*Reader]*Reader),
		lock:     &sync.RWMutex{},
	}
	return LogBuf
}

func (r *Ring) rdrsLog(v *Entry) {
	if r.rdrs != nil {
		r.lock.Lock()
		for _, rdr := range r.rdrs {
			if rdr.ring != nil {
				rdr.ring.processLog(v)
			}
		}
		r.lock.Unlock()
	}
}

func (r *Ring) rdrsQuit() {
	if r.rdrs != nil {
		r.lock.Lock()
		for _, rdr := range r.rdrs {
			rdr.Close()
		}
		r.lock.Unlock()
	}
}

func (r *Ring) run() {
	for {
		select {
		case v := <-r.logInCh:
			r.processLog(v)
			r.rdrsLog(v)

		case errc := <-r.quitCh:
			r.rdrsQuit()
			close(r.logInCh)
			r.logInCh = nil
			close(r.logOutCh)
			r.logOutCh = nil
			errc <- nil
			return
		}
	}
}

func (r *Ring) processLog(v *Entry) {
	r.lock.Lock()
	select {
	case r.logOutCh <- v:
	default:
		<-r.logOutCh
		r.logOutCh <- v
	}
	r.lock.Unlock()
}

// Add a log entry to the ring.
func (r *Ring) Add(v *Entry) {
	if r.logInCh != nil {
		r.logInCh <- v
	}
}

// Close the ring buffer
func (r *Ring) Close() error {
	errc := make(chan error)
	r.quitCh <- errc
	return <-errc
}

// Reader is an "instance" of a log buffer
type Reader struct {
	ring *Ring
	mu   sync.RWMutex
}

func NewReader(ioSize int) *Reader {
	return logBuffer.NewReader(ioSize)
}

// NewReader creates a instance of a Reader of ioSize.
func (r *Ring) NewReader(ioSize int) *Reader {
	rdr := &Reader{
		ring: initLogRing(ioSize, ioSize),
	}

	if rdr.ring == nil {
		panic("rdr.ring is nil")
	}

	// add to array we need to service.
	r.lock.Lock()
	r.rdrs[rdr] = rdr
	r.lock.Unlock()

	// start the reader ring buffer
	go rdr.serv()
	return rdr
}

func (rdr *Reader) getRing() *Ring {
	rdr.mu.RLock()
	x := rdr.ring
	rdr.mu.RUnlock()
	return x
}

func (rdr *Reader) serv() {
	r := rdr.getRing()
	if r == nil {
		return
	}

	for {
		select {
		case v := <-r.logInCh:
			r.processLog(v)

		case <-r.quitCh:
			close(r.logOutCh)
			return
		}
	}
}

// GetRdrChan returns a channel of log *Entry.
func (rdr *Reader) GetRdrChan() chan *Entry {
	if r := rdr.getRing(); r != nil {
		return r.logOutCh
	}
	return nil
}

// Close closes the reader.
func (rdr *Reader) Close() error {
	var err error
	rdr.mu.Lock()
	defer rdr.mu.Unlock()
	if rdr.ring != nil {
		// Remove reader from the ring
		rdr.ring.lock.Lock()
		delete(rdr.ring.rdrs, rdr)
		rdr.ring.lock.Unlock()
		rdr.ring = nil
	}
	return err
}
