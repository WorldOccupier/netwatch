package inputs

import (
	"strings"
	"fmt"
	"os"
	"time"
)

var (
	directory = ".saves/"
	fileName = "save.csv"
	filePath = directory + fileName
)

func (m model) save() {
	os.MkdirAll(directory, 0755)
	_, err := os.Stat(filePath)
	var file *os.File
	if os.IsNotExist(err) {
		file, err = os.Create(filePath)
		if err != nil {
			fmt.Println("Failed to create save file")
		}
	} else {
		file, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Failed to open save file")
		}
	}
	defer file.Close()
	var csvLine strings.Builder
	for i := range m.inputs {
		csvLine.WriteString(m.inputs[i].Value())
		csvLine.WriteString(",")
	}
	csvLine.WriteString(time.Now().String())
	_, err = fmt.Fprintln(file, csvLine.String())
	if err != nil {
		fmt.Println(err)
	}
}

// TODO: load latest save file entry on startup
