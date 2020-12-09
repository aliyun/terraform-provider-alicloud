package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPvtzZoneRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPvtzZoneRecordCreate,
		Read:   resourceAlicloudPvtzZoneRecordRead,
		Update: resourceAlicloudPvtzZoneRecordUpdate,
		Delete: resourceAlicloudPvtzZoneRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"resource_record": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "CNAME", "MX", "TXT", "PTR", "SRV"}, false),
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 50),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("type").(string) != "MX"
				},
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
			"record_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudPvtzZoneRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := pvtz.CreateAddZoneRecordRequest()
	request.RegionId = client.RegionId

	if v, ok := d.GetOk("resource_record"); ok && v.(string) != "" {
		request.Rr = v.(string)
	}

	if v, ok := d.GetOk("type"); ok && v.(string) != "" {
		request.Type = v.(string)
	}

	if v, ok := d.GetOk("value"); ok && v.(string) != "" {
		request.Value = v.(string)
	}

	if v, ok := d.GetOk("zone_id"); ok && v.(string) != "" {
		request.ZoneId = v.(string)
	}

	if v, ok := d.GetOk("priority"); ok && v != nil {
		request.Priority = requests.NewInteger(d.Get("priority").(int))
	}

	if v, ok := d.GetOk("ttl"); ok && v != nil {
		request.Ttl = requests.NewInteger(d.Get("ttl").(int))
	}

	// API AddZoneRecord has a throttling limitation 20qps which one use only can send 20 requests in one second.
	var raw interface{}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.AddZoneRecord(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ServiceUnavailable, ThrottlingUser, "System.Busy"}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {

		if IsExpectedErrors(err, []string{"Record.Invalid.Conflict"}) {
			req := pvtz.CreateDescribeZoneRecordsRequest()
			req.RegionId = client.RegionId
			req.ZoneId = request.ZoneId
			req.Keyword = request.Rr
			req.PageSize = requests.NewInteger(PageSizeXLarge)
			req.PageNumber = requests.NewInteger(1)
			for {
				err = resource.Retry(5*time.Minute, func() *resource.RetryError {
					raw, err = client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
						return pvtzClient.DescribeZoneRecords(req)
					})
					if err != nil {
						if IsExpectedErrors(err, []string{ServiceUnavailable, ThrottlingUser, "System.Busy"}) {
							time.Sleep(5 * time.Second)
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(req.GetActionName(), raw, req.RpcRequest, req)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, "alicloud_pvtz_zone_record", request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				response, _ := raw.(*pvtz.DescribeZoneRecordsResponse)
				if response != nil && len(response.Records.Record) > 0 {
					for _, rec := range response.Records.Record {
						if rec.Rr == request.Rr && rec.Type == request.Type && rec.Value == request.Value {
							d.SetId(fmt.Sprintf("%d%s%s", rec.RecordId, COLON_SEPARATED, request.ZoneId))
							return resourceAlicloudPvtzZoneRecordUpdate(d, meta)
						}
					}
				}
				if len(response.Records.Record) < PageSizeXLarge {
					break
				}

				if page, err := getNextpageNumber(req.PageNumber); err != nil {
					return WrapError(err)
				} else {
					req.PageNumber = page
				}
			}
		}

		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pvtz_zone_record", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*pvtz.AddZoneRecordResponse)

	d.SetId(fmt.Sprintf("%d%s%s", response.RecordId, COLON_SEPARATED, request.ZoneId))

	return resourceAlicloudPvtzZoneRecordUpdate(d, meta)
}

func resourceAlicloudPvtzZoneRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)
	update := false

	request := pvtz.CreateUpdateZoneRecordRequest()
	recordIdStr, _, _ := getRecordIdAndZoneId(d, meta)
	recordId, _ := strconv.Atoi(recordIdStr)
	request.RecordId = requests.NewInteger(recordId)
	request.Rr = d.Get("resource_record").(string)
	request.Type = d.Get("type").(string)
	request.Value = d.Get("value").(string)

	if d.HasChange("type") {
		update = true
	}

	if d.HasChange("value") {
		update = true
	}

	if d.HasChange("priority") {
		request.Priority = requests.NewInteger(d.Get("priority").(int))
		update = true
	}

	if d.HasChange("ttl") {
		request.Ttl = requests.NewInteger(d.Get("ttl").(int))
		update = true
	}

	if update && !d.IsNewResource() {
		request.RegionId = client.RegionId
		if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
				return pvtzClient.UpdateZoneRecord(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{ServiceUnavailable, ThrottlingUser, "System.Busy"}) {
					time.Sleep(5 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			d.SetPartial("type")
			d.SetPartial("value")
			d.SetPartial("priority")
			d.SetPartial("ttl")
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	if d.HasChange("remark") {
		request := pvtz.CreateUpdateRecordRemarkRequest()
		request.RegionId = client.RegionId
		recordIdStr, _, _ := getRecordIdAndZoneId(d, meta)
		recordId, _ := strconv.Atoi(recordIdStr)
		request.RecordId = requests.NewInteger(recordId)
		request.Remark = d.Get("remark").(string)
		raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.UpdateRecordRemark(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("remark")
	}
	d.Partial(false)
	return resourceAlicloudPvtzZoneRecordRead(d, meta)

}

func resourceAlicloudPvtzZoneRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	object, err := pvtzService.DescribePvtzZoneRecord(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("record_id", object.RecordId)
	d.Set("zone_id", parts[1])
	d.Set("resource_record", object.Rr)
	d.Set("type", object.Type)
	d.Set("value", object.Value)
	d.Set("ttl", object.Ttl)
	d.Set("priority", object.Priority)
	d.Set("status", object.Status)
	d.Set("remark", object.Remark)

	return nil
}

func resourceAlicloudPvtzZoneRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}

	request := pvtz.CreateDeleteZoneRecordRequest()
	request.RegionId = client.RegionId
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	recordId, err := strconv.Atoi(parts[0])
	if err != nil {
		return WrapError(err)
	}
	request.RecordId = requests.NewInteger(recordId)

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.DeleteZoneRecord(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "System.Busy"}) {
				time.Sleep(time.Duration(2) * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"Zone.NotExists", "ZoneVpc.NotExists.VpcId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(pvtzService.WaitForPvtzZoneRecord(d.Id(), Deleted, DefaultTimeout))
}

func getRecordIdAndZoneId(d *schema.ResourceData, meta interface{}) (string, string, error) {
	return splitRecordIdAndZoneId(d.Id())
}

func splitRecordIdAndZoneId(s string) (string, string, error) {
	parts := strings.Split(s, string(COLON_SEPARATED))
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id")
	}
	return parts[0], parts[1], nil
}
