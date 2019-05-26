package main

import (
	"time"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	rpio "github.com/stianeikeland/go-rpio"
)

// Reference code in Python
// https://github.com/modmypi/Jam-HAT

const (
	PinRed1 int = 5
	PinRed2 int = 6
	PinOrange1 int = 12
	PinOrange2 int = 13
	PinGreen1 int = 16
	PinGreen2 int = 17
	PinButton1 int = 19
	PinButton2 int = 18
	PinBuzzer int = 20		//Can't use PWM on this pin, will need to buzz manually by looping in a Go routine
)

func main() {

	//Setup GPIO
	err := rpio.Open()
	if (err != nil) {
		return
	}
	defer rpio.Close()

	var leds = [6]int { PinRed1, PinRed2, PinOrange1, PinOrange2, PinGreen1, PinGreen2 }
	var ledPosition = -1	//Default to all off

	// Intercept CTRL-C and cleanup GPIO
	captureExit(leds)

	// Setup BLUE button
	var blueButtonPin = rpio.Pin(PinButton1)
	var blueButtonPressed = rpio.Low;
	blueButtonPin.Mode(rpio.Input)

	// Setup RED button
	var redButtonPin = rpio.Pin(PinButton2)
	var redButtonPressed = rpio.Low;
	redButtonPin.Mode(rpio.Input)

	// Setup BUZZER
	var buzzerPin = rpio.Pin(PinBuzzer)
	buzzerPin.Output()
	var buzzChannel = make(chan bool)		//Will be used in Go routine buzzMonitor()
	defer close(buzzChannel)

	// Kick off the buzzer go routine *this does the buzzing if RED button pressed or not)
	go buzzMonitor(buzzerPin, buzzChannel)

	// Setup LED pins as OUTPUTs
	for ledCount := 0; ledCount < 6; ledCount++ {
		pin := rpio.Pin(leds[ledCount])
		pin.Output()
	}

	// Lets do some processing
	for {

		//Get the press state of BLUE button
		var blueState = blueButtonPin.Read()
				
		//Check for BLUE button state change
		if (blueState != blueButtonPressed) {
			blueButtonPressed = blueState;  //Store new state of the BLUE button

			// Some feedback
			if (blueButtonPressed == rpio.High) {
				fmt.Println("Blue - Pressed")
			} else {
				fmt.Println("Blue - Released")
			}
			
			//Only change if button pressed down, not on release
			if (blueButtonPressed != rpio.Low) {

				//Calculate the next LED position
				ledPosition = ledPosition + 1;
				if (ledPosition > 5) {
					ledPosition = 0
				}

				//Clear all LEDs
				clearLEDs(leds)
 
				//Light up the current LED
				setLED(ledPosition, leds)
			}
		}

		//Get the press state of RED button
		var redState = redButtonPin.Read()

		//Check for RED button state change
		if (redState != redButtonPressed) {
			redButtonPressed = redState;  //Store new state of the RED button

			// Some feedback
			if (redButtonPressed == rpio.High) {
				fmt.Println("Red - Pressed")
				buzzChannel <- true			//Send the button state to the channel
			} else {
				fmt.Println("Red - Released")
				buzzChannel <- false		//Send the button state to the channel
			}
		}

		//Delay for a bit
		time.Sleep(50 * time.Millisecond)
	}
}

// clearLEDs turns off all LEDs on the JamHat
func clearLEDs(leds [6]int) {
	for ledCount := 0; ledCount < 6; ledCount++ {
		pin := rpio.Pin(leds[ledCount])
		pin.Low()
	}
}

// setLED turn on the specified LED on the JamHat
func setLED(position int, leds [6]int) {
	pin := rpio.Pin(leds[position])
	pin.High()
}

// buzzMonitor switches on buzzer if RED button pressed.
func buzzMonitor(pin rpio.Pin, buzzChannel <-chan bool) {

	var canBuzz = false

	for {
		
		select {
		
		case value := <- buzzChannel:
			canBuzz = value
			
		default:

			if (canBuzz)  {
				pin.High()
				time.Sleep(5 * time.Millisecond)
				pin.Low()
				time.Sleep(5 * time.Millisecond)
			} else {
				time.Sleep(200 * time.Millisecond)
			}
		}
	}
}

// captureExit captures Interrupt and SIGTERM signals to handle program exit
func captureExit(leds [6]int) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		for range signalChan {
			clearLEDs(leds)
			rpio.Close()
			os.Exit(1)
		}
	}()
}

