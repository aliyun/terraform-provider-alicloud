package ackone

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Service wraps the AckOne API calls used by the resources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays
// unexported, so callers outside the package (including the flat
// layout's acceptance tests) build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeCluster <<< Encapsulated get interface for AckOne Cluster.

func (s *Service) DescribeCluster(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeHubClusterDetails"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ClusterId"] = id

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("adcp", "2022-01-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"Cluster.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Cluster", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Cluster", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Cluster", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) ClusterStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCluster(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

func (s *Service) DescribeMembershipAttachment(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeManagedClusters"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	clusterId := strings.Split(id, ":")[0]
	subClusterId := strings.Split(id, ":")[1]
	query["ClusterId"] = clusterId

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(15*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("adcp", "2022-01-01", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		clusters := response["Clusters"].([]interface{})
		if len(clusters) == 0 {
			return retry.RetryableError(fmt.Errorf("No managed cluster found"))
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"Cluster.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Cluster", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	clusters := response["Clusters"].([]interface{})

	found := false
	for _, cluster := range clusters {
		managedCluster := cluster.(map[string]interface{})
		clusterId := managedCluster["Cluster"].(map[string]interface{})["ClusterID"].(string)
		if clusterId == subClusterId {
			found = true
			break
		}
	}
	if !found {
		return object, errs.WrapErrorf(errs.NotFoundErr("Attachement", id), errs.NotFoundMsg, response)
	}

	v := map[string]interface{}{
		"cluster_id":     clusterId,
		"sub_cluster_id": subClusterId,
	}
	return v, nil
}

// DescribeCluster >>> Encapsulated.
