package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func ReadJSONBody(response *http.Response, dto interface{}) (err error) {
	if response.Header.Get("Content-Type") != "application/json" {
		return errors.New("Content-Type != application/json")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("cant read request.Body: %w", err)
	}
	defer func() {
		errdefer := response.Body.Close()
		if errdefer != nil {
			err = errdefer
		}
	}()

	err = json.Unmarshal(body, &dto)
	if err != nil {
		return err
	}
	return nil
}

func WriteJSONBody(response http.ResponseWriter, dto interface{}) (err error) {
	response.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	_, err = response.Write(body)
	if err != nil {
		return err
	}

	return nil
}
