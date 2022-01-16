package gakujo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/earlgray283/gakujo-google-calendar/gakujo/scrape"
)

type Client struct {
	http  *http.Client
	token string // org.apache.struts.taglib.html.TOKEN
}

func NewClient() *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		http: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Jar:     jar,
			Timeout: 5 * time.Minute,
		},
	}
}

// search a cookie "JSESSIONID" from c.jar
// if not found, return ""
func (c *Client) SessionID() string {
	u, _ := url.Parse("https://gakujo.shizuoka.ac.jp")
	for _, cookie := range c.http.Jar.Cookies(u) {
		if cookie.Name == "JSESSIONID" {
			return cookie.Value
		}
	}
	return ""
}

// fetch page which needs org.apache.struts.taglib.html.TOKEN and save its token
func (c *Client) FetchPage(url string, datas url.Values) ([]byte, error) {
	datas.Set("org.apache.struts.taglib.html.TOKEN", c.token)

	resp, err := c.postForm(url, datas)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response status was %d(expext %d)", resp.StatusCode, http.StatusOK)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	token, err := scrape.ApacheToken(io.NopCloser(bytes.NewReader(b)))
	if err != nil {
		// fetchPage では必ず apache Token が含まれるページを取得するはず
		return nil, err
	}
	c.token = token

	return b, nil
}

// http.Get wrapper
func (c *Client) getWithReferer(url, referer string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Referer", referer)
	return c.http.Do(req)
}

// http.PostForm wrapper
func (c *Client) postForm(url string, datas url.Values) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(datas.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
