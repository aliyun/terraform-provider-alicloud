package alicloud

import (
	"regexp"

	ots "github.com/alibabacloud-go/tablestore-20201209/v3/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudOtsInstanceAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOtsInstanceAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateOTSInstanceName,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpc_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudOtsInstanceAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	instanceName := d.Get("instance_name").(string)
	allVpcs, err := otsService.ListOtsInstanceVpc(instanceName)
	if err != nil {
		return WrapError(err)
	}

	var filteredVpcs []*ots.ListVpcInfoByInstanceResponseBodyVpcInfos
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		r := regexp.MustCompile(v.(string))
		for _, vpc := range allVpcs {
			if r.MatchString(*vpc.InstanceVpcName) {
				filteredVpcs = append(filteredVpcs, vpc)
			}
		}
	} else {
		filteredVpcs = allVpcs[:]
	}
	return otsAttachmentsDescriptionAttributes(d, filteredVpcs, meta)
}

func otsAttachmentsDescriptionAttributes(d *schema.ResourceData, vpcInfos []*ots.ListVpcInfoByInstanceResponseBodyVpcInfos, meta interface{}) error {
	var ids []string
	var names []string
	var vpcIds []string
	var s []map[string]interface{}

	instanceName := d.Get("instance_name").(string)
	for _, vpc := range vpcInfos {
		mapping := map[string]interface{}{
			"id":            instanceName,
			"domain":        vpc.Domain,
			"endpoint":      vpc.Endpoint,
			"region":        vpc.RegionNo,
			"instance_name": instanceName,
			"vpc_name":      vpc.InstanceVpcName,
			"vpc_id":        vpc.VpcId,
		}
		names = append(names, *vpc.InstanceVpcName)
		ids = append(ids, instanceName)
		vpcIds = append(vpcIds, *vpc.VpcId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("attachments", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("vpc_ids", vpcIds); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
