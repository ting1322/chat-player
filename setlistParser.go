package main

import (
	"bufio"
	"encoding/json"
	"os"
	"regexp"
	"strconv"
)

type VideoTimeStamp struct {
	TimeInMs int `json:"time_in_ms"`
	Title string  `json:"title"`
}

func convertSetlist2Json(filename string) (string, error) {
	if len(filename) == 0 {
		return "", nil
	}
	var timeList []VideoTimeStamp = make([]VideoTimeStamp, 0)
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	re, err := regexp.Compile(`^\s*(\d+\. )?(?P<H>\d\d?):(?P<M>\d\d):(?P<S>\d\d)( ~ \d\d?:\d\d:\d\d)? (?P<T>.+)`)
	if err != nil {
		return "", err
	}
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		if err != nil {
			return "", err
		}
		matches := re.FindSubmatch(line)
		if len(matches) == 0 {
			continue
		}
		hour, _ := strconv.Atoi(string(matches[re.SubexpIndex("H")]))
		minute, _ := strconv.Atoi(string(matches[re.SubexpIndex("M")]))
		second, _ := strconv.Atoi(string(matches[re.SubexpIndex("S")]))
		title := string(matches[re.SubexpIndex("T")])
		total_ms := ((((hour * 60) + minute) * 60) + second) * 1000
		timeList = append(timeList, VideoTimeStamp{ total_ms, title })
	}
	jsondata, err := json.Marshal(timeList)
	return string(jsondata), nil
}