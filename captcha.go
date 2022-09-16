package captcha

import "strings"

type Captcha struct {
	Driver Driver
	Store  Store
}

func NewCaptcha(driver Driver, store Store) *Captcha {
	return &Captcha{Driver: driver, Store: store}
}

func (c *Captcha) Generate() (id, b64s string, err error) {
	id, content, answer := c.Driver.GenerateIdQuestionAnswer()
	item, err := c.Driver.DrawCaptcha(content)
	if err != nil {
		return "", "", err
	}
	err = c.Store.Set(id, answer)
	if err != nil {
		return "", "", err
	}
	b64s = item.EncodeB64string()
	return
}


func (c *Captcha) Verify(id, answer string, clear bool) (match bool) {
	vv := c.Store.Get(id, clear)
	vv = strings.TrimSpace(vv)
	return vv == strings.TrimSpace(answer)
}
