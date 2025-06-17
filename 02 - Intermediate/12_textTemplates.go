package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/template"
)

func main() {
	// tmpl := template.New("example")

	tmpl, err := template.New("example").Parse("Welcome, {{.name}}! How are you doing?\n")
	if err != nil {
		panic(err)
	}

	// Define data for the welcome message template
	data := map[string]string{
		"name": "John",
	}

	if err = tmpl.Execute(os.Stdout, data); err != nil {
		panic(err)
	}

	tmpl = template.Must(template.New("example").Parse("Welcome, {{.name}}!\n"))
	data["name"] = "Nobody"
	tmpl.Execute(os.Stdout, data)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your name")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	templates := map[string]string{
		"welcome":      "Welcome, {{.name}}",
		"notification": "{{.name}}, you have a new notification: {{.notification}}",
		"error":        "Oops, an error occurred: {{.errorMessage}}!",
	}

	parsedTemplates := make(map[string]*template.Template)
	for name, tmp := range templates {
		parsedTemplates[name] = template.Must(template.New(name).Parse(tmp))
	}
}
