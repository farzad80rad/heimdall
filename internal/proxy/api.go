package proxy

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Proxy interface {
	Proxy(ctx *gin.Context) error
	Ping(url string) bool
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
