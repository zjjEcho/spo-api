package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

func init() {
	saltCookieJar, _ = cookiejar.New(nil)
}

var (
	saltCookieJar *cookiejar.Jar
)

type Response struct {
	Code int
	Data interface{}
}

type SaltClient struct {
	Address    string
	UserName   string
	Passwd     string
	HttpClient *http.Client
}

func NewSaltClient(a, u, p string) *SaltClient {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}

	return &SaltClient{
		Address:    a,
		UserName:   u,
		Passwd:     p,
		HttpClient: &http.Client{Jar: saltCookieJar, Transport: tr},
	}
}

func (s *SaltClient) DoRequest(url, method string, buf *bytes.Buffer) (*Response, error) {
	req, err := http.NewRequest(method, "https://"+s.Address+"/"+url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := new(Response)
	ret.Code = resp.StatusCode
	if len(data) > 0 {
		err = json.Unmarshal(data, &ret.Data)
		if err != nil {
			ret.Data = string(data)
		}
	}

	return ret, nil
}

func (s *SaltClient) Auth() (*Response, error) {
	data := map[string]string{
		"username": s.UserName,
		"password": s.Passwd,
		"eauth":    "pam",
	}

	jdata, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Write(jdata)
	resp, err := s.DoRequest("/login", "POST", &buf)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

/*type paramIP struct {
	Client    string   `json:"client"`
	Tgt       string   `json:"tgt"`
	Expr_form string   `json:"expr_form"`
	Fun       string   `json:"fun"`
	Arg       []string `json:"arg"`
}*/

func (s *SaltClient) GetIP(servers string) (*Response, error) {
	data := map[string]string{
		"client":    "local",
		"tgt":       servers,
		"expr_form": "list",
		"fun":       "grains.item",
		"arg":       "ipv4",
	}

	jdata, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Write(jdata)
	resp, err := s.DoRequest("", "POST", &buf)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// func main() {
// 	salt := NewSaltClient("192.168.11.23:8000", "saltapi", "saltapi")
// 	resp, err := salt.Auth()
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println(resp)
// 	}
// }
