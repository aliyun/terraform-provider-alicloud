package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strings"
	"time"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudKVStoreZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKVStroeZonesRead,

		Schema: map[string]*schema.Schema{
			"multi": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      common.PostPaid,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAlicloudKVStroeZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := r_kvstore.CreateDescribeAvailableResourceRequest()
	request.RegionId = client.RegionId
	request.InstanceChargeType = d.Get("instance_charge_type").(string)

	var response = &r_kvstore.DescribeAvailableResourceResponse{}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRkvClient(func(client *r_kvstore.Client) (interface{}, error) {
			return client.DescribeAvailableResource(request)
		})
		if err != nil {
			if IsExceptedError(err, Throttling) {
				time.Sleep(time.Duration(5) * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response = raw.(*r_kvstore.DescribeAvailableResourceResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kvstore_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	if len(response.AvailableZones.AvailableZone) <= 0 {
		return WrapError(fmt.Errorf("[ERROR] There is no available zone for KVStrore."))
	}

	var zoneIds []string

	for _, r := range response.AvailableZones.AvailableZone {
		if d.Get("multi").(bool) {
			if strings.Contains(r.ZoneId, MULTI_IZ_SYMBOL) && r.RegionId == string(client.Region) {
				zoneIds = append(zoneIds, r.ZoneId)
			}
		} else {
			if !strings.Contains(r.ZoneId, MULTI_IZ_SYMBOL) {
				zoneIds = append(zoneIds, r.ZoneId)
			}
		}
	}

	d.SetId(dataResourceIdHash(zoneIds))

	if err := d.Set("ids", zoneIds); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), zoneIds)
	}

	return nil
}
