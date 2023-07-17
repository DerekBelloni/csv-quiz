package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var quizRecords = make([]Quiz, 0)

const duration = 10000 * time.Millisecond

type Quiz struct {
	Question string
	Answer int
}

func readCsvFile() [][]string {
	csv_file, err := os.Open("quiz_problems.csv")
	if err != nil {
		log.Fatal("Unable to read input file", err)
	}
	defer csv_file.Close()
	
	csvReader := csv.NewReader(csv_file)
	csvReader.Comma = ','
	csvReader.LazyQuotes = true
	csvReader.ReuseRecord = true
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Could not read from CSV file: ", err)
	}

	return records
}

func createQuiz(records [][]string) []Quiz {
		for _, record := range records {
			answerInt, _ := strconv.Atoi(record[1])
			quizRecords = append(quizRecords, Quiz{
				Question: record[0],
				Answer: answerInt,
			})
		}
	return quizRecords
}

func administerQuiz(quiz []Quiz, qChannel chan<- []int) {
	quizAnswers := make([]int, len(quiz) - 1)
	fmt.Printf("Thank you for deciding to take out quiz!\n")

	timer := time.NewTimer(duration)

	for i, question := range quiz[1:] {
		select {
		case <-timer.C:
			qChannel <- quizAnswers

		default:
			fmt.Println("Question: ", question.Question)
			fmt.Println("Your answer: ")
			fmt.Scan(&quizAnswers[i])
		}
	}

	qChannel <- quizAnswers
}

func checkAnswers(quizAnswers []int) int {
	quizScore := 0
	for i := 1; i < len(quizRecords); i++ {
		if quizAnswers[i - 1] == quizRecords[i].Answer {
			quizScore++
		}
	}

	return quizScore
}


func main() {
	qChannel := make(chan []int)

	records := readCsvFile()
	quiz := createQuiz(records)


	go administerQuiz(quiz, qChannel)

	quizAnswers := <- qChannel
	quizScore := checkAnswers(quizAnswers)
	fmt.Printf("Your score was %v/13", quizScore)

	close(qChannel)
}

