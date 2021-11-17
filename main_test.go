package main

import (
	"math/rand"
	"regexp"
	"testing"
	"time"
)

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
