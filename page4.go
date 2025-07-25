// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package invidious

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/antoniszymanski/option-go"
	"github.com/cli/browser"
	"github.com/go-json-experiment/json"
)

func (c *Client) AuthorizeToken(req *AuthorizeTokenRequest) (err error) {
	query := make(url.Values, 3)
	query.Set("scopes", strings.Join(req.Scopes, ","))
	query.Set("callback_url", "http://localhost:8080")
	query.Set("expire", itoa(req.Expire.Unix()))
	url := c.InstanceURL + "/authorize_token" + "?" + query.Encode()
	if err = browser.OpenURL(url); err != nil {
		return
	}
	c.RawToken, err = getToken()
	return
}

type AuthorizeTokenRequest struct {
	Scopes []string
	Expire time.Time
}

func getToken() (string, error) {
	srv := http.Server{Addr: ":8080", ReadHeaderTimeout: 10 * time.Second}
	mux := http.NewServeMux()
	var token string
	mux.HandleFunc("/{$}",
		//nolint:errcheck
		func(w http.ResponseWriter, r *http.Request) {
			token = r.URL.Query().Get("token")
			if token == "" {
				w.Write([]byte("Error: 'token' parameter is missing from the URL"))
				return
			}
			w.Write([]byte("Success! Now you can close this page"))
			go func() {
				<-r.Context().Done()
				srv.Close()
			}()
		},
	)
	srv.Handler = mux

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return "", err
	}
	return token, nil
}

func (c *Client) Feed(req *FeedRequest) (*FeedResponse, error) {
	query := make(url.Values)
	if req.MaxResults.IsSome() {
		query.Set("max_results", itoa(req.MaxResults.Unwrap()))
	}
	if req.Page.IsSome() {
		query.Set("page", itoa(req.Page.Unwrap()))
	}
	var resp FeedResponse
	if err := c.call(&requestConfig{
		Method: "GET",
		Path:   "/api/v1/auth/feed",
		Auth:   true,
		Query:  query,
		Output: &resp,
	}); err != nil {
		return nil, err
	}
	return &resp, nil
}

type FeedRequest struct {
	MaxResults option.Option[int32]
	Page       option.Option[int32]
}

type FeedResponse struct {
	Notifications []struct {
		Type            string `json:"type"` // "shortVideo"
		Title           string `json:"title"`
		VideoId         string `json:"videoId"`
		VideoThumbnails []struct {
			Quality string `json:"quality"`
			Url     string `json:"url"`
			Width   int64  `json:"width"`
			Height  int64  `json:"height"`
		} `json:"videoThumbnails"`
		LengthSeconds int64  `json:"lengthSeconds"`
		Author        string `json:"author"`
		AuthorId      string `json:"authorId"`
		AuthorUrl     string `json:"authorUrl"`
		Published     int64  `json:"published"`
		PublishedText string `json:"publishedText"`
		ViewCount     int64  `json:"viewCount"`
	} `json:"notifications"`
	Videos []struct {
		Type            string `json:"type"` // "shortVideo"
		Title           string `json:"title"`
		VideoId         string `json:"videoId"`
		VideoThumbnails []struct {
			Quality string `json:"quality"`
			Url     string `json:"url"`
			Width   int64  `json:"width"`
			Height  int64  `json:"height"`
		} `json:"videoThumbnails"`
		LengthSeconds int64  `json:"lengthSeconds"`
		Author        string `json:"author"`
		AuthorId      string `json:"authorId"`
		AuthorUrl     string `json:"authorUrl"`
		Published     int64  `json:"published"`
		PublishedText string `json:"publishedText"`
		ViewCount     int64  `json:"viewCount"`
	} `json:"videos"`
}

func (c *Client) Playlists() (PlaylistsResponse, error) {
	var resp PlaylistsResponse
	if err := c.call(&requestConfig{
		Method: "GET",
		Path:   "/api/v1/auth/playlists",
		Auth:   true,
		Output: &resp,
	}); err != nil {
		return nil, err
	}
	return resp, nil
}

