package utilfunc

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	DEFAULT_LIMIT = 10
)

func Retry(n int, fn func() error) error {
	start := 0
	var err error
	for {
		if start >= n {
			return err
		}

		err = fn()
		if err != nil {
			start++
		} else {
			return nil
		}
	}
}

func Validate(val any) error {
	v := reflect.ValueOf(val)
	err := make([]string, 0)

	for i := range v.NumField() {
		// field value
		field := v.Field(i)

		fmt.Println("field value", field)

		tag := v.Type().Field(i).Tag.Get("validate")
		if tag == "" {
			continue
		}

		rules := strings.SplitSeq(tag, ",")

		for rule := range rules {
			fieldName := v.Type().Field(i).Name

			switch {
			case rule == "required":
				if len(field.String()) == 0 {
					err = append(err, fmt.Sprintf("%v is required", fieldName))
				}
			case strings.HasPrefix(rule, "max="):
				max, _ := strconv.Atoi(strings.TrimPrefix(rule, "max="))
				if len(field.String()) > max {
					err = append(err,
						fmt.Sprintf("%s should be less than %d characters", fieldName, max))
				}
			case strings.HasPrefix(rule, "min="):
				min, _ := strconv.Atoi(strings.TrimPrefix(rule, "min="))
				if len(field.String()) < min {
					err = append(err,
						fmt.Sprintf("%s should be more than %d characters", fieldName, min))
				}
			case rule == "task_status":
				// if fVal != database.TaskStatusCompleted || field != database.TaskStatusInProgress || field != database.TaskStatusInProgress {
				// }
			}
		}
	}

	if len(err) > 0 {
		return errors.New(strings.Join(err, ", "))
	}
	return nil
}

// offset , limit
func GetOffsetAndLimitFromReq(r *http.Request) (int32, int32) {
	var offset int32

	limit_str := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limit_str)
	if err != nil {
		limit = DEFAULT_LIMIT
	}

	page_no_str := r.URL.Query().Get("page")
	page_no, err := strconv.Atoi(page_no_str)
	if err != nil {
		page_no = 0
		offset = int32(page_no * limit)
	} else {
		offset = int32((page_no - 1) * limit)
	}

	return int32(offset), int32(limit)
}
