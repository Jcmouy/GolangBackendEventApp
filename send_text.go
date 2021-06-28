package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func send(mobile string, otp int) {

	accountSid := "AC8ec5f9d8ec76e517ebf279f05b4bdeda"
	authToken := "c2678ebef0064591b7ad7c4c65feafd8"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	//mobile_without_first_letter := mobile[1:9]
	mobile_without_first_letter := mobile
	otp_string := strconv.Itoa(otp)

	otp_prefix := ":"

	//fmt.Println("+598" + mobile_without_first_letter)
	fmt.Println(mobile_without_first_letter)
	/*
		quotes := [5]string{"XDFZXFXDFXDFXDFXDF",
			"DQWEQWEQWEQWEQWEQW",
			"GTRGRTGRTGRTGRTGGT",
			"LKANSDLASLDKAKLJAD",
			"NVKJNDFGVNVNDFLGDL"}
	*/

	//quotes := "Hello! Welcome to EventApp. Your OPT is: " + otp_string

	quotes := "Hello! Welcome to EventApp. Your OPT is " + otp_prefix + " " + otp_string

	rand.Seed(time.Now().Unix())

	msgData := url.Values{}
	//msgData.Set("To", "+59899946874")
	//msgData.Set("To", "+598"+mobile_without_first_letter)
	msgData.Set("To", mobile_without_first_letter)
	msgData.Set("From", "+18566197912")
	//msgData.Set("Body", quotes[rand.Intn(len(quotes))])
	msgData.Set("Body", quotes)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println(resp.Status)
	}

}
