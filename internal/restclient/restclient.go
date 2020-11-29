package restclient

import (
	"encoding/json"
	"errors"

	resty "github.com/go-resty/resty/v2"
)

func Get(url string, v interface{}) error {
	rc := resty.New()

	req := rc.R()

	r, err := req.Get(url)

	if err != nil {
		return err
	}

	if r.StatusCode()/100 != 2 {
		return errors.New("Non 2xx response")
	}

	if err = json.Unmarshal(r.Body(), v); err != nil {
		return err
	}
	return nil
}
