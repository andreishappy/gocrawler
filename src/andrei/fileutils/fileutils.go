package fileutils

import (
	"encoding/json"
	"os"
	"log"
)

func WriteToFileInJson(filename string, result interface{}) {
	j, err := json.MarshalIndent(result, "", "  ")
	if err != nil{
		log.Printf("Could not marshal json for %s", filename)
		return
	}

	file, fileError := os.Create(filename)
	defer file.Close()

	if fileError != nil {
		log.Printf("Could not open file %s", filename)
		return
	}

	file.Write(j)
	file.Sync()
	log.Printf("Wrote result to %s", filename)
}
