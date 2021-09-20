package networking

import (
    "bytes"
    "context"
    "net/http"
    "crypto/tls"
    "time"
    "net"
    "strconv"
    "io/ioutil"
    "unicode/utf8"
)

//Download results < 5MB
const MAX_DOWNLOAD_SIZE = 5242880

type Client struct {
    ctx context.Context
    client *http.Client
}

func NewClient(ctx context.Context, timeout int, followRedirect bool) *Client {
    client := &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse
        },
        Timeout: time.Duration(time.Duration(timeout) * time.Second), // conf.Timeout
        Transport: &http.Transport{
            MaxIdleConns:        1000,
            MaxIdleConnsPerHost: 100,
            MaxConnsPerHost:     100,
            DialContext: (&net.Dialer{
                Timeout: time.Duration(time.Duration(timeout) * time.Second), // conf.Timeout
            }).DialContext,
            TLSHandshakeTimeout: time.Duration(time.Duration(timeout) * time.Second), // conf.Timeout
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true,
                Renegotiation:      tls.RenegotiateOnceAsClient,
            },
        }}
    if followRedirect {
        client.CheckRedirect = nil
    }
    return &Client{ctx:ctx, client:client}
}

func (c *Client) Execute(req *Request) (Response, error) {
    var httpreq *http.Request
    var err error

    data := bytes.NewReader(req.Data)
    httpreq, err = http.NewRequestWithContext(c.ctx, req.Method, req.Url, data)
    if err != nil {
        return Response{}, err
    }

    // set default User-Agent header if not present
    if _, ok := req.Headers["User-Agent"]; !ok {
        req.Headers["User-Agent"] = "Mozilla/5.0 (iPhone; CPU iPhone OS 12_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko)"
    }

    // Handle Go http.Request special cases
    if _, ok := req.Headers["Host"]; ok {
        httpreq.Host = req.Headers["Host"]
    }

    req.Host = httpreq.Host
    httpreq = httpreq.WithContext(c.ctx)
    for k, v := range req.Headers {
        httpreq.Header.Set(k, v)
    }

    httpresp, err := c.client.Do(httpreq)

    // TODO: handle retry
    if err != nil {
        return Response{}, err
    }

    resp := NewResponse(httpresp, req)
    defer httpresp.Body.Close()

    // Check if we should download the resource or not
    size, err := strconv.Atoi(httpresp.Header.Get("Content-Length"))
    if err == nil {
        resp.ContentLength = int64(size)
    }

    if (req.IgnoreBody) || (size > MAX_DOWNLOAD_SIZE) {
        resp.Cancelled = true
        return resp, nil
    }
    if respbody, err := ioutil.ReadAll(httpresp.Body); err == nil {
        resp.ContentLength = int64(utf8.RuneCountInString(string(respbody)))
        resp.Data = respbody
    }

    return resp, nil
}