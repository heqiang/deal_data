package util

import (
	"fmt"
	"github.com/avast/retry-go"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRetry(t *testing.T) {
	url := "http://example.com"
	var body []byte

	_ = retry.Do(
		func() error {
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			return nil
		},
	)

	fmt.Println(body)
}
