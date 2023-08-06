package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

/// e2e level

func Test_main(t *testing.T) {
	tests := []struct {
		targetDate string
		arrSkipped []bool // その日がスキップ対象(祝日or休日)かどうか
	}{
		{
			targetDate: "2023-02",
			arrSkipped: []bool{
				false, false, false, true,
				true, false, false, false, false, false, true,
				true, false, false, false, false, false, true,
				true, false, false, false, true, false, true,
				true, false, false,
			},
		},
		{
			targetDate: "2023-07",
			arrSkipped: []bool{
				true,                                          // 1st week
				true, false, false, false, false, false, true, // 2nd week
				true, false, false, false, false, false, true, // 3rd week
				true, true, false, false, false, false, true, // 4th week(includes 海の日)
				true, false, false, false, false, false, true, // 5th week
				true, false, // 6th week
			},
		},
	}
	// setup config
	_, err := os.Stat("./config.yml")
	if err != nil {
		data := "from:\n  hour: 9\n  min: 40\nto:\n  hour: 19\n  min: 0"
		err := ioutil.WriteFile("config.yml", []byte(data), 644)
		if err != nil {
			fmt.Println(err)
		}
	}
	for _, tt := range tests {
		t.Run(tt.targetDate, func(t *testing.T) {
			os.Args = []string{"cmd", tt.targetDate}
			main()
			fromClipboard, _ := clipboard.ReadAll()
			result := strings.Split(fromClipboard, "\n")
			if len(result) != len(tt.arrSkipped) {
				t.Fatalf("length not matched: %s", tt.targetDate)
			}
			for i, expected := range tt.arrSkipped {
				if (result[i] == "\t\t") != expected {
					t.Fatalf("not matched: %s %d %v", result[i], i, expected)
				}
			}
		})
	}

	// TODO: ほんとは tear down 的に config 消したほうがいいんだが，僕以外で特に困らないので，いったんこのままとする・・
}

/// UT

func Test_genRandomizeTime(t *testing.T) {
	type args struct {
		hour int
		min  int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid",
			args: args{hour: 0, min: 0},
		},
	}
	r := regexp.MustCompile("[0-2][0-9]:[0-5][0-9]:00")
	rand.Seed(time.Now().UnixNano())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := genRandomizeTime(tt.args.hour, tt.args.min)
			m := r.MatchString(got)
			if m != true {
				t.Errorf("genRandomizeTime() = %v", got)
			}
		})
	}
}
