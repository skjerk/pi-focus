# Pi-focus

This application provides a web-based interface to control the focus
of a webcam on a Raspberry pi.
Very useful on e.g. OctoPrint, where autofocus is not desirable

## Building

The easiest way to build the application is to cross-compile it from a
desktop PC

Go can very easily cross-compile to the Raspberry Pi and to build pi-focus you need to install the Go tools on your pc/mac:

[https://golang.org/doc/install](https://golang.org/doc/install)

To to build the application for the
Raspberry Pi, clone the pi-focus files to your machine and in the pi-focus direcory use the following command:

`env GOOS=linux GOARCH=arm GOARM=7 go build`

When the compilation is done, you can transfer the application to your Pi using the deploy.sh script.

You will need to edit the deploy-script so that it uses the correct IP-address of your Pi

## Deployment

When you have done that, just enter the following command. It will ask for the password for the pi user on the Pi, which will be `raspberry` unless you have changed it.

`./deploy.sh`

## Execution

Now you can ssh to your Pi and start the application:

`./pi-focus`

When the application is running, you can enter the url <ip-of-raspberry-pi>:1080 (if using octopi, try octopi.local:1080) in your browser and then control the focus of your camera.

## Setting up automatic start

To enable pi-focus to always run on your pi, a systemd definition is provided.
You can activate it like this:
```
cd ~
sudo cp pi-focus.service /etc/systemd/system
sudo systemctl start pi-focus.service
sudo systemctl enable pi-focus
```

This will ensure that pi-focus always runs when you boot up your Pi.
