package privatezone

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceEndpointCreate,
		Read:   resourceEndpointRead,
		Update: resourceEndpointUpdate,
		Delete: resourceEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"endpoint_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_configs": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 6,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_region_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddResolverEndpoint"
	request := make(map[string]interface{})
	var err error
	request["Name"] = d.Get("endpoint_name")
	for k, ipConfig := range d.Get("ip_configs").(*schema.Set).List() {
		ipConfigArg := ipConfig.(map[string]interface{})
		request[fmt.Sprintf("IpConfig.%d.VSwitchId", k+1)] = ipConfigArg["vswitch_id"]
		request[fmt.Sprintf("IpConfig.%d.AzId", k+1)] = ipConfigArg["zone_id"]
		request[fmt.Sprintf("IpConfig.%d.CidrBlock", k+1)] = ipConfigArg["cidr_block"]
		request[fmt.Sprintf("IpConfig.%d.Ip", k+1)] = ipConfigArg["ip"]
	}

	request["Lang"] = "en"
	request["SecurityGroupId"] = d.Get("security_group_id")
	request["VpcId"] = d.Get("vpc_id")
	request["VpcRegionId"] = d.Get("vpc_region_id")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ResolverEndpoint.Operation.ExceedsFrequency"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_pvtz_endpoint", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["EndpointId"]))
	pvtzService := Service{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"SUCCESS"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, pvtzService.EndpointStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceEndpointRead(d, meta)
}
func resourceEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := Service{client}
	object, err := pvtzService.DescribeEndpoint(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pvtz_endpoint pvtzService.DescribeEndpoint Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("endpoint_name", object["Name"])
	ipConfigsSli := make([]map[string]interface{}, 0)
	if ipConfigs, ok := object["IpConfigs"].([]interface{}); ok {
		for _, ipConfigArgs := range ipConfigs {
			ipConfigArg := ipConfigArgs.(map[string]interface{})
			ipConfigsMap := make(map[string]interface{})
			ipConfigsMap["zone_id"] = ipConfigArg["AzId"]
			ipConfigsMap["cidr_block"] = ipConfigArg["CidrBlock"]
			ipConfigsMap["ip"] = ipConfigArg["Ip"]
			ipConfigsMap["vswitch_id"] = ipConfigArg["VSwitchId"]
			ipConfigsSli = append(ipConfigsSli, ipConfigsMap)
		}
	}
	d.Set("ip_configs", ipConfigsSli)
	d.Set("security_group_id", object["SecurityGroupId"])
	d.Set("status", object["Status"])
	d.Set("vpc_id", object["VpcId"])
	d.Set("vpc_region_id", object["VpcRegionId"])
	return nil
}
func resourceEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	request := map[string]interface{}{
		"EndpointId": d.Id(),
	}

	update := false
	if d.HasChange("endpoint_name") {
		update = true
		if v, ok := d.GetOk("endpoint_name"); ok {
			request["Name"] = v
		}
	}

	if d.HasChange("ip_configs") {
		update = true
		for k, ipConfig := range d.Get("ip_configs").(*schema.Set).List() {
			ipConfigArg := ipConfig.(map[string]interface{})
			request[fmt.Sprintf("IpConfig.%d.VSwitchId", k+1)] = ipConfigArg["vswitch_id"]
			request[fmt.Sprintf("IpConfig.%d.AzId", k+1)] = ipConfigArg["zone_id"]
			request[fmt.Sprintf("IpConfig.%d.CidrBlock", k+1)] = ipConfigArg["cidr_block"]
			request[fmt.Sprintf("IpConfig.%d.Ip", k+1)] = ipConfigArg["ip"]
		}
	}
	request["Lang"] = "en"
	if update {
		action := "UpdateResolverEndpoint"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, false)
			if err != nil {

				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		pvtzService := Service{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"SUCCESS"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, pvtzService.EndpointStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}
	return resourceEndpointRead(d, meta)
}
func resourceEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteResolverEndpoint"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"EndpointId": d.Id(),
	}

	request["Lang"] = "en"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ResolverEndpoint.NotExists"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
