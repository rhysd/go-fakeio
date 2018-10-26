package fakeio

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestFakeStdout(t *testing.T) {
	stdin, stderr := os.Stdin, os.Stderr

	f := Stdout()
	defer f.Restore()

	if os.Stdin != stdin {
		t.Error("stdin was unexpectedly modified", os.Stdin)
	}
	if os.Stderr != stderr {
		t.Error("stderr was unexpectedly modified", os.Stderr)
	}

	fmt.Print("Hello")
	fmt.Fprint(os.Stderr, "Bye")

	s, err := f.String()
	if err != nil {
		t.Fatal(err)
	}
	if s != "Hello" {
		t.Fatalf("want 'Hello' but have '%s'", s)
	}
	if f.Err() != nil {
		t.Fatal(f.Err())
	}
}

func TestFakeStderr(t *testing.T) {
	stdin, stdout := os.Stdin, os.Stdout

	f := Stderr()
	defer f.Restore()

	if os.Stdin != stdin {
		t.Error("stdin was unexpectedly modified", os.Stdin)
	}
	if os.Stdout != stdout {
		t.Error("stdout was unexpectedly modified", os.Stdout)
	}

	fmt.Print("Hello")
	fmt.Fprint(os.Stderr, "Bye")

	s, err := f.String()
	if err != nil {
		t.Fatal(err)
	}
	if s != "Bye" {
		t.Fatalf("want 'Bye' but have '%s'", s)
	}
	if f.Err() != nil {
		t.Fatal(f.Err())
	}
}

func TestFakeStdin(t *testing.T) {
	stderr, stdout := os.Stderr, os.Stdout

	f := Stdin("hello!\n").Stdin("how are you?\n").Stdin("bye!\n")
	defer f.Restore()

	if os.Stdout != stdout {
		t.Error("stdout was unexpectedly modified", os.Stdout)
	}
	if os.Stderr != stderr {
		t.Error("stderr was unexpectedly modified", os.Stderr)
	}

	r := bufio.NewReader(os.Stdin)
	have := make([]string, 0, 3)
	for i := 0; i < 3; i++ {
		s, err := r.ReadString('\n')
		if err != nil {
			t.Fatal(err)
		}
		have = append(have, s)
	}

	want := []string{
		"hello!\n",
		"how are you?\n",
		"bye!\n",
	}
	if !reflect.DeepEqual(have, want) {
		t.Fatalf("want: '%#v' but have '%#v'", want, have)
	}
	if f.Err() != nil {
		t.Fatal(f.Err())
	}
}

func TestFakeAll(t *testing.T) {
	for i, f := range []func() *FakedIO{
		func() *FakedIO { return Stdout().Stderr().Stdin("hello\n") },
		func() *FakedIO { return Stdout().Stdin("hello\n").Stderr() },
		func() *FakedIO { return Stderr().Stdout().Stdin("hello\n") },
		func() *FakedIO { return Stderr().Stdin("hello\n").Stdout() },
		func() *FakedIO { return Stdin("hello\n").Stdout().Stderr() },
		func() *FakedIO { return Stdin("hello\n").Stderr().Stdout() },
		func() *FakedIO { return StdinBytes([]byte("hello\n")).Stderr().Stdout() },
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			fake := f()
			defer fake.Restore()
			fmt.Fprintln(os.Stderr, "bar!")
			fmt.Println("foo!")
			in, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				t.Fatal(err)
			}
			if in != "hello\n" {
				t.Errorf("want: 'hello\\n' but have '%#v'", in)
			}
			have, err := fake.String()
			if err != nil {
				t.Fatal(err)
			}
			want := "bar!\nfoo!\n"
			if have != want {
				t.Fatalf("want '%#v' but have '%#v'", want, have)
			}
		})
	}
}

func TestRestore(t *testing.T) {
	stdin, stderr, stdout := os.Stdin, os.Stderr, os.Stdout
	f := Stdout().Stderr().Stdin("hello")
	f.Restore()
	if os.Stdin != stdin {
		t.Error("stdin was not restored")
	}
	if os.Stderr != stderr {
		t.Error("stderr was not restored")
	}
	if os.Stdout != stdout {
		t.Error("stdout was not restored")
	}
}

