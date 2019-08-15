package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/nas"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudMountTargets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMountTargetRead,

		Schema: map[string]*schema.Schema{
			"file_system_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"access_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mount_target_domain": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'mount_target_domain' has been deprecated from provider version 1.53.0. New field 'ids' replaces it.",
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// mount_gatgets values
			"targets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount_target_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
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
						"access_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudMountTargetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := nas.CreateDescribeMountTargetsRequest()
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	request.FileSystemId = d.Get("file_system_id").(string)

	var allMt []nas.DescribeMountTargetsMountTarget1

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[vv.(string)] = vv.(string)
		}
	}
	invoker := NewInvoker()
	for {
		var raw interface{}
		if err := invoker.Run(func() error {
			rsp, err := client.WithNasClient(func(nasClient *nas.Client) (interface{}, error) {
				return nasClient.DescribeMountTargets(request)
			})
			raw = rsp
			return err
		}); err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nas_mount_targets", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		resp, _ := raw.(*nas.DescribeMountTargetsResponse)

		if resp == nil || len(resp.MountTargets.MountTarget) < 1 {
			break
		}
		for _, mt := range resp.MountTargets.MountTarget {
			if v, ok := d.GetOk("type"); ok && mt.NetworkType != v.(string) {
				continue
			}
			if v, ok := d.GetOk("vpc_id"); ok && mt.VpcId != v.(string) {
				continue
			}
			if v, ok := d.GetOk("vswitch_id"); ok && mt.VswId != v.(string) {
				continue
			}
			if v, ok := d.GetOk("access_group_name"); ok && mt.AccessGroup != v.(string) {
				continue
			}
			if v, ok := d.GetOk("mount_target_domain"); ok && mt.MountTargetDomain != v.(string) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[mt.MountTargetDomain]; !ok {
					continue
				}
			}
			allMt = append(allMt, mt)
		}
		if len(resp.MountTargets.MountTarget) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	return MountTargetDescriptionAttributes(d, allMt, meta)
}

func MountTargetDescriptionAttributes(d *schema.ResourceData, nasSetTypes []nas.DescribeMountTargetsMountTarget1, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, mt := range nasSetTypes {
		mapping := map[string]interface{}{
			"id":                  mt.MountTargetDomain,
			"type":                mt.NetworkType,
			"vpc_id":              mt.VpcId,
			"vswitch_id":          mt.VswId,
			"access_group_name":   mt.AccessGroup,
			"mount_target_domain": mt.MountTargetDomain,
		}
		ids = append(ids, mt.MountTargetDomain)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("targets", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
