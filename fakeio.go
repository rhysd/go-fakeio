package fakeio

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// FakedIO represents state of faking stdout/stderr/stdin.
// Restore() must be called finally to restore the state.
type FakedIO struct {
	stdout      *os.File
	stderr      *os.File
	outWriter   *os.File
	outReader   *os.File
	stdin       *os.File
	stdinWriter *os.File
	result      []byte
	err         error
}

func (fake *FakedIO) fakeOutput(name string, out *os.File) *os.File {
	if fake.err != nil {
		return out
	}
	if fake.outWriter != nil {
		return fake.outWriter
	}
	r, w, err := os.Pipe()
	if err != nil {
		fake.err = fmt.Errorf("Cannot create %s pipe: %s", name, err)
		return out
	}
	fake.outWriter = w
	fake.outReader = r
	return w
}

// Stdout replace stdout with faked output buffer
func (fake *FakedIO) Stdout() *FakedIO {
	if fake.stdout == nil {
		fake.stdout = os.Stdout
		os.Stdout = fake.fakeOutput("stdout", os.Stdout)
	}
	return fake
}

// Stderr replace stderr with faked output buffer
func (fake *FakedIO) Stderr() *FakedIO {
	if fake.stderr == nil {
		fake.stderr = os.Stderr
		os.Stderr = fake.fakeOutput("stderr", os.Stderr)
	}
	return fake
}

func (fake *FakedIO) writeToStdin(in []byte) {
	if _, err := fake.stdinWriter.Write(in); err != nil {
		fake.err = fmt.Errorf("Cannot write to piped stdin: %s", err)
	}
}

// StdinBytes sets input buffer for stdin with bytes
func (fake *FakedIO) StdinBytes(in []byte) *FakedIO {
	if fake.err != nil {
		return fake
	}

	if fake.stdinWriter != nil {
		fake.writeToStdin(in)
		return fake
	}

	r, w, err := os.Pipe()
	if err != nil {
		fake.err = fmt.Errorf("Cannot create stdin pipe: %s", err)
		return fake
	}

	fake.stdinWriter = w
	fake.stdin = os.Stdin
	os.Stdin = r
	fake.writeToStdin(in)
	return fake
}

// Stdin sets given string as faked input for stdin
func (fake *FakedIO) Stdin(in string) *FakedIO {
	return fake.StdinBytes([]byte(in))
}

// Restore restores faked stdin/stdout/stderr. This must be called finally.
func (fake *FakedIO) Restore() {
	if fake.outReader != nil {
		fake.outReader.Close()
		fake.outReader = nil
		fake.outWriter.Close()
		fake.outWriter = nil
	}
	if fake.stdinWriter != nil {
		fake.stdinWriter.Close()
		fake.stdinWriter = nil
	}
	if fake.stdout != nil {
		os.Stdout = fake.stdout
		fake.stdout = nil
	}
	if fake.stderr != nil {
		os.Stderr = fake.stderr
		fake.stderr = nil
	}
	if fake.stdin != nil {
		os.Stdin = fake.stdin
		fake.stdin = nil
	}
}

// Do runs predicate f and returns output as string
func (fake *FakedIO) Do(f func()) (string, error) {
	defer fake.Restore()
	f()
	return fake.String()
}

// Read reads bytes from buffer while faking stdout/stderr
func (fake *FakedIO) Read(p []byte) (int, error) {
	if fake.err != nil {
		return 0, fake.err
	}
	if fake.outReader == nil {
		return 0, errors.New("stdout nor stderr was not faked")
	}
	return fake.outReader.Read(p)
}

// Bytes returns buffer as []byte while faking stdout/stderr
func (fake *FakedIO) Bytes() ([]byte, error) {
	if fake.err != nil {
		return nil, fake.err
	}
	if fake.result != nil {
		return fake.result, nil
	}
	if fake.outWriter != nil {
		fake.outWriter.Close()
	}
	fake.result, fake.err = ioutil.ReadAll(fake)
	return fake.result, fake.err
}

// String returns buffer as string while faking stdout/stderr
func (fake *FakedIO) String() (string, error) {
	b, err := fake.Bytes()
	return string(b), err
}

// CloseStdin closes faked stdin
func (fake *FakedIO) CloseStdin() *FakedIO {
	if fake.err != nil {
		return fake
	}
	if fake.stdinWriter == nil {
		fake.err = errors.New("Cannot close stdin before faking it")
		return fake
	}
	fake.stdinWriter.Close()
	return fake
}

// Err returns error which occurred while setting faked stdin/stdout/stderr
func (fake *FakedIO) Err() error {
	return fake.err
}

// Stdout starts to fake stdout and returns fakedIO object to restore input/output finally
func Stdout() *FakedIO {
	f := &FakedIO{}
	return f.Stdout()
}

// Stderr starts to fake stderr and returns fakedIO object to restore input/output finally
func Stderr() *FakedIO {
	f := &FakedIO{}
	return f.Stderr()
}

// StdinBytes sets given bytes as faked stdin input and returns fakedIO object to restore input/output finally
func StdinBytes(in []byte) *FakedIO {
	f := &FakedIO{}
	return f.StdinBytes(in)
}

// Stdin sets given string as faked stdin input and returns fakedIO object to restore input/output finally
func Stdin(in string) *FakedIO {
	f := &FakedIO{}
	return f.Stdin(in)
}
