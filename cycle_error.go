package antidig

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/3timeslazy/anti-dig/internal/digreflect"
)

type cycleErrPathEntry struct {
	Key  key
	Func *digreflect.Func
}

type errCycleDetected struct {
	Path  []cycleErrPathEntry
	scope *Scope
}

var _ digError = errCycleDetected{}

func (e errCycleDetected) Error() string {
	// We get something like,
	//
	//   [scope "foo"]
	//   func(*bar) *foo provided by "path/to/package".NewFoo (path/to/file.go:42)
	//   	depends on func(*baz) *bar provided by "another/package".NewBar (somefile.go:1)
	//   	depends on func(*foo) baz provided by "somepackage".NewBar (anotherfile.go:2)
	//   	depends on func(*bar) *foo provided by "path/to/package".NewFoo (path/to/file.go:42)
	//
	b := new(bytes.Buffer)

	if name := e.scope.name; len(name) > 0 {
		fmt.Fprintf(b, "[scope %q]\n", name)
	}
	for i, entry := range e.Path {
		if i > 0 {
			b.WriteString("\n\tdepends on ")
		}
		fmt.Fprintf(b, "%v provided by %v", entry.Key, entry.Func)
	}
	return b.String()
}

func (e errCycleDetected) writeMessage(w io.Writer, v string) {
	fmt.Fprint(w, e.Error())
}

func (e errCycleDetected) Format(w fmt.State, c rune) {
	formatError(e, w, c)
}

// IsCycleDetected returns a boolean as to whether the provided error indicates
// a cycle was detected in the container graph.
func IsCycleDetected(err error) bool {
	return errors.As(err, &errCycleDetected{})
}
