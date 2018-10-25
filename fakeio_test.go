package fakeio

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestFakeStdout(t *testing.T) {
	f := Stdout()
	defer f.Restore()

	fmt.Print("Hello")
	fmt.Fprint(os.Stderr, "Bye")

	s, err := f.String()
	if err != nil {
		t.Fatal(err)
	}
	if s != "Hello" {
		t.Fatalf("want 'Hello' but have '%s'", s)
	}
}

func TestFakeStderr(t *testing.T) {
	f := Stderr()
	defer f.Restore()

	fmt.Print("Hello")
	fmt.Fprint(os.Stderr, "Bye")

	s, err := f.String()
	if err != nil {
		t.Fatal(err)
	}
	if s != "Bye" {
		t.Fatalf("want 'Bye' but have '%s'", s)
	}
}

func TestFakeStdin(t *testing.T) {
	f := Stdin("hello!\n").Stdin("how are you?\n").Stdin("bye!\n")
	defer f.Restore()

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
}

func TestFakeAll(t *testing.T) {
	for i, f := range []func() *fakedIO{
		func() *fakedIO { return Stdout().Stderr().Stdin("hello\n") },
		func() *fakedIO { return Stdout().Stdin("hello\n").Stderr() },
		func() *fakedIO { return Stderr().Stdout().Stdin("hello\n") },
		func() *fakedIO { return Stderr().Stdin("hello\n").Stdout() },
		func() *fakedIO { return Stdin("hello\n").Stdout().Stderr() },
		func() *fakedIO { return Stdin("hello\n").Stderr().Stdout() },
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
		t.Error("stdin was not restured")
	}
	if os.Stderr != stderr {
		t.Error("stderr was not restured")
	}
	if os.Stdout != stdout {
		t.Error("stdout was not restured")
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
