package networking

type Request struct {
    Method   string
    Host     string
    Url      string
    Headers  map[string]string
    Data     []byte
    Input    map[string][]byte
    Position int
    IgnoreBody bool
}

func NewRequest(method, url string, ignoreBody bool) *Request {
    var req Request
    req.Method = method
    req.Url = url
    req.Headers = make(map[string]string)
    req.IgnoreBody = ignoreBody
    return &req
}