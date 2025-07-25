// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package invidious

import (
	"net/url"

	"github.com/antoniszymanski/option-go"
)

func (c *Client) Stats() (*StatsResponse, error) {
	var resp StatsResponse
	if err := c.call(&requestConfig{
		Method: "GET",
		Path:   "/api/v1/stats",
		Output: &resp,
	}); err != nil {
		return nil, err
	}
	return &resp, nil
}

type StatsResponse struct {
	Version  string `json:"version"`
	Software struct {
		Name    string `json:"name"` // "invidious"
		Version string `json:"version"`
		Branch  string `json:"branch"`
	} `json:"software"`
	OpenRegistrations bool `json:"openRegistrations"`
	Usage             struct {
		Users struct {
			Total          int32 `json:"total"`
			ActiveHalfyear int32 `json:"activeHalfyear"`
			ActiveMonth    int32 `json:"activeMonth"`
		} `json:"users"`
	} `json:"usage"`
	Metadata struct {
		UpdatedAt              int64 `json:"updatedAt"`
		LastChannelRefreshedAt int64 `json:"lastChannelRefreshedAt"`
	} `json:"metadata"`
	Playback struct {
		TotalRequests      option.Option[int32]   `json:"totalRequests"`
		SuccessfulRequests option.Option[int32]   `json:"successfulRequests"`
		Ratio              option.Option[float32] `json:"ratio"`
	} `json:"playback"`
}

func (c *Client) Video(req VideoRequest) (*VideoResponse, error) {
	query := make(url.Values)
	if req.Region != "" {
		query.Set("region", req.Region)
	}
	var resp VideoResponse
	if err := c.call(&requestConfig{
		Method: "GET",
		Path:   "/api/v1/videos/" + req.Id,
		Query:  query,
		Output: &resp,
	}); err != nil {
		return nil, err
	}
	return &resp, nil
}

type VideoRequest struct {
	Id     string
	Region string // ISO 3166 country code (default: "US")
}

type VideoResponse struct {
	Type            string `json:"type"` // "video"|"published"
	Title           string `json:"title"`
	VideoId         string `json:"videoId"`
	VideoThumbnails []struct {
		Quality string `json:"quality"`
		Url     string `json:"url"`
		Width   int32  `json:"width"`
		Height  int32  `json:"height"`
	} `json:"videoThumbnails"`
	Storyboards []struct {
		Url              string `json:"url"`
		TemplateUrl      string `json:"templateUrl"`
		Width            int32  `json:"width"`
		Height           int32  `json:"height"`
		Count            int32  `json:"count"`
		Interval         int32  `json:"interval"`
		StoryboardWidth  int32  `json:"storyboardWidth"`
		StoryboardHeight int32  `json:"storyboardHeight"`
		StoryboardCount  int32  `json:"storyboardCount"`
	} `json:"storyboards"`
	Description      string   `json:"description"`
	DescriptionHtml  string   `json:"descriptionHtml"`
	Published        int64    `json:"published"`
	PublishedText    string   `json:"publishedText"`
	Keywords         []string `json:"keywords"`
	ViewCount        int64    `json:"viewCount"`
	LikeCount        int32    `json:"likeCount"`
	DislikeCount     int32    `json:"dislikeCount"`
	Paid             bool     `json:"paid"`
	Premium          bool     `json:"premium"`
	IsFamilyFriendly bool     `json:"isFamilyFriendly"`
	AllowedRegions   []string `json:"allowedRegions"`
	Genre            string   `json:"genre"`
	GenreUrl         string   `json:"genreUrl"`
	Author           string   `json:"author"`
	AuthorId         string   `json:"authorId"`
	AuthorUrl        string   `json:"authorUrl"`
	AuthorThumbnails []struct {
		Url    string `json:"url"`
		Width  int32  `json:"width"`
		Height int32  `json:"height"`
	} `json:"authorThumbnails"`
	SubCountText      string                `json:"subCountText"`
	LengthSeconds     int32                 `json:"lengthSeconds"`
	AllowRatings      bool                  `json:"allowRatings"`
	Rating            float32               `json:"rating"`
	IsListed          bool                  `json:"isListed"`
	LiveNow           bool                  `json:"liveNow"`
	IsPostLiveDvr     bool                  `json:"isPostLiveDvr"`
	IsUpcoming        bool                  `json:"isUpcoming"`
	DashUr            string                `json:"dashUr"`
	PremiereTimestamp option.Option[int64]  `json:"premiereTimestamp"`
	HlsUrl            option.Option[string] `json:"hlsUrl"`
	AdaptiveFormats   []struct {
		Index             string                `json:"index"`
		Bitrate           string                `json:"bitrate"`
		Init              string                `json:"init"`
		Url               string                `json:"url"`
		Itag              string                `json:"itag"`
		Type              string                `json:"type"`
		Clen              string                `json:"clen"`
		Lmt               string                `json:"lmt"`
		ProjectionType    string                `json:"projectionType"`
		Container         string                `json:"container"`
		Encoding          string                `json:"encoding"`
		QualityLabel      option.Option[string] `json:"qualityLabel"`
		Resolution        option.Option[string] `json:"resolution"`
		Fps               int32                 `json:"fps"`
		Size              option.Option[string] `json:"size"`
		TargetDurationsec option.Option[int64]  `json:"targetDurationsec"`
		MaxDvrDurationSec option.Option[int64]  `json:"maxDvrDurationSec"`
		AudioQuality      option.Option[string] `json:"audioQuality"`
		AudioSampleRate   option.Option[string] `json:"audioSampleRate"`
		AudioChannels     option.Option[string] `json:"audioChannels"`
		ColorInfo         option.Option[string] `json:"colorInfo"`
		CaptionTrack      option.Option[string] `json:"captionTrack"`
	} `json:"adaptiveFormats"`
	FormatStreams []struct {
		Url          string                `json:"url"`
		Itag         string                `json:"itag"`
		Type         string                `json:"type"`
		Quality      string                `json:"quality"`
		Bitrate      option.Option[string] `json:"bitrate"`
		Container    string                `json:"container"`
		Encoding     string                `json:"encoding"`
		QualityLabel string                `json:"qualityLabel"`
		Resolution   string                `json:"resolution"`
		Size         string                `json:"size"`
	} `json:"formatStreams"`
	Captions []struct {
		Label         string `json:"label"`
		Language_code string `json:"language_Code"`
		Url           string `json:"url"`
	} `json:"captions"`
	MusicTracks []struct {
		Song    string `json:"song"`
		Artist  string `json:"artist"`
		Album   string `json:"album"`
		License string `json:"license"`
	} `json:"musicTracks"`
	RecommendedVideos []struct {
		VideoId         string `json:"videoId"`
		Title           string `json:"title"`
		VideoThumbnails []struct {
			Quality string `json:"quality"`
			Url     string `json:"url"`
			Width   int32  `json:"width"`
			Height  int32  `json:"height"`
		} `json:"videoThumbnails"`
		Author           string                `json:"author"`
		AuthorUrl        string                `json:"authorUrl"`
		AuthorId         option.Option[string] `json:"authorId"`
		AuthorVerified   bool                  `json:"authorVerified"`
		AuthorThumbnails []struct {
			Url    string `json:"url"`
			Width  int32  `json:"width"`
			Height int32  `json:"height"`
		} `json:"authorThumbnails"`
		LengthSeconds int32  `json:"lengthSeconds"`
		ViewCount     int64  `json:"viewCount"`
		ViewCountText string `json:"viewCountText"`
	} `json:"recommendedVideos"`
}
