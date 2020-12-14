package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudNasMountTargets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudNasMountTargetsRead,
		Schema: map[string]*schema.Schema{
			"access_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mount_target_domain": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'mount_target_domain' has been deprecated from provider version 1.53.0. New field 'ids' replaces it.",
			},
			"type": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'type' has been deprecated from provider version 1.95.0. New field 'network_type' replaces it.",
			},
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Active", "Inactive", "Pending"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"targets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount_target_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudNasMountTargetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := nas.CreateDescribeMountTargetsRequest()
	if v, ok := d.GetOk("file_system_id"); ok {
		request.FileSystemId = v.(string)
	}
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	var objects []nas.MountTarget

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	var response *nas.DescribeMountTargetsResponse
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			raw, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
				return nasClient.DescribeMountTargets(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceUnavailable", "Throttling"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw)
			response, _ = raw.(*nas.DescribeMountTargetsResponse)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nas_mount_targets", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		for _, item := range response.MountTargets.MountTarget {
			if v, ok := d.GetOk("access_group_name"); ok && v.(string) != "" && item.AccessGroup != v.(string) {
				continue
			}
			if v, ok := d.GetOk("mount_target_domain"); ok && v.(string) != "" && item.MountTargetDomain != v.(string) {
				continue
			}
			if v, ok := d.GetOk("type"); ok && v.(string) != "" && item.NetworkType != v.(string) {
				continue
			}
			if v, ok := d.GetOk("network_type"); ok && v.(string) != "" && item.NetworkType != v.(string) {
				continue
			}
			if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" && item.VpcId != v.(string) {
				continue
			}
			if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" && item.VswId != v.(string) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[item.MountTargetDomain]; !ok {
					continue
				}
			}
			if statusOk && status != "" && status != item.Status {
				continue
			}
			objects = append(objects, item)
		}
		if len(response.MountTargets.MountTarget) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"access_group_name":   object.AccessGroup,
			"id":                  object.MountTargetDomain,
			"mount_target_domain": object.MountTargetDomain,
			"network_type":        object.NetworkType,
			"type":                object.NetworkType,
			"status":              object.Status,
			"vpc_id":              object.VpcId,
			"vswitch_id":          object.VswId,
		}
		ids = append(ids, object.MountTargetDomain)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("targets", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
