package functions

import (
	"fmt"
	"github.com/filirnd/owbot/models/telegram"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"
)

const dhcpFile = "/tmp/dhcp.leases"

var functionMap = make(map[string]func(update telegram.TgUpdate) (string, error))

// FUNCTIONS MANAGEMENT

func InitFunctions() {
	functionMap["/clients"] = getClients
	functionMap["/reboot"] = reboot
}

// Return error if not exists function or other function errors
func ExecuteFunction(name string, update telegram.TgUpdate) (string, error) {
	if fun, exist := functionMap[name]; exist {
		return fun(update)
	} else {
		return "", fmt.Errorf("Cannot found this command " + name)
	}
}

// FUNCTIONS

func reboot(update telegram.TgUpdate) (string,error) {
	cmd := exec.Command("reboot")
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("Reboot error " + err.Error())
		return "",err
	}

	return "", nil
}

func getClients(update telegram.TgUpdate) (string, error) {
	l, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	toRet := "Connected clients:\n"
	macList := make([]string, 0)
	for _, iface := range l {
		cmd := exec.Command("iw", "dev", iface.Name, "station", "dump")
		cmdOut, err := cmd.Output()
		if err != nil {
			fmt.Println(iface.Name + " error " + err.Error())
		} else {
			stationsString := string(cmdOut)
			stationsLines := strings.Split(stationsString, "\n")
			for _, line := range stationsLines {
				if strings.Contains(line, "Station ") {
					splitStationMac := strings.Split(line, " ")
					mac := splitStationMac[1]
					macList = append(macList, mac)
				}
			}
		}
	}
	if len(macList) != 0 {
		dhcpBytes, err := ioutil.ReadFile(dhcpFile)
		if err != nil {
			return "", err
		}
		dhcpListString := string(dhcpBytes)
		dhcpList := strings.Split(dhcpListString, "\n")
		for _, mac := range macList {
			for _, dhcpLine := range dhcpList {
				if strings.Contains(dhcpLine, mac) {
					dhcpLineSplit := strings.Split(dhcpLine, " ")
					if dhcpLineSplit[3] != "*" {
						toRet += dhcpLineSplit[3] + " - " + dhcpLineSplit[2] + " - " + dhcpLineSplit[1] + "\n"
					}
				}
			}
		}
	}
	return toRet, nil
}
