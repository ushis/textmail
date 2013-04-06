# Textmail

Textmail is a [Go](http://golang.org/) library to produce nice plain text
from HTML. It uses the [gokogiri](https://github.com/moovweb/gokogiri) library
to parse the input.

Browse the testdata directory to see some input/output samples.

## Use it

A simple command line tool could look like this.

```Go
package main

import (
  "fmt"
  "io/ioutil"
  "github.com/ushis/textmail"
  "os"
)

func main() {
  in, err := ioutil.ReadAll(os.Stdin)

  if err != nil {
    panic(err)
  }
  out, err := textmail.Format(in)

  if err != nil {
    panic(err)
  }
  fmt.Println(string(out))
}
```

## Hack it

Clone the repo and run the tests to verify your work.

```
go get github.com/moovweb/gokogiri
git clone git://github.com/ushis/textmail.git
cd textmail
go test
```

## License (MIT)

Copyright (c) 2013 ushi

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
