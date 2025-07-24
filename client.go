// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package invidious

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Client struct {
	InstanceURL string
	RawToken    string
	UserAgent   string
	HTTPClient  *http.Client
}

func NewClient(instanceURL string) *Client {
	return &Client{InstanceURL: instanceURL}
}

type requestConfig struct {
	Method string
	Path   string
	Auth   bool
	Query  url.Values
	Input  any
	Output any
}

var opts = json.JoinOptions(
	json.WithMarshalers(json.MarshalToFunc(
		func(enc *jsontext.Encoder, t time.Time) error {
			return enc.WriteToken(jsontext.Int(t.Unix()))
		},
	)),
	json.WithUnmarshalers(json.UnmarshalFromFunc(
		func(dec *jsontext.Decoder, t *time.Time) error {
			token, err := dec.ReadToken()
			if err != nil {
				return err
			}
			if kind := token.Kind(); kind != '0' {
				return errors.New("invalid JSON token kind: " + kind.String())
			}
			*t = time.Unix(token.Int(), 0)
			return nil
		},
	)),
)

func (c *Client) call(config *requestConfig) error {
	var query string
	if len(config.Query) > 0 {
		query = "?" + config.Query.Encode()
	}

	var body io.Reader
	if config.Input != nil {
		bodyData, err := json.Marshal(config.Input, opts)
		if err != nil {
			return err
		}
		body = bytes.NewReader(bodyData)
	}

	req, err := http.NewRequest(
		config.Method,
		c.InstanceURL+config.Path+query,
		body,
	)
	if err != nil {
		return err
	}

	if config.Auth {
		req.Header.Set("Authorization", "Bearer "+c.RawToken)
	}
	if config.Input != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	} else {
		req.Header.Set("User-Agent", pkgPath+" "+Version())
	}

	var resp *http.Response
	if c.HTTPClient != nil {
		resp, err = c.HTTPClient.Do(req)
	} else {
		resp, err = http.DefaultClient.Do(req)
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return newError(resp)
	}

	if config.Output != nil {
		err = json.UnmarshalRead(resp.Body, config.Output, opts)
	}
	return err
}

func newError(resp *http.Response) (e Error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		body = nil
	}
	err = json.Unmarshal(body, &e, json.RejectUnknownMembers(true))
	if err != nil {
		e = Error{Message: bytes2string(body)}
	}
	e.StatusCode = resp.StatusCode
	return
}

type Error struct {
	StatusCode int    `json:"-"`
	Message    string `json:"error"`
}

func (e Error) Error() string {
	statusText := http.StatusText(e.StatusCode)
	sz := 3 + 1 + len(statusText)
	if e.Message != "" {
		sz += 3 + quotedLen(e.Message)
	}
	dst := make([]byte, 0, sz)
	dst = appendInt(dst, e.StatusCode)
	dst = append(dst, ' ')
	dst = append(dst, statusText...)
	if e.Message != "" {
		dst = append(dst, " - "...)
		dst = appendQuote(dst, e.Message)
	}
	return bytes2string(dst)
}
