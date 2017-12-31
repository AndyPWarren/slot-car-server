# Leds on the RaspberryPi

A server that allows communications with the RaspberryPi GPIO pins, specifically to control the brightness of LED's. The server creates a socket so brightness change events can be streamed to the Pi from a client. 

## Build
First clone the project to your machine via the usual git commands

Then build a binary for the raspberry simply run:
```
make linuxarm
```
Then deploy it to the Pi using SSH:
```
make deploy
```

ssh into the raspberry pi. The binary will have been placed in the home directory.

Start the server, passing in arguments to configure the Led pins. Pin configuration is in the format `<color>=<bcm-pin-no>` with a space separating the pins, for example:
```
./leds.linux-arm red=17 yellow=27
```
The server will have started and be listening on port `9090`

An array of the he configured leds is available at `<raspberrypi-ip>:9090/leds`

Brightness commands can now be sent to the socket in the format `<color>=<brightness>` where color is the led you want to control and brightness is a number between 0 and 1. `0` = off, `1` = max brightness

