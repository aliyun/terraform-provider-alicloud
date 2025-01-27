package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDdosbgpInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDdosbgpInstanceCreate,
		Read:   resourceAlicloudDdosbgpInstanceRead,
		Update: resourceAlicloudDdosbgpInstanceUpdate,
		Delete: resourceAlicloudDdosbgpInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Required:     false,
				Default:      string(Enterprise),
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Enterprise", "Professional"}, false),
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Required:     false,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"base_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  20,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"ip_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPv4", "IPv6"}, false),
			},
			"ip_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"period": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:     true,
				Default:      12,
			},
			"normal_bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudDdosbgpInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ddosbgpInstanceType := "1"
	if d.Get("type").(string) == string(Professional) {
		ddosbgpInstanceType = "0"
	}

	ddosbgpInstanceIpType := "v4"
	if d.Get("ip_type").(string) == string(IPv6) {
		ddosbgpInstanceIpType = "v6"
	}

	var response map[string]interface{}
	var err error
	var endpoint string
	action := "CreateInstance"
	request := make(map[string]interface{})

	request["ProductCode"] = "ddos"
	request["ProductType"] = "ddosbgp"
	request["SubscriptionType"] = "Subscription"

	parameterMapList := make([]map[string]interface{}, 0)
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "Region",
		"Value": client.RegionId,
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "Type",
		"Value": ddosbgpInstanceType,
	})

	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "IpType",
		"Value": ddosbgpInstanceIpType,
	})

	if v, ok := d.GetOk("base_bandwidth"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "BaseBandwidth",
			"Value": v,
		})
	}

	if v, ok := d.GetOk("normal_bandwidth"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "NormalBandwidth",
			"Value": v,
		})
	}

	if v, ok := d.GetOk("bandwidth"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Bandwidth",
			"Value": v,
		})
	}

	if v, ok := d.GetOk("ip_count"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "IpCount",
			"Value": v,
		})
	}

	request["Parameter"] = parameterMapList
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	} else {
		request["Period"] = 1
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, false, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddosbgp_instance", action, AlibabaCloudSdkGoERROR)
	}

	response = response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(response["InstanceId"]))

	return resourceAlicloudDdosbgpInstanceUpdate(d, meta)
}

func resourceAlicloudDdosbgpInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddosbgpService := DdosbgpService{client}
	object, err := ddosbgpService.DescribeDdosbgpInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ddosbgp_instance DescribeInstanceList Failed!!! %s", err)
			d.SetId("")
			return nil
		}

		return WrapError(err)
	}

	specInfo, err := ddosbgpService.DescribeDdosbgpInstanceSpec(d.Id())
	if err != nil {
		return WrapError(err)
	}

	ddosbgpInstanceType := string(Enterprise)
	if fmt.Sprint(object["InstanceType"]) == "0" {
		ddosbgpInstanceType = string(Professional)
	}

	d.Set("name", object["Remark"])
	d.Set("bandwidth", specInfo["PackConfig"].(map[string]interface{})["Bandwidth"])
	d.Set("base_bandwidth", specInfo["PackConfig"].(map[string]interface{})["PackBasicThre"])
	d.Set("normal_bandwidth", specInfo["PackConfig"].(map[string]interface{})["NormalBandwidth"])
	d.Set("ip_type", object["IpType"])
	d.Set("ip_count", specInfo["PackConfig"].(map[string]interface{})["IpSpec"])
	d.Set("type", ddosbgpInstanceType)

	return nil
}

func resourceAlicloudDdosbgpInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false
	action := "ModifyRemark"
	request := map[string]interface{}{
		"InstanceId": d.Id(),
		"RegionId":   client.RegionId,
		"Remark":     d.Get("name"),
	}

	if d.HasChange("name") {
		update = true
	}
	if update {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ddosbgp", "2018-07-20", action, nil, request, false)
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

	return resourceAlicloudDdosbgpInstanceRead(d, meta)
}

func resourceAlicloudDdosbgpInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AlicloudDdosbgpInstance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
