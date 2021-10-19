package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

// genRandomizeTime returns formatted working time, like '09:45:00'
func genRandomizeTime(hour int, min int) string {
	m := rand.Intn(20) + min
	return fmt.Sprintf("%02d:%02d:00", hour, m)
}

// nationalHolidays returns int array which means national holidays in target month
// Internally uses holidays-jp api(https://holidays-jp.github.io)
func nationalHolidays(current time.Time) []int {
	var ret []int
	api := "https://holidays-jp.github.io/api/v1/date.json"
	res, err := http.Get(api)
	if err != nil {
		fmt.Println("[Warning] Cannot get national holidays!!")
		return ret
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Closing http response is failed")
		}
	}(res.Body)

	var holidays map[string]string
	body, _ := ioutil.ReadAll(res.Body)
	if err != json.Unmarshal(body, &holidays) {
		fmt.Println("[Warning] Cannot get national holidays!! IO error")
		return ret
	}

	for k := range holidays {
		// key maybe formatted like 'yyyy-mm-dd'
		layout := "2006-01-02"
		t, e := time.Parse(layout, k)
		if e != nil {
			continue
		}
		if t.Year() == current.Year() && t.Month() == current.Month() {
			ret = append(ret, t.Day())
		}
	}
	return ret
}

func main() {
	// seed
	rand.Seed(time.Now().UnixNano())

	n := time.Now()
	y, m, _ := n.Date()

	nationalHolidays := nationalHolidays(n)

	var data []string
	for i := 1; i < 31; i++ {
		current := time.Date(y, m, i, 0, 0, 0, 0, time.Local)
		if current.Month() > m {
			break
		}
		if current.Weekday() == time.Sunday || current.Weekday() == time.Saturday {
			data = append(data, "\t\t")
			continue
		}
		isHoliday := false
		for _, d := range nationalHolidays {
			if d == current.Day() {
				isHoliday = true
			}
		}
		if isHoliday == true {
			data = append(data, "\t\t")
			continue
		}
		data = append(data, fmt.Sprintf("%s\t%s\t%s", genRandomizeTime(8, 40), genRandomizeTime(18, 0), "1:00:00"))
	}

	fmt.Println("Copied working times in this month, lets copy into Google Spreadsheet")
	err := clipboard.WriteAll(strings.Join(data, "\n"))
	if err != nil {
		fmt.Println("Coping into clipboard is failed")
		return
	}
}
