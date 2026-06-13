package models

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
	csvLine.WriteString(m.currency.GetActiveValue())
	csvLine.WriteString(valueSeparator)
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

func loadLatestSavedDate() (string, []string) {
	os.MkdirAll(directory, 0755)
	_, err := os.Stat(filePath)
	var file *os.File
	if os.IsNotExist(err) {
		// take a deep breadth, it's okay to skip
		return "", []string{}
	} else if err == nil {
		file, err = os.Open(filePath)
		if err != nil {
			fmt.Println("Failed to open save file")
			return "", []string{}
		}
		defer file.Close()
		line := getLastLineFromFile(file)
		lineData := strings.Split(line, valueSeparator)
		if len(lineData) <= 2 {
			return "", []string{}
		}

		return lineData[0], []string{lineData[1], lineData[2]}
	}

	return "", []string{}
}

func getLastLineFromFile(file *os.File) string {
	scanner := bufio.NewScanner(file)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
	}
	return line
}
