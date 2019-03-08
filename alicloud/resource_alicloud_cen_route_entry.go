package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCenRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	cenId := d.Get("instance_id").(string)
	vtbId := d.Get("route_table_id").(string)
	cidr := d.Get("cidr_block").(string)
	childInstanceId, childInstanceType, err := cenService.CreateCenRouteEntryParas(vtbId)
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
		_, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.PublishRouteEntries(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{OperationBlocking, InvalidStateForOperationMsg}) {
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

	err = cenService.WaitForRouterEntryPublished(d.Id(), PUBLISHED, DefaultCenTimeout)
	if err != nil {
		return fmt.Errorf("Timeout when WaitForCenAvailable")
	}

	return resourceAlicloudCenRouteEntryRead(d, meta)
}

func resourceAlicloudCenRouteEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}

	parts := strings.Split(d.Id(), COLON_SEPARATED)
	if len(parts) != 3 {
		return fmt.Errorf("invalid resource id")
	}
	cenId := parts[0]

	resp, err := cenService.DescribePublishedRouteEntriesById(d.Id())
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
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}

	cenId := d.Get("instance_id").(string)
	vtbId := d.Get("route_table_id").(string)
	cidr := d.Get("cidr_block").(string)
	childInstanceId, childInstanceType, err := cenService.CreateCenRouteEntryParas(vtbId)
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
		_, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.WithdrawPublishedRouteEntries(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{NotFoundRoute, InstanceNotExistMsg}) {
				return nil
			} else if IsExceptedErrors(err, []string{InvalidCenInstanceStatus, InternalError}) {
				return resource.RetryableError(fmt.Errorf("Withdraw CEN route entries timeout and got an error: %#v", err))
			}

			return resource.NonRetryableError(fmt.Errorf("Withdraw CEN route entries timeout and got an error: %#v", err))
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("Withdraw CEN route entry timeout, CEN %s child instance %s vtb %s cidr %s, error info: %#v.", cenId, childInstanceId, vtbId, cidr, err)
	}

	if err := cenService.WaitForRouterEntryPublished(d.Id(), NOPUBLISHED, DefaultCenTimeout); err != nil {
		return fmt.Errorf("Timeout when WaitForRouterEntriesNoPublished")
	}

	return nil
}
