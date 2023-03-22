package hydra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/p4gefau1t/trojan-go/log"
	"io"
	"net/http"
)

type RequestClient struct {
	ctx          context.Context
	baseUrl      string
	accessToken  string
	refreshToken string
	username     string
	password     string
}

func (client *RequestClient) buildUrl(resource string) string {
	url := fmt.Sprintf("%s%s", client.baseUrl, resource)
	return url
}

func (client *RequestClient) RequestSync(
	method string,
	path string,
	params map[string]interface{},
	data map[string]interface{},
	withToken bool) map[string]interface{} {
	url := client.buildUrl(path)
	reqBodyBytes, _ := json.Marshal(data)
	req, _ := http.NewRequest(method, url, bytes.NewReader(reqBodyBytes))

	if method != http.MethodGet {
		req.Header.Set("content-type", "application/json;charset=utf-8")
	}

	if withToken {
		req.Header.Set("x-token", client.accessToken)
	}

	sender := http.Client{}
	resp, err := sender.Do(req)
	if err != nil {
		fmt.Println("HTTP call failed:", err)
		return nil
	}
	respBytes, _ := io.ReadAll(resp.Body)
	log.Debugf("response %s %s from %s", http.MethodPost, url, respBytes)
	var result map[string]interface{}
	json.Unmarshal(respBytes, &result)
	return result
}

func (client *RequestClient) RefreshToken() {
	data := make(map[string]interface{})
	data["refreshToken"] = client.refreshToken
	var result = client.RequestSync(http.MethodPost, "/auth/refreshToken", nil, data, true)
	if result == nil {
		return
	}
	client.accessToken = result["accessToken"].(string)
	client.refreshToken = result["refreshToken"].(string)
	log.Info("hydra token refreshed successfully.")
}

func (client *RequestClient) Login() {
	data := make(map[string]interface{})
	data["username"] = client.username
	data["password"] = client.password
	var result = client.RequestSync(http.MethodPost, "/auth/login", nil, data, false)
	if result == nil {
		return
	}
	client.accessToken = result["accessToken"].(string)
	client.refreshToken = result["refreshToken"].(string)
	log.Info("hydra authenticated successfully.")
}

func (client *RequestClient) UpdateTraffic(hash string, sent uint64, recv uint64) map[string]interface{} {
	data := make(map[string]interface{})
	data["hash"] = hash
	data["sent"] = sent
	data["recv"] = recv
	var result = client.RequestSync(http.MethodPost, "/api/subscriptions/updateTraffic", nil, data, true)
	if result == nil {
		return nil
	}
	log.Debugf("%s", result)
	return result
}

func (client *RequestClient) PullSubscriptions() map[string]interface{} {
	var result = client.RequestSync(http.MethodGet, "/api/subscriptions", nil, nil, true)
	if result == nil {
		return nil
	}
	log.Debugf("%s", result)
	return result
}

func NewRequestClient(ctx context.Context, baseUrl string, username string, password string) (*RequestClient, error) {
	client := &RequestClient{
		ctx:      ctx,
		baseUrl:  baseUrl,
		username: username,
		password: password,
	}
	return client, nil
}
