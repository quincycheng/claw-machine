package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/quincycheng/claw-machine/util"
	"github.com/tidwall/gjson"
)

const authnPathStandalone = "/PasswordVault/API/auth/Cyberark/Logon/"
const groupPathStandalone = "/PasswordVault/api/UserGroups?filter=groupType eq Vault"
const userPathStandalone = "/PasswordVault/api/Users?filter=userType&extendedDetails=true"
const userDetailStandalone = "/PasswordVault/API/Users/"
const groupDetailStandalone = "/PasswordVault/API/UserGroups/"

func CyberArkAuthnStandalone(config util.Config) (string, error) {
	url := config.From.Url + authnPathStandalone
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf("{\"username\": \"%s\",\"password\": \"%s\",\"concurrentSession\": \"true\"}", config.From.User, config.From.Password))
	//	payload := strings.NewReader("{\"username\": \"" + config.From.User + "\",\"password\": \"" + config.From.Password + "\",\"concurrentSession\": \"true\"}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if gjson.Get(string(body), "ErrorCode").String() != "" {
		return "", fmt.Errorf(string(body))
	}

	return util.TrimQuotes(string(body)), nil
}

func GetAllGroupsStandalone(config util.Config, token string, url string) (string, error) {
	if url == "INIT" {
		url = config.From.Url + groupPathStandalone
	}
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func GetAllUsersStandalone(config util.Config, token string, url string) (string, error) {
	if url == "INIT" {
		url = config.From.Url + userPathStandalone
	}
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func GetUserDetailStandalone(config util.Config, token string, userId string) (string, error) {
	url := config.From.Url + userDetailStandalone + userId + "/"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func GetGroupDetailStandalone(config util.Config, token string, groupId string) (string, error) {
	url := config.From.Url + groupDetailStandalone + groupId + "/"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
