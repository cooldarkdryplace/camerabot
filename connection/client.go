package connection

import "net/http"

type Client interface {
	Do(req *http.Request) (resp *http.Response, err error)
	Get(path string) (resp *http.Response, err error)
}
type HttpClient struct {
	Impl *http.Client
}

func (hc *HttpClient) Do(req *http.Request) (resp *http.Response, err error) {
	return hc.Impl.Do(req)
}

func (hc *HttpClient) Get(path string) (resp *http.Response, err error) {
	return hc.Impl.Get(path)
}
