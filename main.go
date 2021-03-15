package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	i := 0
	for {
		log.Printf("iteration %d", i)
		i++
		cmd := exec.Command("docker", "run", "djs55/netperf:latest", "-H", "host.docker.internal", "-l", "600")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		if err := RestartDesktop(); err != nil {
			log.Fatal(err)
		}
	}
}
