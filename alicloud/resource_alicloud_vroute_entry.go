package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunRouteEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunRouteEntryCreate,
		Read:   resourceAliyunRouteEntryRead,
		Delete: resourceAliyunRouteEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"router_id": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Attribute router_id has been deprecated and suggest removing it from your template.",
			},
			"route_table_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_cidrblock": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"nexthop_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRouteEntryNextHopType,
			},
			"nexthop_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunRouteEntryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	rtId := d.Get("route_table_id").(string)
	cidr := d.Get("destination_cidrblock").(string)
	nt := d.Get("nexthop_type").(string)
	ni := d.Get("nexthop_id").(string)

	table, err := meta.(*AliyunClient).QueryRouteTableById(rtId)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error query route table: %#v", err)
	}

	err = resource.Retry(3*time.Minute, func() *resource.RetryError {

		if err := client.WaitForAllRouteEntries(rtId, Available, DefaultTimeout); err != nil {
			return resource.NonRetryableError(fmt.Errorf("WaitFor route entries got error: %#v", err))
		}

		args, err := buildAliyunRouteEntryArgs(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Building CreateRouteEntryArgs got an error: %#v", err))
		}

		if _, err := client.vpcconn.CreateRouteEntry(args); err != nil {
			// Route Entry does not support concurrence when creating or deleting it;
			// Route Entry does not support creating or deleting within 5 seconds frequently
			// It must ensure all the route entries and vswitches' status must be available before creating or deleting route entry.
			if IsExceptedError(err, TaskConflict) || IsExceptedError(err, IncorrectRouteEntryStatus) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(fmt.Errorf("Create route entry timeout and got an error: %#v", err))
			}
			if IsExceptedError(err, RouterEntryConflictDuplicated) {
				en, err := client.QueryRouteEntry(rtId, cidr, nt, ni)
				if err != nil {
					return resource.NonRetryableError(err)
				}
				return resource.NonRetryableError(fmt.Errorf("The route entry %s has already existed. "+
					"Please import it using ID '%s:%s:%s:%s:%s' or specify a new 'destination_cidrblock' and try again.",
					en.DestinationCidrBlock, en.RouteTableId, table.VRouterId, en.DestinationCidrBlock, en.NextHopType, ni))
			}
			return resource.NonRetryableError(fmt.Errorf("Creating Route entry got an error: %#v", err))
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Create Vroute Entry got an error :%#v", err)
	}
	// route_table_id:router_id:destination_cidrblock:nexthop_type:nexthop_id

	d.SetId(rtId + ":" + table.VRouterId + ":" + cidr + ":" + nt + ":" + ni)

	if err := client.WaitForAllRouteEntries(rtId, Available, DefaultTimeout); err != nil {
		return fmt.Errorf("WaitFor route entry got error: %#v", err)
	}
	return resourceAliyunRouteEntryRead(d, meta)
}

func resourceAliyunRouteEntryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parts := strings.Split(d.Id(), ":")
	rtId := parts[0]
	rId := parts[1]
	cidr := parts[2]
	nexthop_type := parts[3]
	nexthop_id := parts[4]

	en, err := client.QueryRouteEntry(rtId, cidr, nexthop_type, nexthop_id)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error route entry: %#v", err)
	}

	d.Set("router_id", rId)
	d.Set("route_table_id", en.RouteTableId)
	d.Set("destination_cidrblock", en.DestinationCidrBlock)
	d.Set("nexthop_type", en.NextHopType)
	d.Set("nexthop_id", en.InstanceId)
	return nil
}

func resourceAliyunRouteEntryDelete(d *schema.ResourceData, meta interface{}) error {
	args, err := buildAliyunRouteEntryDeleteArgs(d, meta)

	if err != nil {
		return err
	}
	client := meta.(*AliyunClient)
	parts := strings.Split(d.Id(), ":")
	rtId := parts[0]
	cidr := parts[2]
	nexthop_type := parts[3]
	nexthop_id := parts[4]

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		en, err := client.QueryRouteEntry(rtId, cidr, nexthop_type, nexthop_id)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error route entry: %#v", err))
		}

		if en.Status != string(Available) {
			return resource.RetryableError(fmt.Errorf("Delete route entry timeout and got an error: %#v.", err))
		}

		if _, err := client.vpcconn.DeleteRouteEntry(args); err != nil {
			if IsExceptedError(err, TaskConflict) || IsExceptedError(err, IncorrectRouteEntryStatus) ||
				IsExceptedError(err, RouterEntryForbbiden) || IsExceptedError(err, UnknownError) {
				// Route Entry does not support creating or deleting within 5 seconds frequently
				time.Sleep(5 * time.Second)
				return resource.RetryableError(fmt.Errorf("Delete route entry timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Deleting RouteEntry got an error: %#v", err))
		}

		return nil
	})
}

func buildAliyunRouteEntryArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateRouteEntryRequest, error) {

	request := vpc.CreateCreateRouteEntryRequest()
	request.RouteTableId = d.Get("route_table_id").(string)
	request.DestinationCidrBlock = d.Get("destination_cidrblock").(string)

	if v := d.Get("nexthop_type").(string); v != "" {
		request.NextHopType = v
	}

	if v := d.Get("nexthop_id").(string); v != "" {
		request.NextHopId = v
	}

	return request, nil
}

func buildAliyunRouteEntryDeleteArgs(d *schema.ResourceData, meta interface{}) (*vpc.DeleteRouteEntryRequest, error) {

	request := vpc.CreateDeleteRouteEntryRequest()
	request.RouteTableId = d.Get("route_table_id").(string)
	request.DestinationCidrBlock = d.Get("destination_cidrblock").(string)

	if v := d.Get("destination_cidrblock").(string); v != "" {
		request.DestinationCidrBlock = v
	}

	if v := d.Get("nexthop_id").(string); v != "" {
		request.NextHopId = v
	}

	return request, nil
}
