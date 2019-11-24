package main

import (
	"os/exec"
)

func openPowerSocket(){
	var path string = "https://maker.ifttt.com/trigger/open/with/key/" + IFTT_KEY
	cmd := exec.Command("curl", "-X", "GET", path)
	_ = cmd.Run()
}

func closePowerSocket(){
	var path string = "https://maker.ifttt.com/trigger/close/with/key/" + IFTT_KEY
	cmd := exec.Command("curl", "-X", "GET", path)
	_ = cmd.Run()
}

