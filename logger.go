package rapidRoot

import (
	"fmt"
	"io"
	"os"
	"time"
)

var log = newLogger(os.Stdout)

const (
	timeFormat          = "[2006/01/02 15:04:05]"
	terminalInfoFormat  = "%s INFO: %s\n"
	terminalWarnFormat  = "%s WARN: %s\n"
	terminalErrorFormat = "%s ERROR: %s | function caller: %s\n"
)

const (
	reset = "\033[0m"

	black        = "\033[30m"
	red          = "\033[31m"
	green        = "\033[32m"
	yellow       = "\033[33m"
	orange       = "\033[34m"
	magenta      = "\033[35m"
	cyan         = "\033[36m"
	lightGray    = "\033[37m"
	darkGray     = "\033[90m"
	lightRed     = "\033[91m"
	lightGreen   = "\033[92m"
	lightYellow  = "\033[93m"
	lightOrange  = "\033[94m"
	lightMagenta = "\033[95m"
	lightCyan    = "\033[96m"
	white        = "\033[97m"
)

type logger struct {
	output io.Writer
}

func newLogger(output io.Writer) *logger {
	return &logger{output: output}
}

func (l *logger) log(msg string) {
	fmt.Fprint(l.output, msg)
}

func colorize(colorCode string, s string) string {
	return fmt.Sprintf("%s%s%s", colorCode, s, reset)
}

func (l *logger) info(msg string) {
	formattedMsg := fmt.Sprintf(terminalInfoFormat, time.Now().Format(timeFormat), msg)
	if l.output == os.Stdout {
		formattedMsg = colorize(lightCyan, formattedMsg)
	}
	l.log(formattedMsg)
}

func (l *logger) warn(msg string) {
	formattedMsg := fmt.Sprintf(terminalWarnFormat, time.Now().Format(timeFormat), msg)
	if l.output == os.Stdout {
		formattedMsg = colorize(lightOrange, formattedMsg)
	}
	l.log(formattedMsg)
}

func (l *logger) error(msg, handlerName string) {
	formattedMsg := fmt.Sprintf(terminalErrorFormat, time.Now().Format(timeFormat), msg, handlerName)
	if l.output == os.Stdout {
		formattedMsg = colorize(red, formattedMsg)
	}
	l.log(formattedMsg)
}

func (l *logger) fatal(err error) {
	l.log(colorize(lightRed, err.Error()+"\n"))
	os.Exit(1)
}

func (l *logger) logRequest(path, method string, status int) {
	logString := fmt.Sprintf("%s %s %s %d\n", time.Now().Format(timeFormat), method, path, status)
	if l.output == os.Stdout {
		switch {
		case status >= 100 && status < 200:
			logString = colorize(yellow, logString)
		case status >= 200 && status < 300:
			logString = colorize(green, logString)
		case status >= 300 && status < 400:
			logString = colorize(orange, logString)
		case status >= 400 && status < 500:
			logString = colorize(cyan, logString)
		case status >= 500 && status < 600:
			logString = colorize(red, logString)
		default:
			logString = colorize(white, logString)
		}
	}
	l.log(logString)
}

func SetOutput(w io.Writer) {
	log.output = w
}

func logHandlers(s string) {
	colorized := colorize(lightCyan, s)
	log.log(colorized)
}
