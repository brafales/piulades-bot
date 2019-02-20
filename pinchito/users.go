package pinchito

type TgPinchitoUser struct {
	TgUsername string
	PinId      int
	PinNick    string
}

var tgPinchitoUsers []TgPinchitoUser

func init() {
	tgPinchitoUsers = []TgPinchitoUser{

		{"parap", 1, "parap"},
		{"rofirrim", 2, "rofi"},
		{"SiR_RiS", 36, "[RiS]"},
		{"frikjan", 4, "Freakhand"},
		{"The_Bell", 18, "The_Bell"},
		{"hexxx", 24, "hex"},
		{"ciddx", 19, "Cid"},
		{"javiKug", 30, "KuG"},
		{"mixtura", 8, "mixtura"},
		{"santiee", 13, "Santiee"},
		{"hydex86", 5, "HyDe"},
		{"n0seres", 41, "n0seres"},
		{"", 3, "CHeRNoBiL"},
		{"", 6, "KoLiFLoR"},
		{"", 7, "Sharek"},
		{"", 9, "^smith"},
		{"", 10, "Menxu"},
		{"", 11, "ReiVaX18"},
		{"", 12, "yorx"},
		{"", 14, "Enehy"},
		{"", 15, "dikFIB"},
		{"", 16, "gOBl1N"},
		{"", 17, "Scenix"},
		{"", 20, "{nimfa}"},
		{"", 21, "Kayl"},
		{"", 22, "graz"},
		{"", 23, "CrC"},
		{"", 25, "[TUmor]"},
		{"", 26, "basted"},
		{"", 27, "SLaKeLS"},
		{"", 28, "AC_DC"},
		{"", 29, "STeArN"},
		{"", 31, "Yak-3"},
		{"", 34, "NeoRob"},
		{"", 32, "{xtreme}"},
		{"", 33, "^KreatoR"},
		{"", 35, "qiz"},
		{"", 37, "{_SaLeM_"},
		{"", 38, "FCS"},
		{"", 39, "DeNeA"},
	}
}
