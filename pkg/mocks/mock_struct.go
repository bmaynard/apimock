package mocks

type Meta struct {
	StatusCode    int    `json:"status_code"`
	RequestPath   string `json:"request_path"`
	RequestMethod string `json:"method"`
}

type MockResponse struct {
	Response interface{} `json:"response"`
	Meta     Meta        `json:"meta"`
}
