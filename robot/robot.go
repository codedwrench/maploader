package robot

import (
	"errors"
	"fmt"
	"log"
	"maploader/util"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type Robot struct {
	model            string
	Serials          []string
	MapFolder        string
	MapToPlannerFile string
	mapFiles         []string
	restartProcesses []Process
}

var CurrentRobot *Robot
const serialFile = "/oem/midea/device.sn8"
var slotRegex = regexp.MustCompile(`([a-zA-Z_0]+)\d+`)

func (r *Robot) MapFiles() []string {
	mapFullPaths := make([]string, len(r.mapFiles))

	for i, mapFile := range r.mapFiles {
    	mapFullPaths[i] = path.Join(r.MapFolder, mapFile)
	}

	return mapFullPaths
}

func (r *Robot) MapFilesForSlot(slot int) []string {
	mapFullPaths := make([]string, len(r.mapFiles))

	for i, mapFile := range r.mapFiles {
		mapFileForNewSlot := slotRegex.ReplaceAllString(mapFile, "${1}"+strconv.Itoa(slot))
    	mapFullPaths[i] = path.Join(r.MapFolder, mapFileForNewSlot)
	}

	return mapFullPaths
}

func DetectRobot() error {
	hostname, err := os.Hostname()

	util.CheckAndHandleError(err)

	var possibleRobots []Robot

	for _, cRobot := range robots {
		if cRobot.model == hostname {
			possibleRobots = append(possibleRobots, cRobot)
		}
	}

	if len(possibleRobots) == 0 {
		return fmt.Errorf("hostname %q is not 'midea', are you sure this is a Midea robot?", hostname)
	}

	rawSerial, err := os.ReadFile(serialFile)
	if err != nil {
		return fmt.Errorf("could not read robot serial from %s: %w; are you sure this is a Midea robot?", serialFile, err)
	}
	serial := strings.TrimSpace(string(rawSerial))

	for _, cRobot := range possibleRobots {
		for _, cSerial := range cRobot.Serials {
			if cSerial == serial && ValidateRobot(cRobot) {
				CurrentRobot = &cRobot
				log.Printf("Using configuration for model: %s", CurrentRobot.model)
				return nil
			}
		}
	}

	return errors.New("robot model could not be determined or is not supported")
}

func ValidateRobot(robot Robot) bool {
	return util.CheckDirectoryExists(robot.MapFolder) &&
		util.CheckFilesExists(robot.MapFiles()...)
}
