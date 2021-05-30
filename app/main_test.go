package main

import (
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
    client := new(http.Client)
    testServer := httptest.NewServer(http.HandlerFunc(helloHandler))
    defer testServer.Close()

    req, _ := http.NewRequest("GET", testServer.URL+"/hello", nil)

    resp, _ := client.Do(req)
    respBody, _ := ioutil.ReadAll(resp.Body)

    assert.Equal(t, http.StatusOK, resp.StatusCode)
    assert.Equal(t, "Hello!", string(respBody))
}