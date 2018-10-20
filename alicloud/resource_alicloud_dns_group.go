package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/dns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudDnsGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDnsGroupCreate,
		Read:   resourceAlicloudDnsGroupRead,
		Update: resourceAlicloudDnsGroupUpdate,
		Delete: resourceAlicloudDnsGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudDnsGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	args := &dns.AddDomainGroupArgs{
		GroupName: d.Get("name").(string),
	}

	raw, err := client.WithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
		return dnsClient.AddDomainGroup(args)
	})
	if err != nil {
		return fmt.Errorf("AddDomainGroup got a error: %#v", err)
	}
	response, _ := raw.(*dns.AddDomainGroupResponse)
	d.SetId(response.GroupId)
	d.Set("name", response.GroupName)
	return resourceAlicloudDnsGroupUpdate(d, meta)
}

func resourceAlicloudDnsGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)
	args := &dns.UpdateDomainGroupArgs{
		GroupId: d.Id(),
	}

	if d.HasChange("name") && !d.IsNewResource() {
		d.SetPartial("name")
		args.GroupName = d.Get("name").(string)
		_, err := client.WithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.UpdateDomainGroup(args)
		})
		if err != nil {
			return fmt.Errorf("UpdateDomainGroup got an error: %#v", err)
		}
	}

	d.Partial(false)
	return resourceAlicloudDnsGroupRead(d, meta)
}

func resourceAlicloudDnsGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := &dns.DescribeDomainGroupsArgs{
		KeyWord: d.Get("name").(string),
	}

	raw, err := client.WithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
		return dnsClient.DescribeDomainGroups(args)
	})
	if err != nil {
		return err
	}
	groups, _ := raw.([]dns.DomainGroupType)
	if groups == nil || len(groups) <= 0 {
		return fmt.Errorf("No domain groups found.")
	}
	for _, v := range groups {
		if v.GroupName == d.Get("name").(string) {
			d.Set("name", v.GroupName)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceAlicloudDnsGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := &dns.DeleteDomainGroupArgs{
		GroupId: d.Id(),
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithDnsClient(func(dnsClient *dns.Client) (interface{}, error) {
			return dnsClient.DeleteDomainGroup(args)
		})
		if err != nil {
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code == FobiddenNotEmptyGroup {
				return resource.RetryableError(fmt.Errorf("The domain group can’t be deleted because it is not empty - trying again after it empty."))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting group %s: %#v", d.Id(), err))
		}
		return nil
	})
}
