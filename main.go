package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
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

type Config struct {
	From struct {
		Hour int
		Min  int
	}
	To struct {
		Hour int
		Min  int
	}
}

var config Config

func readConfig() error {
	// viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return errors.New("cannot find config")
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return errors.New("invalid config")
	}
	return nil
}

func yearAndMonth(arg string) (time.Time, error) {
	var t time.Time
	if arg == "" {
		t = time.Now()
	} else {
		yearAndMonth := strings.Split(arg, "-")
		if len(yearAndMonth) != 2 {
			return t, errors.New("invalid argument. plz enter yyyy-mm for target month")
		}
		year, err := strconv.Atoi(yearAndMonth[0])
		if err != nil {
			return t, err
		}
		month, err := strconv.Atoi(yearAndMonth[1])
		if err != nil {
			return t, err
		}
		t = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	}
	println(t.Format(time.UnixDate))
	return t, nil
}

func main() {
	// seed
	rand.Seed(time.Now().UnixNano())

	err := readConfig()
	if err != nil {
		println(err.Error())
		return
	}

	arg := ""
	if len(os.Args[1:]) > 0 {
		arg = os.Args[1]
	}
	n, err := yearAndMonth(arg)
	if err != nil {
		println(err.Error())
		return
	}

	nationalHolidays := nationalHolidays(n)

	var data []string
	y := n.Year()
	m := n.Month()
	for i := 1; i <= 31; i++ {
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
		data = append(data, fmt.Sprintf("%s\t%s\t%s", genRandomizeTime(config.From.Hour, config.From.Min), genRandomizeTime(config.To.Hour, config.To.Min), "1:00:00"))
	}

	targetYM := n.Format("2006/01")
	fmt.Printf("Generated working times in %s, lets copy into Google Spreadsheet\n", targetYM)
	err = clipboard.WriteAll(strings.Join(data, "\n"))
	if err != nil {
		fmt.Println("Coping into clipboard is failed")
		return
	}
}
