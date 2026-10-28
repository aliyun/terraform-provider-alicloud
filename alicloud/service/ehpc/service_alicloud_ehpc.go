package ehpc

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type EhpcService struct {
	client *connectivity.AliyunClient
}

// NewEhpcService is the exported constructor used by the package's external
// tests (package ehpc_test); production wiring reaches the service through
// inline EhpcService{client} construction.
func NewEhpcService(client *connectivity.AliyunClient) *EhpcService {
	return &EhpcService{client: client}
}

func (s *EhpcService) DescribeEhpcJobTemplate(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListJobTemplates"
	request := map[string]interface{}{
		"PageSize":   helper.PageSizeLarge,
		"PageNumber": 1,
	}
	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcGet("EHPC", "2018-04-12", action, request, nil)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Templates.JobTemplates", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Templates.JobTemplates", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ehpc", id), errs.NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["Id"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Ehpc", id), errs.NotFoundWithResponse, response)
	}
	return
}

func (s *EhpcService) DescribeEhpcCluster(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeCluster"
	request := map[string]interface{}{
		"ClusterId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("EHPC", "2018-04-12", action, request, nil)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ClusterNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Ehpc:Cluster", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.ClusterInfo", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ClusterInfo", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *EhpcService) EhpcClusterStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEhpcCluster(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

// intBetween is the package-private validator migrated verbatim from the root
// alicloud package (alicloud/validators.go). It differs from the exported
// validators.IntBetween, which adds the SkipResourceSchemaValidation branch and
// a trailing hint to the range-violation error; ehpc's ha_enable schema relies
// on this no-skip variant, so the exact behaviour is preserved here.
func intBetween(min, max int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(int)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be int", k))
			return
		}

		if v < min || v > max {
			es = append(es, fmt.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v))
			return
		}

		return
	}
}
