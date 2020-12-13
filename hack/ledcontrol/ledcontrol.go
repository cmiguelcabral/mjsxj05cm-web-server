package ledcontrol

import (
	"io/ioutil"
	"strconv"
)

const ID = "led-control"
const FriendlyName = "Led Control"

const serviceName = "led-control-led-control"

type Led struct {
	Led string `json:"led"`
	Pin string `json:"pin"`
}

type LedStatus struct {
	Led   string `json:"led"`
	Power bool   `json:"power"`
}

var leds = []Led{
	Led{
		Led: "blue",
		Pin: "/sys/class/gpio/gpio76/value",
	},
	Led{
		Led: "yellow",
		Pin: "/sys/class/gpio/gpio77/value",
	},
}

func SetLed(led LedStatus) bool {
	var value []byte
	if !led.Power {
		value := []byte(strconv.Itoa(0))
	} else {
		value := []byte(strconv.Itoa(1))
	}

	for _, v := range leds {
		if v.Led == led.Led {
			err := ioutil.WriteFile(v.Pin, value, 0644)
			if err != nil {
				return false
			}
		}
	}
	return true
}

func GetLedStatus(led string) LedStatus {
	var ledStatus LedStatus
	for _, v := range leds {
		if v.Led == led {
			ledStatus.Led = v.Led
			dat := ioutil.ReadFile(v.Pin)
			switch int(dat[0]) {
			case 0:
				ledStatus.Power = false
			case 1:
				ledStatus.Power = true
			}
		}
	}
	return ledStatus
}
