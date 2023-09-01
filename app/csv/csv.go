package csv

import (
	"encoding/csv"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"log"
	"os"
	"sismo-datagroup-service/app/model"
	"time"
)

func ParseCSV(path string) []model.DataGroupRecord {
	// Open CSV file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// Create a CSV reader
	reader := csv.NewReader(file)
	var list []model.DataGroupRecord
	// Iterate over CSV records
	for {
		// Read each record from CSV
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		fmt.Println("record[0]", record[0])
		fmt.Println("record[1]", record[1])
		doc := model.DataGroupRecord{
			Id:        primitive.NewObjectID(),
			Account:   record[0],
			Value:     record[1],
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			ExpiredAt: time.Now(),
			// Map additional fields from the record as needed
		}
		list = append(list, doc)
	}
	return list
}
