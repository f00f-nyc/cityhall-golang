package cityhall

import (
	"io/ioutil"
	"encoding/json"
	"io"
	"crypto/md5"
	"fmt"
)

func isResponseOkay(body io.ReadCloser) error {
	resp_bytes, err := ioutil.ReadAll(body)
	body.Close()
	if err != nil {
		return err
	}

	type response struct {
		Response, Message string
	}
	var resp response
	if err = json.Unmarshal(resp_bytes, &resp); err != nil {
		return err
	}

	if resp.Response != "Ok" {
		err = cityhallError(resp.Message)
	}
	return err
}

func getValueFromResponse(body io.ReadCloser) (string, error) {
	resp_bytes, err := ioutil.ReadAll(body)
	body.Close()
	if err != nil {
		return "", err
	}
	var request interface{}
	if err = json.Unmarshal(resp_bytes, &request); err != nil {
		return "", err
	}
	m := request.(map[string]interface{})

	response, response_ok := m["Response"]
	message, message_ok := m["Message"]
	value, value_ok := m["value"]

	if response_ok && response == "Ok" && value_ok {
		return value.(string), nil
	} else if response_ok && message_ok {
		return "", cityhallError(message.(string))
	}

	return "", cityhallError(fmt.Sprintf("Response from server is incomplete: %s", resp_bytes))
}

func hash(password string) string {
	if password == "" {
		return ""
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}
