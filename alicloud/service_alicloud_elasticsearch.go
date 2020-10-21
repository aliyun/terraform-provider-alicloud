package alicloud

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ElasticsearchService struct {
	client *connectivity.AliyunClient
}

func (s *ElasticsearchService) DescribeElasticsearchInstance(id string) (*elasticsearch.DescribeInstanceResponse, error) {
	response := &elasticsearch.DescribeInstanceResponse{}
	request := elasticsearch.CreateDescribeInstanceRequest()
	request.RegionId = s.client.RegionId
	request.InstanceId = id
	request.SetContentType("application/json")

	invoker := NewInvoker()
	err := invoker.Run(func() error {
		raw, err := s.client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.DescribeInstance(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"InstanceNotFound"}) {
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

func (s *ElasticsearchService) ElasticsearchRetryFunc(wait func(), errorCodeList []string, do func(*elasticsearch.Client) (interface{}, error)) (interface{}, error) {
	var raw interface{}
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithElasticsearchClient(do)

		if err != nil {
			if IsExpectedErrors(err, errorCodeList) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	return raw, WrapError(err)
}

func (s *ElasticsearchService) TriggerNetwork(d *schema.ResourceData, content map[string]interface{}, meta interface{}) error {
	data, err := json.Marshal(content)
	if err != nil {
		return WrapError(err)
	}
	request := elasticsearch.CreateTriggerNetworkRequest()
	request.ClientToken = buildClientToken(request.GetActionName())
	request.RegionId = s.client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	errorCodeList := []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction", "InternalServerError"}
	raw, err := s.ElasticsearchRetryFunc(wait, errorCodeList, func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.TriggerNetwork(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, s.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func (s *ElasticsearchService) ModifyWhiteIps(d *schema.ResourceData, content map[string]interface{}, meta interface{}) error {
	data, err := json.Marshal(content)
	if err != nil {
		return WrapError(err)
	}
	request := elasticsearch.CreateModifyWhiteIpsRequest()
	request.ClientToken = buildClientToken(request.GetActionName())
	request.RegionId = s.client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	errorCodeList := []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}
	raw, err := s.ElasticsearchRetryFunc(wait, errorCodeList, func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.ModifyWhiteIps(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, s.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func (s *ElasticsearchService) DescribeElasticsearchTags(id string) (tags map[string]string, err error) {
	resourceIds, err := json.Marshal([]string{id})
	if err != nil {
		tmp := make(map[string]string)
		return tmp, WrapError(err)
	}

	request := elasticsearch.CreateListTagResourcesRequest()
	request.RegionId = s.client.RegionId
	request.ResourceIds = string(resourceIds)
	request.ResourceType = strings.ToUpper(string(TagResourceInstance))

	raw, err := s.client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.ListTagResources(request)
	})
	if err != nil {
		tmp := make(map[string]string)
		return tmp, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	response, _ := raw.(*elasticsearch.ListTagResourcesResponse)

	return s.tagsToMap(response.TagResources.TagResource), nil
}

func (s *ElasticsearchService) tagsToMap(tagSet []elasticsearch.TagResourceItem) (tags map[string]string) {
	result := make(map[string]string)
	for _, t := range tagSet {
		if !elasticsearchTagIgnored(t.TagKey, t.TagValue) {
			result[t.TagKey] = t.TagValue
		}
	}

	return result
}

func (s *ElasticsearchService) diffElasticsearchTags(oldTags, newTags map[string]interface{}) (remove []string, add []map[string]string) {
	for k, _ := range oldTags {
		remove = append(remove, k)
	}
	for k, v := range newTags {
		tag := map[string]string{
			"key":   k,
			"value": v.(string),
		}

		add = append(add, tag)
	}
	return
}

func (s *ElasticsearchService) getActionType(actionType bool) string {
	if actionType == true {
		return string(OPEN)
	} else {
		return string(CLOSE)
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
	request.ClientToken = buildClientToken(request.GetActionName())
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

func updateInstanceTags(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})
	remove, add := elasticsearchService.diffElasticsearchTags(o, n)

	// 对系统 Tag 进行过滤
	removeTagKeys := make([]string, 0)
	for _, v := range remove {
		if !elasticsearchTagIgnored(v, "") {
			removeTagKeys = append(removeTagKeys, v)
		}
	}
	if len(removeTagKeys) > 0 {
		tagKeys, err := json.Marshal(removeTagKeys)
		if err != nil {
			return WrapError(err)
		}

		resourceIds, err := json.Marshal([]string{d.Id()})
		if err != nil {
			return WrapError(err)
		}
		request := elasticsearch.CreateUntagResourcesRequest()
		request.RegionId = client.RegionId
		request.TagKeys = string(tagKeys)
		request.ResourceType = strings.ToUpper(string(TagResourceInstance))
		request.ResourceIds = string(resourceIds)
		request.SetContentType("application/json")

		raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.UntagResources(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	}

	if len(add) > 0 {
		content := make(map[string]interface{})
		content["ResourceIds"] = []string{d.Id()}
		content["ResourceType"] = strings.ToUpper(string(TagResourceInstance))
		content["Tags"] = add
		data, err := json.Marshal(content)
		if err != nil {
			return WrapError(err)
		}

		request := elasticsearch.CreateTagResourcesRequest()
		request.RegionId = client.RegionId
		request.SetContent(data)
		request.SetContentType("application/json")
		raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			return elasticsearchClient.TagResources(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	}

	return nil
}

func updateInstanceChargeType(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	content := make(map[string]interface{})
	content["paymentType"] = strings.ToLower(d.Get("instance_charge_type").(string))
	if d.Get("instance_charge_type").(string) == string(PrePaid) {
		paymentInfo := make(map[string]interface{})
		if d.Get("period").(int) >= 12 {
			paymentInfo["duration"] = d.Get("period").(int) / 12
			paymentInfo["pricingCycle"] = string(Year)
		} else {
			paymentInfo["duration"] = d.Get("period").(int)
			paymentInfo["pricingCycle"] = string(Month)
		}

		content["paymentInfo"] = paymentInfo
	}

	data, err := json.Marshal(content)
	if err != nil {
		return WrapError(err)
	}

	request := elasticsearch.CreateUpdateInstanceChargeTypeRequest()
	request.ClientToken = buildClientToken(request.GetActionName())
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdateInstanceChargeType(request)
	})

	time.Sleep(10 * time.Second)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	return nil
}

func renewInstance(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	content := make(map[string]interface{})
	if d.Get("period").(int) >= 12 {
		content["duration"] = d.Get("period").(int) / 12
		content["pricingCycle"] = string(Year)
	} else {
		content["duration"] = d.Get("period").(int)
		content["pricingCycle"] = string(Month)
	}

	data, err := json.Marshal(content)
	if err != nil {
		return WrapError(err)
	}

	request := elasticsearch.CreateRenewInstanceRequest()
	request.ClientToken = buildClientToken(request.GetActionName())
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.RenewInstance(request)
	})

	time.Sleep(10 * time.Second)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RoaRequest, request)

	return nil
}

func updateDataNodeAmount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	content := make(map[string]interface{})
	content["nodeAmount"] = d.Get("data_node_amount").(int)

	data, err := json.Marshal(content)
	if err != nil {
		return WrapError(err)
	}
	request := elasticsearch.CreateUpdateInstanceRequest()
	request.ClientToken = buildClientToken(request.GetActionName())
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	errorCodeList := []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}
	raw, err := elasticsearchService.ElasticsearchRetryFunc(wait, errorCodeList, func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdateInstance(request)
	})

	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
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
	request.ClientToken = buildClientToken(request.GetActionName())
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	errorCodeList := []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}
	raw, err := elasticsearchService.ElasticsearchRetryFunc(wait, errorCodeList, func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdateInstance(request)
	})

	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
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

func updateMasterNode(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	content := make(map[string]interface{})
	if d.Get("master_node_spec") != nil {
		master := make(map[string]interface{})
		master["spec"] = d.Get("master_node_spec").(string)
		master["amount"] = "3"
		master["diskType"] = "cloud_ssd"
		master["disk"] = "20"
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
	request.ClientToken = buildClientToken(request.GetActionName())
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	errorCodeList := []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}
	raw, err := elasticsearchService.ElasticsearchRetryFunc(wait, errorCodeList, func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdateInstance(request)
	})

	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
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

	password := d.Get("password").(string)

	kmsPassword := d.Get("kms_encrypted_password").(string)

	if password == "" && kmsPassword == "" {
		return WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
	}

	if password != "" {
		d.SetPartial("password")
		content["esAdminPassword"] = password
	} else {
		kmsService := KmsService{meta.(*connectivity.AliyunClient)}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		content["esAdminPassword"] = decryptResp.Plaintext
		d.SetPartial("kms_encrypted_password")
		d.SetPartial("kms_encryption_context")
	}

	data, err := json.Marshal(content)
	if err != nil {
		return WrapError(err)
	}
	request := elasticsearch.CreateUpdateAdminPasswordRequest()
	request.ClientToken = buildClientToken(request.GetActionName())
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	errorCodeList := []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}
	raw, err := elasticsearchService.ElasticsearchRetryFunc(wait, errorCodeList, func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
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

func updateClientNode(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	content := make(map[string]interface{})
	content["isHaveClientNode"] = true

	spec := make(map[string]interface{})
	spec["spec"] = d.Get("client_node_spec")
	if d.Get("client_node_amount") == nil {
		spec["amount"] = "2"
	} else {
		spec["amount"] = d.Get("client_node_amount")
	}

	spec["disk"] = "20"
	spec["diskType"] = "cloud_efficiency"
	content["clientNodeConfiguration"] = spec

	data, err := json.Marshal(content)
	if err != nil {
		return WrapError(err)
	}
	request := elasticsearch.CreateUpdateInstanceRequest()
	request.ClientToken = buildClientToken(request.GetActionName())
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContent(data)
	request.SetContentType("application/json")

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	errorCodeList := []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}
	raw, err := elasticsearchService.ElasticsearchRetryFunc(wait, errorCodeList, func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.UpdateInstance(request)
	})

	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
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

func openHttps(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	request := elasticsearch.CreateOpenHttpsRequest()
	request.ClientToken = buildClientToken(request.GetActionName())
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContentType("application/json")

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	errorCodeList := []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}
	raw, err := elasticsearchService.ElasticsearchRetryFunc(wait, errorCodeList, func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.OpenHttps(request)
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

func closeHttps(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	request := elasticsearch.CreateCloseHttpsRequest()
	request.ClientToken = buildClientToken(request.GetActionName())
	request.RegionId = client.RegionId
	request.InstanceId = d.Id()
	request.SetContentType("application/json")

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	errorCodeList := []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}
	raw, err := elasticsearchService.ElasticsearchRetryFunc(wait, errorCodeList, func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.CloseHttps(request)
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
