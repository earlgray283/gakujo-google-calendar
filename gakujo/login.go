package gakujo

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/earlgray283/gakujo-google-calendar/gakujo/scrape"
)

func (c *Client) Login(username, password string) error {
	if err := c.fetchGakujoPortalJSESSIONID(); err != nil {
		return err
	}

	if err := c.fetchGakujoRootJSESSIONID(); err != nil {
		return err
	}

	if err := c.preLogin(); err != nil {
		return err
	}
	resp, err := c.shibbolethlogin()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		return &ErrUnexpectedStatus{resp.StatusCode, http.StatusOK}
	}

	// セッションがないとき
	if resp.StatusCode == http.StatusFound {
		loginAPIurl, err := c.fetchLoginAPIurl(resp.Header.Get("Location"))
		if err != nil {
			return err
		}
		if err := c.login("https://idp.shizuoka.ac.jp"+loginAPIurl, username, password); err != nil {
			return err
		}
	}

	return c.initialize()
}

func (c *Client) fetchGakujoPortalJSESSIONID() error {
	resp, err := c.http.Get("https://gakujo.shizuoka.ac.jp/portal/")
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return &ErrUnexpectedStatus{resp.StatusCode, http.StatusOK}
	}

	return nil
}

func (c *Client) fetchGakujoRootJSESSIONID() error {
	unixmilli := time.Now().UnixNano() / 1000000
	resp, err := c.http.Get("https://gakujo.shizuoka.ac.jp/UI/jsp/topPage/topPage.jsp?_=" + strconv.FormatInt(unixmilli, 10))
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return &ErrUnexpectedStatus{resp.StatusCode, http.StatusOK}
	}

	return nil
}

func (c *Client) preLogin() error {
	datas := url.Values{}
	datas.Set("mistakeChecker", "0")

	resp, err := c.http.PostForm("https://gakujo.shizuoka.ac.jp/portal/login/preLogin/preLogin", datas)
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return &ErrUnexpectedStatus{resp.StatusCode, http.StatusOK}
	}

	return nil
}

func (c *Client) fetchLoginAPIurl(SSOSAMLRequestURL string) (string, error) {
	resp, err := c.http.Get(SSOSAMLRequestURL)
	if err != nil {
		return "", err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusFound {
		return "", &ErrUnexpectedStatus{resp.StatusCode, http.StatusOK}
	}
	return resp.Header.Get("Location"), nil
}

func (c *Client) login(uri, username, password string) error {
	htmlReadCloser, err := c.postSSOexecution(uri, username, password)
	if err != nil {
		return err
	}
	relayState, samlResponse, err := scrape.RelayStateAndSAMLResponse(htmlReadCloser)
	if err != nil {
		return err
	}
	htmlReadCloser.Close()
	_, _ = io.Copy(io.Discard, htmlReadCloser)

	location, err := c.fetchSSOinitLoginLocation(relayState, samlResponse)
	if err != nil {
		return err
	}

	resp, err := c.getWithReferer(location, "https://idp.shizuoka.ac.jp/")
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusOK {
		return &ErrUnexpectedStatus{resp.StatusCode, http.StatusOK}
	}

	return nil
}

func (c *Client) postSSOexecution(uri, username, password string) (io.ReadCloser, error) {
	datas := make(url.Values)
	datas.Set("j_username", username)
	datas.Set("j_password", password)
	datas.Set("_eventId_proceed", "")

	resp, err := c.http.PostForm(uri, datas)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &ErrUnexpectedStatus{resp.StatusCode, http.StatusOK}
	}

	return resp.Body, nil
}

func (c *Client) shibbolethlogin() (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodPost, "https://gakujo.shizuoka.ac.jp/portal/shibbolethlogin/shibbolethLogin/initLogin/sso", nil)
	resp, err := c.http.Do(req)
	resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp, err
}

func (c *Client) fetchSSOinitLoginLocation(relayState, samlResponse string) (string, error) {
	datas := make(url.Values)
	datas.Set("RelayState", relayState)
	datas.Set("SAMLResponse", samlResponse)

	resp, err := c.http.PostForm("https://gakujo.shizuoka.ac.jp/Shibboleth.sso/SAML2/POST", datas)
	if err != nil {
		return "", err
	}
	defer func() {
		resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()
	if resp.StatusCode != http.StatusFound {
		return "", &ErrUnexpectedStatus{resp.StatusCode, http.StatusOK}
	}

	return resp.Header.Get("Location"), nil
}

func (c *Client) initialize() error {
	datas := make(url.Values)
	datas.Set("EXCLUDE_SET", "")

	_, err := c.FetchPage("https://gakujo.shizuoka.ac.jp/portal/home/home/initialize", datas)
	if err != nil {
		return err
	}

	return nil
}
