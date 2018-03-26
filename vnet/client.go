package vnet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/varunamachi/vaali/vlog"
)

//Client - client for orek service
type Client struct {
	http.Client
	Address    string
	VersionStr string
	Token      string
	BaseURL    string
}

//NewClient - creates a new rest client
func NewClient(address, versionStr string) *Client {
	return &Client{
		Client: http.Client{
			Timeout: 0,
		},
		Address:    address,
		VersionStr: versionStr,
		BaseURL:    fmt.Sprintf("%s/%s", address, versionStr),
	}
}

//Get - performs a get request
func (client *Client) Get(
	content interface{},
	urlArgs ...string) (err error) {

	var req *http.Request
	var resp *http.Response
	apiURL := client.getURL(urlArgs...)
	req, err = http.NewRequest("GET", apiURL, nil)
	authHeader := fmt.Sprintf("Bearer %s", client.Token)
	req.Header.Add("Authorization", authHeader)
	resp, err = client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		decoder := json.NewDecoder(resp.Body)
		if resp.StatusCode == http.StatusOK {
			err = decoder.Decode(content)
		} else {
			err = handleStatusCode(resp.StatusCode, decoder)
		}
	}
	return err
}

//Delete - performs a delete request
func (client *Client) Delete(
	urlArgs ...string) (err error) {
	var req *http.Request
	var resp *http.Response
	apiURL := client.getURL(urlArgs...)
	req, err = http.NewRequest("DELETE", apiURL, nil)
	authHeader := fmt.Sprintf("Bearer %s", client.Token)
	req.Header.Add("Authorization", authHeader)
	resp, err = client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		decoder := json.NewDecoder(resp.Body)
		err = handleStatusCode(resp.StatusCode, decoder)
	}
	return err
}

//Post - performs a post request
func (client *Client) Post(content interface{},
	urlArgs ...string) (err error) {
	return client.putOrPost("POST", content, urlArgs...)
}

//Put - performs a put request
func (client *Client) Put(content interface{},
	urlArgs ...string) (err error) {
	return client.putOrPost("PUT", content, urlArgs...)
}

func handleStatusCode(statusCode int, decoder *json.Decoder) (err error) {
	var res Result
	if statusCode == http.StatusOK {
		err = decoder.Decode(&res)
		vlog.Info("REST", "%s : %s", res.Op, res.Msg)
	} else if statusCode == http.StatusInternalServerError ||
		statusCode == http.StatusBadRequest ||
		statusCode == http.StatusUnauthorized {
		err = decoder.Decode(&res)
		if err == nil && len(res.Err) != 0 {
			vlog.Error("REST", "%s : %s - %s", res.Op, res.Msg, res.Err)
			err = errors.New(res.Err)
		} else if err != nil {
			vlog.Error("RESTClient", "Result decode failed: ", err)
		}
	} else {
		err = fmt.Errorf("Status Error: %d - %s", statusCode,
			http.StatusText(statusCode))
		// olog.PrintError("REST", err)
	}
	return err
}

func (client *Client) do(req *http.Request) (
	resp *http.Response, err error) {
	authHeader := fmt.Sprintf("Bearer %s", client.Token)
	req.Header.Add("Authorization", authHeader)
	resp, err = client.Do(req)
	return resp, err
}

func (client *Client) getURL(args ...string) (str string) {
	var buffer bytes.Buffer
	buffer.WriteString(client.BaseURL)
	buffer.WriteString("/in/")
	for i := 0; i < len(args); i++ {
		buffer.WriteString(args[i])
		if i < len(args)-1 {
			buffer.WriteString("/")
		}
	}
	str = buffer.String()
	return str
}

func (client *Client) putOrPost(
	method string,
	content interface{},
	urlArgs ...string) (err error) {

	var data []byte
	var resp *http.Response
	data, err = json.Marshal(content)
	apiURL := client.getURL(urlArgs...)
	req, err := http.NewRequest(method, apiURL, bytes.NewBuffer(data))
	authHeader := fmt.Sprintf("Bearer %s", client.Token)
	req.Header.Add("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		decoder := json.NewDecoder(resp.Body)
		err = handleStatusCode(resp.StatusCode, decoder)
	}
	return err
}
