// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudApigSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigSourceCreate,
		Read:   resourceAliCloudApigSourceRead,
		Update: resourceAliCloudApigSourceUpdate,
		Delete: resourceAliCloudApigSourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"association_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"association_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"k8s_source_info": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"nacos_source_info": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudApigSourceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/v1/sources")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	request["gatewayId"] = d.Get("gateway_id")
	request["type"] = d.Get("type")
	k8sSourceConfig := make(map[string]interface{})

	if v := d.Get("k8s_source_info"); !IsNil(v) {
		clusterId1, _ := jsonpath.Get("$[0].cluster_id", v)
		if clusterId1 != nil && clusterId1 != "" {
			k8sSourceConfig["clusterId"] = clusterId1
		}
		if len(k8sSourceConfig) > 0 {
			request["k8sSourceConfig"] = k8sSourceConfig
		}
	}

	nacosSourceConfig := make(map[string]interface{})

	if v := d.Get("nacos_source_info"); !IsNil(v) {
		instanceId1, _ := jsonpath.Get("$[0].instance_id", v)
		if instanceId1 != nil && instanceId1 != "" {
			nacosSourceConfig["instanceId"] = instanceId1
		}
		if len(nacosSourceConfig) > 0 {
			request["nacosSourceConfig"] = nacosSourceConfig
		}
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_source", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.sourceId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudApigSourceRead(d, meta)
}

func resourceAliCloudApigSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	objectRaw, err := apigServiceV2.DescribeApigSource(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_source DescribeApigSource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("association_reason", objectRaw["associationReason"])
	d.Set("association_status", objectRaw["associationStatus"])
	d.Set("create_time", objectRaw["createTimestamp"])
	d.Set("gateway_id", objectRaw["gatewayId"])
	d.Set("resource_group_id", objectRaw["resourceGroupId"])
	d.Set("source_name", objectRaw["name"])
	d.Set("type", objectRaw["type"])
	d.Set("update_time", objectRaw["updateTimestamp"])

	k8SSourceInfoMaps := make([]map[string]interface{}, 0)
	k8SSourceInfoMap := make(map[string]interface{})
	k8SSourceInfoRaw := make(map[string]interface{})
	if objectRaw["k8SSourceInfo"] != nil {
		k8SSourceInfoRaw = objectRaw["k8SSourceInfo"].(map[string]interface{})
	}
	if len(k8SSourceInfoRaw) > 0 {
		k8SSourceInfoMap["cluster_id"] = k8SSourceInfoRaw["clusterId"]

		k8SSourceInfoMaps = append(k8SSourceInfoMaps, k8SSourceInfoMap)
	}
	if err := d.Set("k8s_source_info", k8SSourceInfoMaps); err != nil {
		return err
	}
	nacosSourceInfoMaps := make([]map[string]interface{}, 0)
	nacosSourceInfoMap := make(map[string]interface{})
	nacosSourceInfoRaw := make(map[string]interface{})
	if objectRaw["nacosSourceInfo"] != nil {
		nacosSourceInfoRaw = objectRaw["nacosSourceInfo"].(map[string]interface{})
	}
	if len(nacosSourceInfoRaw) > 0 {
		nacosSourceInfoMap["address"] = nacosSourceInfoRaw["address"]
		nacosSourceInfoMap["cluster_id"] = nacosSourceInfoRaw["clusterId"]
		nacosSourceInfoMap["instance_id"] = nacosSourceInfoRaw["instanceId"]

		nacosSourceInfoMaps = append(nacosSourceInfoMaps, nacosSourceInfoMap)
	}
	if err := d.Set("nacos_source_info", nacosSourceInfoMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudApigSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	action := fmt.Sprintf("/move-resource-group")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["ResourceId"] = StringPointer(d.Id())
	query["RegionId"] = StringPointer(client.RegionId)
	if d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		query["ResourceGroupId"] = StringPointer(v.(string))
	}

	query["ResourceType"] = StringPointer("Source")
	query["Service"] = StringPointer("APIG")
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		apigServiceV2 := ApigServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("resource_group_id"))}, d.Timeout(schema.TimeoutUpdate), 35*time.Second, apigServiceV2.ApigSourceStateRefreshFunc(d.Id(), "resourceGroupId", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudApigSourceRead(d, meta)
}

func resourceAliCloudApigSourceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	sourceId := d.Id()
	action := fmt.Sprintf("/v1/sources/%s", sourceId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
