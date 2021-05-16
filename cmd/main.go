package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/filirnd/owbot/functions"
	"github.com/filirnd/owbot/models/config"
	"github.com/filirnd/owbot/models/telegram"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const version = "0.1.0b"
const configFile = "resources/config.json"
const lastOffsetFile = "resources/lastOffset"
const sleepTimeoutSeconds = 2

var adminId int64
var botToken string

var offset int64 = 0

func main()  {
	start()
}



func start() {
	fmt.Println("")
	fmt.Println("   \\     /                   #### OWbot ####  v. "+version+"\n   _\\___/_\n /______ /|  Yet another telegram bot, but for your router.\n|_°_____|/   Made with <3 by Filirnd (https://github.com/filirnd/)")
	fmt.Println("")
	cfg,err := loadConfig()
	if err != nil {
		fmt.Println("Cannot load config. Error "+err.Error())
		os.Exit(-1)
	}
	err =sendMsg(adminId,"Router started!")
	if err != nil {
		fmt.Println("Error sending message "+err.Error())
	}
	getLastOffsetFromFile()

	msgChan := make(chan string)
	functions.StartAsyncFunctions(&msgChan,cfg)
	go asyncMessageSender(&msgChan)

	functions.InitFunctions()
	go updatesLoop()
	for {
		time.Sleep(time.Second * 10)
	}
}


/**
 Start all async functions for async messages from bot to clients
 */
func asyncMessageSender(msg *chan string){
	for {
		text := <- *msg
		err := sendMsg(adminId,text)
		if err!=nil {
			fmt.Println("Cannot send async msg. Error: "+err.Error())
		}
	}
}


func loadConfig() (config.Config,error){
	cfg,err := config.ConfigFromFile(configFile)
	if err != nil {
		return cfg,err
	}
	adminId = cfg.TgId
	botToken = cfg.TgBotToken
	return cfg,nil
}


/**
 Update lastOffest variable and file, so if process will be rebooted, don't read old processed messages.
 */
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


/**
 Read lastOffset from file and memorize in the variable
 */
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


// Start update loop for read messages from clients.
func updatesLoop() {
	for {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		resp, err := http.Get(
			"https://api.telegram.org/bot" + botToken + "/" + "getUpdates?offset=" + strconv.FormatInt(offset, 10) + "&timeout=10",
		)
 		if err != nil {
			fmt.Println("Get Updates Error " + err.Error())
		}else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Get Updates Error " + err.Error())
			}
			tgUpdateResult := telegram.TgUpdateResult{}
			err = json.Unmarshal(body, &tgUpdateResult)
			if err != nil {
				fmt.Println("Get Updates Unmarshalling Error " + err.Error())
			} else {
				for _, update := range tgUpdateResult.Result {
					if update.UpdateId > offset {
						updateLastOffset(update.UpdateId)
						executeFunctions(update)
					}
				}
			}

		}
		time.Sleep(time.Second * sleepTimeoutSeconds)
	}
}

// Execute function from clients commands
func executeFunctions(update telegram.TgUpdate) {
	if update.Message.From.Id != adminId {
		err := sendMsg(update.Message.From.Id, "This bot is private. You haven't access to this.")
		if err != nil {
			fmt.Println("Cannot send message. Error: " + err.Error())
		}
	} else {
		resp,err := functions.ExecuteFunction(update.Message.Text,update)
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


// Message Sender
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
