# Slot Car Server on the RaspberryPi

A server that facilitates communication with RaspberryPi GPIO pins, specifically to control slot cars but could be used for any application where GPIO pins need to have PWM levels streamed. The server creates a socket which clients can stream events to. 

## Build
First clone the project to your machine via the usual git commands

Build a binary for the raspberry simply run:
```
make linuxarm
```
Deploying SCP's the binary to the raspberrypi. By default it uses the `id_rsa` key and the current `$USER` environment variable. These can be overridden by exporting `PI_USER` and `PI_SSH_KEY` respectively. Then run:
```
make deploy
```
## Run
SSH into the raspberry pi. The binary will have been placed in the users home directory.

Start the server, passing in arguments to map the BCM pin number to a track number. This configuration is in the format `<track-number>=<bcm-pin-no>` with a space separating the pins, for example:
```
./slot-car.linux-arm 1=17 2=27
```

Events can now be sent to the socket server in the format `<track-number>=<level>` where track-number is the track/pin you want to control and level is a number between 0 and 1. `0` = off, `1` = max 

