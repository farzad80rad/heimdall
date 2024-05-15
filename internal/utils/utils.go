package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"heimdall/internal/config"
	heimdallErrors "heimdall/internal/errors"
	"io"
	"net/http"
	"reflect"
	"time"
)

func ValidateBody(c *gin.Context, hostsInfo []config.HostMatchInfo) error {
	if c.Request.Body != nil {
		for _, b := range hostsInfo {
			if b.SupportedType == c.Request.Method {
				if b.RequestBodyCheckConfig != nil {
					body, _ := io.ReadAll(c.Request.Body)
					c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
					err := validateBody(body, b)
					if err != nil {
						return err
					}
				}
				return nil
			}
		}
	}
	return nil
}

func validateBody(body []byte, config config.HostMatchInfo) error {
	var requestBodyMap map[string]interface{}
	err := json.Unmarshal(body, &requestBodyMap)
	if err != nil {
		return errors.Join(heimdallErrors.BadRequest, errors.New("not in json format"))
	}
	for _, field := range config.RequestBodyCheckConfig.MandatoryFields {
		v, found := requestBodyMap[field.Name]
		if !(found && reflect.TypeOf(v).Kind().String() == field.Type) {
			err = errors.Join(heimdallErrors.BadRequest, fmt.Errorf("missing required field %s by type of %s ", field.Name, field.Type))
			return err
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
