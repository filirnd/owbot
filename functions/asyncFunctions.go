package functions

import (
	"bufio"
	"fmt"
	"github.com/filirnd/owbot/models/config"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var asyncFunctionSlice = make([]func(msg *chan string, config config.Config ),0)


// ASYNC FUNCTIONS MANAGEMENT

func StartAsyncFunctions(msg *chan string, config config.Config) {

	// Init functions
	asyncFunctionSlice = append(asyncFunctionSlice, wifiClientsTriggerAsyncFunction)


	// Start functions
	for _, f := range asyncFunctionSlice {
		go f(msg , config )
	}
}


// CUSTOM ASYNC FUNCTIONS
// NOTE: Each functions defined need to have infinite loop with a sleep for remaining live.

func wifiClientsTriggerAsyncFunction(msg *chan string, config config.Config){
	if !config.Async.NewClient {
		fmt.Println("NewClient Async disabled.")
		return
	}
	cmd := exec.Command("/sbin/logread", "-e" ,"associated", "-f")
	// create a pipe for the output of the script
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			// Sending msg to senderChannel
			txt := scanner.Text()
			// Example text "Sat Jan 23 17:06:30 2021 daemon.info hostapd: wlan0: STA e0:dc:ff:f4:cd:f7 IEEE 802.11: associated (aid 8)"
			//Getting mac
			r, _ := regexp.Compile("([0-9a-fA-F]{2}[:]){5}([0-9a-fA-F]{2})")
			mac := r.FindStringSubmatch(txt)
			if len(mac)>0 {
				time.Sleep(time.Second *2 ) // Waiting two seconds, time for writing new client on dhcp file
				dhcpBytes, err := ioutil.ReadFile(dhcpFile)
				if err != nil {
					fmt.Println( err)
				}else {
					toRet := "Connected new client => "
					dhcpListString := string(dhcpBytes)
					dhcpList := strings.Split(dhcpListString, "\n")
					for _, dhcpLine := range dhcpList {
						if strings.Contains(dhcpLine, mac[0]) {
							dhcpLineSplit := strings.Split(dhcpLine, " ")
							if dhcpLineSplit[3] != "*" {
								toRet += dhcpLineSplit[3] + " - " + dhcpLineSplit[2] + " - " + dhcpLineSplit[1] + "\n"
								}
							}
					}
					 if len(toRet) == 0 {
					 	toRet = "Connected new client => "+mac[0]
					 }
					*msg <- toRet
				}
			}
		}
	}()
	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		return
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		return
	}
}