package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	subCommand1 := flag.NewFlagSet("firstSub", flag.ExitOnError)
	subCommand2 := flag.NewFlagSet("secondSub", flag.ExitOnError)

	firstFlag := subCommand1.Bool("processing", false, "command processing status")
	secondFlag := subCommand1.Int("bytes", 1024, "Byte length")

	flagSc2 := subCommand2.String("language", "go", "Enter your language")

	if len(os.Args) < 2 {
		fmt.Println("This program requires additional commands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "firstSub":
		subCommand1.Parse(os.Args[2:])
		fmt.Println("Subcommand1: ")
		fmt.Println("Processing: ", *firstFlag)
		fmt.Println("Bytes: ", *secondFlag)
	case "secondSub":
		subCommand2.Parse(os.Args[2:])
		fmt.Println("Subcommand2: ")
		fmt.Println("Language: ", *flagSc2)
	default:
		fmt.Println("No subcommand entered")
		os.Exit(1)
	}
}
