package motorcontrol

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"../../config"
	"../../service"
)

const ID = "motor-control"
const FriendlyName = "Motor Control"

const serviceName = "motor-control-motor-control"

var motordFolder = config.GetDirectoryPathForHack(ID) + "/bin"
var eventFile = "event"
var positionFile = "position"
var statusFile = "status"

type Motorcontrol struct {
	*service.Service
	SaveConfig         (bool)
	Command            (bool)
	GetCurrentPosition (MotorControlPosition)
}

type MotorControlConfig struct {
	Enable bool `json:"enable"`
}

type MotorControlCommand struct {
	Command string      `json:"command"`
	Args    interface{} `json:"args"`
}

type MotorControlMove struct {
	Direction   string `json:"direction"`
	Orientation string `json:"orientation"`
	Steps       int    `json:"steps"`
}

type MotorControlPosition struct {
	PositionX int `json:"positionX"`
	PositionY int `json:"positionY"`
}

func GetConfiguration() MotorControlConfig {
	var currentConfig MotorControlConfig

	config.Load(ID, &currentConfig)

	return currentConfig
}

func GetServiceStatus() service.ServiceStatus {
	return service.Status(FriendlyName, service.Runit, serviceName)
}

func SaveConfig(newConfig MotorControlConfig) bool {
	var updatedconfig MotorControlConfig

	config.Load(ID, &updatedconfig)

	updatedconfig.Enable = newConfig.Enable

	success := config.Save(ID, updatedconfig)

	if !success {
		return false
	}

	if updatedconfig.Enable {
		config.EnableService(ID)
		service.Restart(service.Runit, serviceName)
	} else {
		config.DisableService(ID)
		service.Stop(service.Runit, serviceName)
	}
	return true
}

func Command(com MotorControlCommand) bool {
	switch com.Command {
	case "move":
		if args, ok := com.Args.(MotorControlMove); ok {
			motorMove(args)
			return true
		}
	case "goto":
		if args, ok := com.Args.(MotorControlPosition); ok {
			motorGoto(MotorControlPosition(args))
			return true
		}
	}
	return false
}

func motorMove(com MotorControlMove) bool {
	f, err := os.Create(motordFolder + "/" + eventFile)
	if err != nil {
		return false
	}
	_, err = f.WriteString(com.Orientation + " " + com.Direction + " " + strconv.Itoa(com.Steps))
	if err != nil {
		f.Close()
		return false
	}
	err = f.Close()
	if err != nil {
		return false
	}
	return true
}

func motorGoto(com MotorControlPosition) bool {
	f, err := os.Create(motordFolder + "/" + eventFile)
	if err != nil {
		return false
	}
	_, err = f.WriteString("goto " + strconv.Itoa(com.PositionX) + " " + strconv.Itoa(com.PositionY))
	if err != nil {
		f.Close()
		return false
	}
	err = f.Close()
	if err != nil {
		return false
	}
	return true
}

func GetCurrentPosition() MotorControlPosition {
	fmt.Println("[MotorControlPosition]	Getting motor position.")
	var currentPosition MotorControlPosition
	file, err := os.Open(motordFolder + "/" + positionFile)
	defer file.Close()
	if err != nil {
		fmt.Println("[MotorControlPosition]	Error reading getting motor position.")
	} else {
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)
		scanner.Scan()
		currentPosition.PositionX, _ = strconv.Atoi(scanner.Text())
		scanner.Scan()
		currentPosition.PositionY, _ = strconv.Atoi(scanner.Text())
	}
	fmt.Println("[MotorControlPosition]	Response: ", currentPosition)
	return currentPosition
}