type PlaylistsResponse []struct {
	Type             string `json:"type"` // "invidiousPlaylist"
	Title            string `json:"title"`
	PlaylistId       string `json:"playlistId"`
	Author           string `json:"author"`
	AuthorId         any    `json:"authorId"`
	AuthorUrl        any    `json:"authorUrl"`
	AuthorThumbnails []any  `json:"authorThumbnails"`
	Description      string `json:"description"`
	DescriptionHtml  string `json:"descriptionHtml"`
	VideoCount       int32  `json:"videoCount"`
	ViewCount        int64  `json:"viewCount"`
	Updated          int64  `json:"updated"`
	IsListed         bool   `json:"isListed"`
	Videos           []struct {
		Title           string `json:"title"`
		VideoId         string `json:"videoId"`
		Author          string `json:"author"`
		AuthorId        string `json:"authorId"`
		AuthorUrl       string `json:"authorUrl"`
		VideoThumbnails []struct {
			Quality string `json:"quality"`
			Url     string `json:"url"`
			Width   int32  `json:"width"`
			Height  int32  `json:"height"`
		} `json:"videoThumbnails"`
		Index         int32  `json:"index"`
		IndexId       string `json:"indexId"`
		LengthSeconds int32  `json:"lengthSeconds"`
	} `json:"videos"`
}

func (c *Client) CreatePlaylist(req *CreatePlaylistRequest) (*CreatePlaylistResponse, error) {
	var resp CreatePlaylistResponse
	if err := c.call(&requestConfig{
		Method: "POST",
		Path:   "/api/v1/auth/playlists",
		Auth:   true,
		Input:  req,
		Output: &resp,
	}); err != nil {
		return nil, err
	}
	return &resp, nil
}

type CreatePlaylistRequest struct {
	Title   string  `json:"title"`
	Privacy Privacy `json:"privacy"`
}

type Privacy string

const (
	Public   Privacy = "public"
	Unlisted Privacy = "unlisted"
	Private  Privacy = "private"
)

type CreatePlaylistResponse struct {
	Title      string `json:"title"`
	PlaylistId string `json:"playlistId"`
}

func (c *Client) Playlist(id string) (*PlaylistResponse, error) {
	var resp PlaylistResponse
	if err := c.call(&requestConfig{
		Method: "GET",
		Path:   "/api/v1/auth/playlists/" + id,
		Auth:   true,
		Output: &resp,
	}); err != nil {
		return nil, err
	}
	return &resp, nil
}

type PlaylistResponse struct {
	Title            string `json:"title"`
	PlaylistId       string `json:"playlistId"`
	Author           string `json:"author"`
	AuthorId         string `json:"authorId"`
	AuthorThumbnails []struct {
		Url    string `json:"url"`
		Width  string `json:"width"`
		Height string `json:"height"`
	} `json:"authorThumbnails"`
	Description     string `json:"description"`
	DescriptionHtml string `json:"descriptionHtml"`
	VideoCount      int32  `json:"videoCount"`
	ViewCount       int64  `json:"viewCount"`
	ViewCountText   string `json:"viewCountText"`
	Updated         int64  `json:"updated"`
	Videos          []struct {
		Title           string `json:"title"`
		VideoId         string `json:"videoId"`
		Author          string `json:"author"`
		AuthorId        string `json:"authorId"`
		AuthorUrl       string `json:"authorUrl"`
		VideoThumbnails []struct {
			Quality string `json:"quality"`
			Url     string `json:"url"`
			Width   int32  `json:"width"`
			Height  int32  `json:"height"`
		} `json:"videoThumbnails"`
		Index         int32  `json:"index"`
		IndexId       string `json:"indexId"`
		LengthSeconds int32  `json:"lengthSeconds"`
	} `json:"videos"`
}

