package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// DescribeCmsEntityGroup <<< Encapsulated get interface for Cms EntityGroup.
func (s *CmsServiceV2) DescribeCmsEntityGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	entityGroupId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/entity-groups/%s", entityGroupId)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
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
		if IsExpectedErrors(err, []string{"404", "500"}) {
			return object, WrapErrorf(NotFoundErr("EntityGroup", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

// DescribeCmsEntityGroup >>> Encapsulated.
