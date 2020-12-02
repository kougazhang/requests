package requests

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "time"
)

type Request struct {
    http.Request
    URL   string
    Retry *Retry
}

type Retry struct {
    Times    int
    Interval time.Duration
}

func (r Retry) do(fun func(method, url string, body io.Reader) ([]byte, error), method, url string, body io.Reader) (
    res []byte, err error) {

    for i := 0; i < r.Times; i++ {
        res, err = fun(method, url, body)
        if err == nil {
            break
        }
        time.Sleep(r.Interval)
    }
    return res, err
}

func (r *Request) AddQuery(queries map[string]string) error {
    u, err := url.Parse(r.URL)
    if err != nil {
        return err
    }

    values := u.Query()

    for k, v := range queries {
        values.Add(k, v)
    }

    r.URL = u.String()

    return nil
}

func (r *Request) AddHeader(headers map[string]string) {
    if r.Header == nil {
        r.Header = make(http.Header)
    }

    for k, v := range headers {
        r.Header.Add(k, v)
    }
}

func (r Request) do(method, url string, body io.Reader) ([]byte, error) {
    req, err := http.NewRequest(method, url, body)
    if err != nil {
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer req.Body.Close()

    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode/100 != 2 {
        return nil, fmt.Errorf("StatusCode %d, respBody %s", resp.StatusCode, string(respBody))
    }

    return respBody, nil
}

func (r Request) Do(method, url string, body io.Reader) ([]byte, error) {
    if r.Retry != nil {
        return r.Retry.do(r.do, method, url, body)
    }
    return r.do(method, url, body)
}

func (r Request) Get() ([]byte, error) {
    return r.Do(http.MethodGet, r.URL, nil)
}

func (r Request) Post(body io.Reader) ([]byte, error) {
    return r.Do(http.MethodPost, r.URL, body)
}

func (r Request) PostJson(v interface{}) ([]byte, error) {
    payload, err := json.Marshal(v)
    if err != nil {
        return nil, err
    }

    r.AddHeader(map[string]string{
        "Content-Type": "application/json",
    })

    return r.Post(bytes.NewReader(payload))
}
