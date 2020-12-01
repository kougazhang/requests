package requests

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
)

type Request struct {
    http.Request
    URL string
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

func (r Request) Get() (*http.Response, error) {
    req, err := http.NewRequest(http.MethodGet, r.URL, nil)
    if err != nil {
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode/100 != 2 {
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return nil, err
        }
        return nil, fmt.Errorf("StatusCode %d, respBody %s", resp.StatusCode, string(body))
    }

    return resp, nil
}
