package pinchito

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

func (l *Log) PrettyText() string {
	return l.Text
}
