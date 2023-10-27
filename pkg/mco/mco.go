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
		logger.Error(err, "failed marshalling request")
		return err
	}
	url := url(controller)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		logger.Error(err, "failed building request")
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	logger.Info("sending request", "url", url, "body", requestJSON)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err, "failed sending request")
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err, "reading body")
		return err
	}

	if resp.StatusCode >= 300 {
		err := errors.New("failed to send request")
		logger.Error(err, "failed request", "status", resp.StatusCode, "body", body)
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		err = json.Unmarshal(body, result)
		if err != nil {
			logger.Error(err, "failed parsing response", "body", body)
			return err
		}
	}

	logger.Info("got response", "body", result)

	return nil
}