func (c *Client) UpdatePlaylist(req *UpdatePlaylistRequest) error {
	return c.call(&requestConfig{
		Method: "PATCH",
		Path:   "/api/v1/auth/playlists/" + req.Id,
		Auth:   true,
		Input:  req,
	})
}

type UpdatePlaylistRequest struct {
	Id          string                 `json:"-"`
	Title       option.Option[string]  `json:"title"`
	Description option.Option[string]  `json:"description"`
	Privacy     option.Option[Privacy] `json:"privacy"`
}

func (c *Client) DeletePlaylist(id string) error {
	return c.call(&requestConfig{
		Method: "DELETE",
		Path:   "/api/v1/auth/playlists/" + id,
		Auth:   true,
	})
}

func (c *Client) AddVideo(req *AddVideoRequest) (*AddVideoResponse, error) {
	var resp AddVideoResponse
	if err := c.call(&requestConfig{
		Method: "POST",
		Path:   "/api/v1/auth/playlists/" + req.PlaylistId + "/videos",
		Auth:   true,
		Input:  req,
		Output: &resp,
	}); err != nil {
		return nil, err
	}
	return &resp, nil
}

type AddVideoRequest struct {
	PlaylistId string `json:"-"`
	VideoId    string `json:"videoId"`
}

type AddVideoResponse struct {
	Title           string `json:"title"`
	VideoId         string `json:"videoId"`
	Author          string `json:"author"`
	AuthorId        string `json:"authorId"`
	AuthorUrl       string `json:"authorUrl"`
	VideoThumbnails []struct {
		Quality string `json:"quality"`
		Url     string `json:"url"`
		Width   int32  `json:"width"`
		Height  int32  `json:"height"`
	} `json:"videoThumbnails"`
}

func (c *Client) DeleteVideo(req *DeleteVideoRequest) error {
	return c.call(&requestConfig{
		Method: "DELETE",
		Path:   "/api/v1/auth/playlists/" + req.PlaylistId + "/videos/" + req.IndexId,
		Auth:   true,
	})
}

type DeleteVideoRequest struct {
	PlaylistId string
	IndexId    string
}

func (c *Client) Preferences() (*PreferencesResponse, error) {
	var resp PreferencesResponse
	if err := c.call(&requestConfig{
		Method: "GET",
		Path:   "/api/v1/auth/preferences",
		Auth:   true,
		Output: &resp,
	}); err != nil {
		return nil, err
	}
	return &resp, nil
}

type PreferencesResponse struct {
	Annotations           bool     `json:"annotations"`
	AnnotationsSubscribed bool     `json:"annotations_subscribed"`
	Autoplay              bool     `json:"autoplay"`
	Captions              []string `json:"captions"`
	Comments              []string `json:"comments"`
	Continue              bool     `json:"continue"`
	ContinueAutoplay      bool     `json:"continue_autoplay"`
	DarkMode              string   `json:"dark_mode"`
	LatestOnly            bool     `json:"latest_only"`
	Listen                bool     `json:"listen"`
	Local                 bool     `json:"local"`
	Locale                string   `json:"locale"`
	MaxResults            int32    `json:"max_results"`
	NotificationsOnly     bool     `json:"notifications_only"`
	PlayerStyle           string   `json:"player_style"`
	Quality               string   `json:"quality"`
	DefaultHome           string   `json:"default_home"`
	FeedMenu              []string `json:"feed_menu"`
	RelatedVideos         bool     `json:"related_videos"`
	Sort                  string   `json:"sort"`
	Speed                 float64  `json:"speed"`
	ThinMode              bool     `json:"thin_mode"`
	UnseenOnly            bool     `json:"unseen_only"`
	VideoLoop             bool     `json:"video_loop"`
	Volume                uint8    `json:"volume"`
}

// func (c *Client) UpdatePreferences() {
// }

