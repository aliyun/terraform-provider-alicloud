package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"strings"
	"time"
)

type NATGatewayServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeNATGatewaySnatEntry <<< Encapsulated get interface for NATGateway SnatEntry.

func (s *NATGatewayServiceV2) DescribeNATGatewaySnatEntry(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "DescribeSnatTableEntries"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SnatEntryId"] = parts[1]
	request["SnatTableId"] = parts[0]
	request["RegionId"] = client.RegionId

wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidSnatTableId.NotFound", "InvalidSnatEntryId.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("SnatEntry", id)), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.SnatTableEntries.SnatTableEntry[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SnatTableEntries.SnatTableEntry[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(Error(GetNotFoundMessage("SnatEntry", id)), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *NATGatewayServiceV2) NATGatewaySnatEntryStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNATGatewaySnatEntry(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeNATGatewaySnatEntry >>> Encapsulated.
