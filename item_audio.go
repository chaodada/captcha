
package captcha

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
)

type ItemAudio struct {
	answer      string
	body        *bytes.Buffer
	digitSounds [][]byte
}


func newAudio(id string, digits []byte, lang string) *ItemAudio {
	a := new(ItemAudio)

	if sounds, ok := digitSounds[lang]; ok {
		a.digitSounds = sounds
	} else {
		a.digitSounds = digitSounds["en"]
	}
	numsnd := make([][]byte, len(digits))
	for i, n := range digits {
		snd := a.randomizedDigitSound(n)
		setSoundLevel(snd, 1.5)
		numsnd[i] = snd
	}
	intervals := make([]int, len(digits)+1)
	intdur := 0
	for i := range intervals {
		dur := randIntRange(sampleRate, sampleRate*2)
		intdur += dur
		intervals[i] = dur
	}
	bg := a.makeBackgroundSound(a.longestDigitSndLen()*len(digits) + intdur)
	sil := makeSilence(sampleRate / 5)
	bufcap := 3*len(beepSound) + 2*len(sil) + len(bg) + len(endingBeepSound)
	a.body = bytes.NewBuffer(make([]byte, 0, bufcap))
	a.body.Write(beepSound)
	a.body.Write(sil)
	a.body.Write(beepSound)
	a.body.Write(sil)
	a.body.Write(beepSound)
	pos := intervals[0]
	for i, v := range numsnd {
		mixSound(bg[pos:], v)
		pos += len(v) + intervals[i+1]
	}
	a.body.Write(bg)
	a.body.Write(endingBeepSound)
	return a
}

func (a *ItemAudio) encodedLen() int {
	return len(waveHeader) + 4 + a.body.Len()
}

func (a *ItemAudio) makeBackgroundSound(length int) []byte {
	b := a.makeWhiteNoise(length, 4)
	for i := 0; i < length/(sampleRate/10); i++ {
		snd := reversedSound(a.digitSounds[rand.Intn(10)])
		place := rand.Intn(len(b) - len(snd))
		setSoundLevel(snd, randFloat64Range(0.04, 0.08))
		mixSound(b[place:], snd)
	}
	return b
}

func (a *ItemAudio) randomizedDigitSound(n byte) []byte {
	s := a.randomSpeed(a.digitSounds[n])
	setSoundLevel(s, randFloat64Range(0.85, 1.2))
	return s
}

func (a *ItemAudio) longestDigitSndLen() int {
	n := 0
	for _, v := range a.digitSounds {
		if n < len(v) {
			n = len(v)
		}
	}
	return n
}

func (a *ItemAudio) randomSpeed(b []byte) []byte {
	pitch := randFloat64Range(0.95, 1.1)
	return changeSpeed(b, pitch)
}

func (a *ItemAudio) makeWhiteNoise(length int, level uint8) []byte {
	noise := randBytes(length)
	adj := 128 - level/2
	for i, v := range noise {
		v %= level
		v += adj
		noise[i] = v
	}
	return noise
}


func (a *ItemAudio) WriteTo(w io.Writer) (n int64, err error) {
	bodyLen := uint32(a.body.Len())
	paddedBodyLen := bodyLen
	if bodyLen%2 != 0 {
		paddedBodyLen++
	}
	totalLen := uint32(len(waveHeader)) - 4 + paddedBodyLen
	header := make([]byte, len(waveHeader)+4)
	copy(header, waveHeader)
	binary.LittleEndian.PutUint32(header[4:], totalLen)
	binary.LittleEndian.PutUint32(header[len(waveHeader):], bodyLen)
	nn, err := w.Write(header)
	n = int64(nn)
	if err != nil {
		return
	}
	n, err = a.body.WriteTo(w)
	n += int64(nn)
	if err != nil {
		return
	}
	if bodyLen != paddedBodyLen {
		w.Write([]byte{0})
		n++
	}
	return
}

func (a *ItemAudio) EncodeB64string() string {
	var buf bytes.Buffer
	if _, err := a.WriteTo(&buf); err != nil {
		panic(err)
	}
	return fmt.Sprintf("data:%s;base64,%s", MimeTypeAudio, base64.StdEncoding.EncodeToString(buf.Bytes()))

}
