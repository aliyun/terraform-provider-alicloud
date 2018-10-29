package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudSlbAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSlbAclCreate,
		Read:   resourceAlicloudSlbAclRead,
		Update: resourceAlicloudSlbAclUpdate,
		Delete: resourceAlicloudSlbAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"ip_version": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      IPVersion4,
				ValidateFunc: validateAllowedStringValue([]string{string(IPVersion4), string(IPVersion6)}),
			},
			"entry_list": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entry": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"comment": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				MaxItems: 300,
				MinItems: 0,
			},
		},
	}
}

func resourceAlicloudSlbAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	name := strings.Trim(d.Get("name").(string), " ")
	ip_version := d.Get("ip_version").(string)

	request := slb.CreateCreateAccessControlListRequest()
	request.AclName = name
	request.AddressIPVersion = ip_version

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.CreateAccessControlList(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{SlbAclInvalidActionRegionNotSupport, SlbAclNumberOverLimit}) {
			return fmt.Errorf("CreateAccessControlList got an error: %#v", err)
		}

		return fmt.Errorf("CreateAccessControlList got an unknown error: %#v", err)
	}
	response, _ := raw.(*slb.CreateAccessControlListResponse)

	d.SetId(response.AclId)
	return resourceAlicloudSlbAclUpdate(d, meta)
}

func resourceAlicloudSlbAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	request := slb.CreateDescribeAccessControlListAttributeRequest()
	request.AclId = d.Id()
	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeAccessControlListAttribute(request)
	})
	if err != nil {
		if IsExceptedError(err, SlbAclNotExists) {
			d.SetId("")
			return nil
		}
		return err
	}
	acl, _ := raw.(*slb.DescribeAccessControlListAttributeResponse)

	d.Set("name", acl.AclName)
	d.Set("ip_version", acl.AddressIPVersion)
	if err := d.Set("entry_list", slbService.FlattenSlbAclEntryMappings(acl.AclEntrys.AclEntry)); err != nil {
		return err
	}

	return nil
}

func resourceAlicloudSlbAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("name") {
		request := slb.CreateSetAccessControlListAttributeRequest()
		request.AclId = d.Id()
		request.AclName = d.Get("name").(string)
		_, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetAccessControlListAttribute(request)
		})
		if err != nil {

			return fmt.Errorf("SetAccessControlListAttribute set %s  name %s got an error: %#v",
				d.Id(), request.AclName, err)

		}
		d.SetPartial("name")
	}

	if d.HasChange("entry_list") {
		o, n := d.GetChange("entry_list")
		oe := o.(*schema.Set)
		ne := n.(*schema.Set)
		remove := oe.Difference(ne).List()
		add := ne.Difference(oe).List()

		if len(remove) > 0 {
			if err := slbService.SlbRemoveAccessControlListEntry(remove, d.Id()); err != nil {
				return err
			}
		}

		if len(add) > 0 {
			if err := slbService.SlbAddAccessControlListEntry(add, d.Id()); err != nil {
				return err
			}
		}

		d.SetPartial("entry_list")
	}

	d.Partial(false)

	return resourceAlicloudSlbAclRead(d, meta)
}

func resourceAlicloudSlbAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := slb.CreateDeleteAccessControlListRequest()
		request.AclId = d.Id()
		_, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteAccessControlList(request)
		})
		if err != nil {
			if IsExceptedError(err, SlbAclNotExists) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("DeleteAccessControlList %s got an error: %#v.", d.Id(), err))
		}

		req := slb.CreateDescribeAccessControlListAttributeRequest()
		req.AclId = d.Id()
		_, err = client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeAccessControlListAttribute(req)
		})
		if err != nil {
			if IsExceptedError(err, SlbAclNotExists) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("While DeleteAccessControlList，DescribeAccessControlListAttribute %s got an error: %#v.", d.Id(), err))
		}

		return resource.RetryableError(fmt.Errorf("DeleteAccessControlList %s timeout.", d.Id()))
	})
}
