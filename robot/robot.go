package robot

import (
	"errors"
	"log"
	"maploader/util"
	"os"
	"strings"
)

const serialFile = "/oem/midea/device.sn8"

type Robot struct {
	model            string
	Serials          []string
	MapFolder        string
	mapFiles         []string
	restartProcesses []Process
}

func (r *Robot) MapFilesAndFolders() []string {
	return append(r.MapFolders, r.mapFiles...)
}

var CurrentRobot *Robot

func DetectRobot() {
	hostname, err := os.Hostname()

	util.CheckAndHandleError(err)

	var possibleRobots []Robot

	for _, cRobot := range robots {
		if cRobot.model == hostname {
			possibleRobots = append(possibleRobots, cRobot)
		}
	}

	if len(possibleRobots) == 0 {
		log.Printf("Hostname is not 'midea', are you sure this is a Midea robot?")
    	CurrentRobot = nil, errors.New("Robot model could not be determined or is not supported")
        return
	}

	serial, err := os.ReadFile(serialFile)
	if err != nil {
		log.Printf("Could not read robot serial from %s: %v; are you sure this is a Midea robot?", serialFile, err)
		CurrentRobot = nil, errors.New("Robot model could not be determined or is not supported")
		return
	}

	for _, cRobot := range possibleRobots {
		for _, cSerial := range cRobot.Serials {
			if cSerial == serial {
				CurrentRobot = ValidateRobot(cRobot)
				log.Printf("Using configuration for model: %s", CurrentRobot.model)
				break
			}
		}
	    if (CurrentRobot != nil) {
	    	break
	    }
	}
}

func ValidateRobot(robot Robot) bool {

	if util.CheckDirectoryExists(cRobot.MapFolder) &&
		util.CheckFilesExists(cRobot.mapFiles...) {
		return true
	}

	return false
}
