package forward

import (
	"bytes"
	"log"
	"net/http"
)

type Forwarder struct {
	name     string
	endpoint string
}

func (f *Forwarder) Write(p []byte) (n int, err error) {

	buf := bytes.NewBuffer(p)

	resp, err := http.Post(f.endpoint, "text/plain", buf)
	if err != nil {
		log.Printf("error on %s forward: %v", f.name, err)
		return 0, err
	}
	log.Printf("forward %s -- %s", f.name, resp.Status)
	return len(p), nil

}

func NewForwarder(name, endpoint string) *Forwarder {
	return &Forwarder{
		name:     name,
		endpoint: endpoint,
	}
}
