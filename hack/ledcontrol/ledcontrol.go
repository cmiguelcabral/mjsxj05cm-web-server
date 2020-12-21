package ledcontrol

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/stianeikeland/go-rpio"
)

const ID = "led-control"

const FriendlyName = "Led Control"

const serviceName = "led-control-led-control"

type Led struct {
	Name  string `json:"led"`
	Pin   string `json:"pin"`
	Power bool   `json:"power"`
	Error string `json:"error"`
}

var leds = map[string]Led{
	"blue": {
		Name: "blue",
		Pin:  "/sys/class/gpio/gpio76",
	},
	"yellow": {
		Name: "yellow",
		Pin:  "/sys/class/gpio/gpio77",
	},
}

func SetLed(setLed Led) bool {
	var led Led = leds[setLed.Name]
	var valueString string
	if !led.Power {
		valueString = strconv.Itoa(0)

	} else {
		valueString = strconv.Itoa(1)
	}
	err := ioutil.WriteFile(led.Pin+"/value", []byte(valueString), 0644)
	fmt.Println("[Led] Changing", led.Name, "Led state to", led.Power)
	if err != nil {
		fmt.Println("Error changing Led.")
		return false
	}
	return true
}

func BlinkLed(led string, count int) {
	err := rpio.Open()
	if err != nil {
		os.Exit(1)
	}
	defer rpio.Close()

	pin := rpio.Pin(76)
	pin.Mode(rpio.Pwm)
	pin.Freq(64000)
	pin.DutyCycle(0, 32)
	// the LED will be blinking at 2000Hz
	// (source frequency divided by cycle length => 64000/32 = 2000)

	// five times smoothly fade in and out
	for i := 0; i < count; i++ {
		for i := uint32(0); i < 32; i++ { // increasing brightness
			pin.DutyCycle(i, 32)
			time.Sleep(time.Second / 32)
		}
		for i := uint32(32); i > 0; i-- { // decreasing brightness
			pin.DutyCycle(i, 32)
			time.Sleep(time.Second / 32)
		}
	}
}

func GetLedStatus(getLed string) Led {
	var led Led
	led = leds[getLed]
	fmt.Println("[Led] Getting %v Led state.", led.Name)
	file, err := os.Open(led.Pin + "/value")
	defer file.Close()
	if err != nil {
		fmt.Println("[Led] Error getting Led state.")
		led.Error = "Error reading Led state."
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			switch scanner.Text() {
			case "0":
				led.Power = false
			case "1":
				led.Power = true
			}
		}
	}
	fmt.Println("[Led] Response: %+v", led)
	return led
}
