package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func main() {
	nodelete := flag.Bool("nodelete", false, "do not delete the 'bin' directory")
	flag.Parse()

	// Deferred removal of the 'bin' directory if it exists
	defer func() {
		if *nodelete {
			fmt.Printf("Skipping removal of 'bin' directory\n")
			return
		}
		fmt.Printf("Removing 'bin' directory\n")
		if _, err := os.Stat("bin"); err == nil {
			err = os.RemoveAll("bin") // Remove the 'bin' directory
			if err != nil {
				fmt.Printf("Failed to remove 'bin' directory: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("os.Stat(\"bin\") failed: %v\n", err)
		}
	}()

	// Run go list to find all main packages
	cmd := exec.Command("go", "list", "-f", `{{if eq .Name "main"}}{{.ImportPath}}{{end}}`, "./...")
	var out bytes.Buffer
	cmd.Stdout = &out
	fmt.Printf("Listing main packages...\n")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to list main packages: %v\n", err)
		os.Exit(1)
	}

	// Read the output and filter the packages
	var tobuild []string
	for _, line := range strings.Split(out.String(), "\n") {
		if line != "" && !strings.Contains(line, "buildall") {
			tobuild = append(tobuild, line)
		}
	}

	// semaphore to limit the number of concurrent builds
	smphr := make(chan struct{}, runtime.NumCPU())
	var wg sync.WaitGroup
	for _, pkg := range tobuild {
		wg.Add(1)
		go func(pkg string) {
			defer wg.Done()
			smphr <- struct{}{}        // acquire a semaphore
			defer func() { <-smphr }() // release the semaphore
			// put all executables in the same directory
			exePath := "bin/" + filepath.Base(pkg) + ".exe" // Put executables in a 'bin' directory

			cmd := exec.Command("go", "build", "-ldflags", "-H=windowsgui", "-o", exePath, pkg)
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Printf("%s failed: %v\n", pkg, err)
				return
			}
			fmt.Printf("%s ok\n", pkg)
		}(pkg)
	}
	wg.Wait()
}
