package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aclgo/grpc-mail/internal/models"
	"github.com/aclgo/grpc-mail/pkg/logger"
)

type e2eHTTPTest struct {
}

func Newe2eTest() *e2eHTTPTest {
	return &e2eHTTPTest{}
}

func RunHTTP(addrServer string, logger logger.Logger) {
	params := models.NewMailBody(
		"from",
		"to",
		"subject",
		"body",
		"template",
		"service_name",
	)

	js, err := json.Marshal(params)
	if err != nil {
		logger.Errorf("Run.Marshal: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", addrServer, bytes.NewReader(js))
	if err != nil {
		logger.Errorf("Run.NewRequestWithContext: %v", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorf("Run.Default.Client.Do: %v", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusOK {
		logger.Errorf("ERROR e2e test http endpoint")
		return
	}

	logger.Info("TEST e2e HTTP PASS")

}
