# Pi-focus

This application provides a web-based interface to control the focus
of a webcam on a Raspberry pi.
Very useful on e.g. OctoPrint, where autofocus is not desirable

## Building

The easiest way to build the application is to cross-compile it from a
desktop PC

Go can very easily cross-compile, and to build the application for the
Raspberry Pi, use the following command on your PC or Mac:

`env GOOS=linux GOARCH=arm GOARM=7 go build`

When the compilation is done, you can transfer the application to your Pi using the deploy.sh script.
You will need to edit the deploy-script so that it uses the correct IP-address of your Pi

## Deployment

When you have done that, just enter the following command. It will ask for the password for the pi user on the Pi, which will be `raspberry` unless you have changed it.

`./deploy.sh`

## Execution

Now you can ssh to your Pi and start the application:

`./pi-focus`

When the application is running, you can enter the url <ip-of-raspberry-pi>:1080 in your browser and then control the focus of your camera.

I use OctoPi to view the stream from the camera, but any streaming application should do the job.
