package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"maploader/config"
	"maploader/robot"
	"maploader/tar"
	"maploader/util"
	"path/filepath"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	stateTopic      string
	currentMapTopic string
	saveTopic       string
	loadTopic       string
	setTopic        string
)

var maploaderDir string
var currentMap string

var rotationKeepMaps int

var messageCurrentMapTopicHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	if string(msg.Payload()) != currentMap {
		log.Printf("Loaded current map from status topic")
		currentMap = string(msg.Payload())
	}
}

var messageSaveTopicHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	saveMap(client, string(msg.Payload()))
}

var messageLoadTopicHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	loadMap(client, string(msg.Payload()))
}

var messageSetTopicHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	saveMap(client, currentMap)
	loadMap(client, string(msg.Payload()))
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
	publishState(client, "idle")

	subscriptions := []struct {
		Topic   string
		Handler mqtt.MessageHandler
	}{
		{currentMapTopic, messageCurrentMapTopicHandler},
		{saveTopic, messageSaveTopicHandler},
		{loadTopic, messageLoadTopicHandler},
		{setTopic, messageSetTopicHandler},
	}
	for _, sub := range subscriptions {
		token := client.Subscribe(sub.Topic, 1, sub.Handler)
		token.Wait()
		log.Printf("Subscribed to topic: %s", sub.Topic)
	}

	if currentMap == "" {
		currentMap = config.Getenv("DEFAULT_MAP_NAME", "main")
	}
	publishCurrentMap(client)
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

func initMqttTopics() {
	topicPrefix := config.MqttTopicPrefix()
	stateTopic = fmt.Sprintf("%s/%s/maploader/status", topicPrefix, config.MqttIdentifier())
	currentMapTopic = fmt.Sprintf("%s/%s/maploader/map", topicPrefix, config.MqttIdentifier())
	saveTopic = currentMapTopic + "/save"
	loadTopic = currentMapTopic + "/load"
	setTopic = currentMapTopic + "/set"

}

func main() {
	maploaderDir = config.Getenv("MAPLOADER_DIR", "/userdata/maploader")
	rotationKeepMaps = config.RotationKeepMaps()

	util.InitLogging(maploaderDir + "/log")
	config.InitConfig(config.Getenv("VALETUDO_CONFIG_PATH", "/userdata/valetudo_config.json"))
	initMqttTopics()

	err := robot.DetectRobot()
	if err != nil {
		log.Fatal(err)
	}

	var broker = config.MqttHost()
	var port = config.MqttPort()
	var protocol = "tcp"

	opts := mqtt.NewClientOptions()
	if config.MqttTLSEnabled() {
		protocol = "ssl"
		opts.TLSConfig = new(tls.Config)
		opts.TLSConfig.InsecureSkipVerify = config.MqttTLSIgnoreCertificateErrors()
		if len(config.MqttTLSCA()) > 0 {
			opts.TLSConfig.RootCAs = x509.NewCertPool()
			opts.TLSConfig.RootCAs.AppendCertsFromPEM([]byte(config.MqttTLSCA()))
		}
	}
	opts.AddBroker(fmt.Sprintf("%s://%s:%d", protocol, broker, port))

	opts.SetClientID(fmt.Sprintf("maploader_%s", config.MqttIdentifier())) // Valetudo uses a generated ID but that's not available to us
	opts.SetUsername(config.MqttUsername())
	opts.SetPassword(config.MqttPassword())
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	opts.WillTopic = stateTopic
	opts.WillEnabled = true
	opts.WillPayload = []byte("offline")
	opts.WillRetained = true

	client := mqtt.NewClient(opts)
	retry := time.NewTicker(5 * time.Second)
	for range retry.C {
		token := client.Connect()
		token.Wait()
		error := token.Error()
		if error != nil {
			log.Printf("MQTT connection failed: %s\n", error.Error())
		} else {
			retry.Stop()
			break
		}
	}
	for {
		time.Sleep(time.Second)
	}
}

