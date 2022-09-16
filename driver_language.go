package captcha

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/golang/freetype/truetype"
)

var langMap = map[string][]int{
	"latin":  {0x0000, 0x007f},
	"zh":     {0x4e00, 0x9fa5},
	"ko":     {12593, 12686},
	"jp":     {12449, 12531},
	"ru":     {1025, 1169},
	"th":     {0x0e00, 0x0e7f},
	"greek":  {0x0380, 0x03ff},
	"arabic": {0x0600, 0x06ff},
	"hebrew": {0x0590, 0x05ff},
}

func generateRandomRune(size int, code string) string {
	lang, ok := langMap[code]
	if !ok {
		fmt.Sprintf("can not font language of %s", code)
		lang = langMap["latin"]
	}
	start := lang[0]
	end := lang[1]
	randRune := make([]rune, size)
	for i := range randRune {
		idx := rand.Intn(end-start) + start
		randRune[i] = rune(idx)
	}
	return string(randRune)
}

type DriverLanguage struct {
	Height int
	Width int

	NoiseCount int

	ShowLineOptions int

	Length int

	BgColor *color.RGBA

	fontsStorage FontsStorage

	Fonts        []*truetype.Font
	LanguageCode string
}

func NewDriverLanguage(height int, width int, noiseCount int, showLineOptions int, length int, bgColor *color.RGBA, fontsStorage FontsStorage, fonts []*truetype.Font, languageCode string) *DriverLanguage {
	return &DriverLanguage{Height: height, Width: width, NoiseCount: noiseCount, ShowLineOptions: showLineOptions, Length: length, BgColor: bgColor, fontsStorage: fontsStorage, Fonts: fonts, LanguageCode: languageCode}
}

func (d *DriverLanguage) GenerateIdQuestionAnswer() (id, content, answer string) {
	id = RandomId()
	content = generateRandomRune(d.Length, d.LanguageCode)
	return id, content, content
}

func (d *DriverLanguage) DrawCaptcha(content string) (item Item, err error) {
	var bgc color.RGBA
	if d.BgColor != nil {
		bgc = *d.BgColor
	} else {
		bgc = RandLightColor()
	}
	itemChar := NewItemChar(d.Width, d.Height, bgc)

	if d.ShowLineOptions&OptionShowHollowLine == OptionShowHollowLine {
		itemChar.drawHollowLine()
	}

	if d.ShowLineOptions&OptionShowSlimeLine == OptionShowSlimeLine {
		itemChar.drawSlimLine(3)
	}

	if d.ShowLineOptions&OptionShowSineLine == OptionShowSineLine {
		itemChar.drawSineLine()
	}

	if d.NoiseCount > 0 {
		noise := RandText(d.NoiseCount, TxtNumbers+TxtAlphabet+",.[]<>")
		err = itemChar.drawNoise(noise, fontsAll)
		if err != nil {
			return
		}
	}

	err = itemChar.drawText(content, []*truetype.Font{fontChinese})
	if err != nil {
		return
	}

	return itemChar, nil
}
