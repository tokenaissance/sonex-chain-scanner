package eth_client

import "net/http"

type transport struct {
	Transport http.RoundTripper
	OnSend    func(req *http.Request)
	OnFailure func(req *http.Request, err error)
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	if t.OnSend != nil {
		t.OnSend(req)
	}

	resp, err = transport.RoundTrip(req)

	if err != nil && t.OnFailure != nil {
		t.OnFailure(req, err)
	}

	return resp, err
}
