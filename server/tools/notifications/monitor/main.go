package main

import (
	"encoding/json"
	"errors"
	"time"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/tom-draper/api-analytics/server/database"
	"github.com/tom-draper/api-analytics/server/email"

	"github.com/joho/godotenv"
)

var url string = "https://apianalytics-server.com/api/"

func TestNewUser() error {
	response, err := http.Get(url + "generate-api-key")
	if err != nil {
		return err
	} else if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("status code: %d", response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	sb := string(body)
	apiKey := sb[1 : len(sb)-1]
	if len(apiKey) != 36 {
		return errors.New(fmt.Sprintf("uuid value returned is invalid"))
	}
	
	err = database.DeleteUser(apiKey)
	if err != nil {
		return err
	}
	return nil
}

func TestFetchData() error {
	client := http.Client{}
	apiKey := getTestAPIKey()
	request, err := http.NewRequest("GET", url+"data", nil)
	if err != nil {
		return err
	}

	request.Header = http.Header{
		"X-AUTH-TOKEN": {apiKey},
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	} else if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("status code: %d", response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var data interface{}
	err = json.Unmarshal(body, &data)
	return err
}

func TestFetchDashboardData() error {
	userID := getTestUserID()
	response, err := http.Get(url + "requests/" + userID)
	if err != nil {
		return err
	} else if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("status code: %d", response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var data interface{}
	err = json.Unmarshal(body, &data)
	return err
}

func TestFetchUserID() error {
	apiKey := getTestAPIKey()
	response, err := http.Get(url + "user-id/" + apiKey)
	if err != nil {
		return err
	} else if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("status code: %d", response.StatusCode))
	}
	
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	sb := string(body)
	userID := sb[1 : len(sb)-1]
	if len(userID) != 36 {
		return errors.New(fmt.Sprintf("uuid value returned is invalid"))
	}
	return nil
}

func TestFetchMonitorPings() error {
	userID := getTestUserID()
	response, err := http.Get(url + "monitor/pings/" + userID)
	if err != nil {
		return err
	} else if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("status code: %d", response.StatusCode))
	}
	
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	
	var data interface{}
	err = json.Unmarshal(body, &data)
	return err
}

func getTestAPIKey() string {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	apiKey := os.Getenv("MONITOR_API_KEY")
	return apiKey
}

func getTestUserID() string {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	userID := os.Getenv("MONITOR_USER_ID")
	return userID
}

func buildEmailBody(newUser error, fetchDashboardData error, fetchData error) string {
	var body strings.Builder
	body.WriteString(fmt.Sprintf("Failure detected at %v\n", time.Now()))
	if newUser != nil {
		body.WriteString(fmt.Sprintf("Error when creating new user: %s\n", newUser.Error()))
	}
	if fetchDashboardData != nil {
		body.WriteString(fmt.Sprintf("Error when fetching dashboard data: %s\n", fetchDashboardData.Error()))
	}
	if fetchData != nil {
		body.WriteString(fmt.Sprintf("Error when fetching API data: %s\n", fetchData.Error()))
	}
	return body.String()
}

func main() {
	newUserSuccessful := TestNewUser()
	fetchDashboardDataSuccessful := TestFetchDashboardData()
	fetchDataSuccessful := TestFetchData()
	newUserSuccessful = errors.New("test error")
	if newUserSuccessful != nil || fetchDashboardDataSuccessful != nil || fetchDataSuccessful != nil {
		address := email.GetEmailAddress()
		body := buildEmailBody(newUserSuccessful, fetchDashboardDataSuccessful, fetchDataSuccessful)
		err := email.SendEmail("Failure detected at API Analytics", body, address)
		if err != nil {
			panic(err)
		}
	}
}