func publishCurrentMap(client mqtt.Client) {
	token := client.Publish(currentMapTopic, 0, true, currentMap)
	token.Wait()
	time.Sleep(time.Second)
}

func publishState(client mqtt.Client, status string) {
	token := client.Publish(stateTopic, 0, true, status)
	token.Wait()
	time.Sleep(time.Second)
}

// Save the current map as the given name
func saveMap(client mqtt.Client, mapName string) {
	if mapName == "" {
		mapName = currentMap
	}

	publishState(client, "saving_map")
	log.Printf("Saving current map as %s\n", mapName)

	err := util.RotateFile(rotationKeepMaps, filepath.Join(maploaderDir, mapName), "tar.gz")
	checkAndHandleErrorWithMqtt(err, client)

	err = tar.Tar(fmt.Sprintf("%s/%s.tar.gz", maploaderDir, mapName), robot.CurrentRobot.MapFiles()...)
	checkAndHandleErrorWithMqtt(err, client)

	currentMap = mapName
	publishCurrentMap(client)
	publishState(client, "idle")
}

// TODO: Handle slots, and backups of slots, this tar.gz business might not work so well
// Idea:
// We make the assumption that 2 slots are always the same (1 and something else). (To be tested)
// We can restore a tar gz to 2 slots with the same date, that would be our current slot selection.
// We might have to change `map_id` on the other slot, to its corresponding number, it seems slot 1 will fix itself.
// This might only happen on reboot, but we will see.
//
// We should also use its own backup files where possible, but lets start listing those in general first.
// So basically there would be 3 things this would handle, Slots, internal backups and the maploader backups.
// Tiny bit complicated, but we'll get there.

// Load the map of the given name
func loadMap(client mqtt.Client, mapName string) {
	if mapName == "" {
		mapName = currentMap
	}

	mapFileToLoadMatches, err := filepath.Glob(filepath.Join(maploaderDir, mapName+".tar.gz"))
	checkAndHandleErrorWithMqtt(err, client)

	publishState(client, "loading_map")
	log.Printf("Attempting to load map %s\n", mapName)
	log.Println("stopping processes")
	robot.StopProcesses()

	log.Println("removing current map files")
	// TODO: No not this
	//util.RemoveDirContents(robot.CurrentRobot.MapFolders...)
	// FIXME: not this either, we should take our slot in account, but this is a start
	// Going to try to get general maploader functionality working first, where it uses its own backups
	for _, mapFile := range robot.CurrentRobot.MapFiles() {
		// If this fails we have a problem maybe, maybe not
		RemoveFile(mapFile)
    }

	// This seems to be some cache it will come back after restarting the process
	RemoveFile(robot.CurrentRobot.MapToPlannerFile)

	if len(mapFileToLoadMatches) == 0 {
		log.Println("requested map does not exist, loading blank map")
	} else {
		log.Println("loading map from archive")
		// TODO: Hmmmmmmm, i'm not sure if i'm a fan of this, since we only have 1 map folder, this could probably be relative?
		err = tar.Untar(mapFileToLoadMatches[0], "/")
		checkAndHandleErrorWithMqtt(err, client)
	}

	log.Println("map load complete, syncing files")
	robot.ExcuteCmd("sync")

	log.Println("restarting processes")
	robot.StartProcesses()
	log.Printf("After process restart")
	currentMap = mapName
	publishCurrentMap(client)
	log.Printf("After map publish")
	publishState(client, "starting_services")
	log.Printf("Wait for processes")
	robot.WaitForProcesses()

	publishState(client, "idle")
	log.Printf("before sound playback")
	robot.PlayMapLoadedSound(currentMap)
}

func checkAndHandleErrorWithMqtt(err error, client mqtt.Client) {
	if err != nil {
		util.CheckAndHandleError(err)
		publishState(client, "error")
	}
}
