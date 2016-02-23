package cityhall

import (
	"encoding/json"
	"crypto/md5"
	"fmt"
)

func isResponseOkay(body []byte) error {
	type response struct {
		Response, Message string
	}
	var resp response
	if err := json.Unmarshal(body, &resp); err != nil {
		return err
	}

	if resp.Response != "Ok" {
		return cityhallError(resp.Message)
	}
	return nil
}

func getValueFromResponse(body []byte) (string, error) {
	type response struct {
		Response string
		Message string
		Value string
	}
	var ret response
	if err := json.Unmarshal(body, &ret); err != nil {
		return "", err
	}

	if ret.Response == "Ok" {
		return ret.Value, nil
	}

	if len(ret.Message) > 0 {
		return "", cityhallError(ret.Message)
	}

	return "", cityhallError(fmt.Sprintf("Response from server is incomplete: %s", body))
}

func hash(password string) string {
	if password == "" {
		return ""
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}
