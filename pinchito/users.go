package pinchito

type TgPinchitoUser struct {
	TgUsername string
	PinId      int
	PinNick    string
}

var tgPinchitoUsers []TgPinchitoUser

func init() {
	tgPinchitoUsers = []TgPinchitoUser{

		TgPinchitoUser{"PaRaP", 1, "PaRaP"},
		TgPinchitoUser{"rofirrim", 2, "rofi"},
		TgPinchitoUser{"SiR_RiS", 36, "[RiS]"},
		TgPinchitoUser{"frikjan", 4, "Freakhand"},
		TgPinchitoUser{"The_Bell", 18, "The_Bell"},
		TgPinchitoUser{"hexxx", 24, "hex"},
		TgPinchitoUser{"ciddx", 19, "Cid"},
		TgPinchitoUser{"javiKug", 30, "KuG"},
		TgPinchitoUser{"mixtura", 8, "mixtura"},
		TgPinchitoUser{"santiee", 13, "Santiee"},
		TgPinchitoUser{"hydex86", 5, "HyDe"},
		TgPinchitoUser{"", 3, "CHeRNoBiL"},
		TgPinchitoUser{"", 6, "KoLiFLoR"},
		TgPinchitoUser{"", 7, "Sharek"},
		TgPinchitoUser{"", 9, "^smith"},
		TgPinchitoUser{"", 10, "Menxu"},
		TgPinchitoUser{"", 11, "ReiVaX18"},
		TgPinchitoUser{"", 12, "yorx"},
		TgPinchitoUser{"", 14, "Enehy"},
		TgPinchitoUser{"", 15, "dikFIB"},
		TgPinchitoUser{"", 16, "gOBl1N"},
		TgPinchitoUser{"", 17, "Scenix"},
		TgPinchitoUser{"", 20, "{nimfa}"},
		TgPinchitoUser{"", 21, "Kayl"},
		TgPinchitoUser{"", 22, "graz"},
		TgPinchitoUser{"", 23, "CrC"},
		TgPinchitoUser{"", 25, "[TUmor]"},
		TgPinchitoUser{"", 26, "basted"},
		TgPinchitoUser{"", 27, "SLaKeLS"},
		TgPinchitoUser{"", 28, "AC_DC"},
		TgPinchitoUser{"", 29, "STeArN"},
		TgPinchitoUser{"", 31, "Yak-3"},
		TgPinchitoUser{"", 34, "NeoRob"},
		TgPinchitoUser{"", 32, "{xtreme}"},
		TgPinchitoUser{"", 33, "^KreatoR"},
		TgPinchitoUser{"", 35, "qiz"},
		TgPinchitoUser{"", 37, "{_SaLeM_"},
		TgPinchitoUser{"", 38, "FCS"},
		TgPinchitoUser{"", 39, "DeNeA"},
	}

}
