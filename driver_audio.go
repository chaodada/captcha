package captcha

type DriverAudio struct {
	Length int
	Language string
}

var DefaultDriverAudio = NewDriverAudio(6, "en")

func NewDriverAudio(length int, language string) *DriverAudio {
	return &DriverAudio{Length: length, Language: language}
}

func (d *DriverAudio) DrawCaptcha(content string) (item Item, err error) {
	digits := stringToFakeByte(content)
	audio := newAudio("", digits, d.Language)
	return audio, nil
}

func (d *DriverAudio) GenerateIdQuestionAnswer() (id, q, a string) {
	id = RandomId()
	digits := randomDigits(d.Length)
	a = parseDigitsToString(digits)
	return id, a, a
}
