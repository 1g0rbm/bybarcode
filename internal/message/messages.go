package message

import "strings"

func OnStartMessage() string {
	msg := `
Привет!
Отправь мне штрихкод товара и я пришлю тебе информацию о нем.
Чтобы открыть веб-приложение воспользуйся командой:
/open
`
	return strings.Trim(msg, " ")
}
