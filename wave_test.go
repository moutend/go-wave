package wave

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

func TestNew(t *testing.T) {
	var a1, a2 *WAVE
	var err error

	if a1, err = New(44100, 16, 2); err != nil {
		t.Fatal(err)
	}
	if a1.FormatTag != WAVE_FORMAT_PCM {
		t.Fatalf("FormatTag should be %d but got %d", WAVE_FORMAT_PCM, a1.FormatTag)
	}
	if a2, err = New(96000, 32, 2); err != nil {
		t.Fatal(err)
	}
	if a2.FormatTag != WAVE_FORMAT_EXTENSIBLE {
		t.Fatalf("FormatTag should be %d but got %d", WAVE_FORMAT_EXTENSIBLE, a1.FormatTag)
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
				if bits == 16 {
					if audio.FormatTag != WAVE_FORMAT_PCM {
						t.Fatalf("format tag should be %d but got %d (%s)", WAVE_FORMAT_PCM, audio.FormatTag, filename)
					}
				}
				if bits == 24 || bits == 32 {
					if audio.FormatTag != WAVE_FORMAT_EXTENSIBLE {
						t.Fatalf("format tag should be %d (%s)", WAVE_FORMAT_EXTENSIBLE, filename)
					}
				}
			}
		}
	}
	return
}

func TestBytes_16bit(t *testing.T) {
	var audio *WAVE
	var actualBytes, expectedBytes []byte
	var err error

	if audio, err = OpenFile("./testdata/sawtooth.wav"); err != nil {
		t.Fatal(err)
	}
	if expectedBytes, err = ioutil.ReadFile("./testdata/sawtooth.wav"); err != nil {
		t.Fatal(err)
	}

	actualBytes = audio.Bytes()
	sizeOfExpectedBytes := len(expectedBytes)
	sizeOfActualBytes := len(actualBytes)

	if sizeOfExpectedBytes != sizeOfActualBytes {
		t.Fatalf("expected: %d actual: %d", sizeOfExpectedBytes, sizeOfActualBytes)
	}
	for i, b := range expectedBytes {
		if b != actualBytes[i] {
			t.Fatalf("[%v] expected: %v actual: %v\n", i, b, actualBytes[i])
		}
	}
	return
}

func TestBytes_32bit(t *testing.T) {
	var audio *WAVE
	var actualBytes, expectedBytes []byte
	var err error

	if audio, err = OpenFile("./testdata/sine.wav"); err != nil {
		t.Fatal(err)
	}
	if expectedBytes, err = ioutil.ReadFile("./testdata/sine.wav"); err != nil {
		t.Fatal(err)
	}

	actualBytes = audio.Bytes()
	sizeOfExpectedBytes := len(expectedBytes)
	sizeOfActualBytes := len(actualBytes)

	if sizeOfExpectedBytes != sizeOfActualBytes {
		t.Fatalf("expected: %d actual: %d", sizeOfExpectedBytes, sizeOfActualBytes)
	}
	for i, b := range expectedBytes {
		if b != actualBytes[i] {
			t.Fatalf("[%v] expected: %v actual: %v\n", i, b, actualBytes[i])
		}
	}
	return
}

func TestRead_16bit(t *testing.T) {
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

func TestRead_32bit(t *testing.T) {
	var audio *WAVE
	var rawdata []byte
	var buf []byte
	var err error

	if audio, err = OpenFile("./testdata/sine.wav"); err != nil {
		t.Fatal(err)
	}
	if rawdata, err = ioutil.ReadFile("./testdata/sine.raw"); err != nil {
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

func TestWrite_16bit(t *testing.T) {
	var n int64
	var err error

	src, _ := OpenFile("./testdata/sawtooth.wav")
	dest, _ := New(src.SamplesPerSec, src.BitsPerSample, src.Channels)

	if n, err = io.Copy(dest, src); err != nil {
		t.Fatal(err)
	}
	if n != int64(src.DataSize) {
		t.Errorf("expect: %v actual: %v", n, dest.DataSize)
	}
	return
}

func TestWrite_32bit(t *testing.T) {
	var n int64
	var err error

	src, _ := OpenFile("./testdata/sine.wav")
	dest, _ := New(src.SamplesPerSec, src.BitsPerSample, src.Channels)

	if n, err = io.Copy(dest, src); err != nil {
		t.Fatal(err)
	}
	if n != int64(src.DataSize) {
		t.Errorf("expect: %v actual: %v", n, dest.DataSize)
	}
	return
}
