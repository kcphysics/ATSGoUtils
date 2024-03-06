package main

import "log"

func main() {
	atsDataFilename := "./atsdata.json"
	atsData, err := ParseATSDataFromFile(atsDataFilename)
	if err != nil {
		panic(err)
	}
	log.Printf("%#v", atsData)
	// log.Println("Completed")
}
