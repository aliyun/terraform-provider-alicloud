package privatezone

import (
	"context"
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

func ResourceZoneAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceZoneAttachmentCreate,
		Read:   resourceZoneAttachmentRead,
		Update: resourceZoneAttachmentUpdate,
		Delete: resourceZoneAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, v interface{}) error {
			if d.HasChange("vpcs") {
				d.SetNewComputed("vpc_ids")
			} else if d.HasChange("vpc_ids") {
				d.SetNewComputed("vpcs")
			}
			return nil
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_client_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpcs": {
				Type:          schema.TypeSet,
				Computed:      true,
				ConflictsWith: []string{"vpc_ids"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"vpcs"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceZoneAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "BindZoneVpc"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	if v, ok := d.GetOk("user_client_ip"); ok {
		request["UserClientIp"] = v
	}

	vpcIdMap := make(map[string]string)
	if v, ok := d.GetOk("vpcs"); ok {
		vpcs := make([]map[string]interface{}, len(v.(*schema.Set).List()))
		for i, j := range v.(*schema.Set).List() {
			vpcs[i] = make(map[string]interface{})
			regionId := j.(map[string]interface{})["region_id"]
			if regionId == "" {
				regionId = client.RegionId
			}
			vpcs[i]["RegionId"] = regionId
			vpcIdMap[j.(map[string]interface{})["vpc_id"].(string)] = j.(map[string]interface{})["vpc_id"].(string)
			vpcs[i]["VpcId"] = j.(map[string]interface{})["vpc_id"]
		}
		request["Vpcs"] = vpcs

	} else {
		vpcIds := d.Get("vpc_ids").(*schema.Set).List()
		vpcs := make([]map[string]interface{}, len(vpcIds))
		for i, j := range vpcIds {
			vpcs[i] = make(map[string]interface{})
			vpcs[i]["RegionId"] = client.RegionId
			vpcs[i]["VpcId"] = j.(string)
			vpcIdMap[j.(string)] = j.(string)
		}
		request["Vpcs"] = vpcs
	}

	request["ZoneId"] = d.Get("zone_id")
	wait := helper.IncrementalWait(3*time.Second, 2*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"System.Busy", "Throttling.User", "ServiceUnavailable", "Zone.NotExists"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_pvtz_zone_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["ZoneId"]))
	pvtzService := Service{client}
	if err := pvtzService.WaitForZoneAttachment(d.Id(), vpcIdMap, DefaultTimeout); err != nil {
		return errs.WrapError(err)
	}

	return resourceZoneAttachmentRead(d, meta)
}
func resourceZoneAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := Service{client}
	object, err := pvtzService.DescribeZoneAttachment(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_private_zone_zone_attachment pvtzService.DescribeZoneAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("zone_id", d.Id())
	vpcIds := make([]string, 0)

	vpc := make([]map[string]interface{}, 0)
	if vpcList, ok := object["BindVpcs"].(map[string]interface{})["Vpc"].([]interface{}); ok {
		for _, v := range vpcList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"region_id": m1["RegionId"],
					"vpc_id":    m1["VpcId"],
				}
				vpcIds = append(vpcIds, m1["VpcId"].(string))
				vpc = append(vpc, temp1)

			}
		}
	}
	if err := d.Set("vpcs", vpc); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("vpc_ids", vpcIds); err != nil {
		return errs.WrapError(err)
	}
	return nil
}
func resourceZoneAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	if d.HasChange("vpcs") || d.HasChange("vpc_ids") {
		request := map[string]interface{}{
			"ZoneId": d.Id(),
		}
		vpcIdMap := make(map[string]string)
		vpcs := make([]map[string]interface{}, len(d.Get("vpcs").(*schema.Set).List()))
		if d.HasChange("vpcs") && len(vpcs) != 0 {
			for i, j := range d.Get("vpcs").(*schema.Set).List() {
				vpcs[i] = make(map[string]interface{})
				regionId := j.(map[string]interface{})["region_id"]
				if regionId == "" {
					regionId = client.RegionId
				}
				vpcs[i]["RegionId"] = regionId
				vpcs[i]["VpcId"] = j.(map[string]interface{})["vpc_id"]
				vpcIdMap[j.(map[string]interface{})["vpc_id"].(string)] = j.(map[string]interface{})["vpc_id"].(string)
			}
		} else {
			vpcs = make([]map[string]interface{}, len(d.Get("vpc_ids").(*schema.Set).List()))
			for i, j := range d.Get("vpc_ids").(*schema.Set).List() {
				vpcs[i] = make(map[string]interface{})
				vpcs[i]["RegionId"] = client.RegionId
				vpcs[i]["VpcId"] = j.(string)
				vpcIdMap[j.(string)] = j.(string)
			}
		}
		request["Vpcs"] = vpcs
		if _, ok := d.GetOk("lang"); ok {
			request["Lang"] = d.Get("lang")
		}
		if _, ok := d.GetOk("user_client_ip"); ok {
			request["UserClientIp"] = d.Get("user_client_ip")
		}
		action := "BindZoneVpc"
		wait := helper.IncrementalWait(3*time.Second, 2*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, false)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"System.Busy", "Throttling.User", "ServiceUnavailable", "Zone.NotExists"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		pvtzService := Service{client}
		if err := pvtzService.WaitForZoneAttachment(d.Id(), vpcIdMap, DefaultTimeout); err != nil {
			return errs.WrapError(err)
		}
	}
	return resourceZoneAttachmentRead(d, meta)
}
func resourceZoneAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "BindZoneVpc"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"ZoneId": d.Id(),
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("user_client_ip"); ok {
		request["UserClientIp"] = v
	}
	vpcs := make([]map[string]interface{}, 0)
	request["Vpcs"] = vpcs
	wait := helper.IncrementalWait(3*time.Second, 2*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"System.Busy", "Throttling.User"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"Zone.NotExists", "ZoneVpc.NotExists.VpcId"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	pvtzService := Service{client}
	return errs.WrapError(pvtzService.WaitForZoneAttachmentStatus(d.Id(), Deleted, DefaultTimeout))
}
