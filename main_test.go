package requests

import (
    "fmt"
    "testing"
)

func TestRequest_Get(t *testing.T) {
    r := Request{
        URL: "http://www.baidu.com",
    }
    resp, err := r.Get()
    if err != nil {
        t.Fatal(err)
    }
    fmt.Println(string(resp))
}
