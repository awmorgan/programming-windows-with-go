package main

import (
	"log"
	"os"
	"os/exec"
	"sync"
)

func main() {
	tobuild := []string{
		"x/ch1/hello",
		"x/ch2/scrnsize",
		"x/ch3/hellowin",
		"x/ch4/sysmets1",
		"x/ch4/sysmets2",
	}

	var wg sync.WaitGroup
	for _, pkg := range tobuild {
		wg.Add(1)
		go func(pkg string) {
			defer wg.Done()

			// Run "go build -v" command and capture stderr
			cmd := exec.Command("go", "build", "-v", pkg)
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.Fatalf("Build failed for %s: %v", pkg, err)
			}
		}(pkg)
	}
	wg.Wait()
}
