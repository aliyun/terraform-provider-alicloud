// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudHbrPolicyBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudHbrPolicyBindingCreate,
		Read:   resourceAliCloudHbrPolicyBindingRead,
		Update: resourceAliCloudHbrPolicyBindingUpdate,
		Delete: resourceAliCloudHbrPolicyBindingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"advanced_options": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"udm_detail": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"exclude_disk_id_list": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"destination_kms_key_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"disk_id_list": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_source_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"exclude": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"include": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_binding_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"UDM_ECS", "NAS", "OSS", "File", "ECS_FILE"}, true),
			},
			"speed_limit": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudHbrPolicyBindingCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreatePolicyBindings"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["PolicyId"] = d.Get("policy_id")

	policyBindingListLocalMaps := make([]map[string]interface{}, 0)
	policyBindingListLocalMap := make(map[string]interface{})
	policyBindingListLocalMap["DataSourceId"] = d.Get("data_source_id")
	policyBindingListLocalMap["SourceType"] = d.Get("source_type")
	if d.HasChange("disabled") {
		policyBindingListLocalMap["Disabled"] = d.Get("disabled")
	}

	if v, ok := d.GetOk("source"); ok {
		policyBindingListLocalMap["Source"] = v
	}

	if v, ok := d.GetOk("policy_binding_description"); ok {
		policyBindingListLocalMap["PolicyBindingDescription"] = v
	}

	if v, ok := d.GetOk("include"); ok {
		policyBindingListLocalMap["Include"] = v
	}

	if v, ok := d.GetOk("exclude"); ok {
		policyBindingListLocalMap["Exclude"] = v
	}

	if v, ok := d.GetOk("speed_limit"); ok {
		policyBindingListLocalMap["SpeedLimit"] = v
	}

	if _, ok := d.GetOk("advanced_options"); ok {
		objectDataLocalMap := make(map[string]interface{})
		if v := d.Get("advanced_options"); v != nil {
			udmDetail := make(map[string]interface{})
			nodeNative, _ := jsonpath.Get("$[0].udm_detail[0].disk_id_list", d.Get("advanced_options"))
			if nodeNative != nil && nodeNative != "" {
				udmDetail["DiskIdList"] = nodeNative
			}
			nodeNative1, _ := jsonpath.Get("$[0].udm_detail[0].destination_kms_key_id", v)
			if nodeNative1 != nil && nodeNative1 != "" {
				udmDetail["DestinationKmsKeyId"] = nodeNative1
			}
			nodeNative2, _ := jsonpath.Get("$[0].udm_detail[0].exclude_disk_id_list", d.Get("advanced_options"))
			if nodeNative2 != nil && nodeNative2 != "" {
				udmDetail["ExcludeDiskIdList"] = nodeNative2
			}

			objectDataLocalMap["UdmDetail"] = udmDetail
			policyBindingListLocalMap["AdvancedOptions"] = convertMapToJsonStringIgnoreError(objectDataLocalMap)
		}
	}
	policyBindingListLocalMaps = append(policyBindingListLocalMaps, policyBindingListLocalMap)
	request["PolicyBindingList"], _ = convertListMapToJsonString(policyBindingListLocalMaps)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), query, request, &runtime)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_hbr_policy_binding", action, AlibabaCloudSdkGoERROR)
	}

	PolicyBindingListSourceType := d.Get("source_type")
	PolicyBindingListDataSourceId := d.Get("data_source_id")
	d.SetId(fmt.Sprintf("%v:%v:%v", request["PolicyId"], PolicyBindingListSourceType, PolicyBindingListDataSourceId))

	return resourceAliCloudHbrPolicyBindingRead(d, meta)
}

func resourceAliCloudHbrPolicyBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	hbrServiceV2 := HbrServiceV2{client}

	objectRaw, err := hbrServiceV2.DescribeHbrPolicyBinding(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_hbr_policy_binding DescribeHbrPolicyBinding Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreatedTime"])
	d.Set("disabled", objectRaw["Disabled"])
	d.Set("exclude", objectRaw["Exclude"])
	d.Set("include", objectRaw["Include"])
	d.Set("policy_binding_description", objectRaw["PolicyBindingDescription"])
	d.Set("source", objectRaw["Source"])
	d.Set("speed_limit", objectRaw["SpeedLimit"])
	d.Set("data_source_id", objectRaw["DataSourceId"])
	d.Set("policy_id", objectRaw["PolicyId"])
	d.Set("source_type", objectRaw["SourceType"])

	advancedOptionsMaps := make([]map[string]interface{}, 0)
	advancedOptionsMap := make(map[string]interface{})
	udmDetail1RawObj, _ := jsonpath.Get("$.AdvancedOptions.UdmDetail", objectRaw)
	udmDetail1Raw := make(map[string]interface{})
	if udmDetail1RawObj != nil {
		udmDetail1Raw = udmDetail1RawObj.(map[string]interface{})
	}
	if len(udmDetail1Raw) > 0 {
		udmDetailMaps := make([]map[string]interface{}, 0)
		udmDetailMap := make(map[string]interface{})
		if len(udmDetail1Raw) > 0 {
			udmDetailMap["destination_kms_key_id"] = udmDetail1Raw["DestinationKmsKeyId"]

			diskIdList1Raw := make([]interface{}, 0)
			if udmDetail1Raw["DiskIdList"] != nil {
				diskIdList1Raw = udmDetail1Raw["DiskIdList"].([]interface{})
			}

			udmDetailMap["disk_id_list"] = diskIdList1Raw
			excludeDiskIdList1Raw := make([]interface{}, 0)
			if udmDetail1Raw["ExcludeDiskIdList"] != nil {
				excludeDiskIdList1Raw = udmDetail1Raw["ExcludeDiskIdList"].([]interface{})
			}

			udmDetailMap["exclude_disk_id_list"] = excludeDiskIdList1Raw
			udmDetailMaps = append(udmDetailMaps, udmDetailMap)
		}
		advancedOptionsMap["udm_detail"] = udmDetailMaps
		advancedOptionsMaps = append(advancedOptionsMaps, advancedOptionsMap)
	}
	d.Set("advanced_options", advancedOptionsMaps)

	parts := strings.Split(d.Id(), ":")
	d.Set("policy_id", parts[0])
	d.Set("source_type", parts[1])
	d.Set("data_source_id", parts[2])

	return nil
}

func resourceAliCloudHbrPolicyBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "UpdatePolicyBinding"
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["PolicyId"] = parts[0]
	query["DataSourceId"] = parts[2]
	query["SourceType"] = parts[1]
	if d.HasChange("disabled") {
		update = true
		request["Disabled"] = d.Get("disabled")
	}

	if d.HasChange("source") {
		update = true
		request["Source"] = d.Get("source")
	}

	if d.HasChange("policy_binding_description") {
		update = true
		request["PolicyBindingDescription"] = d.Get("policy_binding_description")
	}

	if d.HasChange("include") {
		update = true
		request["Include"] = d.Get("include")
	}

	if d.HasChange("exclude") {
		update = true
		request["Exclude"] = d.Get("exclude")
	}

	if d.HasChange("speed_limit") {
		update = true
		request["SpeedLimit"] = d.Get("speed_limit")
	}

	if d.HasChange("advanced_options") {
		update = true
		objectDataLocalMap := make(map[string]interface{})
		if v := d.Get("advanced_options"); v != nil {
			udmDetail := make(map[string]interface{})
			nodeNative, _ := jsonpath.Get("$[0].udm_detail[0].disk_id_list", d.Get("advanced_options"))
			if nodeNative != nil && nodeNative != "" {
				udmDetail["DiskIdList"] = nodeNative
			}
			nodeNative1, _ := jsonpath.Get("$[0].udm_detail[0].destination_kms_key_id", v)
			if nodeNative1 != nil && nodeNative1 != "" {
				udmDetail["DestinationKmsKeyId"] = nodeNative1
			}
			nodeNative2, _ := jsonpath.Get("$[0].udm_detail[0].exclude_disk_id_list", d.Get("advanced_options"))
			if nodeNative2 != nil && nodeNative2 != "" {
				udmDetail["ExcludeDiskIdList"] = nodeNative2
			}

			objectDataLocalMap["UdmDetail"] = udmDetail
			request["AdvancedOptions"] = convertMapToJsonStringIgnoreError(objectDataLocalMap)
		}
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), query, request, &runtime)

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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudHbrPolicyBindingRead(d, meta)
}

func resourceAliCloudHbrPolicyBindingDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeletePolicyBinding"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewHbrClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["PolicyId"] = parts[0]
	request["SourceType"] = parts[1]
	request["DataSourceIds"] = "[\"" + parts[2] + "\"]"

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-08"), StringPointer("AK"), query, request, &runtime)

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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
