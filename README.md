# JamHat-rpi

Simple program written in Go for the JamHat from ModMyPi.  The code is meant to be simple and readable to get a basic understanding of how to interace with the JamHat.  GPIO pin edge detection was not used for button presses, nor was PWM used to buzz the buzzer.

NOTE: As the buzzer uses GPIO Pin 20, PWM code not be used to make the buzzer buzz, so we just used basic on/off with sleeps.

Product from ModMyPi
<a href="https://www.modmypi.com/raspberry-pi/led-displays-and-drivers-1034/jam-hat" target="_blank">JAM HAT (LED & Buzzer Board)</a>

Python source code for JamHat
<a href="https://github.com/modmypi/Jam-HAT" target="_blank">https://github.com/modmypi/Jam-HAT</a>


## Usage ##

1) go get github.com/adrianh-za/jamhat-rpi
2) browse to $/go/src/github.com/adrianh-za/jamhat-rpi
3) sudo -E go run jamhat.go


## Acknowledgements ##

This project uses the GPIO library written and maintained by <b><a href="https://github.com/stianeikeland" target="_blank">Stian Eikeland</a></b>


## Gits ##

https://github.com/stianeikeland/go-rpio
