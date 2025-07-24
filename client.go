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
	URL   string
	Token string
}

const ua = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"

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

func (c *Client) call(config requestConfig) error {
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
		c.URL+config.Path+query,
		body,
	)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", ua)
	if config.Auth {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	if config.Input != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
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
	if e.Message == "" {
		return "error " + itoa(e.StatusCode) + " - " + http.StatusText(e.StatusCode)
	} else {
		return "error " + itoa(e.StatusCode) + " - " + http.StatusText(e.StatusCode) + " - " + quote(e.Message)
	}
}
