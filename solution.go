package kata

func BoolToWord(word bool) string {
	result := ""
	switch word {
	case true:
		result = "Yes"

	case false:
		result = "No"

	}
	return result
}
