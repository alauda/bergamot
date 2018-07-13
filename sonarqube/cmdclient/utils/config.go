package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ParseConfigContent parse config content to map, config file content like
// name=abc
// age=14
// hosturl=localhost
// port=9009
func ParseConfigContent(filePath string) (map[string]string, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("Not found task data file %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Open %s error : ", err.Error())
	}
	defer file.Close()
	r := bufio.NewReader(file)
	taskData, err := parseLineData(r)
	if err != nil {
		return nil, fmt.Errorf("Parse Task Data Error： %s", err.Error())
	}
	return taskData, nil
}

func parseLineData(r *bufio.Reader) (map[string]string, error) {
	data := map[string]string{}
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return data, fmt.Errorf("parse line data: %s error : %s", line, err.Error())
		}
		line = strings.TrimRight(line, "\n")
		if strings.TrimSpace(line) == "" {
			continue
		}
		line = strings.TrimRight(line, "=")
		index := strings.Index(line, "=")
		if index == -1 {
			data[line] = ""
		} else {
			data[line[0:index]] = line[index+1:]
		}
	}
	return data, nil
}
