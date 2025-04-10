package steam

import (
	"encoding/json"
	"io"
	"net/url"
	"strconv"
)

const (
	apiUpToDateCheck = APIBaseUrl + "/ISteamApps/UpToDateCheck/v1?"
)

func (session *Session) GetRequiredSteamAppVersion(appID int) (int, error) {
	resp, err := session.client.Get(apiUpToDateCheck + url.Values{
		"appid":   {strconv.Itoa(appID)},
		"version": {"0"},
	}.Encode())
	if resp != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
	}

	if err != nil {
		return 0, err
	}

	type UpToDateCheckResponse struct {
		RequiredVersion int `json:"required_version"`
	}

	type Response struct {
		Inner UpToDateCheckResponse `json:"response"`
	}

	var response Response
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}
	return response.Inner.RequiredVersion, nil
}
