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
				ValidateFunc: validateIntegerInRange(1, 10),
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
	client := meta.(*AliyunClient)
	conn := client.pvtzconn

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

	resp, err := conn.AddZoneRecord(args)

	if err != nil {
		return fmt.Errorf("AddZoneRecord got a error: %#v", err)
	}
	if resp == nil {
		return fmt.Errorf("AddZoneRecord got a nil response: %#v", resp)
	}

	d.SetId(strconv.Itoa(resp.RecordId) + ":" + args.ZoneId)

	return resourceAlicloudPvtzZoneRecordUpdate(d, meta)
}

func resourceAlicloudPvtzZoneRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)

	attributeUpdate := false

	args := pvtz.CreateUpdateZoneRecordRequest()
	recordIdStr, _, _ := getRecordIdAndZoneId(d, meta)
	recordId, _ := strconv.Atoi(recordIdStr)
	args.RecordId = requests.NewInteger(recordId)
	args.Rr = d.Get("resource_record").(string)

	if d.HasChange("type") {
		d.SetPartial("type")
		args.Type = d.Get("type").(string)

		attributeUpdate = true
	}

	if d.HasChange("value") {
		d.SetPartial("value")
		args.Value = d.Get("value").(string)

		attributeUpdate = true
	}

	if d.HasChange("priority") {
		d.SetPartial("priority")
		args.Priority = requests.NewInteger(d.Get("priority").(int))

		attributeUpdate = true
	}

	if d.HasChange("ttl") {
		d.SetPartial("ttl")
		args.Ttl = requests.NewInteger(d.Get("ttl").(int))

		attributeUpdate = true
	}

	if attributeUpdate {
		if _, err := meta.(*AliyunClient).pvtzconn.UpdateZoneRecord(args); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAlicloudPvtzZoneRecordRead(d, meta)

}

func resourceAlicloudPvtzZoneRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	recordIdStr, zoneId, _ := getRecordIdAndZoneId(d, meta)
	recordId, e := strconv.Atoi(recordIdStr)
	if e != nil {
		return e
	}

	record, err := client.DescribeZoneRecord(recordId, zoneId)
	if err != nil {
		if NotFoundError(e) {
			d.SetId("")
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
	client := meta.(*AliyunClient)
	conn := client.pvtzconn

	request := pvtz.CreateDeleteZoneRecordRequest()
	recordIdStr, zoneId, _ := getRecordIdAndZoneId(d, meta)
	recordId, err := strconv.Atoi(recordIdStr)
	if err != nil {
		return err
	}
	request.RecordId = requests.NewInteger(recordId)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DeleteZoneRecord(request)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error deleting zone record failed: %#v", err))
		}

		if _, e := client.DescribeZoneRecord(recordId, zoneId); e != nil {
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
