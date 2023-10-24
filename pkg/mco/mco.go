package mco

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"k8s.io/klog/v2"
)

const Port = 7070
const Host = "http://localhost"

func url(controller string) string {
	return fmt.Sprintf("%s:%d/%s", Host, Port, controller)
}

func Send(ctx context.Context, controller string, request interface{}, result interface{}) error {
	logger := klog.FromContext(ctx)

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return err
	}
	url := url(controller)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		logger.Info("failed building request")
		logger.Error(err, "building request")
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	logger.Info("sending request", "url", url, "body", requestJSON)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Info("failed sending request")
		logger.Error(err, "sending request")
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err, "reading body")
		return err
	}

	if resp.StatusCode >= 300 {
		logger.Info("failed request", "status", resp.StatusCode, "body", body)
		return errors.New("failed to send request")
	}

	if resp.StatusCode != http.StatusNoContent {
		err = json.Unmarshal(body, result)
		if err != nil {
			logger.Error(err, "failed parsing response")
			return err
		}
	}

	logger.Info("got response", "body_raw", body, "body", result)

	return nil
}
