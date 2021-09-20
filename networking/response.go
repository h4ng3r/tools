package networking

import (
    "net/http"
)

type Response struct {
    StatusCode    int64
    Headers       map[string][]string
    Data          []byte
    ContentLength int64
    Cancelled     bool
    Request       *Request
}

func NewResponse(httpresp *http.Response, req *Request) Response {
    var resp Response
    resp.Request = req
    resp.StatusCode = int64(httpresp.StatusCode)
    resp.Headers = httpresp.Header
    resp.Cancelled = false
    return resp
}