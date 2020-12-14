package ledcontrol

import (
	"io/ioutil"
	"strconv"
)

const ID = "led-control"

const FriendlyName = "Led Control"

const serviceName = "led-control-led-control"

type Led struct {
	Name string `json:"led"`
	Pin  string `json:"pin"`
}

type LedStatus struct {
	Name  string `json:"led"`
	Power bool   `json:"power"`
}

var leds = map[string]Led{
	"blue": {
		Name: "blue",
		Pin:  "/sys/class/gpio/gpio76/value",
	},
	"yellow": {
		Name: "yellow",
		Pin:  "/sys/class/gpio/gpio77/value",
	},
}

func SetLed(setLed LedStatus) bool {
	var value []byte
	if !setLed.Power {
		value = []byte(strconv.Itoa(0))
	} else {
		value = []byte(strconv.Itoa(1))
	}
	err := ioutil.WriteFile(leds[setLed.Name].Pin, value, 0644)
	if err != nil {
		return false
	}
	return true
}

func GetLedStatus(getLed string) LedStatus {
	var ledStatus LedStatus
	ledStatus.Name = leds[getLed].Name
	dat, err := ioutil.ReadFile(leds[getLed].Pin)
	if err != nil {
		switch int(dat[0]) {
		case 0:
			ledStatus.Power = false
		case 1:
			ledStatus.Power = true
		}
	}
	return ledStatus
}
