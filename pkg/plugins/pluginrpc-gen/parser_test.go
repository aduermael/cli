package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const testFixture = "fixtures/foo.go"

func TestParseEmptyInterface(t *testing.T) {
	pkg, err := Parse(testFixture, "Fooer")
	if err != nil {
		t.Fatal(err)
	}

	assertName(t, "foo", pkg.Name)
	assertNum(t, 0, len(pkg.Functions))
}

func TestParseNonInterfaceType(t *testing.T) {
	_, err := Parse(testFixture, "wobble")
	if _, ok := err.(errUnexpectedType); !ok {
		t.Fatal("expected type error when parsing non-interface type")
	}
}

func TestParseWithOneFunction(t *testing.T) {
	pkg, err := Parse(testFixture, "Fooer2")
	if err != nil {
		t.Fatal(err)
	}

	assertName(t, "foo", pkg.Name)
	assertNum(t, 1, len(pkg.Functions))
	assertName(t, "Foo", pkg.Functions[0].Name)
	assertNum(t, 0, len(pkg.Functions[0].Args))
	assertNum(t, 0, len(pkg.Functions[0].Returns))
}

func TestParseWithMultipleFuncs(t *testing.T) {
	pkg, err := Parse(testFixture, "Fooer3")
	if err != nil {
		t.Fatal(err)
	}

	assertName(t, "foo", pkg.Name)
	assertNum(t, 6, len(pkg.Functions))

	f := pkg.Functions[0]
	assertName(t, "Foo", f.Name)
	assertNum(t, 0, len(f.Args))
	assertNum(t, 0, len(f.Returns))

	f = pkg.Functions[1]
	assertName(t, "Bar", f.Name)
	assertNum(t, 1, len(f.Args))
	assertNum(t, 0, len(f.Returns))
	arg := f.Args[0]
	assertName(t, "a", arg.Name)
	assertName(t, "string", arg.ArgType)

	f = pkg.Functions[2]
	assertName(t, "Baz", f.Name)
	assertNum(t, 1, len(f.Args))
	assertNum(t, 1, len(f.Returns))
	arg = f.Args[0]
	assertName(t, "a", arg.Name)
	assertName(t, "string", arg.ArgType)
	arg = f.Returns[0]
	assertName(t, "err", arg.Name)
	assertName(t, "error", arg.ArgType)

	f = pkg.Functions[3]
	assertName(t, "Qux", f.Name)
	assertNum(t, 2, len(f.Args))
	assertNum(t, 2, len(f.Returns))
	arg = f.Args[0]
	assertName(t, "a", f.Args[0].Name)
	assertName(t, "string", f.Args[0].ArgType)
	arg = f.Args[1]
	assertName(t, "b", arg.Name)
	assertName(t, "string", arg.ArgType)
	arg = f.Returns[0]
	assertName(t, "val", arg.Name)
	assertName(t, "string", arg.ArgType)
	arg = f.Returns[1]
	assertName(t, "err", arg.Name)
	assertName(t, "error", arg.ArgType)

	f = pkg.Functions[4]
	assertName(t, "Wobble", f.Name)
	assertNum(t, 0, len(f.Args))
	assertNum(t, 1, len(f.Returns))
	arg = f.Returns[0]
	assertName(t, "w", arg.Name)
	assertName(t, "*wobble", arg.ArgType)

	f = pkg.Functions[5]
	assertName(t, "Wiggle", f.Name)
	assertNum(t, 0, len(f.Args))
	assertNum(t, 1, len(f.Returns))
	arg = f.Returns[0]
	assertName(t, "w", arg.Name)
	assertName(t, "wobble", arg.ArgType)
}

func TestParseWithUnamedReturn(t *testing.T) {
	_, err := Parse(testFixture, "Fooer4")
	if !strings.HasSuffix(err.Error(), errBadReturn.Error()) {
		t.Fatalf("expected ErrBadReturn, got %v", err)
	}
}

func TestEmbeddedInterface(t *testing.T) {
	pkg, err := Parse(testFixture, "Fooer5")
	if err != nil {
		t.Fatal(err)
	}

	assertName(t, "foo", pkg.Name)
	assertNum(t, 2, len(pkg.Functions))

	f := pkg.Functions[0]
	assertName(t, "Foo", f.Name)
	assertNum(t, 0, len(f.Args))
	assertNum(t, 0, len(f.Returns))

	f = pkg.Functions[1]
	assertName(t, "Boo", f.Name)
	assertNum(t, 2, len(f.Args))
	assertNum(t, 2, len(f.Returns))

	arg := f.Args[0]
	assertName(t, "a", arg.Name)
	assertName(t, "string", arg.ArgType)

	arg = f.Args[1]
	assertName(t, "b", arg.Name)
	assertName(t, "string", arg.ArgType)

	arg = f.Returns[0]
	assertName(t, "s", arg.Name)
	assertName(t, "string", arg.ArgType)

	arg = f.Returns[1]
	assertName(t, "err", arg.Name)
	assertName(t, "error", arg.ArgType)
}

func assertName(t *testing.T, expected, actual string) {
	if expected != actual {
		fatalOut(t, fmt.Sprintf("expected name to be `%s`, got: %s", expected, actual))
	}
}

func assertNum(t *testing.T, expected, actual int) {
	if expected != actual {
		fatalOut(t, fmt.Sprintf("expected number to be %d, got: %d", expected, actual))
	}
}

func fatalOut(t *testing.T, msg string) {
	_, file, ln, _ := runtime.Caller(2)
	t.Fatalf("%s:%d: %s", filepath.Base(file), ln, msg)
}
