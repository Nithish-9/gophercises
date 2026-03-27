package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func main() {
	fileName := flag.String("file", "problems.csv", "CSV file name")
	timeLimit := flag.Int("limit", 30, "quiz time limit")
	flag.Parse()

	if !strings.Contains(*fileName, ".csv") {
		fmt.Println("Only CSV files are allowed")
		return
	}
	if *timeLimit <= 0 {
		fmt.Println("Quiz time limit should be positive")
		return
	}

	records := readCsvFile(*fileName)
	hash := make(map[string]string)
	for _, row := range records {
		_, ok := hash[row[0]]
		if !ok {
			hash[row[0]] = row[1]
		}
	}

	totalQuestions := len(hash)
	correct := 0
	wrong := 0
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()

quizLoop:
	for question, correctAnswer := range hash {
		fmt.Printf("%s : ", question)

		userChan := make(chan string)

		go func() {
			input, _ := reader.ReadString('\n')
			userAnswer := strings.TrimSpace(input)
			userChan <- userAnswer
		}()

		select {
		case <-timer.C:
			break quizLoop
		case userAnswer := <-userChan:
			if userAnswer == correctAnswer {
				correct++
			} else {
				wrong++
			}
		}

	}

	fmt.Println()
	fmt.Printf("\n")
	fmt.Printf("Total number of questions : %d\n", totalQuestions)
	fmt.Printf("Total number of correct answers : %d\n", correct)
	fmt.Printf("Total number of wrong answers : %d\n", wrong)

}
