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

			// Build the executable path
			exePath := pkg + ".exe"

			// Run "go build -v" command specifying output directly
			cmd := exec.Command("go", "build", "-v", "-o", exePath, pkg)
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.Printf("Build failed for %s: %v", pkg, err)
				return
			}
			log.Printf("Successfully built %s", pkg)

			// If build was successful, remove the executable
			if err := os.Remove(exePath); err != nil {
				log.Printf("Failed to remove %s: %v", exePath, err)
			} else {
				log.Printf("Removed %s", exePath)
			}
		}(pkg)
	}
	wg.Wait()
}
