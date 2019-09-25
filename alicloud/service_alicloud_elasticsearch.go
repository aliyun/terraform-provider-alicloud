package alicloud

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type ElasticsearchService struct {
	client *connectivity.AliyunClient
}

func (s *ElasticsearchService) DescribeElasticsearchInstance(id string) (response *elasticsearch.DescribeInstanceResponse, err error) {
	request := elasticsearch.CreateDescribeInstanceRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id
	request.SetContentType("application/json")

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, err := s.client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.DescribeInstance(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{ESInstanceNotFound}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}

			return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)

		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		response, _ = raw.(*elasticsearch.DescribeInstanceResponse)
		if response.Result.InstanceId != id {
			return WrapErrorf(Error(GetNotFoundMessage("Elasticsearch Instance", id)), NotFoundMsg, ProviderERROR)
		}

		return nil
	})

	return response, WrapError(err)
}

func (s *ElasticsearchService) ElasticsearchStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeElasticsearchInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.Result.Status == failState {
				return object, object.Result.Status, WrapError(Error(FailedToReachTargetStatus, object.Result.Status))
			}
		}

		return object, object.Result.Status, nil
	}
}

func updateDescription(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	content := make(map[string]interface{})
	content["description"] = d.Get("description").(string)
	data, err := json.Marshal(content)
	if err != nil {
		return WrapError(err)
	}

	request := elasticsearch.CreateUpdateDescriptionRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdateDescription(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	return nil
}

func updatePrivateWhitelist(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	content := make(map[string]interface{})
	content["esIPWhitelist"] = d.Get("private_whitelist").(*schema.Set).List()
	data, err := json.Marshal(content)
	if err != nil {
		return WrapError(err)
	}

	request := elasticsearch.CreateUpdateWhiteIpsRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdateWhiteIps(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func updatePublicWhitelist(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	content := make(map[string]interface{})
	content["publicIpWhitelist"] = d.Get("public_whitelist").(*schema.Set).List()
	data, err := json.Marshal(content)
	if err != nil {
		return WrapError(err)
	}
	request := elasticsearch.CreateUpdatePublicWhiteIpsRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdatePublicWhiteIps(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func updateDateNodeAmount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	content := make(map[string]interface{})
	content["nodeAmount"] = d.Get("data_node_amount").(int)

	data, err := json.Marshal(content)
	if err != nil {
		return WrapError(err)
	}
	request := elasticsearch.CreateUpdateInstanceRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (resp interface{}, errs error) {
			return elasticsearchClient.UpdateInstance(request)
		})
		if err != nil {
			if IsExceptedError(err, ESConcurrencyConflictError) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		return nil
	})
	if err != nil && !IsExceptedErrors(err, []string{ESMustChangeOneResource, ESCssCheckUpdowngradeError}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
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
	if err != nil {
		return WrapError(err)
	}
	request := elasticsearch.CreateUpdateInstanceRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.UpdateInstance(request)
		})
		if err != nil {
			if IsExceptedError(err, ESConcurrencyConflictError) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		return nil
	})
	if err != nil && !IsExceptedErrors(err, []string{ESMustChangeOneResource, ESCssCheckUpdowngradeError}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
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
	if err != nil {
		return WrapError(err)
	}
	request := elasticsearch.CreateUpdateInstanceRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.UpdateInstance(request)
		})
		if err != nil {
			if IsExceptedError(err, ESConcurrencyConflictError) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		return nil
	})
	if err != nil && !IsExceptedErrors(err, []string{ESMustChangeOneResource, ESCssCheckUpdowngradeError}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func updateKibanaWhitelist(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	content := make(map[string]interface{})
	content["kibanaIPWhitelist"] = d.Get("kibana_whitelist").(*schema.Set).List()
	data, err := json.Marshal(content)
	if err != nil {
		return WrapError(err)
	}
	request := elasticsearch.CreateUpdateKibanaWhiteIpsRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdateKibanaWhiteIps(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func updatePassword(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	content := make(map[string]interface{})
	content["esAdminPassword"] = d.Get("password")
	data, err := json.Marshal(content)
	if err != nil {
		return WrapError(err)
	}
	request := elasticsearch.CreateUpdateAdminPasswordRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdateAdminPassword(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func getChargeType(paymentType string) string {
	if strings.ToLower(paymentType) == strings.ToLower(string(PostPaid)) {
		return string(PostPaid)
	} else {
		return string(PrePaid)
	}
}

func filterWhitelist(destIPs []string, localIPs *schema.Set) []string {
	var whitelist []string
	if destIPs != nil {
		for _, ip := range destIPs {
			if (ip == "::1" || ip == "::/0" || ip == "127.0.0.1" || ip == "0.0.0.0/0") && !localIPs.Contains(ip) {
				continue
			}
			whitelist = append(whitelist, ip)
		}
	}
	return whitelist
}
