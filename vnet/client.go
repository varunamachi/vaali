package vnet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/varunamachi/vaali/vsec"
)

//ResultReader - proto result, use it to read the Result with proper data struct
type ResultReader struct {
	RawData []byte
	Err     error
	Res     Result
}

//NewResultReader - creates a new result-reader from response body
func NewResultReader(r *http.Response) (reader *ResultReader) {
	reader = &ResultReader{}
	loc, _ := r.Location()
	url := "N/A"
	if loc != nil {
		url = loc.String()
	}
	if r.StatusCode == http.StatusNotFound {
		reader.Err = fmt.Errorf("Not found. URL: %s", url)
	} else if r.StatusCode == http.StatusUnauthorized {
		reader.Err = fmt.Errorf("Not logged in. URL: %s", url)
	} else if r.StatusCode == http.StatusForbidden {
		reader.Err = fmt.Errorf("Unauthorized access. URL: %s", url)
	} else {
		defer r.Body.Close()
		reader.RawData, reader.Err = ioutil.ReadAll(r.Body)
	}
	return reader
}

//Read - read result data from reader. The provided data object will be populatd
//with result's data field
func (rr *ResultReader) Read(data interface{}) (err error) {
	if rr.Err == nil {
		rr.Res = Result{
			Data: data,
		}
		err = json.Unmarshal(rr.RawData, &rr.Res)
	} else {
		err = rr.Err
	}
	if err == nil {
		if rr.Res.Status == http.StatusInternalServerError ||
			rr.Res.Status == http.StatusBadRequest {
			err = fmt.Errorf("%s : %s - %s",
				rr.Res.Op,
				rr.Res.Msg,
				rr.Res.Err)
		}
	}
	return err
}

//Finish - decodes the server response and returns error if it failed. Use this
//method if data is not expected from server call
func (rr *ResultReader) Finish() (err error) {
	if rr.Err == nil {
		rr.Res = Result{}
		err = json.Unmarshal(rr.RawData, &rr.Res)
	} else {
		err = rr.Err
	}
	if err == nil {
		if rr.Res.Status == http.StatusNotFound ||
			rr.Res.Status == http.StatusInternalServerError ||
			rr.Res.Status == http.StatusBadRequest ||
			rr.Res.Status == http.StatusUnauthorized {
			err = fmt.Errorf("%s : %s - %s",
				rr.Res.Op,
				rr.Res.Msg,
				rr.Res.Err)
		}
	}
	return err
}

//Client - client for orek service
type Client struct {
	http.Client
	Address    string
	VersionStr string
	BaseURL    string
	Token      string
	User       *vsec.User
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
	access vsec.AuthLevel,
	urlArgs ...string) (rr *ResultReader) {
	var req *http.Request
	var resp *http.Response
	var err error
	apiURL := client.CreateURL(access, urlArgs...)
	req, err = http.NewRequest("GET", apiURL, nil)
	authHeader := fmt.Sprintf("Bearer %s", client.Token)
	req.Header.Add("Authorization", authHeader)
	resp, err = client.Do(req)
	if err == nil {
		rr = NewResultReader(resp)
	} else {
		rr = &ResultReader{
			Err: err,
		}
	}
	return rr
}

//Delete - performs a delete request
func (client *Client) Delete(
	access vsec.AuthLevel,
	urlArgs ...string) (rr *ResultReader) {
	var req *http.Request
	var resp *http.Response
	var err error
	apiURL := client.CreateURL(access, urlArgs...)
	req, err = http.NewRequest("DELETE", apiURL, nil)
	authHeader := fmt.Sprintf("Bearer %s", client.Token)
	req.Header.Add("Authorization", authHeader)
	resp, err = client.Do(req)
	if err == nil {
		rr = NewResultReader(resp)
	} else {
		rr = &ResultReader{
			Err: err,
		}
	}
	return rr
}

//Post - performs a post request
func (client *Client) Post(
	content interface{},
	access vsec.AuthLevel,
	urlArgs ...string) (rr *ResultReader) {
	return client.putOrPost("POST", access, content, urlArgs...)
}

//Put - performs a put request
func (client *Client) Put(
	content interface{},
	access vsec.AuthLevel,
	urlArgs ...string) (rr *ResultReader) {
	return client.putOrPost("PUT", access, content, urlArgs...)
}

//CreateURL - constructs URL from base URL, access level and the given
//path components
func (client *Client) CreateURL(
	access vsec.AuthLevel,
	args ...string) (str string) {
	var buffer bytes.Buffer
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
	buffer.WriteString("/")
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

//Login - login to a vaali based service with userID and password. If successful
//client will have the session information and can perform REST calls that needs
//authentication
func (client *Client) Login(userID, password string) (err error) {
	data := make(map[string]string)
	data["userID"] = userID
	data["password"] = password
	loginResult := struct {
		Token string     `json:"token"`
		User  *vsec.User `json:"user"`
	}{}
	rr := client.Post(data, vsec.Public, "login")
	err = rr.Read(&loginResult)
	if err == nil {
		client.Token = loginResult.Token
		client.User = loginResult.User
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

func (client *Client) putOrPost(
	method string,
	access vsec.AuthLevel,
	content interface{},
	urlArgs ...string) (rr *ResultReader) {
	var data []byte
	var resp *http.Response
	var err error
	data, err = json.Marshal(content)
	apiURL := client.CreateURL(access, urlArgs...)
	req, err := http.NewRequest(method, apiURL, bytes.NewBuffer(data))
	authHeader := fmt.Sprintf("Bearer %s", client.Token)
	req.Header.Add("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err == nil {
		rr = NewResultReader(resp)
	} else {
		rr = &ResultReader{
			Err: err,
		}
	}
	return rr
}
