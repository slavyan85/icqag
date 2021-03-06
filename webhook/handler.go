package webhook

import (
	"bytes"
	"github.com/labstack/echo"
)

func (p *Provider) handleMessage(c echo.Context) error {
	payload, err := p.payloadBySourceName(c.Param("source"))
	if err != nil {
		return err
	}
	messageString, err := payload.Parse(c.Request(), p.instance.Logger)
	if err != nil {
		return err
	}
	for _, msg := range wrapMessage(messageString, 4000) {
		message := p.Bot.NewTextMessage(c.Param("target"), msg)
		err = message.Send()
		if err != nil {
			return err
		}
	}
	return nil
}

func wrapMessage(message string, maxLen int) []string {
	runes := bytes.Runes([]byte(message))
	runesCount := len(runes)
	if runesCount <= maxLen {
		return []string{message}
	}
	var result []string
	for i := 0; i < runesCount; i += maxLen {
		var s string
		start := i
		stop := i + maxLen
		if stop > runesCount {
			stop = runesCount
		}
		for _, r := range runes[start:stop] {
			s += string(r)
		}
		result = append(result, s)
	}
	return result
}
