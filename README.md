# go-wave

[![CircleCI](https://circleci.com/gh/moutend/go-wave/tree/master.svg?style=svg&circle-token=e7f4cd48b17c809c0351ad5d116f3bb4a5c8a16e)][status]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GitHub release](http://img.shields.io/github/release/moutend/go-wave.svg?style=flat-square)][release]

[status]: https://circleci.com/gh/moutend/go-wave/tree/master
[license]: https://github.com/moutend/go-wave/blob/master/LICENSE
[release]: https://github.com/moutend/go-wave/releases

`go-wave` reads and writes wave (.wav) file.

# Example

The following example concatinates `input1.wav` and `input2.wav` into `output.wav`. Note that the example assumes that the two input files have same sample rate, bit depth and channels.

```go
package main

import (
	"io"
	"io/ioutil"
	"log"

	"github.com/moutend/go-wave"
)

func main() {
	var err error
	if err = run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	var input1, input2, output *wave.WAVE

	if input1, err = wave.OpenFile("./input1.wav"); err != nil {
		return
	}
	if input2, err = wave.OpenFile("./input2.wav"); err != nil {
		return
	}
	if output, err = wave.New(input1.SamplesPerSec, input1.BitsPerSample, input1.Channels); err != nil {
		return
	}

	io.Copy(output, input1)
	io.Copy(output, input2)

	return ioutil.WriteFile("./output.wav", output.Bytes(), 0644)
}
```

## Contributing

1. Fork ([https://github.com/moutend/go-wca/fork](https://github.com/moutend/go-wca/fork))
1. Create a feature branch
1. Add changes
1. Run `go fmt` and `go test`
1. Commit your changes
1. Open a new Pull Request

## Author

[Yoshiyuki Koyanagi](https://github.com/moutend)
