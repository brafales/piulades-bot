package pinchito


type TgPinchitoUser struct {
	PinId int
	TgId  int
	Nick  string
}

var tgPinchitoUsers []TgPinchitoUser

func InitUsers() {
	tgPinchitoUsers = []TgPinchitoUser{
		TgPinchitoUser{1, 5774355, "[RiS]"},
		TgPinchitoUser{2, 5774356, "Rofi"},
	}

}