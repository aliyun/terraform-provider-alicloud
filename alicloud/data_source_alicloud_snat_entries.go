package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudSnatEntries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSnatEntriesRead,

		Schema: map[string]*schema.Schema{
			"snat_table_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snat_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// the snat_entry resource id is spliced from snat_table_id and snat_entry_id, but,this id refers to snat_entry_id
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				Computed: true,
			},

			// Computed values
			"entries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snat_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceAlicloudSnatEntriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := vpc.CreateDescribeSnatTableEntriesRequest()
	request.RegionId = string(client.Region)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	request.SnatTableId = d.Get("snat_table_id").(string)
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	var allSnatEntries []vpc.SnatTableEntry
	invoker := NewInvoker()
	for {
		var raw interface{}
		if err := invoker.Run(func() error {
			response, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeSnatTableEntries(request)
			})
			raw = response
			return err
		}); err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_snat_entries", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeSnatTableEntriesResponse)
		if len(response.SnatTableEntries.SnatTableEntry) < 1 {
			break
		}

		for _, entries := range response.SnatTableEntries.SnatTableEntry {
			if snat_ip, ok := d.GetOk("snat_ip"); ok && entries.SnatIp != snat_ip.(string) {
				continue
			}
			if source_cidr, ok := d.GetOk("source_cidr"); ok && entries.SourceCIDR != source_cidr.(string) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[entries.SnatEntryId]; !ok {
					continue
				}
			}
			allSnatEntries = append(allSnatEntries, entries)
		}

		if len(response.SnatTableEntries.SnatTableEntry) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return SnatEntriesDecriptionAttributes(d, allSnatEntries, meta)
}

func SnatEntriesDecriptionAttributes(d *schema.ResourceData, entries []vpc.SnatTableEntry, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, entry := range entries {
		mapping := map[string]interface{}{
			"id":          entry.SnatEntryId,
			"snat_ip":     entry.SnatIp,
			"source_cidr": entry.SourceCIDR,
			"status":      entry.Status,
		}
		ids = append(ids, entry.SnatEntryId)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("entries", s); err != nil {
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
