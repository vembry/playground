package http

type BaseResponse[C any] struct {
	Error  string `json:"error"`
	Object C      `json:"object"`
}
