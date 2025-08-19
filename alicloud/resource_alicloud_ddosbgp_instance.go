package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDdosBgpInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDdosBgpInstanceCreate,
		Read:   resourceAliCloudDdosBgpInstanceRead,
		Update: resourceAliCloudDdosBgpInstanceUpdate,
		Delete: resourceAliCloudDdosBgpInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(26 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"base_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  20,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"ip_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"IPv4", "IPv6"}, false),
			},
			"normal_bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"period": {
				Type:         schema.TypeInt,
				ValidateFunc: IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:     true,
				Default:      12,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Enterprise", "Professional"}, false),
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field `name` has been deprecated from provider version 1.259.0. New field `instance_name` instead.",
			},
		},
	}
}

func resourceAliCloudDdosBgpInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	parameterMapList := make([]map[string]interface{}, 0)
	if v, ok := d.GetOk("type"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Type",
			"Value": convertDdosBgpInstanceInstanceTypeRequest(v),
		})
	} else {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Type",
			"Value": "1",
		})
	}
	if v, ok := d.GetOk("ip_type"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "IpType",
			"Value": convertDdosBgpInstanceInstanceListIpTypeRequest(v),
		})
	}
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "Region",
		"Value": client.RegionId,
	})
	if v, ok := d.GetOk("normal_bandwidth"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "NormalBandwidth",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("ip_count"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "IpCount",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("base_bandwidth"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "BaseBandwidth",
			"Value": v,
		})
	}
	if v, ok := d.GetOk("bandwidth"); ok {
		parameterMapList = append(parameterMapList, map[string]interface{}{
			"Code":  "Bandwidth",
			"Value": v,
		})
	}
	request["Parameter"] = parameterMapList

	request["ProductCode"] = "ddos"
	request["ProductType"] = "ddosbgp"
	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = v
	} else {
		request["Period"] = 1
	}
	request["SubscriptionType"] = "Subscription"
	var endpoint string
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
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

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	ddosBgpServiceV2 := DdosBgpServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"1"}, d.Timeout(schema.TimeoutCreate), 50*time.Second, ddosBgpServiceV2.DdosBgpInstanceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudDdosBgpInstanceUpdate(d, meta)
}

func resourceAliCloudDdosBgpInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddosBgpServiceV2 := DdosBgpServiceV2{client}

	objectRaw, err := ddosBgpServiceV2.DescribeDdosBgpInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ddosbgp_instance DescribeDdosBgpInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_name", objectRaw["Remark"])
	d.Set("ip_type", convertDdosBgpInstanceInstanceListIpTypeResponse(objectRaw["IpType"]))
	d.Set("resource_group_id", objectRaw["ResourceGroupId"])
	d.Set("status", objectRaw["Status"])
	d.Set("type", convertDdosBgpInstanceInstanceListInstanceTypeResponse(objectRaw["InstanceType"]))
	d.Set("name", objectRaw["Remark"])

	objectRaw, err = ddosBgpServiceV2.DescribeInstanceListTagResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", objectRaw)
	d.Set("tags", tagsToMap(tagsMaps))

	objectRaw, err = ddosBgpServiceV2.DescribeInstanceDescribeInstanceSpecs(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("bandwidth", objectRaw["Bandwidth"])
	d.Set("base_bandwidth", objectRaw["PackBasicThre"])
	d.Set("ip_count", objectRaw["IpSpec"])
	d.Set("normal_bandwidth", objectRaw["NormalBandwidth"])

	return nil
}

func resourceAliCloudDdosBgpInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "MoveResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId"] = d.Id()
	request["ResourceRegionId"] = client.RegionId
	if d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["ResourceType"] = "INSTANCE"
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ddosbgp", "2018-07-20", action, query, request, true)
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
	update = false
	action = "ModifyRemark"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("instance_name") {
		update = true

		request["Remark"] = d.Get("instance_name")
	}

	if d.HasChange("name") {
		update = true

		request["Remark"] = d.Get("name")
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ddosbgp", "2018-07-20", action, query, request, true)
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

	if d.HasChange("tags") {
		ddosBgpServiceV2 := DdosBgpServiceV2{client}
		if err := ddosBgpServiceV2.SetResourceTags(d, "INSTANCE"); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAliCloudDdosBgpInstanceRead(d, meta)
}

func resourceAliCloudDdosBgpInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Instance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func convertDdosBgpInstanceInstanceListIpTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "v4":
		return "IPv4"
	case "v6":
		return "IPv6"
	}
	return source
}

func convertDdosBgpInstanceInstanceListIpTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "IPv4":
		return "v4"
	case "IPv6":
		return "v6"
	}
	return source
}

func convertDdosBgpInstanceInstanceListInstanceTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "0":
		return "Professional"
	case "1":
		return "Enterprise"
	}
	return source
}

func convertDdosBgpInstanceInstanceTypeRequest(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "Enterprise":
		return "1"
	case "Professional":
		return "0"
	}
	return source
}
