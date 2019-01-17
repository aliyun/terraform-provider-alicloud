package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type ElasticsearchService struct {
	client *connectivity.AliyunClient
}

func (s *ElasticsearchService) DescribeInstance(instanceId string) (v elasticsearch.DescribeInstanceResponse, err error) {
	request := elasticsearch.CreateDescribeInstanceRequest()
	request.InstanceId = instanceId
	request.SetContentType("application/json")

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.DescribeInstance(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{ESInstanceNotFound}) {
				return GetNotFoundErrorFromString(GetNotFoundMessage("elasticsearch", instanceId))
			}

			return err
		}

		resp, _ := raw.(*elasticsearch.DescribeInstanceResponse)
		if resp == nil || resp.Result.InstanceId != instanceId {
			return GetNotFoundErrorFromString(GetNotFoundMessage("elasticsearch", instanceId))
		}

		v = *resp
		return nil
	})

	return
}

func (s *ElasticsearchService) WaitForElasticsearchInstance(instanceId string, status []ElasticsearchStatus, timeout int) error {

	for _, elasticsearchStatus := range status {
		for {
			if resp, err := s.DescribeInstance(instanceId); err == nil {
				if resp.Result.Status == string(elasticsearchStatus) {
					break
				}
			}

			if timeout <= 0 {
				return common.GetClientErrorFromString(fmt.Sprintf("Timeout for %s", string(elasticsearchStatus)))
			}

			timeout = timeout - DefaultIntervalLong
			time.Sleep(DefaultIntervalLong * time.Second)
		}
	}

	return nil
}

func updateDescription(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	content := make(map[string]interface{})
	content["description"] = d.Get("description").(string)
	data, err := json.Marshal(content)

	request := elasticsearch.CreateUpdateDescriptionRequest()
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	_, err = client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdateDescription(request)
	})

	return BuildWrapError(request.GetActionName(), d.Id(), AlibabaCloudSdkGoERROR, err)
}

func updatePrivateWhitelist(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	content := make(map[string]interface{})
	content["esIPWhitelist"] = d.Get("private_whitelist").(*schema.Set).List()
	data, err := json.Marshal(content)

	request := elasticsearch.CreateUpdateWhiteIpsRequest()
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	if _, err = client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdateWhiteIps(request)
	}); err != nil {
		return BuildWrapError(request.GetActionName(), d.Id(), AlibabaCloudSdkGoERROR, err)
	}

	return elasticsearchService.WaitForElasticsearchInstance(d.Id(), []ElasticsearchStatus{ElasticsearchStatusActive}, WaitInstanceActiveTimeout)
}

func updatePublicWhitelist(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	content := make(map[string]interface{})
	content["publicIpWhitelist"] = d.Get("public_whitelist").(*schema.Set).List()
	data, err := json.Marshal(content)

	request := elasticsearch.CreateUpdatePublicWhiteIpsRequest()
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	if _, err = client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdatePublicWhiteIps(request)
	}); err != nil {
		return BuildWrapError(request.GetActionName(), d.Id(), AlibabaCloudSdkGoERROR, err)
	}

	return elasticsearchService.WaitForElasticsearchInstance(d.Id(), []ElasticsearchStatus{ElasticsearchStatusActive}, WaitInstanceActiveTimeout)
}

func updateDateNodeAmount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	content := make(map[string]interface{})
	content["nodeAmount"] = d.Get("data_node_amount").(int)

	data, err := json.Marshal(content)

	request := elasticsearch.CreateUpdateInstanceRequest()
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	if _, err = client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (resp interface{}, errs error) {
		return elasticsearchClient.UpdateInstance(request)
	}); err != nil {
		return BuildWrapError(request.GetActionName(), d.Id(), AlibabaCloudSdkGoERROR, err)
	}

	return elasticsearchService.WaitForElasticsearchInstance(d.Id(), []ElasticsearchStatus{ElasticsearchStatusActivating, ElasticsearchStatusActive}, WaitInstanceActiveTimeout)
}

func updateDataNodeSpec(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	content := make(map[string]interface{})
	spec := make(map[string]interface{})

	spec["spec"] = d.Get("data_node_spec")
	spec["disk"] = d.Get("data_node_disk_size")
	spec["diskType"] = d.Get("data_node_disk_type")
	content["nodeSpec"] = spec
	data, err := json.Marshal(content)

	request := elasticsearch.CreateUpdateInstanceRequest()
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	if _, err = client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdateInstance(request)
	}); err != nil {
		return BuildWrapError(request.GetActionName(), d.Id(), AlibabaCloudSdkGoERROR, err)
	}

	return elasticsearchService.WaitForElasticsearchInstance(d.Id(), []ElasticsearchStatus{ElasticsearchStatusActivating, ElasticsearchStatusActive}, WaitInstanceActiveTimeout)
}

func updateMasterNode(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	content := make(map[string]interface{})
	if d.Get("master_node_spec") != nil {
		master := make(map[string]interface{})
		master["spec"] = d.Get("master_node_spec").(string)
		master["amount"] = MasterNodeAmount
		master["diskType"] = MasterNodeDiskType
		master["disk"] = MasterNodeDisk
		content["masterConfiguration"] = master
		content["advancedDedicateMaster"] = true
	} else {
		content["advancedDedicateMaster"] = false
	}

	data, err := json.Marshal(content)
	request := elasticsearch.CreateUpdateInstanceRequest()
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	if _, err = client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdateInstance(request)
	}); err != nil {
		return BuildWrapError(request.GetActionName(), d.Id(), AlibabaCloudSdkGoERROR, err)
	}

	return elasticsearchService.WaitForElasticsearchInstance(d.Id(), []ElasticsearchStatus{ElasticsearchStatusActivating, ElasticsearchStatusActive}, WaitInstanceActiveTimeout)
}

func updateKibanaWhitelist(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	content := make(map[string]interface{})
	content["kibanaIPWhitelist"] = d.Get("kibana_whitelist").(*schema.Set).List()
	data, err := json.Marshal(content)

	request := elasticsearch.CreateUpdateKibanaWhiteIpsRequest()
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	if _, err = client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdateKibanaWhiteIps(request)
	}); err != nil {
		return BuildWrapError(request.GetActionName(), d.Id(), AlibabaCloudSdkGoERROR, err)
	}

	return elasticsearchService.WaitForElasticsearchInstance(d.Id(), []ElasticsearchStatus{ElasticsearchStatusActive}, WaitInstanceActiveTimeout)
}

func updatePassword(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	content := make(map[string]interface{})
	content["esAdminPassword"] = d.Get("password")
	data, err := json.Marshal(content)

	request := elasticsearch.CreateUpdateAdminPasswordRequest()
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	if _, err = client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdateAdminPassword(request)
	}); err != nil {
		return BuildWrapError(request.GetActionName(), d.Id(), AlibabaCloudSdkGoERROR, err)
	}

	return elasticsearchService.WaitForElasticsearchInstance(d.Id(), []ElasticsearchStatus{ElasticsearchStatusActive}, WaitInstanceActiveTimeout)
}

func getChargeType(paymentType string) string {
	if strings.ToLower(paymentType) == strings.ToLower(string(PostPaid)) {
		return string(PostPaid)
	} else {
		return string(PrePaid)
	}
}
