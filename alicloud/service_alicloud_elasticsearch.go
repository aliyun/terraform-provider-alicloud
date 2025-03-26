package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ElasticsearchService struct {
	client *connectivity.AliyunClient
}

func (s *ElasticsearchService) DescribeElasticsearchInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeInstance"
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGetWithApiName("elasticsearch", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s", id), nil, nil, nil)
		if err != nil {
			if IsExpectedErrors(err, []string{"QPS Limit Exceeded"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			addDebug(action, response, nil)
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	addDebug(action, response, nil)
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Result", response)
	}
	object = v.(map[string]interface{})
	if (object["instanceId"].(string)) != id {
		return object, WrapErrorf(NotFoundErr("Elasticsearch Instance", id), NotFoundWithResponse, response)
	}

	return object, WrapError(err)
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
			if object["status"].(string) == failState {
				return object, object["status"].(string), WrapError(Error(FailedToReachTargetStatus, object["status"].(string)))
			}
		}

		return object, object["status"].(string), nil
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
	var response map[string]interface{}
	var err error
	client := s.client
	action := "TriggerNetwork"
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.RoaPostWithApiName("elasticsearch", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s/actions/network-trigger", d.Id()), requestQuery, nil, content, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction", "InternalServerError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	addDebug(action, response, content)
	if err != nil {
		if IsExpectedErrors(err, []string{"RepetitionOperationError"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, s.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func (s *ElasticsearchService) getKibanaPvlNetworkInfo(id string) (interface{}, error) {
	listKibanaPvlNetworkReq := requests.RoaRequest{}
	listKibanaPvlNetworkReq.InitWithApiInfo("elasticsearch", "2017-06-13", "ListKibanaPvlNetwork", "/openapi/instances/[InstanceId]/actions/get-kibana-private", "elasticsearch", "openAPI")
	listKibanaPvlNetworkReq.Method = requests.GET
	listKibanaPvlNetworkReq.RegionId = s.client.RegionId
	listKibanaPvlNetworkReq.PathParams["InstanceId"] = id
	listKibanaPvlNetworkReq.SetContentType("application/json")
	listKibanaPvlNetworkResp := responses.BaseResponse{}

	invoker := NewInvoker()
	err := invoker.Run(func() error {
		raw, err := s.client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
			err := elasticsearchClient.DoAction(&listKibanaPvlNetworkReq, &listKibanaPvlNetworkResp)
			return listKibanaPvlNetworkResp, err
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"InstanceNotFound"}) {
				return WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
			}
			return WrapErrorf(err, DefaultErrorMsg, id, listKibanaPvlNetworkReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(listKibanaPvlNetworkReq.GetActionName(), raw, listKibanaPvlNetworkReq, listKibanaPvlNetworkReq)
		listKibanaPvlNetworkResp, _ = raw.(responses.BaseResponse)
		return nil
	})

	instanceMap := jsonToMap(listKibanaPvlNetworkResp.GetHttpContentString())
	resultMap := instanceMap["Result"]
	return resultMap, err
}

func (s *ElasticsearchService) updateKibanaPrivatePvlNetwork(d *schema.ResourceData, content map[string]interface{}, meta interface{}) error {
	data, err := json.Marshal(content)
	if err != nil {
		return WrapError(err)
	}
	updateKibanaPvlNetworkReq := requests.RoaRequest{}
	updateKibanaPvlNetworkReq.InitWithApiInfo("elasticsearch", "2017-06-13", "UpdateKibanaPvlNetwork", "/openapi/instances/[InstanceId]/actions/update-kibana-private", "elasticsearch", "openAPI")
	updateKibanaPvlNetworkReq.Method = requests.POST
	updateKibanaPvlNetworkReq.RegionId = s.client.RegionId
	updateKibanaPvlNetworkReq.PathParams["InstanceId"] = d.Id()
	updateKibanaPvlNetworkReq.SetContent(data)
	updateKibanaPvlNetworkReq.SetContentType("application/json")

	pvlInfoInterface, err := s.getKibanaPvlNetworkInfo(d.Id())
	if err != nil {
		return WrapErrorf(err, "get kibana pvl info error %s", d.Id())
	}

	pvlInfoArr := pvlInfoInterface.([]interface{})

	if len(pvlInfoArr) == 0 {
		return WrapErrorf(err, "get kibana pvl info empty %s", d.Id())
	}

	pvlInfo := pvlInfoArr[0]
	pvlId := pvlInfo.(map[string]interface{})["pvlId"].(string)
	updateKibanaPvlNetworkReq.QueryParams["pvlId"] = pvlId

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	errorCodeList := []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}
	raw, err := s.ElasticsearchRetryFunc(wait, errorCodeList, func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		enableResponse := responses.BaseResponse{}
		err := elasticsearchClient.DoAction(&updateKibanaPvlNetworkReq, &enableResponse)
		return enableResponse, err
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), updateKibanaPvlNetworkReq.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(updateKibanaPvlNetworkReq.GetActionName(), raw, updateKibanaPvlNetworkReq, updateKibanaPvlNetworkReq)
	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, s.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func (s *ElasticsearchService) enableKibanaPrivatePvlNetwork(d *schema.ResourceData, content map[string]interface{}, meta interface{}) error {
	data, err := json.Marshal(content)
	if err != nil {
		return WrapError(err)
	}
	enableKibanaPvlNetworkReq := requests.RoaRequest{}
	enableKibanaPvlNetworkReq.InitWithApiInfo("elasticsearch", "2017-06-13", "EnableKibanaPvlNetwork", "/openapi/instances/[InstanceId]/actions/enable-kibana-private", "elasticsearch", "openAPI")
	enableKibanaPvlNetworkReq.Method = requests.POST

	enableKibanaPvlNetworkReq.RegionId = s.client.RegionId
	enableKibanaPvlNetworkReq.PathParams["InstanceId"] = d.Id()
	enableKibanaPvlNetworkReq.SetContent(data)
	enableKibanaPvlNetworkReq.SetContentType("application/json")

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	errorCodeList := []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}
	raw, err := s.ElasticsearchRetryFunc(wait, errorCodeList, func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		enableResponse := responses.BaseResponse{}

		err := elasticsearchClient.DoAction(&enableKibanaPvlNetworkReq, &enableResponse)
		return enableResponse, err
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), enableKibanaPvlNetworkReq.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(enableKibanaPvlNetworkReq.GetActionName(), raw, enableKibanaPvlNetworkReq, enableKibanaPvlNetworkReq)

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, s.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func (s *ElasticsearchService) disableKibanaPrivatePvlNetwork(d *schema.ResourceData, content map[string]interface{}, meta interface{}) error {
	data, err := json.Marshal(content)
	if err != nil {
		return WrapError(err)
	}
	disableRequest := requests.RoaRequest{}
	disableRequest.InitWithApiInfo("elasticsearch", "2017-06-13", "DisableKibanaPvlNetwork", "/openapi/instances/[InstanceId]/actions/disable-kibana-private", "elasticsearch", "openAPI")
	disableRequest.Method = requests.POST

	disableRequest.RegionId = s.client.RegionId
	disableRequest.PathParams["InstanceId"] = d.Id()
	disableRequest.SetContent(data)
	disableRequest.SetContentType("application/json")

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	errorCodeList := []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}
	raw, err := s.ElasticsearchRetryFunc(wait, errorCodeList, func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		enableResponse := responses.BaseResponse{}

		err := elasticsearchClient.DoAction(&disableRequest, &enableResponse)
		return enableResponse, err
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), disableRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(disableRequest.GetActionName(), raw, disableRequest, disableRequest)

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, s.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func (s *ElasticsearchService) ModifyWhiteIps(d *schema.ResourceData, content map[string]interface{}, meta interface{}) error {
	var response map[string]interface{}
	var err error
	client := s.client
	action := "ModifyWhiteIps"
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.RoaPostWithApiName("elasticsearch", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s/actions/modify-white-ips", d.Id()), requestQuery, nil, content, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction", "InternalServerError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	addDebug(action, response, nil)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, s.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
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

	addDebug(request.GetActionName(), raw, request.RoaRequest, request)
	if err != nil {
		tmp := make(map[string]string)
		return tmp, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

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
	for k := range oldTags {
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
	action := "UpdateDescription"

	content := make(map[string]interface{})
	content["description"] = d.Get("description").(string)
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.RoaPostWithApiName("elasticsearch", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s/description", d.Id()), requestQuery, nil, content, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"GetCustomerLabelFail"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
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

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
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

		addDebug(request.GetActionName(), raw, request.RoaRequest, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	return nil
}

func updateInstanceChargeType(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	action := "UpdateInstanceChargeType"

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
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}
	response, err := client.RoaPostWithApiName("elasticsearch", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s/actions/convert-pay-type", d.Id()), requestQuery, nil, content, false)

	time.Sleep(10 * time.Second)

	addDebug(action, response, content)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), response, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func renewInstance(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	action := "RenewInstance"

	content := make(map[string]interface{})
	if d.Get("period").(int) >= 12 {
		content["duration"] = d.Get("period").(int) / 12
		content["pricingCycle"] = string(Year)
	} else {
		content["duration"] = d.Get("period").(int)
		content["pricingCycle"] = string(Month)
	}
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}
	response, err := client.RoaPostWithApiName("elasticsearch", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s/actions/renew", d.Id()), requestQuery, nil, content, false)

	time.Sleep(10 * time.Second)

	addDebug(action, response, content)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), response, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func setRenewalInstance(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "SetRenewal"
	var renewalResponse map[string]interface{}
	var err error
	var endpoint string
	setRenewalReq := map[string]interface{}{
		"InstanceIDs":      d.Id(),
		"ProductCode":      "elasticsearch",
		"ProductType":      "elasticsearchpre",
		"SubscriptionType": "Subscription",
	}
	if client.IsInternationalAccount() {
		setRenewalReq["ProductType"] = "elasticsearchpre_intl"
	}

	if _, ok := d.GetOk("auto_renew_duration"); !ok && d.Get("renew_status").(string) == "AutoRenewal" {
		return WrapError(fmt.Errorf("UpdateInstance error: auto_renew_duration is null"))
	}

	if _, ok := d.GetOk("renewal_duration_unit"); !ok && d.Get("renew_status").(string) == "AutoRenewal" {
		return WrapError(fmt.Errorf("CreateInstance error: renewal_duration_unit is null"))
	}

	if v, ok := d.GetOk("renew_status"); ok {
		setRenewalReq["RenewalStatus"] = v.(string)
	}

	if v, ok := d.GetOk("auto_renew_duration"); ok {
		setRenewalReq["RenewalPeriod"] = v.(int)
	}

	if v, ok := d.GetOk("renewal_duration_unit"); ok {
		setRenewalReq["RenewalPeriodUnit"] = v.(string)
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		renewalResponse, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, setRenewalReq, false, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				setRenewalReq["ProductType"] = "elasticsearchpre_intl"
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, renewalResponse, setRenewalReq)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

func updateDataNodeAmount(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}
	var err error
	action := "UpdateInstance"

	var response map[string]interface{}
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}
	content := make(map[string]interface{})
	content["nodeAmount"] = d.Get("data_node_amount").(int)

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.RoaPutWithApiName("elasticsearch", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s", d.Id()), requestQuery, nil, content, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction", "InstanceDuplicateScheduledTask"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	addDebug(action, response, content)
	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

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
	var err error
	action := "UpdateInstance"

	var response map[string]interface{}
	content := make(map[string]interface{})
	spec := make(map[string]interface{})
	spec["spec"] = d.Get("data_node_spec")
	spec["disk"] = d.Get("data_node_disk_size")
	spec["diskType"] = d.Get("data_node_disk_type")

	if v, ok := d.GetOkExists("data_node_disk_performance_level"); ok {
		spec["performanceLevel"] = v.(string)
	}
	content["nodeSpec"] = spec
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.RoaPutWithApiName("elasticsearch", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s", d.Id()), requestQuery, nil, content, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction", "InstanceDuplicateScheduledTask"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	addDebug(action, response, content)
	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func updateMasterNode(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}
	var err error
	action := "UpdateInstance"

	var response map[string]interface{}
	content := make(map[string]interface{})

	var masterDiskType string
	if v, ok := d.GetOkExists("master_node_disk_type"); ok {
		masterDiskType = v.(string)
	} else {
		masterDiskType = "cloud_ssd"
	}
	if d.Get("master_node_spec") != nil {
		master := make(map[string]interface{})
		master["spec"] = d.Get("master_node_spec").(string)
		master["amount"] = "3"
		master["diskType"] = masterDiskType
		master["disk"] = "20"
		content["masterConfiguration"] = master
		content["advancedDedicateMaster"] = true
	} else {
		content["advancedDedicateMaster"] = false
	}
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.RoaPutWithApiName("elasticsearch", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s", d.Id()), requestQuery, nil, content, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction", "InstanceDuplicateScheduledTask"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, content)

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 1*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func updatePassword(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}
	var err error
	action := "UpdateAdminPassword"

	var response map[string]interface{}
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
		content["esAdminPassword"] = decryptResp
		d.SetPartial("kms_encrypted_password")
		d.SetPartial("kms_encryption_context")
	}
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.RoaPostWithApiName("elasticsearch", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s/admin-pwd", d.Id()), requestQuery, nil, content, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction", "InstanceDuplicateScheduledTask"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	addDebug(action, response, content)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 120*time.Second, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
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
			whitelist = append(whitelist, ip)
		}
	}
	return whitelist
}

func updateClientNode(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}
	var err error
	action := "UpdateInstance"

	var response map[string]interface{}
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
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.RoaPutWithApiName("elasticsearch", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s", d.Id()), requestQuery, nil, content, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction", "InstanceDuplicateScheduledTask"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	addDebug(action, response, content)
	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func updateWarmNode(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	content := make(map[string]interface{})
	warm := make(map[string]interface{})
	if d.Get("warm_node_spec") != nil {
		warm["spec"] = d.Get("warm_node_spec").(string)
		warm["amount"] = d.Get("warm_node_amount")
		warm["diskType"] = "cloud_efficiency"
		warm["disk"] = d.Get("warm_node_disk_size")
		warm["diskEncryption"] = d.Get("warm_node_disk_encrypted")
		content["warmNodeConfiguration"] = warm
	} else {
		content["warmNodeConfiguration"] = warm
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

func updateKibanaNode(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}
	var err error
	action := "UpdateInstance"

	var response map[string]interface{}
	content := make(map[string]interface{})
	content["haveKibana"] = true

	spec := make(map[string]interface{})
	spec["spec"] = d.Get("kibana_node_spec")
	spec["amount"] = "1"
	spec["disk"] = "0"
	content["kibanaConfiguration"] = spec
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.RoaPutWithApiName("elasticsearch", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s", d.Id()), requestQuery, nil, content, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction", "InstanceDuplicateScheduledTask"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, content)
		return nil
	})

	addDebug(action, response, content)
	if err != nil && !IsExpectedErrors(err, []string{"MustChangeOneResource", "CssCheckUpdowngradeError"}) {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

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
	var err error
	action := "OpenHttps"

	var response map[string]interface{}
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.RoaPostWithApiName("elasticsearch", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s/actions/open-https", d.Id()), requestQuery, nil, nil, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	addDebug(action, response, nil)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 90*time.Second, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func closeHttps(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}
	var err error
	action := "CloseHttps"

	var response map[string]interface{}
	requestQuery := map[string]*string{
		"clientToken": StringPointer(buildClientToken(action)),
	}

	// retry
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := client.RoaPostWithApiName("elasticsearch", "2017-06-13", action, fmt.Sprintf("/openapi/instances/%s/actions/close-https", d.Id()), requestQuery, nil, nil, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ConcurrencyUpdateInstanceConflict", "InstanceStatusNotSupportCurrentAction"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})

	addDebug(action, response, nil)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"activating"}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 90*time.Second, elasticsearchService.ElasticsearchStateRefreshFunc(d.Id(), []string{"inactive"}))
	stateConf.PollInterval = 5 * time.Second

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func jsonToMap(content string) map[string]interface{} {
	var dataMap map[string]interface{}

	err := json.Unmarshal([]byte(content), &dataMap)
	if err != nil {
		log.Fatalf("parse json to map error: %v", err)
	}

	return dataMap
}
