package helper

func FormatResponse(msg string, data any) map[string]any{
	response := map[string]any{}
	response["message"] = msg
	if data != nil{
		response["data"] = data
	}
	return response
}