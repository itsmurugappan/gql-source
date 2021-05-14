package adapter

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"knative.dev/pkg/logging"
)

const (
	userSecret                   = "GQL_OIDC_USER"
	passwordSecret               = "GQL_OIDC_PASSWORD"
	clientIDSecret               = "GQL_OIDC_CLIENT"
	tokenUrl                     = "GQL_AUTH_TOKEN_URL"
	timeOut        time.Duration = 5 * time.Second
)

type authClient interface {
	setAuthHeader(ctx context.Context, header http.Header)
}

func newAuthClient(ctx context.Context) authClient {
	switch {
	case os.Getenv("GQL_OIDC_AUTH_ENABLE") == "true":
		return newKCClient(ctx)
		// implement other auth methods here
	}
	return nil
}

type kcClient struct {
	client   *http.Client
	authData url.Values
	tokenUrl string
}

func newKCClient(ctx context.Context) *kcClient {
	c := &http.Client{
		Timeout: timeOut,
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		},
	}
	data := url.Values{}
	data.Set("username", getSecret(ctx, userSecret))
	data.Set("password", getSecret(ctx, passwordSecret))
	data.Set("grant_type", "password")
	data.Set("client_id", getSecret(ctx, clientIDSecret))
	return &kcClient{
		c,
		data,
		getSecret(ctx, tokenUrl),
	}
}

// fetchToken fetches token from keycloak endpoint
func (c *kcClient) setAuthHeader(ctx context.Context, header http.Header) {
	logger := logging.FromContext(ctx)
	res, err := c.client.Do(c.constructTokenRequest(ctx))
	if err != nil {
		logger.Error("error getting token from keycloak", zap.Error(err))
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("error getting body from response", zap.Error(err))
		return
	}
	var token map[string]interface{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		logger.Error("fail to parse response body to get token", zap.Error(err))
		return
	}
	if token["access_token"] != nil {
		header.Add("Authorization", fmt.Sprintf("Bearer %s", token["access_token"].(string)))
	}
	return
}

func (c *kcClient) constructTokenRequest(ctx context.Context) *http.Request {
	r, err := http.NewRequest("POST", c.tokenUrl, strings.NewReader(c.authData.Encode()))
	if err != nil {
		logging.FromContext(ctx).Error("fail to construct request to get token", zap.Error(err))
	}
	r.Close = true
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(c.authData.Encode())))

	return r
}

func getSecret(ctx context.Context, envName string) string {
	secret, present := os.LookupEnv(envName)
	if !present {
		logging.FromContext(ctx).Fatal("Unable to retrieve auth creds")
	}
	return secret
}
