package util

import "net/http"

//HTTPResponseInterceptor intercept http response writer
//need this to intercept http header written inside blevehttp and update echo context value
//eg. getting error code then update echo context response error code
type HTTPResponseInterceptor struct {
	http.ResponseWriter
	StatusCode int
}

//NewHTTPResponseInterceptor create new httpInterceptor
func NewHTTPResponseInterceptor(w http.ResponseWriter) *HTTPResponseInterceptor {
	return &HTTPResponseInterceptor{w, http.StatusOK}
}

//WriteHeader override response WriteHeader
func (i *HTTPResponseInterceptor) WriteHeader(code int) {
	i.StatusCode = code
	i.ResponseWriter.WriteHeader(code)
}
