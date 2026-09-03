// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type PrivatelinkServiceV2 struct {
	client *connectivity.AliyunClient
}

// parsePrivatelinkVpcEndpointServiceUserId parses the resource ID of
// alicloud_privatelink_vpc_endpoint_service_user. Supported formats:
//   - <service_id>:<user_id>            (legacy format for account ID whitelist entries)
//   - <service_id>:<user_id>:<user_arn> (user_id may be empty for ARN whitelist entries;
//     the ARN itself may contain colons, so the remainder after the second colon is
//     treated as the whole UserARN)
func parsePrivatelinkVpcEndpointServiceUserId(id string) (serviceId, userId, userArn string, err error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) < 2 || parts[0] == "" {
		return "", "", "", WrapError(fmt.Errorf("invalid Resource Id %s. Expected <service_id>:<user_id>[:<user_arn>]", id))
	}
	serviceId = parts[0]
	userId = parts[1]
	if len(parts) == 3 {
		userArn = parts[2]
	}
	return serviceId, userId, userArn, nil
}

// formatPrivatelinkUserId normalizes the UserId value returned by the API
// (integer in the API definition) into its plain decimal string form.
func formatPrivatelinkUserId(v interface{}) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case json.Number:
		return t.String()
	case float64:
		return strconv.FormatInt(int64(t), 10)
	case float32:
		return strconv.FormatInt(int64(t), 10)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	default:
		return fmt.Sprint(t)
	}
}

// parsePrivatelinkAccountArnUserId extracts the account ID from the ARN form
// "acs:ram:*:<uid>:*" that the API auto-generates for account ID whitelist
// entries in the UserARNs view. Returns "" when the ARN is not of that form.
func parsePrivatelinkAccountArnUserId(userArn string) string {
	parts := strings.Split(userArn, ":")
	if len(parts) == 5 && parts[0] == "acs" && parts[1] == "ram" && parts[2] == "*" && parts[4] == "*" {
		return parts[3]
	}
	return ""
}

// DescribePrivatelinkVpcEndpointServiceUser <<< Encapsulated get interface for Privatelink VpcEndpointServiceUser.

func (s *PrivatelinkServiceV2) DescribePrivatelinkVpcEndpointServiceUser(id string) (object map[string]interface{}, err error) {
	client := s.client
	serviceId, userId, userArn, err := parsePrivatelinkVpcEndpointServiceUserId(id)
	if err != nil {
		return nil, WrapError(err)
	}

	action := "ListVpcEndpointServiceUsers"
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"ServiceId":  serviceId,
		"MaxResults": PageSizeLarge,
	}
	if userId != "" {
		// The API supports filtering account ID whitelist entries by UserId directly.
		request["UserId"] = userId
	} else {
		// The API does not support filtering by UserARN. Pull the ARN whitelist
		// page by page and match the target ARN locally.
		request["UserListType"] = "UserARNs"
	}

	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	for {
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)

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
			if IsExpectedErrors(err, []string{"EndpointServiceNotFound"}) {
				return object, WrapErrorf(NotFoundErr("VpcEndpointServiceUser", id), NotFoundMsg, response)
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		if userId != "" {
			if users, e := jsonpath.Get("$.Users[*]", response); e == nil {
				arr, _ := users.([]interface{})
				for _, item := range arr {
					if m, ok := item.(map[string]interface{}); ok && formatPrivatelinkUserId(m["UserId"]) == userId {
						return m, nil
					}
				}
			}
		} else {
			if arns, e := jsonpath.Get("$.UserARNs[*]", response); e == nil {
				arr, _ := arns.([]interface{})
				for _, item := range arr {
					if m, ok := item.(map[string]interface{}); ok && fmt.Sprint(m["UserARN"]) == userArn {
						return m, nil
					}
				}
			}
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	return object, WrapErrorf(NotFoundErr("VpcEndpointServiceUser", id), NotFoundMsg, response)
}

func (s *PrivatelinkServiceV2) PrivatelinkVpcEndpointServiceUserStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePrivatelinkVpcEndpointServiceUser(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePrivatelinkVpcEndpointServiceUser >>> Encapsulated.
