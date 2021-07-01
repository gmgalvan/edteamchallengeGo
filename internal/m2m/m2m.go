package m2m

import (
	"io"
	"log"
	"net/http"
)

type M2MCall interface {
	GET(url string) ([]byte, error)
}

type M2M struct {
	client *http.Client
	M2MCall
}

func NewM2M(c *http.Client) *M2M {
	return &M2M{
		client: c,
	}
}

func (m2m *M2M) GET(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error on get response from url %v", url)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error on get body from url %v", url)
		return nil, err
	}
	return body, nil
}
