package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const myId = 0
const botToken = ""

var offset int64 = 0

const lastOffsetFile = "lastOffset"
const sleepTimeoutSeconds = 5

func main() {
	fmt.Println("    *")
	fmt.Println(" *  *  *   -------------------------------------------------------")
	fmt.Println("* OWbot *   Made with <3 by Filirnd (https://github.com/filirnd/)")
	fmt.Println(" *  *  *   -------------------------------------------------------")
	fmt.Println("    *")
	fmt.Println("")
	err:=sendMsg(myId,"Router started!")
	if err != nil {
		fmt.Println("Error sending message "+err.Error())
	}
	getLastOffsetFromFile()
	initFunctions()
	go updatesLoop()

	for {
		time.Sleep(time.Second * 10)
	}
}

func updateLastOffset(newOffset int64) {
	f, err := os.OpenFile(lastOffsetFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(strconv.FormatInt(newOffset, 10))); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	offset = newOffset
}

func getLastOffsetFromFile() {
	if _, err := os.Stat(lastOffsetFile); err == nil { // File Exists
		byteLine, err := ioutil.ReadFile(lastOffsetFile)
		if err != nil {
			fmt.Println("Cannot read file " + lastOffsetFile + ". Error:" + err.Error())
		}
		line := string(byteLine)
		if line != "" {
			offset, err = strconv.ParseInt(line, 10, 64)
			if err != nil {
				offset = 0
				updateLastOffset(offset)
			}
		} else {
			offset = 0
			updateLastOffset(offset)
		}
	} else {
		offset = 0
		updateLastOffset(offset)
	}

}

func updatesLoop() {
	for {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		resp, err := http.Get(
			"https://api.telegram.org/bot" + botToken + "/" + "getUpdates?offset=" + strconv.FormatInt(offset, 10) + "&timeout=10",
		)
		if err != nil {
			fmt.Println("Get Updates Error " + err.Error())
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Get Updates Error " + err.Error())
		}
		tgUpdateResult := TgUpdateResult{}
		err = json.Unmarshal(body, &tgUpdateResult)
		if err != nil {
			fmt.Println("Get Updates Unmarshalling Error " + err.Error())
		}

		latestUpdate := TgUpdate{}
		for _, update := range tgUpdateResult.Result {
			if(update.UpdateId > offset){
				executeFunctions(update)
			}
			latestUpdate = update
		}
		if latestUpdate.UpdateId != offset {
			updateLastOffset(latestUpdate.UpdateId)
		}
		time.Sleep(time.Second * sleepTimeoutSeconds)
	}

}

func executeFunctions(update TgUpdate) {
	if update.Message.From.Id != myId {
		err := sendMsg(update.Message.From.Id, "This bot is private. You haven't access to this.")
		if err != nil {
			fmt.Println("Cannot send message. Error: " + err.Error())
		}
	} else {
		resp,err := executeFunction(update.Message.Text,update)
		if err != nil {
			err := sendMsg(update.Message.From.Id, err.Error())
			if err != nil {
				fmt.Println("Cannot send message. Error: " + err.Error())
			}
		}else{
			err := sendMsg(update.Message.From.Id, resp)
			if err != nil {
				fmt.Println("Cannot send message. Error: " + err.Error())
			}
		}
	}
}

func sendMsg(chatID int64, msg string) error {
	// Convert our custom type into jso	n format
	reqBytes := []byte(fmt.Sprintf("{\"chat_id\":\"%d\", \"text\":\"%s\"}", chatID, msg))

	// Make a request to send our message using the POST method to the telegram bot API
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Post(
		"https://api.telegram.org/bot"+botToken+"/"+"sendMessage", "application/json",
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status" + resp.Status)
	}
	return err
}
