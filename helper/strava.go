package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"Jimandy-Website-Backend/data"
	"Jimandy-Website-Backend/database"
	"Jimandy-Website-Backend/model"
)

var urlMap = map[string]string{
	"Authorization":                "https://www.strava.com/oauth/token?grant_type=authorization_code",
	"getLoggedInAthleteActivities": "https://www.strava.com/api/v3/athlete/activities?per_page=200",
	"getLapsByActivityId":          "https://www.strava.com/api/v3/activities/{id}/laps",
	"refreshToken":                 "https://www.strava.com/api/v3/oauth/token?", // client_id={client_id}&client_secret={client_secret}&grant_type=refresh_token&refresh_token={refresh_token}
}

func FetchStravaApi[T any](method string, url string, accessToken string, result T) (T, error) {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return result, err
	}

	request.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Error making request:", err)
		return result, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return result, err
	}

	return result, nil
}

func AuthorizeStrava(myAthlete model.Athlete) bool {
	authorizationUrl := fmt.Sprintf("%s&client_id=%d&client_secret=%s&code=%s", urlMap["Authorization"], myAthlete.ClientID, myAthlete.ClientSecret, myAthlete.AuthorizationCode)

	tokens, err := FetchStravaApi("POST", authorizationUrl, "", data.Token{})
	if err != nil || tokens.Athlete.AthleteID == 0 {
		return false
	}

	myAthlete.ID = tokens.Athlete.AthleteID
	myAthlete.AccessToken = tokens.AccessToken
	myAthlete.RefreshToken = tokens.RefreshToken

	return database.AddAthlete(&myAthlete)
}

func FetchNewAccessToken(myAthlete model.Athlete) string {
	refreshTokenUrl := fmt.Sprintf("%sclient_id=%d&client_secret=%s&grant_type=refresh_token&refresh_token=%s", urlMap["refreshToken"], myAthlete.ClientID, myAthlete.ClientSecret, myAthlete.RefreshToken)

	tokens, err := FetchStravaApi("POST", refreshTokenUrl, "", data.Token{})
	if err != nil {
		return ""
	}

	_ = database.UpdateAthleteToken(myAthlete.ID, tokens.AccessToken, tokens.RefreshToken)

	return tokens.AccessToken
}

func FetchStravaActivities(myAthlete model.Athlete) error {
	url := urlMap["getLoggedInAthleteActivities"]

	lastActivityDate := database.GetLatestActivityDate(uint64(myAthlete.ID)) // 取得運動員最近一次的跑步活動日期
	// 若 有最後一次跑步活動 則 爬取這之後的活動
	if lastActivityDate.Unix() > 0 {
		url += fmt.Sprintf("&after=%d", lastActivityDate.Unix())
	}

	var activities []data.Activities

	page := 1
	for {
		result, err := FetchStravaApi("GET", fmt.Sprintf("%s&page=%d", url, page), myAthlete.AccessToken, []data.Activities{})
		if err != nil {
			accessToken := FetchNewAccessToken(myAthlete)

			if accessToken == "" {
				return errors.New("Invalid AccessToken")
			}

			myAthlete.AccessToken = accessToken
			result, err = FetchStravaApi("GET", url, myAthlete.AccessToken, []data.Activities{})
			if err != nil {
				return errors.New("Error fetching strava api")
			}
		}

		if len(result) == 0 {
			break
		}
		page += 1
		activities = append(activities, result...)
	}

	if !database.AddActivity(activities) {
		return errors.New("Error Saving activities")
	}

	return nil
}

// 取得活動圈數
func FetchStravaLaps(myAthlete model.Athlete, activityID string) error {
	url := strings.Replace(urlMap["getLapsByActivityId"], "{id}", activityID, 1)

	result, err := FetchStravaApi("GET", url, myAthlete.AccessToken, []data.Lap{})
	if err != nil {
		accessToken := FetchNewAccessToken(myAthlete)
		if accessToken == "" {
			return errors.New("Invalid AccessToken")
		}

		result, err = FetchStravaApi("GET", url, accessToken, []data.Lap{})
		if err != nil {
			return errors.New("Error fetching strava api")
		}
	}

	if !database.AddLaps(result) {
		return errors.New("Error Saving laps error")
	}

	return nil
}
