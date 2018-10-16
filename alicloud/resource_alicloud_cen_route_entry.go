package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudCenRouteEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenRouteEntryCreate,
		Read:   resourceAlicloudCenRouteEntryRead,
		Delete: resourceAlicloudCenRouteEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_table_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCenRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	cenId := d.Get("instance_id").(string)
	vtbId := d.Get("route_table_id").(string)
	cidr := d.Get("cidr_block").(string)
	childInstanceId, childInstanceType, err := client.createCenRouteEntryParas(vtbId)
	if err != nil {
		return fmt.Errorf("Publish CEN route entry encounter an error when query childInstance ID, CEN %s vtb %s cidr %s, error info: %#v.", cenId, vtbId, cidr, err)
	}

	request := cbn.CreatePublishRouteEntriesRequest()
	request.CenId = cenId
	request.ChildInstanceId = childInstanceId
	request.ChildInstanceType = childInstanceType
	request.ChildInstanceRegionId = client.RegionId
	request.ChildInstanceRouteTableId = vtbId
	request.DestinationCidrBlock = cidr

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err = client.cenconn.PublishRouteEntries(request)
		if err != nil {
			if IsExceptedError(err, OperationBlocking) {
				return resource.RetryableError(fmt.Errorf("Publish CEN route entry timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("Publish CEN route entry timeout, CEN %s child instance %s vtb %s cidr %s, error info: %#v.", cenId, childInstanceId, vtbId, cidr, err)
	}

	d.SetId(cenId + COLON_SEPARATED + vtbId + COLON_SEPARATED + cidr)

	err = client.WaitForRouterEntryPublished(d.Id(), PUBLISHED, DefaultCenTimeout)
	if err != nil {
		return fmt.Errorf("Timeout when WaitForCenAvailable")
	}

	return resourceAlicloudCenRouteEntryRead(d, meta)
}

func resourceAlicloudCenRouteEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	parts := strings.Split(d.Id(), COLON_SEPARATED)
	if len(parts) != 3 {
		return fmt.Errorf("invalid resource id")
	}
	cenId := parts[0]

	resp, err := client.DescribePublishedRouteEntriesById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return err
	}

	if resp.PublishStatus == string(NOPUBLISHED) {
		d.SetId("")
		return nil
	}

	d.Set("instance_id", cenId)
	d.Set("route_table_id", resp.ChildInstanceRouteTableId)
	d.Set("cidr_block", resp.DestinationCidrBlock)

	return nil
}

func resourceAlicloudCenRouteEntryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	cenId := d.Get("instance_id").(string)
	vtbId := d.Get("route_table_id").(string)
	cidr := d.Get("cidr_block").(string)
	childInstanceId, childInstanceType, err := client.createCenRouteEntryParas(vtbId)
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return fmt.Errorf("Withdraw CEN route entry encounter an error when query childInstance ID, CEN %s vtb %s cidr %s, error info: %#v.", cenId, vtbId, cidr, err)
	}

	request := cbn.CreateWithdrawPublishedRouteEntriesRequest()
	request.CenId = cenId
	request.ChildInstanceId = childInstanceId
	request.ChildInstanceType = childInstanceType
	request.ChildInstanceRegionId = client.RegionId
	request.ChildInstanceRouteTableId = vtbId
	request.DestinationCidrBlock = cidr

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = client.cenconn.WithdrawPublishedRouteEntries(request)
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidCenInstanceStatus, InternalError}) {
				return resource.RetryableError(fmt.Errorf("Withdraw CEN route entries timeout and got an error: %#v", err))
			}

			return resource.NonRetryableError(fmt.Errorf("Withdraw CEN route entries timeout and got an error: %#v", err))
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("Withdraw CEN route entry timeout, CEN %s child instance %s vtb %s cidr %s, error info: %#v.", cenId, childInstanceId, vtbId, cidr, err)
	}

	if err := client.WaitForRouterEntryPublished(d.Id(), NOPUBLISHED, DefaultCenTimeout); err != nil {
		return fmt.Errorf("Timeout when WaitForRouterEntriesNoPublished")
	}

	return nil
}
