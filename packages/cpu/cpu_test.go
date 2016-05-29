package cpu

import (
	"os"
	"testing"

	i3 "github.com/denbeigh2000/goi3bar"

	"fmt"
)

func TestMain(m *testing.M) {
	loadSource = testingLoadFunc

	os.Exit(m.Run())
}

func TestBasicFormat(t *testing.T) {
	inst := &Cpu{
		WarnThreshold: 0.75,
		CritThreshold: 0.85,
	}

	testingVal = loadInfo{0.1, 0.1, 0.1}
	testingErr = nil

	val, err := inst.Generate()
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}

	for i, out := range val {
		if out.Color != i3.DefaultColors.OK {
			t.Errorf("Colors incorrect: Expected %v, got %v", i3.DefaultColors.OK, out.Color)
		}

		if out.FullText != "0.10" {
			t.Errorf("Values incorrect: Expected %v, got %v", "0.10", out.FullText)
		}

		if i < 2 && out.Separator {
			t.Errorf(
				"Incorrect use of separator on load #%v: Expected %v, got %v",
				i, false, out.Separator,
			)
		}

		if i == 2 && !out.Separator {
			t.Errorf("Expected separator after last load avg")
		}
	}

}

func TestThresholds(t *testing.T) {
	inst := Cpu{WarnThreshold: 0.75, CritThreshold: 0.85}

	testingErr = nil

	for _, in := range []struct {
		info   loadInfo
		colors []string
	}{
		{
			loadInfo{0.1, 0.1, 0.75},
			[]string{i3.DefaultColors.OK, i3.DefaultColors.OK, i3.DefaultColors.Warn},
		},
		{
			loadInfo{0.1, 0.75, 0.1},
			[]string{i3.DefaultColors.OK, i3.DefaultColors.Warn, i3.DefaultColors.OK},
		},
		{
			loadInfo{0.75, 0.1, 0.1},
			[]string{i3.DefaultColors.Warn, i3.DefaultColors.OK, i3.DefaultColors.OK},
		},
		{
			loadInfo{0.85, 0.75, 0.1},
			[]string{i3.DefaultColors.Crit, i3.DefaultColors.Warn, i3.DefaultColors.OK},
		},
		{
			loadInfo{0.6, 0.75, 0.9},
			[]string{i3.DefaultColors.OK, i3.DefaultColors.Warn, i3.DefaultColors.Crit},
		},
	} {
		testingVal = in.info

		result, _ := inst.Generate()
		for i, str := range result {
			if str.Color != in.colors[i] {
				t.Errorf("Incorrect color: expected %v, got %v", in.colors[i], str.Color)
			}
		}
	}
}

var testInst Cpu

func testingLoadFunc() (loadInfo, error) {
	return testingVal, testingErr
}

var testingVal loadInfo
var testingErr = fmt.Errorf("Something went wrong")
