# Pi-focus

This application provides a web-based interface to control the focus
of a webcam on a Raspberry pi.
Very useful on e.g. OctoPrint, where autofocus is not desireable

## Building

The easiest way to build the application is to cross-compile it from a
desktop PC

Go can very easily cross-compile, and to build the application for the
Raspberry Pi, use the following command on your PC or Mac:

`env GOOS=linux GOARCH=arm GOARM=7 go build`

Then transfer the application to your Pi like this:

`tar zcvf pi-focus.tgz pi-focus public templates`

`scp pi-focus.tgz <ip-of-raspberry-pi>.local:` 

Now login to the Pi and unpack the application:

```shell
ssh pi@octopi.local
tar zxvf pi-focus.tgz
./pi-focus
```

