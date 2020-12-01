package requests

import (
    "fmt"
    "io/ioutil"
    "testing"
)

func TestRequest_Get(t *testing.T) {
    r := Request{
        URL:     "http://www.baidu.com",
    }
    resp, err := r.Get()
    if err != nil {
        t.Fatal(err)
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        t.Fatal(err)
    }
    fmt.Println(string(body))
}
