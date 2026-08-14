package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCrInternetEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCrInternetEndpointRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entries": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entry": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudCrInternetEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	instanceId := d.Get("instance_id").(string)
	objectRaw, err := crService.DescribeCrInternetEndpoint(instanceId)
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Data source alicloud_cr_internet_endpoint DescribeCrInternetEndpoint Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.SetId(fmt.Sprint(instanceId))
	d.Set("status", objectRaw["Status"])
	d.Set("instance_id", instanceId)

	aclEntriesRaw := objectRaw["AclEntries"]
	entriesMaps := make([]map[string]interface{}, 0)
	if aclEntriesRaw != nil {
		for _, aclEntriesChildRaw := range convertToInterfaceArray(aclEntriesRaw) {
			aclEntriesChildRaw := aclEntriesChildRaw.(map[string]interface{})
			// Exclude the system-managed loopback default ACL policy
			// (entry 127.0.0.1/32, comment "default") that GetInstanceEndpoint
			// auto-returns once the endpoint is enabled, so the data source
			// exposes only the user-managed entries and stays consistent with
			// the paired resource's state.
			if fmt.Sprint(aclEntriesChildRaw["Entry"]) == "127.0.0.1/32" && fmt.Sprint(aclEntriesChildRaw["Comment"]) == "default" {
				continue
			}
			entriesMap := make(map[string]interface{})
			entriesMap["comment"] = aclEntriesChildRaw["Comment"]
			entriesMap["entry"] = aclEntriesChildRaw["Entry"]
			entriesMaps = append(entriesMaps, entriesMap)
		}
	}
	if err := d.Set("entries", entriesMaps); err != nil {
		return err
	}

	return nil
}
