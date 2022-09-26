package apilayer

type ErrorsData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ResultData struct {
	Result    float64    `json:"result"`
	Success   bool       `json:"success"`
	ErrorData ErrorsData `json:"error"`
}
