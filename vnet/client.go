package vnet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/varunamachi/vaali/vsec"

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
func NewClient(address, appName, versionStr string) *Client {
	return &Client{
		Client: http.Client{
			Timeout: 1 * time.Minute,
		},
		Address:    address,
		VersionStr: versionStr,
		BaseURL: fmt.Sprintf("%s/%s/api/%s",
			address,
			appName,
			versionStr),
	}
}

//Get - performs a get request
func (client *Client) Get(
	content interface{},
	access vsec.AuthLevel,
	urlArgs ...string) (err error) {
	var req *http.Request
	var resp *http.Response
	apiURL := client.CreateURL(access, urlArgs...)
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
	access vsec.AuthLevel,
	urlArgs ...string) (err error) {
	var req *http.Request
	var resp *http.Response
	apiURL := client.CreateURL(access, urlArgs...)
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
func (client *Client) Post(
	content interface{},
	access vsec.AuthLevel,
	urlArgs ...string) (err error) {
	var res Result
	return client.putOrPost("POST", access, content, &res, urlArgs...)
}

//Put - performs a put request
func (client *Client) Put(
	content interface{},
	access vsec.AuthLevel,
	urlArgs ...string) (err error) {
	var res Result
	return client.putOrPost("PUT", access, content, &res, urlArgs...)
}

//CreateURL - constructs URL from base URL, access level and the given
//path components
func (client *Client) CreateURL(
	access vsec.AuthLevel,
	args ...string) (str string) {
	var buffer bytes.Buffer
	// buffer.WriteString("/in/")
	accessStr := ""
	switch access {
	case vsec.Super:
		accessStr = "in/r0/"
	case vsec.Admin:
		accessStr = "in/r1/"
	case vsec.Normal:
		accessStr = "in/r2/"
	case vsec.Monitor:
		accessStr = "in/r3/"
	case vsec.Public:
		accessStr = ""
	}
	buffer.WriteString(client.BaseURL)
	buffer.WriteString(accessStr)
	for i := 0; i < len(args); i++ {
		buffer.WriteString(args[i])
		if i < len(args)-1 {
			buffer.WriteString("/")
		}
	}
	str = buffer.String()
	return str
}

func (client *Client) Login(userID, password string) (err error) {
	data := make(map[string]string)
	data["userID"] = userID
	data["password"] = password
	err = client.Post(data, vsec.Public, "login")

	var req *http.Request

	// req, err = http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	// if err == nil {
	// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// 	var resp *http.Response
	// 	resp, err = client.Do(req)
	// 	if err == nil {
	// 		defer resp.Body.Close()
	// 		decoder := json.NewDecoder(resp.Body)
	// 		if resp.StatusCode == http.StatusOK {
	// 			tmap := make(map[string]string)
	// 			err = decoder.Decode(&tmap)
	// 			client.Token = tmap["token"]
	// 		} else {
	// 			err = handleStatusCode(resp.StatusCode, decoder)
	// 		}
	// 	}
	// }
	// if err != nil {
	// 	olog.PrintError("RESTClient", err)
	// }
	return err
}

func handleStatusCode(statusCode int, decoder *json.Decoder) (err error) {
	var res Result
	// if statusCode == http.StatusOK {
	// 	err = decoder.Decode(&res)
	// 	vlog.Info("REST", "%s : %s", res.Op, res.Msg)
	// } else
	if statusCode == http.StatusInternalServerError ||
		statusCode == http.StatusBadRequest ||
		statusCode == http.StatusUnauthorized {
		err = decoder.Decode(&res)
		if err == nil && len(res.Err) != 0 {
			vlog.Error("REST", "%s : %s - %s", res.Op, res.Msg, res.Err)
			err = errors.New(res.Err)
		} else if err != nil {
			vlog.Error("RESTClient", "Result decode failed: ", err)
		}
	}
	// else {
	// 	err = fmt.Errorf("Status Error: %d - %s", statusCode,
	// 		http.StatusText(statusCode))
	// 	// olog.PrintError("REST", err)
	// }
	return err
}

func (client *Client) do(req *http.Request) (
	resp *http.Response, err error) {
	authHeader := fmt.Sprintf("Bearer %s", client.Token)
	req.Header.Add("Authorization", authHeader)
	resp, err = client.Do(req)
	return resp, err
}

func (client *Client) putOrPost(
	method string,
	access vsec.AuthLevel,
	content interface{},
	resultOut interface{},
	urlArgs ...string) (err error) {
	var data []byte
	var resp *http.Response
	data, err = json.Marshal(content)
	apiURL := client.CreateURL(access, urlArgs...)
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
