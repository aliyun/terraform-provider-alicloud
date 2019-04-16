package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{string(RecordA), string(RecordCNAME),
					string(RecordMX), string(RecordTXT), string(RecordPTR)}),
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
				ValidateFunc: validateIntegerInRange(1, 50),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("type").(string) != string(RecordMX)
				},
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
		},
	}
}

func resourceAlicloudPvtzZoneRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := pvtz.CreateAddZoneRecordRequest()

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
	invoker := PvtzInvoker()
	var raw interface{}
	var err error
	err = invoker.Run(func() error {
		raw, err = client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.AddZoneRecord(request)
		})
		return err
	})

	addDebug(request.GetActionName(), raw)

	if err != nil {

		if IsExceptedErrors(err, []string{RecordInvalidConflict}) {
			req := pvtz.CreateDescribeZoneRecordsRequest()
			req.ZoneId = request.ZoneId
			req.Keyword = request.Rr
			req.PageSize = requests.NewInteger(PageSizeXLarge)
			req.PageNumber = requests.NewInteger(1)
			for {
				var raw interface{}
				var err error
				if err = invoker.Run(func() error {
					raw, err = client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
						return pvtzClient.DescribeZoneRecords(req)
					})

					addDebug(req.GetActionName(), raw)

					return err
				}); err != nil {
					return WrapErrorf(err, DefaultErrorMsg, "alicloud_pvtz_zone_record", request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				response, _ := raw.(*pvtz.DescribeZoneRecordsResponse)
				if response != nil && len(response.Records.Record) > 0 {
					for _, rec := range response.Records.Record {
						if rec.Rr == request.Rr && rec.Type == request.Type && rec.Value == request.Value {
							d.SetId(fmt.Sprintf("%d%s%s", rec.RecordId, COLON_SEPARATED, request.ZoneId))
							return resourceAlicloudPvtzZoneRecordRead(d, meta)
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

	return resourceAlicloudPvtzZoneRecordRead(d, meta)
}

func resourceAlicloudPvtzZoneRecordUpdate(d *schema.ResourceData, meta interface{}) error {

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

	if update {
		client := meta.(*connectivity.AliyunClient)
		invoker := PvtzInvoker()

		if err := invoker.Run(func() error {
			raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
				return pvtzClient.UpdateZoneRecord(request)
			})
			addDebug(request.GetActionName(), raw)
			return err
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudPvtzZoneRecordRead(d, meta)

}

func resourceAlicloudPvtzZoneRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	object, err := pvtzService.DescribeZoneRecord(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("record_id", object.RecordId)
	d.Set("zone_id", parts[1])
	d.Set("resource_record", object.Rr)
	d.Set("type", object.Type)
	d.Set("value", object.Value)
	d.Set("ttl", object.Ttl)
	d.Set("priority", object.Priority)
	d.Set("status", object.Status)

	return nil
}

func resourceAlicloudPvtzZoneRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}

	request := pvtz.CreateDeleteZoneRecordRequest()
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
			if IsExceptedErrors(err, []string{ZoneNotExists, ZoneVpcNotExists}) {
				return nil
			}
			if IsExceptedErrors(err, []string{PvtzThrottlingUser, PvtzSystemBusy}) {
				time.Sleep(time.Duration(2) * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(request.GetActionName(), raw)
		return nil
	})

	if err != nil {
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
