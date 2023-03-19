package orion

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type O struct {
	breaches []Breach
	handler  BreachHandler
}

func New(handler BreachHandler) *O {
	return &O{
		handler: handler,
	}
}

func (o *O) NewBreach(err error) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("whoops")
	}
	f := runtime.FuncForPC(pc)
	b := Breach{
		Err:  err,
		F:    f,
		File: file,
		Line: line,
	}
	o.breaches = append(o.breaches, b)
	if o.handler != nil {
		o.handler(b)
	}
}

type Breach struct {
	Err  error
	F    *runtime.Func
	File string
	Line int
}

type BreachHandler func(breach Breach)

func TestBreachHandler(t *testing.T) BreachHandler {
	return testBreachHandler{t: t}.BreachHandler
}

type testBreachHandler struct {
	t *testing.T
}

func (tbh testBreachHandler) BreachHandler(b Breach) {
	s := tbh.decorate(b, fmt.Sprintf("%+v", b))
	_, _ = fmt.Fprintf(os.Stdout, s)
	tbh.t.Fail()
}

func (tbh testBreachHandler) decorate(b Breach, s string) string {
	file := filepath.Base(b.File)
	// Truncate file name at last file name separator.
	if index := strings.LastIndex(file, "/"); index >= 0 {
		file = file[index+1:]
	} else if index = strings.LastIndex(file, "\\"); index >= 0 {
		file = file[index+1:]
	}

	buf := new(bytes.Buffer)
	// Every line is indented at least one tab.
	buf.WriteByte('\t')
	_, _ = fmt.Fprintf(buf, "%s:%d: ", file, b.Line)

	s = tbh.escapeFormatString(s)

	lines := strings.Split(s, "\n")
	if l := len(lines); l > 1 && lines[l-1] == "" {
		lines = lines[:l-1]
	}
	for i, line := range lines {
		if i > 0 {
			// Second and subsequent lines are indented an extra tab.
			buf.WriteString("\n\t\t")
		}
		// // expand arguments (if $ARGS is present)
		// if strings.Contains(line, "$ARGS") {
		// 	args, _ := loadArguments(path, lineNumber)
		// 	line = strings.Replace(line, "$ARGS", args, -1)
		// }
		buf.WriteString(line)
	}
	buf.WriteString("\n")
	return buf.String()
}

// escapeFormatString escapes strings for use in formatted functions like Sprintf.
func (testBreachHandler) escapeFormatString(fmt string) string {
	return strings.Replace(fmt, "%", "%%", -1)
}

func LogBreachHandler() BreachHandler {
	return func(b Breach) {
		log.Printf("%+v", b)
	}
}