func TestRepeatStdout(t *testing.T) {
	f := Stdout().Stdout().Stdout()
	defer f.Restore()
	want := "hello, world\n"
	fmt.Print(want)
	have, err := f.String()
	if err != nil {
		t.Fatal(err)
	}
	if want != have {
		t.Fatalf("want '%#v' but have '%#v'", want, have)
	}
}

func TestRepeatStderr(t *testing.T) {
	f := Stderr().Stderr().Stderr()
	defer f.Restore()
	want := "goodbye, world\n"
	fmt.Fprint(os.Stderr, want)
	have, err := f.String()
	if err != nil {
		t.Fatal(err)
	}
	if want != have {
		t.Fatalf("want '%#v' but have '%#v'", want, have)
	}
}

func TestRepeatGetResult(t *testing.T) {
	f := Stdout()
	defer f.Restore()
	want := "hello, world\n"
	fmt.Print(want)
	for i := 0; i < 5; i++ {
		have, err := f.String()
		if err != nil {
			t.Fatal(err)
		}
		if want != have {
			t.Fatalf("want '%#v' but have '%#v'", want, have)
		}
	}
}

func TestFakeNotSetButGetResult(t *testing.T) {
	f := Stdin("hello")
	defer f.Restore()

	_, err := f.String()
	if err == nil {
		t.Fatal("Error did not occur")
	}
	if !strings.Contains(err.Error(), "stdout nor stderr was not faked") {
		t.Fatal("Unexpected error:", err)
	}
}

func TestDo(t *testing.T) {
	have, err := Stdin("hello\n").Stderr().Stdout().Do(func() {
		fmt.Fprintln(os.Stderr, "bar!")
		fmt.Println("foo!")
		in, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			t.Fatal(err)
		}
		if in != "hello\n" {
			t.Errorf("want: 'hello\\n' but have '%#v'", in)
		}
	})
	if err != nil {
		t.Fatal(err)
	}

	want := "bar!\nfoo!\n"
	if have != want {
		t.Fatalf("want '%#v' but have '%#v'", want, have)
	}
}

func TestDoRestore(t *testing.T) {
	stdin, stderr, stdout := os.Stdin, os.Stderr, os.Stdout

	if _, err := Stdin("hello\n").Stderr().Stdout().Do(func() {}); err != nil {
		t.Fatal(err)
	}

	if os.Stdin != stdin {
		t.Error("stdin was not restored")
	}
	if os.Stderr != stderr {
		t.Error("stderr was not restored")
	}
	if os.Stdout != stdout {
		t.Error("stdout was not restored")
	}
}

func TestDoRestoreOnPanic(t *testing.T) {
	stdin, stderr, stdout := os.Stdin, os.Stderr, os.Stdout

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Did not panic")
		}

		if os.Stdin != stdin {
			t.Error("stdin was not restored")
		}
		if os.Stderr != stderr {
			t.Error("stderr was not restored")
		}
		if os.Stdout != stdout {
			t.Error("stdout was not restored")
		}
	}()

	Stdin("hello\n").Stderr().Stdout().Do(func() {
		panic("oops!!")
	})
}

func TestCloseStdin(t *testing.T) {
	want := "hello"
	fake := Stdin(want).CloseStdin()
	defer fake.Restore()

	// When stdin is not closed, this line blocks forever
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		t.Fatal(err)
	}

	have := string(b)
	if want != have {
		t.Fatalf("want '%#v' but have '%#v'", want, have)
	}
}

func TestRestoreNothing(t *testing.T) {
	stdin, stderr, stdout := os.Stdin, os.Stderr, os.Stdout
	f := &FakedIO{}
	f.Restore()
	if os.Stdin != stdin {
		t.Error("stdin was modified", os.Stdin)
	}
	if os.Stderr != stderr {
		t.Error("stderr was modified", os.Stderr)
	}
	if os.Stdout != stdout {
		t.Error("stdout was modified", os.Stdout)
	}
}
