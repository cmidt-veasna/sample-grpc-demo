package query

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

type SampleFilter struct {
	Name      string
	Age       int
	Status    int
	CreatedAt string
	UpdatedAt string
}

var sampleFilters []*SampleFilter

func init() {
	status := 0
	createDate := ""
	updateDate := ""
	for i := 0; i < 100; i++ {
		switch {
		case i < 20:
			status = 1
			createDate = time.Now().Add(-time.Hour * 64).Format(time.RFC3339Nano)
			updateDate = time.Now().Format(time.RFC3339Nano)

		case i < 40:
			status = 2
			createDate = time.Now().Add(-time.Hour * 100).Format(time.RFC3339Nano)
			updateDate = time.Now().Add(-time.Hour * 65).Format(time.RFC3339Nano)

		case i < 60:
			status = 3
			createDate = time.Now().Add(-time.Hour * 120).Format(time.RFC3339Nano)
			updateDate = time.Now().Add(-time.Hour * 101).Format(time.RFC3339Nano)

		default:
			status = 4
			createDate = time.Now().Format(time.RFC3339Nano)
			updateDate = time.Now().Format(time.RFC3339Nano)

		}
		sampleFilters = append(sampleFilters, &SampleFilter{
			Name:      fmt.Sprintf("sample filter %d", i),
			Age:       i + 10,
			Status:    status,
			CreatedAt: createDate,
			UpdatedAt: updateDate,
		})
	}
}

func TestFilterActionStringMatch(t *testing.T) {
	var result []*SampleFilter
	for j := 1; j < 10; j++ {
		name := fmt.Sprintf("sample filter %d", j)
		action := FilterActionString(name)
		result = nil
		for i := range sampleFilters {
			if action(sampleFilters[i].Name) {
				result = append(result, sampleFilters[i])
			}
		}
		if size := len(result); size != 1 {
			t.Error("index", j, "expected result 1 got", size)
		} else if name != result[0].Name {
			t.Error("index", j, "expected result name", name, "got", result[0].Name)
		}
	}
}

func TestFilterActionStringParam(t *testing.T) {
	var result []*SampleFilter
	// suffix
	for j := 1; j < 10; j++ {
		name := fmt.Sprintf("%%filter %d", j)
		action := FilterActionString(name)
		name = name[1:]
		result = nil
		for i := range sampleFilters {
			if action(sampleFilters[i].Name) {
				result = append(result, sampleFilters[i])
			}
		}
		if size := len(result); size != 1 {
			t.Error("index", j, "expected result 1 got", size)
		} else if !strings.HasSuffix(result[0].Name, name) {
			t.Error("index", j, "expected result name", "'"+result[0].Name+"'", "to end with", "'"+name+"'")
		}
	}

	// prefix
	for j := 1; j < 10; j++ {
		name := fmt.Sprintf("sample filter %d%%", j)
		action := FilterActionString(name)
		name = name[:len(name)-1]
		result = nil
		for i := range sampleFilters {
			if action(sampleFilters[i].Name) {
				result = append(result, sampleFilters[i])
			}
		}
		if size := len(result); size != 11 {
			t.Error("index", j, "expected result 11 got", size)
		} else {
			for _, sf := range result {
				if !strings.HasPrefix(sf.Name, name) {
					t.Error("index", j, "expected result name", result[0].Name, "to end with", name)
				}
			}
		}
	}

	// prefix, escape
	for j := 1; j < 10; j++ {
		name := fmt.Sprintf("%%%%sample filter %d%%", j)
		action := FilterActionString(name)
		name = name[:len(name)-1]
		result = nil
		for i := range sampleFilters {
			if action(sampleFilters[i].Name) {
				result = append(result, sampleFilters[i])
			}
		}
		if size := len(result); size != 0 {
			t.Error("index", j, "expected result 0 got", size)
		}
	}

	// contain
	for j := 1; j < 10; j++ {
		name := fmt.Sprintf("%%filter %d%%", j)
		action := FilterActionString(name)
		name = name[1 : len(name)-1]
		result = nil
		for i := range sampleFilters {
			if action(sampleFilters[i].Name) {
				result = append(result, sampleFilters[i])
			}
		}
		if size := len(result); size != 11 {
			t.Error("index", j, "expected result 11 got", size)
		} else {
			for _, sf := range result {
				if !strings.Contains(sf.Name, name) {
					t.Error("index", j, "expected result name", result[0].Name, "to end with", name)
				}
			}
		}
	}
}

