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
			"resource_record": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{string(RecordA), string(RecordCNAME),
					string(RecordMX), string(RecordTXT), string(RecordPTR)}),
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validateIntegerInRange(1, 50),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("type").(string) != string(RecordMX)
				},
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
		},
	}
}

func resourceAlicloudPvtzZoneRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := pvtz.CreateAddZoneRecordRequest()

	if v, ok := d.GetOk("resource_record"); ok && v.(string) != "" {
		args.Rr = v.(string)
	}

	if v, ok := d.GetOk("type"); ok && v.(string) != "" {
		args.Type = v.(string)
	}

	if v, ok := d.GetOk("value"); ok && v.(string) != "" {
		args.Value = v.(string)
	}

	if v, ok := d.GetOk("zone_id"); ok && v.(string) != "" {
		args.ZoneId = v.(string)
	}

	if v, ok := d.GetOk("priority"); ok && v != nil {
		args.Priority = requests.NewInteger(d.Get("priority").(int))
	}

	if v, ok := d.GetOk("ttl"); ok && v != nil {
		args.Ttl = requests.NewInteger(d.Get("ttl").(int))
	}

	raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
		return pvtzClient.AddZoneRecord(args)
	})

	if err != nil {
		return fmt.Errorf("AddZoneRecord got a error: %#v", err)
	}
	resp, _ := raw.(*pvtz.AddZoneRecordResponse)
	if resp == nil {
		return fmt.Errorf("AddZoneRecord got a nil response: %#v", resp)
	}

	d.SetId(strconv.Itoa(resp.RecordId) + ":" + args.ZoneId)

	return resourceAlicloudPvtzZoneRecordUpdate(d, meta)
}

func resourceAlicloudPvtzZoneRecordUpdate(d *schema.ResourceData, meta interface{}) error {

	attributeUpdate := false

	args := pvtz.CreateUpdateZoneRecordRequest()
	recordIdStr, _, _ := getRecordIdAndZoneId(d, meta)
	recordId, _ := strconv.Atoi(recordIdStr)
	args.RecordId = requests.NewInteger(recordId)
	args.Rr = d.Get("resource_record").(string)
	args.Type = d.Get("type").(string)
	args.Value = d.Get("value").(string)

	if d.HasChange("resource_record") {
		attributeUpdate = true
	}

	if d.HasChange("type") {
		attributeUpdate = true
	}

	if d.HasChange("value") {
		attributeUpdate = true
	}

	if d.HasChange("priority") {
		args.Priority = requests.NewInteger(d.Get("priority").(int))
		attributeUpdate = true
	}

	if d.HasChange("ttl") {
		args.Ttl = requests.NewInteger(d.Get("ttl").(int))
		attributeUpdate = true
	}

	if attributeUpdate {
		client := meta.(*connectivity.AliyunClient)
		_, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.UpdateZoneRecord(args)
		})
		if err != nil {
			return err
		}
	}

	return resourceAlicloudPvtzZoneRecordRead(d, meta)

}

func resourceAlicloudPvtzZoneRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}

	recordIdStr, zoneId, _ := getRecordIdAndZoneId(d, meta)
	recordId, e := strconv.Atoi(recordIdStr)
	if e != nil {
		return e
	}

	record, err := pvtzService.DescribeZoneRecord(recordId, zoneId)
	if err != nil {
		if NotFoundError(e) {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("record_id", record.RecordId)
	d.Set("zone_id", zoneId)
	d.Set("resource_record", record.Rr)
	d.Set("type", record.Type)
	d.Set("value", record.Value)
	d.Set("ttl", record.Ttl)
	d.Set("priority", record.Priority)
	d.Set("status", record.Status)

	return nil
}

func resourceAlicloudPvtzZoneRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}

	request := pvtz.CreateDeleteZoneRecordRequest()
	recordIdStr, zoneId, _ := getRecordIdAndZoneId(d, meta)
	recordId, err := strconv.Atoi(recordIdStr)
	if err != nil {
		return err
	}
	request.RecordId = requests.NewInteger(recordId)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.DeleteZoneRecord(request)
		})

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error deleting zone record failed: %#v", err))
		}

		if _, e := pvtzService.DescribeZoneRecord(recordId, zoneId); e != nil {
			if NotFoundError(e) {
				return nil
			}

			return resource.NonRetryableError(e)
		}

		return nil
	})
}

func getRecordIdAndZoneId(d *schema.ResourceData, meta interface{}) (string, string, error) {
	return splitRecordIdAndZoneId(d.Id())
}

func splitRecordIdAndZoneId(s string) (string, string, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid resource id")
	}
	return parts[0], parts[1], nil
}
