package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	cmd := exec.Command("go", "list", "-f", `{{if eq .Name "main"}}{{.ImportPath}}{{end}}`, "./...")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to list main packages: %v", err)
	}

	// Read the output and filter the packages
	var tobuild []string
	for _, line := range strings.Split(out.String(), "\n") {
		if line != "" && !strings.Contains(line, "x/buildall") {
			tobuild = append(tobuild, line)
		}
	}

	var wg sync.WaitGroup

	for _, pkg := range tobuild {
		wg.Add(1)
		go func(pkg string) {
			defer wg.Done()

			// Build the executable path
			exePath := strings.Replace(pkg, "x/", "", 1) + ".exe"

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
