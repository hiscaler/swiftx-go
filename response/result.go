package response

type Result struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
