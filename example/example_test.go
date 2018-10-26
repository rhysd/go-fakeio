package fakeioexample

import (
	"bufio"
	"fmt"
	"github.com/rhysd/go-fakeio"
	"os"
)

func Example_FakeStdout() {
	f := fakeio.Stdout()

	fmt.Print("Hello")

	s, err := f.String()
	if err != nil {
		f.Restore()
		panic(err)
	}

	// 'defer' is better, but here it's unavailable due to output test
	f.Restore()

	fmt.Println(s)

	// Output:
	// Hello
}

func Example_FakeStderr() {
	f := fakeio.Stderr()

	fmt.Fprint(os.Stderr, "Hello")

	s, err := f.String()
	if err != nil {
		f.Restore()
		panic(err)
	}

	// 'defer' is better, but here it's unavailable due to output test
	f.Restore()

	fmt.Println(s)

	// Output:
	// Hello
}

func Example_FakeStdin() {
	f := fakeio.Stdin("Bye!")

	s, err := bufio.NewReader(os.Stdin).ReadString('!')
	if err != nil {
		f.Restore()
		panic(err)
	}

	// 'defer' is better, but here it's unavailable due to output test
	f.Restore()

	fmt.Println(s)

	// Output:
	// Bye!
}

func Example_FakeAll() {
	f := fakeio.Stdout().Stderr().Stdin("from stdin!")

	fromInput, err := bufio.NewReader(os.Stdin).ReadString('!')
	if err != nil {
		f.Restore()
		panic(err)
	}

	fmt.Println("from stdout!")
	fmt.Fprintln(os.Stderr, "from stderr!")

	fromOutput, err := f.String()
	if err != nil {
		f.Restore()
		panic(err)
	}

	// 'defer' is better, but here it's unavailable due to output test
	f.Restore()

	fmt.Println(fromInput)
	fmt.Println(fromOutput)

	// Output:
	// from stdin!
	// from stdout!
	// from stderr!
	//
}
