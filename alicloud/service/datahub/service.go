// Package datahub contains the datahub service resources and data sources
// migrated from the root alicloud package. The root package cannot be imported
// from a service package, so the encapsulated datahub calls live here.
package datahub

import (
	"strings"
	"time"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// Status mirrors the root package status type used by the WaitFor helpers.
type Status string

const (
	Deleted         = Status("Deleted")
	DefaultTimeout  = 120
	COLON_SEPARATED = ":"
	DoesNotExist    = "does not exist"
)

// DatahubService wraps the datahub API calls used by the resources of this package.
type DatahubService struct {
	client *connectivity.AliyunClient
}

// NewDatahubService is the exported constructor: the client field stays
// unexported, so callers outside the package (including the acceptance
// tests) build values with it.
func NewDatahubService(client *connectivity.AliyunClient) *DatahubService {
	return &DatahubService{client: client}
}

func (s *DatahubService) DescribeDatahubProject(id string) (*datahub.GetProjectResult, error) {
	var requestInfo *datahub.DataHub
	project := &datahub.GetProjectResult{}
	raw, err := s.client.WithDataHubClient(func(dataHubClient datahub.DataHubApi) (interface{}, error) {
		requestInfo = dataHubClient.(*datahub.DataHub)
		return dataHubClient.GetProject(id)
	})
	if err != nil {
		if isDatahubNotExistError(err) {
			return project, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunDatahubSdkGo)
		}
		return project, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetProject", errs.AliyunDatahubSdkGo)
	}
	if logs.DebugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = id
		logs.AddDebug("GetProject", raw, requestInfo, requestMap)
	}
	project, _ = raw.(*datahub.GetProjectResult)
	if project == nil {
		return project, errs.WrapErrorf(errs.NotFoundErr("DatahubProject", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return project, nil
}

func (s *DatahubService) WaitForDatahubProject(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeDatahubProject(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}

		if time.Now().After(deadline) {
			objstringfy, err := convert.ConvertArrayObjectToJsonString(object)
			if err != nil {
				return errs.WrapError(err)
			}
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, objstringfy, id, errs.ProviderERROR)
		}

	}
}

func (s *DatahubService) DescribeDatahubSubscription(id string) (*datahub.GetSubscriptionResult, error) {
	subscription := &datahub.GetSubscriptionResult{}
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return subscription, errs.WrapError(err)
	}
	projectName, topicName, subId := parts[0], parts[1], parts[2]

	var requestInfo *datahub.DataHub

	raw, err := s.client.WithDataHubClient(func(dataHubClient datahub.DataHubApi) (interface{}, error) {
		requestInfo = dataHubClient.(*datahub.DataHub)
		return dataHubClient.GetSubscription(projectName, topicName, subId)
	})
	if err != nil {
		if isDatahubNotExistError(err) {
			return subscription, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunDatahubSdkGo)
		}
		return subscription, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetSubscription", errs.AliyunDatahubSdkGo)
	}
	if logs.DebugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["TopicName"] = topicName
		requestMap["SubId"] = subId
		logs.AddDebug("GetProject", raw, requestInfo, requestMap)
	}
	subscription, _ = raw.(*datahub.GetSubscriptionResult)
	if subscription == nil || subscription.TopicName != topicName || subscription.SubId != subId {
		return subscription, errs.WrapErrorf(errs.NotFoundErr("DatahubSubscription", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return subscription, nil
}

func (s *DatahubService) WaitForDatahubSubscription(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 3)
	if err != nil {
		return errs.WrapError(err)
	}
	topicName, subId := parts[1], parts[2]
	for {
		object, err := s.DescribeDatahubSubscription(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.TopicName == topicName && object.SubId == subId && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.TopicName+":"+object.SubId, parts[1]+":"+parts[2], errs.ProviderERROR)
		}

	}
}

func (s *DatahubService) DescribeDatahubTopic(id string) (*datahub.GetTopicResult, error) {
	topic := &datahub.GetTopicResult{}
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return topic, errs.WrapError(err)
	}
	projectName, topicName := parts[0], parts[1]

	var requestInfo *datahub.DataHub

	raw, err := s.client.WithDataHubClient(func(dataHubClient datahub.DataHubApi) (interface{}, error) {
		requestInfo = dataHubClient.(*datahub.DataHub)
		return dataHubClient.GetTopic(projectName, topicName)
	})
	if err != nil {
		if isDatahubNotExistError(err) {
			return topic, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AliyunDatahubSdkGo)
		}
		return topic, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "GetTopic", errs.AliyunDatahubSdkGo)
	}
	if logs.DebugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["TopicName"] = topicName
		logs.AddDebug("GetTopic", raw, requestInfo, requestMap)
	}
	topic, _ = raw.(*datahub.GetTopicResult)
	if topic == nil {
		return topic, errs.WrapErrorf(errs.NotFoundErr("DatahubTopic", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return topic, nil
}

func (s *DatahubService) WaitForDatahubTopic(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	projectName, topicName := parts[0], parts[1]
	for {
		object, err := s.DescribeDatahubTopic(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object.ProjectName == projectName && object.TopicName == topicName && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object.ProjectName+":"+object.TopicName, id, errs.ProviderERROR)
		}

	}
}

func convUint64ToDate(t uint64) string {
	return time.Unix(int64(t), 0).Format("2006-01-02 15:04:05")
}

func getNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func getRecordSchema(typeMap map[string]interface{}) (recordSchema *datahub.RecordSchema) {
	recordSchema = datahub.NewRecordSchema()

	for k, v := range typeMap {
		recordSchema.AddField(datahub.Field{Name: string(k), Type: datahub.FieldType(v.(string))})
	}

	return recordSchema
}

func isRetryableDatahubError(err error) bool {
	if e, ok := err.(*datahub.DatahubClientError); ok && e.StatusCode >= 500 {
		return true
	}

	return false
}

func isDatahubNotExistError(err error) bool {
	return errs.IsExpectedErrors(err, []string{datahub.NoSuchProject, datahub.NoSuchTopic, datahub.NoSuchShard, datahub.NoSuchSubscription, DoesNotExist})
}

func getDefaultRecordSchemainMap() map[string]interface{} {

	return map[string]interface{}{
		"string_field": "STRING",
	}
}

func recordSchemaToMap(fields []datahub.Field) map[string]string {
	result := make(map[string]string)
	for _, f := range fields {
		result[f.Name] = f.Type.String()
	}

	return result
}

// IsTerraformTestingDatahubObject reports whether name looks like a Terraform
// acceptance test fixture, used by the project sweeper.
func IsTerraformTestingDatahubObject(name string) bool {
	prefixes := []string{
		"tf_testAcc",
		"tf_test_",
		"testAcc",
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			return true
		}
	}

	return false
}