func TestFilterActionDateTime(t *testing.T) {
	// TODO: to implemented
}

func TestFilterActionNumber(t *testing.T) {
	var result []*SampleFilter
	// test equal
	for j := 0; j < 10; j++ {
		age := fmt.Sprintf("%d", 10+j)
		action, err := FilterActionNumber(age)
		if err != nil {
			t.Error("index", j, "expected error nil but got", err)
			continue
		}
		result = nil
		for i := range sampleFilters {
			if action(int64(sampleFilters[i].Age)) {
				result = append(result, sampleFilters[i])
			}
		}
		if size := len(result); size != 1 {
			t.Error("index", j, "expected result 1 got", size)
		} else if result[0].Age != 10+j {
			t.Error("index", j, "expected result age", result[0].Age, "to end with", age)
		}
	}

	// test set
	var ageSets []map[int]bool
	var ages []string
	var expecteds []int

	// expect 3
	ageSets = append(ageSets, map[int]bool{10: true, 11: true, 12: true})
	ages = append(ages, "{10, 11, 12}")
	expecteds = append(expecteds, 3)
	// expect 2
	ageSets = append(ageSets, map[int]bool{50: true, 15: true})
	ages = append(ages, "{140, 50, 15}")
	expecteds = append(expecteds, 2)
	// expect 0
	ageSets = append(ageSets, map[int]bool{})
	ages = append(ages, "{150, 141, 162}")
	expecteds = append(expecteds, 0)
	// expect 6
	ageSets = append(ageSets, map[int]bool{14: true, 18: true, 20: true, 22: true, 109: true, 98: true})
	ages = append(ages, "{14, 18, 20, 22, 109, 98}")
	expecteds = append(expecteds, 6)

	for i := range ages {
		ageSet := ageSets[i]
		age := ages[i]
		expected := expecteds[i]
		action, err := FilterActionNumber(age)
		if err != nil {
			t.Error("expected error nil but got", err)
		} else {
			result = nil
			for i := range sampleFilters {
				if action(int64(sampleFilters[i].Age)) {
					result = append(result, sampleFilters[i])
				}
			}
			if size := len(result); size != expected {
				t.Error("expected result", expected, "got", size)
			} else {
				for _, r := range result {
					if !ageSet[r.Age] {
						t.Error("expected result age", r.Age, "be one of 10, 11 or 12")
					}
				}
			}
		}
	}

	// test interval

	// test status interval
	expecteds = []int{20, 40, 20, 40, 40, 20, 20, 20}
	statusArg := []string{"[1,2[", "[1,2]", "]1,2]", "[4,5]", "[4,5[", "(1, 2]", "[1,2)", "]1,3["}
	for i, sa := range statusArg {
		action, err := FilterActionNumber(sa)
		if err != nil {
			t.Error("expected error nil but got", err)
		} else {
			result = nil
			for i := range sampleFilters {
				if action(int64(sampleFilters[i].Status)) {
					result = append(result, sampleFilters[i])
				}
			}
			if size := len(result); size != expecteds[i] {
				t.Error("expected result", expecteds[i], "got", size)
			}
		}
	}

	// test age interval
	var ageRange [][]int
	var ageIvs []string
	expecteds = nil

	// expect 5
	ageIvs = append(ageIvs, "[10, 14]")
	ageRange = append(ageRange, []int{9, 15})
	expecteds = append(expecteds, 5)
	// expect 0
	ageIvs = append(ageIvs, "[1, 8]")
	ageRange = append(ageRange, []int{})
	expecteds = append(expecteds, 0)
	// expect 7
	ageIvs = append(ageIvs, "[11, 18[")
	ageRange = append(ageRange, []int{10, 18})
	expecteds = append(expecteds, 7)

	for i := range ageIvs {
		ar := ageRange[i]
		age := ageIvs[i]
		expected := expecteds[i]
		action, err := FilterActionNumber(age)
		if err != nil {
			t.Error("expected error nil but got", err)
		} else {
			result = nil
			for i := range sampleFilters {
				if action(int64(sampleFilters[i].Age)) {
					result = append(result, sampleFilters[i])
				}
			}
			if size := len(result); size != expected {
				t.Error("expected result", expected, "got", size)
			} else {
				for _, r := range result {
					if !(ar[0] < r.Age && r.Age < ar[1]) {
						t.Error("expected result age", r.Age, "be in range of", ar[0], "and", ar[1])
					}
				}
			}
		}
	}

}
