`fakeio` pacakge for Go
=======================
[![Linux/macOS Build Status](https://travis-ci.org/rhysd/go-fakeio.svg?branch=master)](https://travis-ci.org/rhysd/go-fakeio)
[![Windows Build status](https://ci.appveyor.com/api/projects/status/5b9t6932m5dt2e23/branch/master?svg=true)](https://ci.appveyor.com/project/rhysd/go-fakeio/branch/master)
[![Documentation](https://godoc.org/github.com/rhysd/go-fakeio?status.svg)](http://godoc.org/github.com/rhysd/go-fakeio)

[`fakeio`](https://github.com/rhysd/go-fakeio) is a small library to fake stdout/stderr/stdin.
This is mainly for unit testing of CLI applications. Please see [documentation](https://godoc.org/github.com/rhysd/go-fakeio)
for more details.

## Installation

```
$ go get github.com/rhysd/go-fakeio
```

## Usage

Basic usage:

```go
import (
    "bufio"
    "github.com/rhysd/go-fakeio"
)

// Fake stdout and input 'hello' to stdin
fake := fakeio.Stdout().Stdin("hello!")
defer fake.Restore()

// Do something...

// "hello!" is stored to variable `i`
i, err := bufio.NewReader(os.Stdin).ReadString('!')

// At this point, this line outputs nothing
fmt.Print("bye")

// "bye" is stored to variable `o`
o, err := fake.String()
```

### Faking stdout/stderr/stdin

Faking stderr:

```go
fake := fakeio.Stderr()
defer fake.Restore()
```

Faking stdin:

```go
fake, err := fakeio.Stdin("hello")
defer fake.Restore()
```

Faking stderr/stdout/stdin

```go
fake := fakeio.Stderr().Stdout().Stdin("Faked input to stdin")
defer fake.Restore()
```

### Read bufferred stdout/stderr

Reading as string:

```go
s, err := fake.String()
if err != nil {
    // Faking IO failed
    panic(err)
}
fmt.Println(s)
```

Reading as bytes:

```go
b, err := fake.Bytes()
if err != nil {
    // Faking IO failed
    panic(err)
}
fmt.Println(b)
```

Reading via `io.Reader` interface:

```go
s := bufio.NewScanner(fake)
for s.Scan() {
    // Reading line by line
    line := s.Text()
    fmt.Println(line)
}
if s.Err() != nil {
    // Error happened while reading
    panic(s.Err)
}
```

### Shortcut

`.Do()` is a shortcut

```go
s, err := fakeio.Stderr().Stdout().Do(func () {
    // Do something

    // Faked stderr and stdout are restored at exit of this scope
})
if err != nil {
    // Faking IO failed
    panic(err)
}
fmt.Println(s)
```

### Examples

Please see [examples](example_test.go) for actual examples.

## Repository

https://github.com/rhysd/go-fakeio

## License

[MIT License](LICENSE.txt)