func (c *Client) Subscriptions() (SubscriptionsResponse, error) {
	var resp SubscriptionsResponse
	if err := c.call(&requestConfig{
		Method: "GET",
		Path:   "/api/v1/auth/subscriptions",
		Auth:   true,
		Output: &resp,
	}); err != nil {
		return nil, err
	}
	return resp, nil
}

type SubscriptionsResponse []struct {
	Author   string `json:"author"`
	AuthorId string `json:"authorId"`
}

func (c *Client) AddSubscription(ucid string) error {
	return c.call(&requestConfig{
		Method: "POST",
		Path:   "/api/v1/auth/subscriptions/" + ucid,
		Auth:   true,
	})
}

func (c *Client) RemoveSubscription(ucid string) error {
	return c.call(&requestConfig{
		Method: "DELETE",
		Path:   "/api/v1/auth/subscriptions/" + ucid,
		Auth:   true,
	})
}

func (c *Client) Tokens() (TokensResponse, error) {
	var resp TokensResponse
	if err := c.call(&requestConfig{
		Method: "GET",
		Path:   "/api/v1/auth/tokens",
		Auth:   true,
		Output: &resp,
	}); err != nil {
		return nil, err
	}
	return resp, nil
}

type TokensResponse []struct {
	Session string `json:"session"`
	Issued  int64  `json:"issued"`
}

func (c *Client) RegisterToken(req *RegisterTokenRequest) (*Token, error) {
	var resp Token
	if err := c.call(&requestConfig{
		Method: "POST",
		Path:   "/api/v1/auth/tokens/register",
		Auth:   true,
		Input:  req,
		Output: &resp,
	}); err != nil {
		return nil, err
	}
	return &resp, nil
}

type RegisterTokenRequest struct {
	Scopes      []string              `json:"scopes"`
	CallbackUrl option.Option[string] `json:"callbackUrl"`
	Expire      time.Time             `json:"expire"`
}

type Token struct {
	Session   string    `json:"session"`
	Scopes    []string  `json:"scopes"`
	Expire    time.Time `json:"expire"`
	Signature string    `json:"signature"`
}

func (t *Token) Encode() (string, error) {
	out, err := json.Marshal(&t, opts)
	if err != nil {
		return "", err
	}
	return url.QueryEscape(bytes2string(out)), nil
}

func ParseToken(in string) (*Token, error) {
	in, err := url.QueryUnescape(in)
	if err != nil {
		return nil, err
	}
	var t Token
	err = json.Unmarshal(string2bytes(in), &t, opts)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (c *Client) RevokeToken(req *RevokeRequest) error {
	return c.call(&requestConfig{
		Method: "POST",
		Path:   "/api/v1/auth/tokens/unregister",
		Auth:   true,
		Input:  req,
	})
}

type RevokeRequest struct {
	Session string `json:"session"`
}

func (c *Client) History(req *HistoryRequest) (HistoryResponse, error) {
	query := make(url.Values)
	if req.MaxResults.IsSome() {
		query.Set("max_results", itoa(req.MaxResults.Unwrap()))
	}
	if req.Page.IsSome() {
		query.Set("page", itoa(req.Page.Unwrap()))
	}
	var resp HistoryResponse
	if err := c.call(&requestConfig{
		Method: "GET",
		Path:   "/api/v1/auth/history",
		Auth:   true,
		Query:  query,
		Output: &resp,
	}); err != nil {
		return nil, err
	}
	return resp, nil
}

type HistoryRequest struct {
	MaxResults option.Option[int32]
	Page       option.Option[int32]
}

type HistoryResponse []string

func (c *Client) AddToHistory(id string) error {
	return c.call(&requestConfig{
		Method: "POST",
		Path:   "/api/v1/auth/history/" + id,
		Auth:   true,
	})
}

func (c *Client) DeleteFromHistory(id string) error {
	return c.call(&requestConfig{
		Method: "DELETE",
		Path:   "/api/v1/auth/history/" + id,
		Auth:   true,
	})
}
