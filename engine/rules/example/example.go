package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/libopenstorage/autopilot/engine/rules"
)

func main() {
	reader, err := os.Open("./Prometheus-Node.csv")
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(reader)

	// Parse Headers, Create Map
	scanner.Scan()
	columns := strings.Split(scanner.Text(), ",")
	fmt.Println(columns)

	/*
		lastID := -1
		for scanner.Scan() {
			id, err := strconv.Atoi(columns[0])
			if err != nil {
				log.Fatalf("ParseInt: %v", err)
			}
			if id <= lastID {
				log.Fatal("😭")
			}
			lastID = id
		}
		if err := scanner.Err(); err != nil {
			log.Fatal("scanner: %v", err)
		}
		log.Println("👍")

		nodeCsv, err := os.Open("./Prometheus-Node.csv")
		if err != nil {
			fmt.Println(err)
		}
		r := csv.NewReader(nodeCsv)

		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(record)
		}
	*/
	rules.ReadRules()
}
