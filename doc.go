/*
Package fakeio is a small library to fake stdout/stderr/stdin mainly for
unit testing.

Following example fakes stdout/stderr/stdin. Set 'from stdin!' as stdin and get
output from both stdout and stderr as string.

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

    // Output: from stdin!
    fmt.Println(fromInput)

    // Output:
    // from stdout!
    // from stderr!
    fmt.Println(fromOutput)

Please read example/example_test.go to see live examples.

If you find some bugs, please report it to repository page: https://github.com/rhysd/go-fakeio
*/
package fakeio
