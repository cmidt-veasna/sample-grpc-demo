package query

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Action func(interface{}) bool

func FilterActionString(cond string) (action Action) {
	if cond == "" {
		return
	}

	cond = strings.TrimSpace(cond)
	last := len(cond) - 1

	ripEscape := func(in string) (bool, string) {
		suffix := false
		lastIndex := len(in) - 1
		if len(in) < 4 {
			switch {
			case in == "%%":
				in = "%"

			case in == "%%%":
				in = "%%"
				suffix = true

			case in == "%%%%":
				in = "%%"
			}
		} else {
			if in[lastIndex] == '%' && in[lastIndex-1] == '%' {
				in = in[:lastIndex-1]
			}
			if in[0] == '%' && in[1] == '%' {
				in = in[1:]
			}
		}
		return suffix, in
	}

	suffix := func(cond string) Action {
		_, comp := ripEscape(cond[1:])
		return func(i interface{}) bool {
			return strings.HasSuffix(i.(string), comp)
		}
	}

	prefix := func(cond string) Action {
		_, comp := ripEscape(cond[:len(cond)-1])
		return func(i interface{}) bool {
			return strings.HasPrefix(i.(string), comp)
		}
	}

	switch {
	case cond[0] == '%' && cond[1] != '%' && cond[last] == '%' && cond[last-1] != '%':
		_, comp := ripEscape(cond[1 : len(cond)-1])
		action = func(i interface{}) bool {
			return strings.Contains(i.(string), comp)
		}

	case cond[last] == '%' && last > 0 && cond[last-1] != '%':
		return prefix(cond)

	case cond[0] == '%' && last > 0 && cond[1] != '%':
		return suffix(cond)

	default:
		// remove escape sign (%)
		suf := false
		if suf, cond = ripEscape(cond); suf {
			suffix(cond)
		} else {
			action = func(i interface{}) bool {
				return i.(string) == cond
			}
		}
	}
	return
}

func FilterActionDateTime(cond string) (action Action, ierr error) {
	if cond == "" {
		return
	}
	// remove all whitespace
	cond = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, cond)

	var to2DateTime = func(cond string) (dateTimes []string, err error) {
		cond = cond[1 : len(cond)-1]
		index := strings.IndexByte(cond, ',')
		if index == 0 || index < 0 {
			err = errors.New(fmt.Sprintf("invalid given date time %s", cond))
			return
		}
		dateTimes = []string{cond[:index], cond[index+1:]}
		return
	}

	switch {
	case strings.HasPrefix(cond, "]") && strings.HasSuffix(cond, "["):
		dts, err := to2DateTime(cond)
		if err != nil {
			ierr = err
			return
		}
		action = func(i interface{}) bool {
			dt := i.(string)
			return dts[0] < dt && dt < dts[1]
		}

	case strings.HasPrefix(cond, "]") && strings.HasSuffix(cond, "]"):
		dts, err := to2DateTime(cond)
		if err != nil {
			ierr = err
			return
		}
		action = func(i interface{}) bool {
			dt := i.(string)
			return dts[0] < dt && dt <= dts[1]
		}

	case strings.HasPrefix(cond, "[") && strings.HasSuffix(cond, "["):
		dts, err := to2DateTime(cond)
		if err != nil {
			ierr = err
			return
		}
		action = func(i interface{}) bool {
			dt := i.(string)
			return dts[0] <= dt && dt < dts[1]
		}

	case strings.HasPrefix(cond, "[") && strings.HasSuffix(cond, "]"):
		cond = cond[1 : len(cond)-1]
		fallthrough

	default:
		action = func(i interface{}) bool {
			return strings.HasPrefix(i.(string), cond)
		}
	}
	return
}

func FilterActionNumber(cond string) (action Action, ierr error) {
	if cond == "" {
		return
	}
	// remove all whitespace
	cond = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, cond)

	var to2Num = func(val string) (nums []int64, err error) {
		// remove first and last sign of math interval
		val = val[1 : len(val)-1]
		// get split index
		index := strings.LastIndexByte(val, ',')
		if index == 0 || index < 0 {
			err = errors.New(fmt.Sprintf("invalid give value %s", val))
			return
		}
		nums = make([]int64, 2)
		if nums[0], err = strconv.ParseInt(val[:index], 10, 64); err == nil {
			nums[1], err = strconv.ParseInt(val[index+1:], 10, 64)
		}
		return
	}

	switch {
	case strings.HasPrefix(cond, "{") && strings.HasSuffix(cond, "}"):
		// math set syntax, see https://en.wikipedia.org/wiki/Set_(mathematics)
		conds := strings.Split(cond[1:len(cond)-1], ",")
		m := make(map[int64]bool)
		for i := range conds {
			n, err := strconv.ParseInt(conds[i], 10, 64)
			if err != nil {
				ierr = err
				return
			}
			m[n] = true
		}
		action = func(i interface{}) bool {
			return m[i.(int64)]
		}

		// below is math interval syntax, see https://en.wikipedia.org/wiki/Interval_(mathematics)
	case strings.HasPrefix(cond, "[") && strings.HasSuffix(cond, "]"):
		if nums, err := to2Num(cond); err != nil {
			ierr = err
			return
		} else {
			action = func(i interface{}) bool {
				num := i.(int64)
				return nums[0] <= num && num <= nums[1]
			}
		}

	case (strings.HasPrefix(cond, "]") || strings.HasPrefix(cond, "(")) && strings.HasSuffix(cond, "]"):
		if nums, err := to2Num(cond); err != nil {
			ierr = err
			return
		} else {
			action = func(i interface{}) bool {
				num := i.(int64)
				return nums[0] < num && num <= nums[1]
			}
		}

	case strings.HasPrefix(cond, "[") && (strings.HasSuffix(cond, "[") || strings.HasSuffix(cond, ")")):
		if nums, err := to2Num(cond); err != nil {
			ierr = err
			return
		} else {
			action = func(i interface{}) bool {
				num := i.(int64)
				return nums[0] <= num && num < nums[1]
			}
		}

	case (strings.HasPrefix(cond, "]") && strings.HasSuffix(cond, "[")) || (strings.HasPrefix(cond, "(") && strings.HasSuffix(cond, ")")):
		if nums, err := to2Num(cond); err != nil {
			ierr = err
			return
		} else {
			action = func(i interface{}) bool {
				num := i.(int64)
				return nums[0] < num && num < nums[1]
			}
		}

	default:
		n, err := strconv.ParseInt(cond, 10, 64)
		if err != nil {
			ierr = err
			return
		}
		action = func(i interface{}) bool {
			return i.(int64) == n
		}
	}
	return
}
