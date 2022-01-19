package assert

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
)

var (
	plumber = fileMetaPlumber{
		plumbing: map[string]*fileMeta{},
	}
)

type fileMetaPlumber struct {
	plumbing map[string]*fileMeta
}

type fileMeta struct {
	src  string
	f    *ast.File
	fset *token.FileSet
}

func (p *fileMetaPlumber) MGetLocateAndArgs(dep int, locations ...int) []string {
	_, file, line, _ := runtime.Caller(dep)
	pc, _, _, _ := runtime.Caller(dep - 1)
	callFuncName := runtime.FuncForPC(pc).Name()
	callFuncName = strings.Split(callFuncName, "[")[0]
	fm, ok := p.plumbing[file]
	if !ok {
		fset := token.NewFileSet()
		content, err := ioutil.ReadFile(file)
		if err != nil {
			panic(fmt.Errorf("read file %v fail, %w", file, err))
		}
		f, err := parser.ParseFile(fset, file, string(content), 0)
		if err != nil {
			panic(fmt.Errorf("parse file %v fail, %w", file, err))
		}
		fm = &fileMeta{
			src:  string(content),
			f:    f,
			fset: fset,
		}
		p.plumbing[file] = fm
	}

	res := []string{file + ":" + strconv.Itoa(line) + ":"}
	ast.Inspect(fm.f, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.CallExpr:
			if line == fm.fset.Position(n.Pos()).Line {
				ce := n.(*ast.CallExpr)
				if len(ce.Args) == 0 {
					return false
				}
				for _, v := range locations {
					if v > len(ce.Args) {
						panic(fmt.Errorf("Func args position not exists, %v", v))
					}
					res = append(res, fm.nodeToString_(ce.Args[v]))
				}
				return false
			}
			return true
		default:
			return true
		}
	})
	return res
}

func (fm *fileMeta) nodeToString_(n ast.Node) string {
	start := fm.fset.Position(n.Pos()).Offset
	stop := fm.fset.Position(n.End()).Offset
	return fm.src[start:stop]
}

func valueToString(v any) string {
	switch x := v.(type) {
	case string:
		return strconv.Quote(x)
	default:
		return fmt.Sprintf("%#v", v)
	}
}

func isEqual(expected, actual any) bool {
	if expected == nil && actual == nil {
		return true
	}

	return reflect.DeepEqual(expected, actual)
}

func Equal[T any](t *testing.T, expected, actual T) {
	if isEqual(expected, actual) {
		return
	}

	// testing.T 报错的时候，会忽略这个文件的函数调用，这样可以在真正调用 assert.Equal 的地方报错
	t.Helper()

	fargs := plumber.MGetLocateAndArgs(2, 1, 2)

	t.Errorf(`
		Expected: %v 

		is equal to %v, 

		but got: %v\n`,
		fargs[2], fargs[1], actual)
}

func True[B ~bool](t *testing.T, actual B) {
	if actual {
		return
	}

	// testing.T 报错的时候，会忽略这个文件的函数调用，这样可以在真正调用 assert.Equal 的地方报错
	t.Helper()

	fargs := plumber.MGetLocateAndArgs(2, 1)

	t.Errorf(`Expected: %v is true, but got false`, fargs[1])
}

func False[B ~bool](t *testing.T, actual B) {
	if !actual {
		return
	}

	// testing.T 报错的时候，会忽略这个文件的函数调用，这样可以在真正调用 assert.Equal 的地方报错
	t.Helper()

	fargs := plumber.MGetLocateAndArgs(2, 1)

	t.Errorf(`Expected: %v is false, but got true`, fargs[1])
}

func Panic(t *testing.T, f func()) {
	if ok, _, _ := didPanic(f); !ok {
		// testing.T 报错的时候，会忽略这个文件的函数调用，这样可以在真正调用 assert.Equal 的地方报错
		t.Helper()

		fargs := plumber.MGetLocateAndArgs(2, 1)
		t.Errorf("Func %v should panic", fargs[1])
	}
}

func NotPanic(t *testing.T, f func()) {
	if ok, err, stack := didPanic(f); ok {
		// testing.T 报错的时候，会忽略这个文件的函数调用，这样可以在真正调用 assert.Equal 的地方报错
		t.Helper()

		fargs := plumber.MGetLocateAndArgs(2, 1)
		t.Errorf("Func %v should not panic, cause: %v", fargs[1], err)
		t.Errorf("Stack: %v", stack)
	}
}

func didPanic(f func()) (p bool, err any, stack string) {
	defer func() {
		if err = recover(); err != nil {
			p = true
			stack = string(debug.Stack())
		}
	}()

	f()
	return
}
