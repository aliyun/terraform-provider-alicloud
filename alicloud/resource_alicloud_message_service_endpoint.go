// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func resourceAliCloudMessageServiceEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMessageServiceEndpointCreate,
		Read:   resourceAliCloudMessageServiceEndpointRead,
		Update: resourceAliCloudMessageServiceEndpointUpdate,
		Delete: resourceAliCloudMessageServiceEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cidr_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:     schema.TypeString,
							Required: true,
						},
						"acl_strategy": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: StringInSlice([]string{"allow"}, false),
						},
					},
				},
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"public"}, false),
			},
		},
	}
}

func resourceAliCloudMessageServiceEndpointCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "EnableEndpoint"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	if v, ok := d.GetOk("endpoint_type"); ok {
		request["EndpointType"] = v
	}
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), query, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_message_service_endpoint", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["EndpointType"]))

	return resourceAliCloudMessageServiceEndpointUpdate(d, meta)
}

func resourceAliCloudMessageServiceEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	messageServiceServiceV2 := MessageServiceServiceV2{client}

	objectRaw, err := messageServiceServiceV2.DescribeMessageServiceEndpoint(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_message_service_endpoint DescribeMessageServiceEndpoint Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	cidrList1Raw := objectRaw["CidrList"]
	cidrListMaps := make([]map[string]interface{}, 0)
	if cidrList1Raw != nil {
		for _, cidrListChild1Raw := range cidrList1Raw.([]interface{}) {
			cidrListMap := make(map[string]interface{})
			cidrListChild1Raw := cidrListChild1Raw.(map[string]interface{})
			cidrListMap["acl_strategy"] = cidrListChild1Raw["AclStrategy"]
			cidrListMap["cidr"] = cidrListChild1Raw["Cidr"]

			cidrListMaps = append(cidrListMaps, cidrListMap)
		}
	}
	if objectRaw["CidrList"] != nil {
		if err := d.Set("cidr_list", cidrListMaps); err != nil {
			return err
		}
	}

	d.Set("endpoint_type", d.Id())

	return nil
}

func resourceAliCloudMessageServiceEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}

	if d.HasChange("cidr_list") {
		oldEntry, newEntry := d.GetChange("cidr_list")
		removed := oldEntry.(*schema.Set).Difference(newEntry.(*schema.Set)).List()
		added := newEntry.(*schema.Set).Difference(oldEntry.(*schema.Set)).List()

		if len(removed) > 0 {
			action := "RevokeEndpointAcl"
			conn, err := client.NewMnsClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["EndpointType"] = d.Id()
			request["RegionId"] = client.RegionId
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)

			allowCidrList := make([]interface{}, 0)
			for _, cidrList := range removed {
				cidrListArg := cidrList.(map[string]interface{})

				if fmt.Sprint(cidrListArg["acl_strategy"]) == "allow" {
					allowCidrList = append(allowCidrList, cidrListArg["cidr"])
				}
			}

			if len(allowCidrList) > 0 {
				request["AclStrategy"] = "allow"
				request["CidrList"] = convertListToCommaSeparate(allowCidrList)

				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), query, request, &runtime)
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
			}
		}

		if len(added) > 0 {
			action := "AuthorizeEndpointAcl"
			conn, err := client.NewMnsClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["EndpointType"] = d.Id()
			request["RegionId"] = client.RegionId
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)

			allowCidrList := make([]interface{}, 0)
			for _, cidrList := range added {
				cidrListArg := cidrList.(map[string]interface{})

				if fmt.Sprint(cidrListArg["acl_strategy"]) == "allow" {
					allowCidrList = append(allowCidrList, cidrListArg["cidr"])
				}
			}

			if len(allowCidrList) > 0 {
				request["AclStrategy"] = "allow"
				request["CidrList"] = convertListToCommaSeparate(allowCidrList)

				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), query, request, &runtime)
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
			}
		}
	}
	return resourceAliCloudMessageServiceEndpointRead(d, meta)
}

func resourceAliCloudMessageServiceEndpointDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DisableEndpoint"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["EndpointType"] = d.Id()
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), query, request, &runtime)

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
