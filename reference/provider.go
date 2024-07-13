package reference

var smsProviders = map[string]bool{
	"Topolo": true,
	"Rond":   true,
	"Kildy":  true,
}

func IsSmsProvider(name string) bool {
	_, found := smsProviders[name]
	return found
}

var mmsProviders = map[string]bool{
	"Topolo": true,
	"Rond":   true,
	"Kildy":  true,
}

func IsMmsProvider(name string) bool {
	_, found := mmsProviders[name]
	return found
}

var vcdProviders = map[string]struct{}{
	"TransparentCalls": {},
	"E-Voice":          {},
	"JustPhone":        {},
}

func IsVcdProvider(name string) bool {
	_, found := vcdProviders[name]
	return found
}

var emailProviders = map[string]struct{}{
	"Gmail":      {},
	"Yahoo":      {},
	"Hotmail":    {},
	"MSN":        {},
	"Orange":     {},
	"Comcast":    {},
	"AOL":        {},
	"Live":       {},
	"RediffMail": {},
	"GMX":        {},
	"Protonmail": {},
	"Yandex":     {},
	"Mail.ru":    {},
}

func IsEmailProvider(name string) bool {
	_, found := emailProviders[name]
	return found
}
