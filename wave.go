// Package wave reads and writes wave (.wav) file.
package wave

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
)

type WAVE struct {
	FormatTag      uint16
	Channels       uint16
	SamplesPerSec  uint32
	AvgBytesPerSec uint32
	BlockAlign     uint16
	BitsPerSample  uint16
	DataSize       uint32
	RawData        []byte
	offset         int
}

func (v *WAVE) Read(p []byte) (n int, err error) {
	frames := len(v.RawData)
	size := len(p)

	for n = 0; n < size; n++ {
		i := v.offset + n
		if i >= frames {
			return n, io.EOF
		}
		p[n] = v.RawData[i]
	}
	v.offset += size
	return
}

func (v *WAVE) Write(b []byte) (n int, err error) {
	size := len(b)

	for n = 0; n < size; n++ {
		v.RawData = append(v.RawData, b[n])
	}
	v.DataSize += uint32(size)
	return
}

func (v *WAVE) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, []byte("RIFF"))
	binary.Write(buf, binary.LittleEndian, uint32(v.DataSize+36))
	binary.Write(buf, binary.BigEndian, []byte("WAVEfmt "))
	binary.Write(buf, binary.LittleEndian, uint32(16))
	binary.Write(buf, binary.LittleEndian, uint16(1))
	binary.Write(buf, binary.LittleEndian, v.Channels)
	binary.Write(buf, binary.LittleEndian, v.SamplesPerSec)
	binary.Write(buf, binary.LittleEndian, v.AvgBytesPerSec)
	binary.Write(buf, binary.LittleEndian, v.BlockAlign)
	binary.Write(buf, binary.LittleEndian, v.BitsPerSample)
	binary.Write(buf, binary.BigEndian, []byte("data"))
	binary.Write(buf, binary.LittleEndian, v.DataSize)
	binary.Write(buf, binary.LittleEndian, v.RawData)

	return buf.Bytes()
}

func OpenFile(path string) (audio *WAVE, err error) {
	var file []byte

	if file, err = ioutil.ReadFile(path); err != nil {
		return
	}

	audio = &WAVE{}
	reader := bytes.NewReader(file)

	binary.Read(io.NewSectionReader(reader, 20, 2), binary.LittleEndian, &audio.FormatTag)
	binary.Read(io.NewSectionReader(reader, 22, 2), binary.LittleEndian, &audio.Channels)
	binary.Read(io.NewSectionReader(reader, 24, 4), binary.LittleEndian, &audio.SamplesPerSec)
	binary.Read(io.NewSectionReader(reader, 28, 4), binary.LittleEndian, &audio.AvgBytesPerSec)
	binary.Read(io.NewSectionReader(reader, 32, 2), binary.LittleEndian, &audio.BlockAlign)
	binary.Read(io.NewSectionReader(reader, 34, 2), binary.LittleEndian, &audio.BitsPerSample)
	binary.Read(io.NewSectionReader(reader, 40, 4), binary.LittleEndian, &audio.DataSize)

	buf := new(bytes.Buffer)
	io.Copy(buf, io.NewSectionReader(reader, 44, int64(audio.DataSize)))
	audio.RawData = buf.Bytes()

	return
}

func New(samplesPerSec uint32, bitsPerSample, channels uint16) (audio *WAVE, err error) {
	audio = &WAVE{}
	audio.SamplesPerSec = samplesPerSec
	audio.Channels = channels
	audio.BitsPerSample = bitsPerSample
	audio.BlockAlign = audio.Channels * audio.BitsPerSample / 8
	audio.AvgBytesPerSec = audio.SamplesPerSec * uint32(audio.BlockAlign)
	audio.RawData = []byte{}
	return
}
