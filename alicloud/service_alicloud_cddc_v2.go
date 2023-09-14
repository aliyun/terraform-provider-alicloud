package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type CddcServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeCddcDedicatedPropreHost <<< Encapsulated get interface for Cddc DedicatedPropreHost.

func (s *CddcServiceV2) DescribeCddcDedicatedPropreHost(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "DescribeMyBaseDedicatedInstances"
	conn, err := client.NewCddcClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DedicatedHostGroupName"] = parts[0]
	query["EcsInstanceIds"] = fmt.Sprintf("[\"%s\"]", parts[1])
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-20"), StringPointer("AK"), query, request, &util.RuntimeOptions{})

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("DedicatedPropreHost", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *CddcServiceV2) CddcDedicatedPropreHostStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCddcDedicatedPropreHost(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		sysDiskCapacity, _ := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(sysDiskCapacity)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeCddcDedicatedPropreHost >>> Encapsulated.
