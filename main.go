package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func ReadCsvFile(filename string) []problem {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("error on reading ", err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	rec, err := r.ReadAll()
	if err != nil {
		log.Fatal("error on parsing csv ", err)
	}
	return parseLines(rec)
}

func parseLines(lines [][]string) []problem {
	records := make([]problem, len(lines))
	for i, l := range lines {
		records[i] = problem{
			q: l[0],
			a: l[1],
		}
	}
	return records
}

func checkValidAnswer(answerSubmitted string, realAnswer string) bool {
	return answerSubmitted == realAnswer
}

type problem struct {
	q, a string
}

func readFlags() (*string, *int) {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	return csvFilename, timeLimit
}

func main() {
	csvFilename, timeLimit := readFlags()
	records := ReadCsvFile(*csvFilename)
	var rightAnswers int
	var pressKey string
	fmt.Println("press any key to start...")
	fmt.Scan(&pressKey)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
problemloop:
	for i, v := range records {
		fmt.Print("Problem #", i+1, ": ", v.q+"=")
		answerCh := make(chan string)
		go func() {
			var test string
			fmt.Scan(&test)
			answerCh <- test
		}()
		select {
		case <-timer.C:
			fmt.Println("\nrun out of time :(")
			break problemloop
		case answer := <-answerCh:
			if checkValidAnswer(answer, v.a) {
				rightAnswers++
			}
		}
	}
	fmt.Println("you scored ", rightAnswers, "out of ", len(records))
}
