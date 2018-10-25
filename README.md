`fakeio` pacakge for Go
=======================

`fakeio` is a small library to fake stdout/stderr/stdin.
This is mainly for unit testing of CLI applications.

## Usage

Basic usage:

```go
import "github.com/rhysd/fakeio"

// Fake stdout and input 'hello' to stdin
fake := fakeio.Stdout().Stdin("hello")
defer fake.Restore()

// Do something...

// Get bufferred stdout as string
s, err := fake.String()
if err != nil {
    panic(err)
}

fmt.Println(s)
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

Please see [examples](example/example_test.go) for actual examples.


## License

[MIT License](LICENSE.txt)

