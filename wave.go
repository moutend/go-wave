// Package wave reads and writes wave (.wav) file.
package wave

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

const (
	WAVE_FORMAT_PCM        = 0x1
	WAVE_FORMAT_EXTENSIBLE = 0xFFFE
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
	if v.FormatTag == WAVE_FORMAT_PCM {
		binary.Write(buf, binary.LittleEndian, uint32(v.DataSize+36))
	} else if v.FormatTag == WAVE_FORMAT_EXTENSIBLE {
		binary.Write(buf, binary.LittleEndian, uint32(v.DataSize+72))
	}
	binary.Write(buf, binary.BigEndian, []byte("WAVEfmt "))
	if v.FormatTag == WAVE_FORMAT_PCM {
		binary.Write(buf, binary.LittleEndian, uint32(16))
	} else if v.FormatTag == WAVE_FORMAT_EXTENSIBLE {
		binary.Write(buf, binary.LittleEndian, uint32(60))
	}
	binary.Write(buf, binary.LittleEndian, v.FormatTag)
	binary.Write(buf, binary.LittleEndian, v.Channels)
	binary.Write(buf, binary.LittleEndian, v.SamplesPerSec)
	binary.Write(buf, binary.LittleEndian, v.AvgBytesPerSec)
	binary.Write(buf, binary.LittleEndian, v.BlockAlign)
	binary.Write(buf, binary.LittleEndian, v.BitsPerSample)
	if v.FormatTag == WAVE_FORMAT_EXTENSIBLE {
		fmt.Println("hoge")
		binary.Write(buf, binary.LittleEndian, uint16(22))      // cbSize
		binary.Write(buf, binary.LittleEndian, uint16(22))      // cbSize
		binary.Write(buf, binary.LittleEndian, v.BitsPerSample) // validBitsPerSample
		binary.Write(buf, binary.LittleEndian, uint16(0))       // samplesPerBlock
		binary.Write(buf, binary.LittleEndian, uint16(0))       // reserved
		binary.Write(buf, binary.LittleEndian, v.Channels)      // channelMask
		guid := [16]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0xaa, 0x00, 0x38, 0x9b, 0x71}
		binary.Write(buf, binary.BigEndian, guid)
		binary.Write(buf, binary.BigEndian, []byte("fact"))
		binary.Write(buf, binary.LittleEndian, uint32(4))
	}
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
	if !(audio.FormatTag == WAVE_FORMAT_PCM || audio.FormatTag == WAVE_FORMAT_EXTENSIBLE) {
		err = fmt.Errorf("UnknownformatError")
		return
	}

	binary.Read(io.NewSectionReader(reader, 22, 2), binary.LittleEndian, &audio.Channels)
	binary.Read(io.NewSectionReader(reader, 24, 4), binary.LittleEndian, &audio.SamplesPerSec)
	binary.Read(io.NewSectionReader(reader, 28, 4), binary.LittleEndian, &audio.AvgBytesPerSec)
	binary.Read(io.NewSectionReader(reader, 32, 2), binary.LittleEndian, &audio.BlockAlign)
	binary.Read(io.NewSectionReader(reader, 34, 2), binary.LittleEndian, &audio.BitsPerSample)

	if audio.FormatTag == WAVE_FORMAT_PCM {
		binary.Read(io.NewSectionReader(reader, 40, 4), binary.LittleEndian, &audio.DataSize)
	} else if audio.FormatTag == WAVE_FORMAT_EXTENSIBLE {
		binary.Read(io.NewSectionReader(reader, 76, 4), binary.LittleEndian, &audio.DataSize)
	}

	buf := new(bytes.Buffer)
	if audio.FormatTag == WAVE_FORMAT_PCM {
		io.Copy(buf, io.NewSectionReader(reader, 44, int64(audio.DataSize)))
	} else if audio.FormatTag == WAVE_FORMAT_EXTENSIBLE {
		io.Copy(buf, io.NewSectionReader(reader, 80, int64(audio.DataSize)))
	}
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
