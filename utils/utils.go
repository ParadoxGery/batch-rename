package utils

func Deleters(text string) string {
	res := ""

	for range text {
		res += "\b"
	}

	return res
}

func Pad(text string, l int) string {
	res := text
	for i := len(text); i < l; i++ {
		res += " "
	}

	return res
}
