package pinchito

import (
	"fmt"
)

type Log struct {
	Id           int
	Text         string
	Protagonista User
	Autor        User
	Titol        string
	Dia          string
	Hora         string
	Nota         float32
}

type User struct {
	Id     int
	Login  string
	Avatar []byte
}

func (l *Log) TelegramText() string {
	if l.Text == "" {
		return "No s'ha trobat cap log amb la cerca proporcionada"
	}
	return fmt.Sprintf("%s\n\n%s\n\n%s", l.Titol, l.Text, l.dateAndAuthor())
}

func (l *Log) dateAndAuthor() string {
	return fmt.Sprintf("Enviat el %s per %s", l.Dia, l.Autor.Login)
}
