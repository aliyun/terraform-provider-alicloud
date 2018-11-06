package alicloud

import (
	"fmt"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudSlbAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSlbAclsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"acls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entry_list": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"entry": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"comment": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MinItems: 0,
						},
						"related_listeners": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"load_balancer_id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"frontend_port": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"protocol": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"acl_type": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							MinItems: 0,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSlbAclsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := slb.CreateDescribeAccessControlListsRequest()

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeAccessControlLists(args)
	})
	if err != nil {
		return fmt.Errorf("DescribeAccessControlLists got an error: %#v", err)
	}
	resp, _ := raw.(*slb.DescribeAccessControlListsResponse)
	if resp == nil {
		return fmt.Errorf("there is no SLB acl. Please change your search criteria and try again")
	}

	var filteredAclsTemp []slb.Acl
	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, acl := range resp.Acls.Acl {
			if r != nil && !r.MatchString(acl.AclName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[acl.AclId]; !ok {
					continue
				}
			}

			filteredAclsTemp = append(filteredAclsTemp, acl)
		}
	} else {
		filteredAclsTemp = resp.Acls.Acl
	}

	if len(filteredAclsTemp) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	return slbAclsDescriptionAttributes(d, filteredAclsTemp, client)
}

func slbAclsDescriptionAttributes(d *schema.ResourceData, acls []slb.Acl, client *connectivity.AliyunClient) error {

	var ids []string
	var s []map[string]interface{}
	slbService := SlbService{client}

	req := slb.CreateDescribeAccessControlListAttributeRequest()
	for _, item := range acls {
		req.AclId = item.AclId
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeAccessControlListAttribute(req)
		})
		if err != nil {
			return fmt.Errorf("DescribeAccessControlListAttribute %s got an error: %#v", req.AclId, err)
		}
		acl, _ := raw.(*slb.DescribeAccessControlListAttributeResponse)

		mapping := map[string]interface{}{
			"id":                acl.AclId,
			"name":              acl.AclName,
			"ip_version":        acl.AddressIPVersion,
			"entry_list":        slbService.FlattenSlbAclEntryMappings(acl.AclEntrys.AclEntry),
			"related_listeners": slbService.flattenSlbRelatedListenerMappings(acl.RelatedListeners.RelatedListener),
		}

		ids = append(ids, acl.AclId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("acls", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
