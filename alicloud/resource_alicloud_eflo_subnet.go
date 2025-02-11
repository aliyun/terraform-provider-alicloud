package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEfloSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEfloSubnetCreate,
		Read:   resourceAlicloudEfloSubnetRead,
		Update: resourceAlicloudEfloSubnetUpdate,
		Delete: resourceAlicloudEfloSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cidr": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"gmt_modified": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"message": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"resource_group_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"subnet_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"subnet_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"type": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"OOB", "LB"}, false),
			},
			"vpd_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"zone_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudEfloSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloService := EfloService{client}
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("vpd_id"); ok {
		request["VpdId"] = v
	}
	request["SubnetName"] = d.Get("subnet_name")
	request["ZoneId"] = d.Get("zone_id")
	if v, ok := d.GetOk("cidr"); ok {
		request["Cidr"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}

	request["ClientToken"] = buildClientToken("CreateSubnet")
	var response map[string]interface{}
	action := "CreateSubnet"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("eflo", "2022-05-30", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eflo_subnet", action, AlibabaCloudSdkGoERROR)
	}
	subnetIdValue, err := jsonpath.Get("$.Content.SubnetId", response)
	if err != nil || subnetIdValue == nil {
		return WrapErrorf(err, IdMsg, "alicloud_eflo_subnet")
	}

	d.SetId(fmt.Sprint(request["VpdId"], ":", subnetIdValue))
	stateConf := BuildStateConf([]string{}, []string{"Available", "Waiting4AssignPort"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, efloService.EfloSubnetStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudEfloSubnetRead(d, meta)
}

func resourceAlicloudEfloSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloService := EfloService{client}

	object, err := efloService.DescribeEfloSubnet(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eflo_subnet efloService.DescribeEfloSubnet Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("vpd_id", parts[0])
	d.Set("subnet_id", parts[1])
	d.Set("cidr", object["Cidr"])
	d.Set("create_time", object["CreateTime"])
	d.Set("gmt_modified", object["GmtModified"])
	d.Set("message", object["Message"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["Status"])
	d.Set("subnet_name", object["SubnetName"])
	d.Set("type", object["Type"])
	d.Set("zone_id", object["ZoneId"])

	return nil
}

func resourceAlicloudEfloSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"VpdId":    parts[0],
		"SubnetId": parts[1],
		"RegionId": client.RegionId,
	}

	if d.HasChange("subnet_name") {
		update = true
		request["SubnetName"] = d.Get("subnet_name")
	}
	request["ZoneId"] = d.Get("zone_id")

	if update {
		action := "UpdateSubnet"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("eflo", "2022-05-30", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudEfloSubnetRead(d, meta)
}

func resourceAlicloudEfloSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloService := EfloService{client}
	var response map[string]interface{}
	var err error
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"VpdId":    parts[0],
		"SubnetId": parts[1],
		"ZoneId":   d.Get("zone_id"),
	}

	action := "DeleteSubnet"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("eflo", "2022-05-30", action, nil, request, false)
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
		if NotFoundError(err) || IsExpectedErrors(err, []string{"1003"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, efloService.EfloSubnetStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
