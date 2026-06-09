package inputs

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	directory = ".saves/"
	fileName = "save.csv"
	filePath = directory + fileName
	valueSeparator = ","
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
	} else if err == nil {
		file, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Failed to open save file")
		}
	}
	defer file.Close()
	var csvLine strings.Builder
	for i := range m.inputs {
		csvLine.WriteString(m.inputs[i].Value())
		csvLine.WriteString(valueSeparator)
	}
	csvLine.WriteString(time.Now().String())
	_, err = fmt.Fprintln(file, csvLine.String())
	if err != nil {
		fmt.Println(err)
	}
}

func (m model) loadLatestSavedDate() {
	os.MkdirAll(directory, 0755)
	_, err := os.Stat(filePath)
	var file *os.File
	if os.IsNotExist(err) {
		// take a deep breadth, it's okay to skip
		return
	} else if err == nil {
		file, err = os.Open(filePath)
		if err != nil {
			fmt.Println("Failed to open save file")
			return
		}
		defer file.Close()
		line := getLastLineFromFile(file)
		lineData := strings.Split(line, valueSeparator)
		if len(lineData) <= len(m.inputs) {
			return
		}

		for i := range m.inputs {
			m.inputs[i].SetValue(lineData[i])
		}
	}
}

func getLastLineFromFile(file *os.File) string {
	scanner := bufio.NewScanner(file)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
	}
	return line
}
