package helper

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// State-wait helpers moved from the root package (alicloud/common.go);
// package alicloud re-exports them as var aliases so the flat layout
// keeps compiling unchanged.

func BuildStateConf(pending, target []string, timeout, delay time.Duration, f retry.StateRefreshFunc) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Refresh:    f,
		Timeout:    timeout,
		Delay:      delay,
		MinTimeout: 3 * time.Second,
	}
}

// If auto renew, the period computed from ComputePeriodByUnit will be changed
// This method used to compute a period accourding to current period and unit
func ComputePeriodByUnit(createTime, endTime interface{}, currentPeriod int, periodUnit string) (int, error) {
	if createTime == nil {
		return 0, errs.WrapError(fmt.Errorf("createTime should not be nil"))
	}
	if endTime == nil {
		return 0, errs.WrapError(fmt.Errorf("endTime should not be nil"))
	}
	var createTimeStr, endTimeStr string
	switch value := createTime.(type) {
	case int64:
		createTimeStr = time.Unix(createTime.(int64), 0).Format(time.RFC3339)
		endTimeStr = time.Unix(endTime.(int64), 0).Format(time.RFC3339)
	case string:
		createTimeStr = createTime.(string)
		endTimeStr = endTime.(string)
	default:
		return 0, errs.WrapError(fmt.Errorf("Unsupported time type: %#v", value))
	}
	// currently, there is time value does not format as standard RFC3339
	UnStandardRFC3339 := "2006-01-02T15:04Z07:00"
	create, err := time.Parse(time.RFC3339, createTimeStr)
	if err != nil {
		log.Printf("Parase the CreateTime %#v failed and error is: %#v.", createTime, err)
		create, err = time.Parse(UnStandardRFC3339, createTimeStr)
		if err != nil {
			return 0, errs.WrapError(err)
		}
	}
	end, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		log.Printf("Parase the EndTime %#v failed and error is: %#v.", endTime, err)
		end, err = time.Parse(UnStandardRFC3339, endTimeStr)
		if err != nil {
			return 0, errs.WrapError(err)
		}
	}
	var period int
	switch periodUnit {
	case "Month":
		period = int(math.Floor(end.Sub(create).Hours() / 24 / 30))
	case "Week":
		period = int(math.Floor(end.Sub(create).Hours() / 24 / 7))
	case "Year":
		period = int(math.Floor(end.Sub(create).Hours() / 24 / 365))
	default:
		err = fmt.Errorf("Unexpected period unit %s", periodUnit)
	}
	// The period at least is 1
	if period < 1 {
		period = 1
	}
	if period > 12 {
		period = 12
	}
	// period can not be modified and if the new period is changed, using the previous one.
	if currentPeriod > 0 && currentPeriod != period {
		period = currentPeriod
	}
	return period, errs.WrapError(err)
}

func CheckWaitForReady(object interface{}, conditions map[string]interface{}) (bool, map[string]interface{}, error) {
	if conditions == nil {
		return false, nil, nil
	}
	objectType := reflect.TypeOf(object)
	objectValue := reflect.ValueOf(object)
	values := make(map[string]interface{})
	for key, value := range conditions {
		if _, ok := objectType.FieldByName(key); ok {
			current := objectValue.FieldByName(key)
			values[key] = current
			if fmt.Sprintf("%v", current) != fmt.Sprintf("%v", value) {
				return false, values, nil
			}
		} else {
			return false, values, errs.WrapError(fmt.Errorf("There is missing attribute %s in the object.", key))
		}
	}
	return true, values, nil
}
