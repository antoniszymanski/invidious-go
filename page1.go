// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package invidious

import "github.com/antoniszymanski/option-go"

func (c *Client) Stats() (StatsResponse, error) {
	var resp StatsResponse
	if err := c.call(requestConfig{
		Method: "GET",
		Path:   "/api/v1/stats",
		Output: &resp,
	}); err != nil {
		return StatsResponse{}, err
	}
	return resp, nil
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
