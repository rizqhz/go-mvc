package helpers

func FormatResponse(message string, data any) map[string]any {
	res := make(map[string]any)
	res["message"] = message
	if data != nil {
		res["data"] = data
	}
	return res
}
