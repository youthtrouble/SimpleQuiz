package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file should be in the format of question, answer")
	timeLimit := flag.Int("limit", 40, "the time limit for the quiz is in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("failed to open the csv file %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("failed to parse the csv file")
	}
	problems := parselines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	<-timer.C

	correct := 0
problemLoop:
	for i, p := range problems {
		fmt.Printf("Problem number %d: %s = ", i+1, p.question)
		solnCh := make(chan string)
		go func() {
			var soln string
			fmt.Scanf("%s\n", &soln)
			solnCh <- soln
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemLoop
		case soln := <-solnCh:
			if soln == p.answer {
				correct++
			}

		}
	}
	fmt.Printf("You Scored %d out of %d correct\n", correct, len(problems))

}

func parselines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}

	}
	return ret
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
