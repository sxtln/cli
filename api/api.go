package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
	"sextillion.io/cli/models"
)

type Api struct {
	authorization string
}

func NewApi(appConfig models.Config) *Api {
	ret := Api{
		authorization: "",
	}
	if appConfig.ApiKey != "" {
		ret.authorization = appConfig.ApiKey
	}
	if appConfig.Token != "" {
		ret.authorization = appConfig.Token
	}
	return &ret
}

func getClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    60 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{},
	}

	client := &http.Client{
		Transport: tr,
	}

	return client

}

func (t *Api) endpoint2Url(ep string) string {
	url, err := url.JoinPath("https://api.sextillion.io", ep)
	if err != nil {
		log.Fatal(err)
	}
	return url
}

func (t *Api) handleResponse(url string, hResp *http.Response) (int, models.JSON, error) {

	defer hResp.Body.Close()
	resBody, err := io.ReadAll(hResp.Body)
	if err != nil {
		return 0, nil, err
	}

	if hResp.StatusCode != 200 {
		return hResp.StatusCode, nil, fmt.Errorf(`failed to call %s %s. status code: %d body: %s`, hResp.Request.Method, url, hResp.StatusCode, string(resBody))
	}
	var json models.JSON

	var contentType = hResp.Header.Get(`Content-Type`)
	if strings.Contains(contentType, `json`) || strings.Contains(contentType, `yaml`) || strings.Contains(contentType, `yml`) {
		if len(resBody) > 0 {
			json, err = t.unmarshalHttpBody(resBody)
			if err != nil {
				return 0, nil, err
			}
		}
	} else {
		json = make(models.JSON)
		json[`body`] = resBody
	}

	return hResp.StatusCode, json, nil
}

func (t *Api) Get(endpoint string) (int, models.JSON, error) {

	var url = t.endpoint2Url(endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, err
	}

	if t.authorization != "" {
		req.Header.Add("Authorization", fmt.Sprintf("bearer %s", t.authorization))
	}
	req.Header.Add("Content-Type", "application/json")

	client := getClient()

	hResp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	return t.handleResponse(url, hResp)
}

func (t *Api) Delete(endpoint string) (int, models.JSON, error) {

	var url = t.endpoint2Url(endpoint)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return 0, nil, err
	}

	if t.authorization != "" {
		req.Header.Add("Authorization", fmt.Sprintf("bearer %s", t.authorization))
	}
	req.Header.Add("Content-Type", "application/json")

	client := getClient()

	hResp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	return t.handleResponse(url, hResp)
}

func (t *Api) Post(endpoint string, body models.JSON) (int, models.JSON, error) {
	var reqBody io.Reader

	if body != nil {
		json, err := json.Marshal(body)
		if err != nil {
			return 0, nil, err
		}
		reqBody = bytes.NewBuffer(json)
	}

	var url = t.endpoint2Url(endpoint)

	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return 0, nil, err
	}

	if t.authorization != "" {
		req.Header.Add("Authorization", fmt.Sprintf("bearer %s", t.authorization))
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if len(body) > 0 {
		req.Header.Set("Content-Length", fmt.Sprint(len(body)))
	}

	client := getClient()

	hResp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	return t.handleResponse(url, hResp)

}

func (t *Api) unmarshalHttpBody(rawJson []byte) (models.JSON, error) {
	var ret map[string]interface{}

	jsonErr := json.Unmarshal(rawJson, &ret)
	if jsonErr != nil {
		yamlErr := yaml.Unmarshal(rawJson, &ret)
		if yamlErr != nil {
			return nil, jsonErr
		}
	}
	return ret, nil
}
