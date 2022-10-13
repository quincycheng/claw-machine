package pkg

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/quincycheng/claw-machine/util"
)

const authnPathISP = "/api/idadmin/oauth2/platformtoken"
const userPathISP = "/api/idadmin/CDirectoryService/GetUsers"
const groupPathISP = "/api/idadmin//UserMgmt/GetUsersRolesAndAdministrativeRights"

func AuthnSharedService(config util.Config) (string, error) {
	theUrl := config.To.Url + authnPathISP
	method := "POST"

	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	client := &http.Client{}
	req, err := http.NewRequest(method, theUrl, strings.NewReader(form.Encode()))

	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Basic "+basicAuth(config.To.User, config.To.Password))
	req.Header.Add("Accept", "*/*")

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(string(body)), &jsonMap)
	return jsonMap["access_token"].(string), nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func GetAllRolesISP(config util.Config, token string) (string, error) {
	theUrl := config.To.Url + groupPathISP
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, theUrl, nil)

	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "*/*")

	res, err := client.Do(req)
	if err != nil {
		//	return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}
	return string(body), nil
}
func GetAllUsersISP(config util.Config, token string) (string, error) {
	theUrl := config.To.Url + userPathISP
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, theUrl, nil)

	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "*/*")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	_ = body
	if err != nil {
		return "", err
	}

	//fmt.Println(string(body))
	return string(body), nil
}
