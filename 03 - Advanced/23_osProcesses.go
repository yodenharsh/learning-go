package main

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("echo", "Hello there")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
	fmt.Println(string(output))

	cmd = exec.Command("grep", "foo")
	cmd.Stdin = strings.NewReader("Hi there\nfoo\nbar")
	output, err = cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
	fmt.Println(string(output))

	cmd = exec.Command("sleep", "1.5")
	err = cmd.Start()

	if err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error waiting for command:", err)
		return
	}

	fmt.Println("sleep process complete")

	cmd = exec.Command("sleep", "60")
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	err = cmd.Process.Kill()
	if err != nil {
		fmt.Println("Error killing command:", err)
		return
	}
	fmt.Println("long 60 second sleep process killed")

	cmd = exec.Command("printenv", "SHELL")
	output, err = cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
	fmt.Printf("SHELL environment variable: %s", output)

	pr, pw := io.Pipe()
	cmd = exec.Command("grep", "foo")
	cmd.Stdin = pr

	go func() {
		defer pw.Close()
		pw.Write([]byte("hello there\nfoo\nbaz"))
	}()

	output, err = cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
	fmt.Printf("Output from grep: %s", output)

	cmd = exec.Command("ls", "-l")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
	fmt.Printf("Output from ls -l:\n%s", output)
}
