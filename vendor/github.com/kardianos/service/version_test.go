package service

import (
	"reflect"
	"testing"
)

func Test_versionCompare(t *testing.T) {
	type args struct {
		v1 []int
		v2 []int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"segment-mismatch-1", args{[]int{0, 0, 0, 0}, []int{0, 0, 0}}, 0, true},
		{"segment-mismatch-2", args{[]int{0, 0, 0}, []int{0, 0, 0, 0}}, 0, true},
		{"equal-to-1", args{[]int{0, 0, 0}, []int{0, 0, 0}}, 0, false},
		{"equal-to-2", args{[]int{0, 0, 1}, []int{0, 0, 1}}, 0, false},
		{"equal-to-3", args{[]int{0, 1, 0}, []int{0, 1, 0}}, 0, false},
		{"equal-to-4", args{[]int{1, 0, 0}, []int{1, 0, 0}}, 0, false},
		{"equal-to-5", args{[]int{2, 2, 2}, []int{2, 2, 2}}, 0, false},
		{"less-than-1", args{[]int{0, 0, 0}, []int{0, 0, 1}}, -1, false},
		{"less-than-2", args{[]int{0, 0, 1}, []int{0, 1, 0}}, -1, false},
		{"less-than-3", args{[]int{0, 1, 0}, []int{1, 0, 0}}, -1, false},
		{"less-than-4", args{[]int{0, 1, 0}, []int{1, 0, 0}}, -1, false},
		{"less-than-5", args{[]int{0, 8, 0}, []int{1, 5, 0}}, -1, false},
		{"greater-than-1", args{[]int{0, 0, 1}, []int{0, 0, 0}}, 1, false},
		{"greater-than-2", args{[]int{0, 1, 0}, []int{0, 0, 1}}, 1, false},
		{"greater-than-3", args{[]int{1, 0, 0}, []int{0, 1, 0}}, 1, false},
		{"greater-than-4", args{[]int{1, 0, 0}, []int{0, 1, 0}}, 1, false},
		{"greater-than-5", args{[]int{1, 5, 0}, []int{0, 8, 0}}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := versionCompare(tt.args.v1, tt.args.v2)
			if (err != nil) != tt.wantErr {
				t.Errorf("versionCompare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("versionCompare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseVersion(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"sanity-check", args{"0.0.0"}, []int{0, 0, 0}},
		{"should-fail", args{"0.zero.0"}, nil},
		{"should-fail-no-semver", args{"0.0.0-test+1"}, nil},
		{"double-digits", args{"5.200.1"}, []int{5, 200, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseVersion(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_versionAtMost(t *testing.T) {
	type args struct {
		version []int
		max     []int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"segment-mismatch-1", args{[]int{0, 0, 0, 0}, []int{0, 0, 0}}, false, true},
		{"segment-mismatch-2", args{[]int{0, 0, 0}, []int{0, 0, 0, 0}}, false, true},
		{"test-1", args{[]int{0, 0, 0}, []int{0, 0, 1}}, true, false},
		{"test-2", args{[]int{0, 1, 0}, []int{0, 1, 0}}, true, false},
		{"test-3", args{[]int{1, 0, 0}, []int{1, 0, 0}}, true, false},
		{"negative-test-1", args{[]int{0, 0, 1}, []int{0, 0, 0}}, false, false},
		{"negative-test-2", args{[]int{0, 1, 0}, []int{0, 0, 1}}, false, false},
		{"negative-test-3", args{[]int{1, 0, 0}, []int{0, 1, 0}}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := versionAtMost(tt.args.version, tt.args.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("versionAtMost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("versionAtMost() = %v, want %v", got, tt.want)
			}
		})
	}
}
