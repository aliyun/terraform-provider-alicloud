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

type AckOneServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeAckOneCluster <<< Encapsulated get interface for AckOne Cluster.

func (s *AckOneServiceV2) DescribeAckOneCluster(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeHubClusterDetails"
	conn, err := client.NewAckoneClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ClusterId"] = id

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), query, request, &runtime)

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
		if IsExpectedErrors(err, []string{"Cluster.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Cluster", id)), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Cluster", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Cluster", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *AckOneServiceV2) AckOneClusterStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAckOneCluster(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
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

func (s *AckOneServiceV2) DescribeAckOneMembershipAttachment(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeManagedClusters"
	conn, err := client.NewAckoneClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	clusterId := strings.Split(id, ":")[0]
	subClusterId := strings.Split(id, ":")[1]
	query["ClusterId"] = clusterId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(15*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-01"), StringPointer("AK"), query, request, &runtime)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		clusters := response["Clusters"].([]interface{})
		if len(clusters) == 0 {
			return resource.RetryableError(fmt.Errorf("No managed cluster found"))
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"Cluster.NotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Cluster", id)), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
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
		return object, WrapErrorf(Error(GetNotFoundMessage("Attachement", id)), NotFoundMsg, response)
	}

	v := map[string]interface{}{
		"cluster_id":     clusterId,
		"sub_cluster_id": subClusterId,
	}
	return v, nil
}

// DescribeAckOneCluster >>> Encapsulated.
