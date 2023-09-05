package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sismo-datagroup-service/app/form"
	"sismo-datagroup-service/app/model"
)

type Users struct {
	Data []UserInfo `json:"data"`
}

type UserInfo struct {
	Username      string `json:"username"`
	PublicMetrics struct {
		FollowersCount int `json:"followers_count"`
		FollowingCount int `json:"following_count"`
		TweetCount     int `json:"tweet_count"`
		ListedCount    int `json:"listed_count"`
	} `json:"public_metrics"`
	Name string `json:"name"`
	ID   string `json:"id"`
}

func parseUserInfo(jsonString string) (*Users, error) {
	//jsonString := `{"data":[{"username":"randomprime","public_metrics":{"followers_count":5,"following_count":54,"tweet_count":14,"listed_count":0},"name":"Bruce Bromley","id":"19370684"}]}`

	var users Users
	err := json.Unmarshal([]byte(jsonString), &users)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}

	//for _, user := range users.Data {
	//	fmt.Println("Username:", user.Username)
	//	fmt.Println("Followers Count:", user.PublicMetrics.FollowersCount)
	//	fmt.Println("Following Count:", user.PublicMetrics.FollowingCount)
	//	fmt.Println("Tweet Count:", user.PublicMetrics.TweetCount)
	//	fmt.Println("Listed Count:", user.PublicMetrics.ListedCount)
	//	fmt.Println("Name:", user.Name)
	//	fmt.Println("ID:", user.ID)
	//	fmt.Println()
	//}
	return &users, nil
}

var client = &http.Client{}

var _ Handler = new(TwitterFollowerHandler)

var TwitterFollowerHandlerName = "TwitterFollower"

type TwitterFollowerHandler struct{}

func (*TwitterFollowerHandler) ValidateRecord(record form.RecordForm, groupMeta *model.DataGroupMate) (string, error) {
	account, err := processRecord(record, groupMeta)
	if err != nil {
		return "", err
	}
	fmt.Println("account", account)
	return account, nil
}

func (*TwitterFollowerHandler) HandlerName() string {
	return TwitterFollowerHandlerName
}

func processRecord(record form.RecordForm, groupMeta *model.DataGroupMate) (string, error) {
	//url := "https://api.twitter.com/2/users/by?usernames=randomprime&user.fields=public_metrics"
	url := fmt.Sprintf("https://api.twitter.com/2/users/by?usernames=%s&user.fields=public_metrics", record.Account)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if record.TwitterAccessToken == "" {
		return "", errors.New("EmptyAccessToken")
	}

	authorization := fmt.Sprintf("Bearer %s", record.TwitterAccessToken)
	req.Header.Add("Authorization", authorization)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("ReadResponseBodyError")
	}
	user, err := parseUserInfo(string(body))
	if err != nil {
		return "", errors.New("ParseTwitterUserInfoError")
	}

	var account string
	if len(user.Data) == 0 {
		return "", errors.New("TwitterUserDataNoExist")
	}

	if user.Data[0].PublicMetrics.FollowersCount > groupMeta.TwitterConfig.Followers {
		account = fmt.Sprintf("twitter:%s", user.Data[0].Username)
		return account, nil
	} else {
		return "", errors.New("FollowersCountTooLow")
	}
}
