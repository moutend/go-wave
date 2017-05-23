package wave

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

func TestBytes(t *testing.T) {
	var audio *WAVE
	var file []byte
	var err error

	if audio, err = OpenFile("./testdata/sawtooth.wav"); err != nil {
		t.Fatal(err)
	}
	if file, err = ioutil.ReadFile("./testdata/sawtooth.wav"); err != nil {
		t.Fatal(err)
	}
	actual := audio.Bytes()

	for i, b := range file {
		if b != actual[i] {
			t.Errorf("[%v] expected: %v actual: %v\n", i, b, actual[i])
		}
	}
	return
}

func TestOpenFile(t *testing.T) {
	var audio *WAVE
	var err error

	for _, samples := range []uint32{44100, 48000, 96000, 192000} {
		for _, bits := range []uint16{16, 24, 32} {
			for _, channels := range []uint16{1, 2} {
				filename := fmt.Sprintf("./testdata/%vHz-%vbit-%vch-empty.wav", samples, bits, channels)
				if audio, err = OpenFile(filename); err != nil {
					t.Fatal(err)
				}
				if audio.SamplesPerSec != samples {
					t.Errorf("expected: %v actual: %v (%v)\n", samples, audio.SamplesPerSec, filename)
				}
				if audio.BitsPerSample != bits {
					t.Errorf("expected: %v actual: %v (%v)\n", bits, audio.BitsPerSample, filename)
				}
				if audio.Channels != channels {
					t.Errorf("expected: %v actual: %v\n (%v)", channels, audio.Channels, filename)
				}
			}
		}
	}
	return
}

func TestRead(t *testing.T) {
	var audio *WAVE
	var rawdata []byte
	var buf []byte
	var err error

	if audio, err = OpenFile("./testdata/sawtooth.wav"); err != nil {
		t.Fatal(err)
	}
	if rawdata, err = ioutil.ReadFile("./testdata/sawtooth.raw"); err != nil {
		t.Fatal(err)
	}
	if buf, err = ioutil.ReadAll(audio); err != nil {
		t.Fatal(err)
	}

	size := len(rawdata)

	for i := 0; i < size; i++ {
		if buf[i] != rawdata[i] {
			t.Fatalf("[%v] expected: %v actual: %v", i, rawdata[i], buf[i])
		}
	}
	return
}

func TestWrite(t *testing.T) {
	var n int64
	var err error

	src, _ := OpenFile("./testdata/sawtooth.wav")
	dest, _ := New(src.SamplesPerSec, src.BitsPerSample, src.Channels)
	if n, err = io.Copy(dest, src); err != nil {
		t.Fatal(err)
	}
	if n != int64(dest.DataSize) {
		t.Errorf("expect: %v actual: %v", n, dest.DataSize)
	}
	return
}
