// Package sitooutils contains utility functions for working with Sitoo.
package sitooutils

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

//BasicAuth - Functon to convert string base64
func BasicAuth(user, pass string) string {
	auth := user + ":" + pass
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

//Node - Function to set node
func Node(accountNo string) string {
	var node string
	if string(accountNo[0:2]) == "90" {
		node = "-sandbox"
	} else if string(accountNo[0:2]) == "91" {
		node = "140"
	} else if string(accountNo[0:3]) == "119" {
		node = "130"
	} else if string(accountNo[0:2]) == "13" {
		node = "130"
	} else if string(accountNo[0:3]) == "100" {
		node = ""
	} else if string(accountNo[0:3]) == "201" {
		node = "201"
	} else if string(accountNo[0:3]) == "202" {
		node = "202"
	} else if string(accountNo[0:3]) == "203" {
		node = "203"
	}
	return node
}

//GetSitoo - Function to GET data from Sitoo
func GetSitoo(baseURL string, endpoint string, account string, password string) []byte {
	accountSplit := strings.Split(account, "-")
	accountNo := accountSplit[0]

	req, err := http.NewRequest("GET", baseURL+accountNo+endpoint, nil)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Basic "+BasicAuth(account, password))
	resp, err := http.DefaultClient.Do(req)

	log.WithFields(log.Fields{
		"requesttype": "GET",
		"account":     account,
		"endpoint":    endpoint,
		"body":        nil,
	}).Debug("Request sent")

	if err != nil {
		log.WithFields(log.Fields{
			"requesttype": "GET/Error",
			"account":     account,
			"endpoint":    endpoint,
			"body":        nil,
			"response":    err,
		}).Error("ERROR")
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 500 || resp.StatusCode == 401 {
		log.WithFields(log.Fields{
			"requesttype": "GET/Response",
			"account":     account,
			"endpoint":    endpoint,
			"statuscode":  resp.StatusCode,
			"response":    string(response),
		}).Error("ERROR")
		os.Exit(1)
	} else if resp.StatusCode == 429 {
		xRateLimitReset := resp.Header.Get("X-Rate-Limit-Reset")
		sleepTime, _ := strconv.Atoi(xRateLimitReset)
		time.Sleep(time.Duration(sleepTime) * time.Second)
		GetSitoo(baseURL,endpoint,account,password)
	} else {
		log.WithFields(log.Fields{
			"requesttype": "GET/Response",
			"account":     account,
			"endpoint":    endpoint,
			"statuscode":  resp.StatusCode,
			"response":    "",
		}).Debug("OK")
		return response
	}
	return response
}

//PostSitoo - Function to POST data to Sitoo
func PostSitoo(baseURL string, endpoint string, account string, password string, json []byte) []byte {
	accountSplit := strings.Split(account, "-")
	accountNo := accountSplit[0]

	req, err := http.NewRequest("POST", baseURL+accountNo+endpoint, bytes.NewBuffer(json))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Basic "+BasicAuth(account, password))
	resp, err := http.DefaultClient.Do(req)

	log.WithFields(log.Fields{
		"requesttype": "POST",
		"account":     account,
		"endpoint":    endpoint,
		"body":        string(json),
	}).Debug("Request sent")

	if err != nil {
		log.WithFields(log.Fields{
			"requesttype": "POST/Error",
			"account":     account,
			"endpoint":    endpoint,
			"body":        json,
			"response":    err,
		}).Error("ERROR")
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		if resp.StatusCode == 429 {
			xRateLimitReset := resp.Header.Get("X-Rate-Limit-Reset")
			sleepTime, _ := strconv.Atoi(xRateLimitReset)
			time.Sleep(time.Duration(sleepTime) * time.Second)
			PostSitoo(baseURL, endpoint, account, password, json)
		} else {
			log.WithFields(log.Fields{
				"requesttype": "POST/Response",
				"account":     account,
				"endpoint":    endpoint,
				"statuscode":  resp.StatusCode,
				"response":    string(response),
			}).Error("ERROR")
			os.Exit(1)
		}

	} else {
		log.WithFields(log.Fields{
			"requesttype": "POST/Response",
			"account":     account,
			"endpoint":    endpoint,
			"statuscode":  resp.StatusCode,
			"response":    string(response),
		}).Debug("OK")
		return response
	}
	return response
}

//PutSitoo - Function to PUT data to Sitoo
func PutSitoo(baseURL string, endpoint string, account string, password string, json []byte) []byte {
	accountSplit := strings.Split(account, "-")
	accountNo := accountSplit[0]

	req, err := http.NewRequest("PUT", baseURL+accountNo+endpoint, bytes.NewBuffer(json))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Basic "+BasicAuth(account, password))
	resp, err := http.DefaultClient.Do(req)

	log.WithFields(log.Fields{
		"requesttype": "PUT",
		"account":     account,
		"endpoint":    endpoint,
		"body":        string(json),
	}).Debug("Request sent")

	if err != nil {
		log.WithFields(log.Fields{
			"requesttype": "PUT/Error",
			"account":     account,
			"endpoint":    endpoint,
			"body":        json,
			"response":    err,
		}).Error("ERROR")
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		if resp.StatusCode == 429 {
			xRateLimitReset := resp.Header.Get("X-Rate-Limit-Reset")
			sleepTime, _ := strconv.Atoi(xRateLimitReset)
			time.Sleep(time.Duration(sleepTime) * time.Second)
			PutSitoo(baseURL,endpoint,account,password,json)
		} else {
			log.WithFields(log.Fields{
				"requesttype": "PUT/Response",
				"account":     account,
				"endpoint":    endpoint,
				"statuscode":  resp.StatusCode,
				"response":    string(response),
			}).Error("ERROR")
		}
	} else {
		log.WithFields(log.Fields{
			"requesttype": "PUT/Response",
			"account":     account,
			"endpoint":    endpoint,
			"statuscode":  resp.StatusCode,
			"response":    string(response),
		}).Debug("OK")
		return response
	}
	return response
}

//DeleteSitoo - Function to DELETE data to Sitoo
func DeleteSitoo(baseURL string, endpoint string, account string, password string) []byte {
	accountSplit := strings.Split(account, "-")
	accountNo := accountSplit[0]

	req, err := http.NewRequest("DELETE", baseURL+accountNo+endpoint, nil)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Basic "+BasicAuth(account, password))
	resp, err := http.DefaultClient.Do(req)

	log.WithFields(log.Fields{
		"requesttype": "DELETE",
		"account":     account,
		"endpoint":    endpoint,
	}).Debug("Request sent")

	if err != nil {
		log.WithFields(log.Fields{
			"requesttype": "DELETE/Error",
			"account":     account,
			"endpoint":    endpoint,
			"response":    err,
		}).Error("ERROR")
	}
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		if resp.StatusCode == 429 {
			xRateLimitReset := resp.Header.Get("X-Rate-Limit-Reset")
			sleepTime, _ := strconv.Atoi(xRateLimitReset)
			time.Sleep(time.Duration(sleepTime) * time.Second)
			DeleteSitoo(baseURL,endpoint,account,password)
		} else {
			log.WithFields(log.Fields{
				"requesttype": "DELETE/Response",
				"account":     account,
				"endpoint":    endpoint,
				"statuscode":  resp.StatusCode,
				"response":    string(response),
			}).Error("ERROR")
		}
	} else {
		log.WithFields(log.Fields{
			"requesttype": "DELETE/Response",
			"account":     account,
			"endpoint":    endpoint,
			"statuscode":  resp.StatusCode,
			"response":    string(response),
		}).Debug("OK")
		return response
	}
	return response
}
