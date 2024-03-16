package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/vltvdnl/Adviser-Bot/lib/e"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host, token string) Client {
	return Client{
		host:     host,
		basePath: NewBasePath(token),
		client:   http.Client{},
	}
}
func NewBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))
	data, err := c.DoRequest("getUpdates", q)

	if err != nil {
		return nil, err
	}

	var res UpdatesResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil

}
func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	_, err := c.DoRequest("sendMessage", q)
	if err != nil {
		return e.Wrap("cant't send message", err)
	}
	return nil
}
func (c *Client) DoRequest(method string, query url.Values) (data []byte, err error) {
	defer func() {
		err = e.Wrap("can't do request", err)
	}()
	const ErrMsg = "can't do request"

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, e.Wrap(ErrMsg, err)
	}
	req.URL.RawQuery = query.Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap(ErrMsg, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap(ErrMsg, err)
	}
	return body, nil

}
