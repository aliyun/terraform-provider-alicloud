package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudApiGatewayTrafficControl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApiGatewayTrafficControlCreate,
		Read:   resourceAliCloudApiGatewayTrafficControlRead,
		Update: resourceAliCloudApiGatewayTrafficControlUpdate,
		Delete: resourceAliCloudApiGatewayTrafficControlDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"traffic_control_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"traffic_control_unit": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"MINUTE", "HOUR", "DAY"}, false),
			},
			"api_default": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"user_default": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"app_default": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"special_policies": {
				Type:      schema.TypeList,
				Computed:  true,
				Sensitive: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"special_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"specials": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"traffic_value": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"special_key": {
										Type:     schema.TypeString,
										Computed: true,
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
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudApiGatewayTrafficControlCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateTrafficControl"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["TrafficControlName"] = d.Get("traffic_control_name")
	request["TrafficControlUnit"] = d.Get("traffic_control_unit")
	request["ApiDefault"] = d.Get("api_default")
	if v, ok := d.GetOk("user_default"); ok {
		request["UserDefault"] = v
	}
	if v, ok := d.GetOk("app_default"); ok {
		request["AppDefault"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_traffic_control", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TrafficControlId"]))

	return resourceAliCloudApiGatewayTrafficControlRead(d, meta)
}

func resourceAliCloudApiGatewayTrafficControlRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apiGatewayServiceV2 := ApiGatewayServiceV2{client}

	objectRaw, err := apiGatewayServiceV2.DescribeApiGatewayTrafficControl(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_api_gateway_traffic_control DescribeApiGatewayTrafficControl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("traffic_control_name", objectRaw["TrafficControlName"])
	d.Set("traffic_control_unit", objectRaw["TrafficControlUnit"])
	d.Set("api_default", objectRaw["ApiDefault"])
	d.Set("user_default", objectRaw["UserDefault"])
	d.Set("app_default", objectRaw["AppDefault"])
	d.Set("description", objectRaw["Description"])
	d.Set("create_time", objectRaw["CreatedTime"])
	d.Set("modified_time", objectRaw["ModifiedTime"])
	d.Set("region_id", client.RegionId)

	specialPoliciesRaw, _ := jsonpath.Get("$.SpecialPolicies.SpecialPolicy", objectRaw)
	specialPoliciesMaps := make([]map[string]interface{}, 0)
	if specialPoliciesRaw != nil {
		if spList, ok := specialPoliciesRaw.([]interface{}); ok {
			for _, spRaw := range spList {
				spItem, ok := spRaw.(map[string]interface{})
				if !ok {
					continue
				}
				specialPolicyMap := make(map[string]interface{})
				specialPolicyMap["special_type"] = spItem["SpecialType"]

				specialsMaps := make([]map[string]interface{}, 0)
				specialsRaw, _ := jsonpath.Get("$.Specials.Special", spItem)
				if specialsRaw != nil {
					if sList, ok := specialsRaw.([]interface{}); ok {
						for _, sRaw := range sList {
							sItem, ok := sRaw.(map[string]interface{})
							if !ok {
								continue
							}
							specialMap := make(map[string]interface{})
							specialMap["traffic_value"] = sItem["TrafficValue"]
							specialMap["special_key"] = sItem["SpecialKey"]
							specialsMaps = append(specialsMaps, specialMap)
						}
					}
				}
				specialPolicyMap["specials"] = specialsMaps
				specialPoliciesMaps = append(specialPoliciesMaps, specialPolicyMap)
			}
		}
	}
	if err := d.Set("special_policies", specialPoliciesMaps); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAliCloudApiGatewayTrafficControlUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	if d.HasChange("traffic_control_name") || d.HasChange("traffic_control_unit") || d.HasChange("api_default") || d.HasChange("user_default") || d.HasChange("app_default") || d.HasChange("description") {
		action := "ModifyTrafficControl"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		query["TrafficControlId"] = d.Id()
		request["TrafficControlName"] = d.Get("traffic_control_name")
		request["TrafficControlUnit"] = d.Get("traffic_control_unit")
		request["ApiDefault"] = d.Get("api_default")
		request["UserDefault"] = d.Get("user_default")
		request["AppDefault"] = d.Get("app_default")
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
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
		d.SetPartial("traffic_control_name")
		d.SetPartial("traffic_control_unit")
		d.SetPartial("api_default")
		d.SetPartial("user_default")
		d.SetPartial("app_default")
		d.SetPartial("description")
	}

	d.Partial(false)

	return resourceAliCloudApiGatewayTrafficControlRead(d, meta)
}

func resourceAliCloudApiGatewayTrafficControlDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DeleteTrafficControl"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["TrafficControlId"] = d.Id()

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"NotFoundTrafficControl"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
