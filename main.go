package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/saulhoward/martian-robots/robot"
)

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "i", "", "path of file containing robot instructions")
	flag.Parse()

	instructions, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	result, err := robot.RunRobots(string(instructions))
	if err != nil {
		log.Fatalf("failed to execute instructions: %v", err)
	}

	for _, p := range result {
		fmt.Println(p)
	}
}
