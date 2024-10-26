package checker

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type HTTPChecker struct {
	timeout time.Duration
}

func NewHTTPChecker(timeout time.Duration) *HTTPChecker {
	return &HTTPChecker{
		timeout: timeout,
	}
}

func (hc *HTTPChecker) Check(url string) (bool, error) {
	client := http.Client{
		Timeout: hc.timeout,
	}
	resp, err := client.Get(url)

	if err != nil {
		return false, fmt.Errorf("не удалось получить доступ к %s: %w", url, err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Ошибка при закрытии ответа от %s: %v", url, err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("сайт %s вернул статус %d", url, resp.StatusCode)
	}

	return true, nil
}
