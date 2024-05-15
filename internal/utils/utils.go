package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"heimdall/internal/config"
	heimdallErrors "heimdall/internal/errors"
	"net/http"
	"reflect"
	"time"
)

func ValidateBody(method string, body []byte, config []config.HostMatchInfo) error {
	var requestBodyMap map[string]interface{}
	err := json.Unmarshal(body, &requestBodyMap)
	if err != nil {
		return errors.Join(heimdallErrors.BadRequest, errors.New("not in json format"))
	}
	for _, methodInfo := range config {
		if methodInfo.SupportedType == method {
			for _, field := range methodInfo.RequestBodyCheckConfig.MandatoryFields {
				v, found := requestBodyMap[field.Name]
				if !(found && reflect.TypeOf(v).Kind().String() == field.Type) {
					err = errors.Join(heimdallErrors.BadRequest, fmt.Errorf("missing required field %s by type of %s ", field.Name, field.Type))
					return err
				}
			}

			return nil
		}
	}
	return nil
}

func DoHTTPGetRequest(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{
		Timeout: time.Minute * 1,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code is not OK: %v", resp.StatusCode)
	}

	return nil
}
