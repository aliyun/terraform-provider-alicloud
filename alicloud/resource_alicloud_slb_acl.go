package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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

	client := meta.(*AliyunClient).slbconn

	name := strings.Trim(d.Get("name").(string), " ")
	ip_version := d.Get("ip_version").(string)

	request := slb.CreateCreateAccessControlListRequest()
	request.AclName = name
	request.AddressIPVersion = ip_version

	response, err := client.CreateAccessControlList(request)
	if err != nil {
		if IsExceptedErrors(err, []string{SlbAclInvalidActionRegionNotSupport, SlbAclNumberOverLimit}) {
			return fmt.Errorf("CreateAccessControlList got an error: %#v", err)
		}

		return fmt.Errorf("CreateAccessControlList got an unknown error: %#v", err)
	}

	d.SetId(response.AclId)
	return resourceAlicloudSlbAclUpdate(d, meta)
}

func resourceAlicloudSlbAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient).slbconn

	request := slb.CreateDescribeAccessControlListAttributeRequest()
	request.AclId = d.Id()
	acl, err := client.DescribeAccessControlListAttribute(request)
	if err != nil {
		if IsExceptedError(err, SlbAclNotExists) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", acl.AclName)
	d.Set("ip_version", acl.AddressIPVersion)
	if err := d.Set("entry_list", flattenSlbAclEntryMappings(acl.AclEntrys.AclEntry)); err != nil {
		return err
	}

	return nil
}

func resourceAlicloudSlbAclUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient).slbconn

	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("name") {
		request := slb.CreateSetAccessControlListAttributeRequest()
		request.AclId = d.Id()
		request.AclName = d.Get("name").(string)
		if _, err := client.SetAccessControlListAttribute(request); err != nil {

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
			if err := slbRemoveAccessControlListEntry(client, remove, d.Id()); err != nil {
				return err
			}
		}

		if len(add) > 0 {
			if err := slbAddAccessControlListEntry(client, add, d.Id()); err != nil {
				return err
			}
		}

		d.SetPartial("entry_list")
	}

	d.Partial(false)

	return resourceAlicloudSlbAclRead(d, meta)
}

func resourceAlicloudSlbAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient).slbconn

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := slb.CreateDeleteAccessControlListRequest()
		request.AclId = d.Id()
		if _, err := client.DeleteAccessControlList(request); err != nil {
			if IsExceptedError(err, SlbAclNotExists) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("DeleteAccessControlList %s got an error: %#v.", d.Id(), err))
		}

		req := slb.CreateDescribeAccessControlListAttributeRequest()
		req.AclId = d.Id()
		if _, err := client.DescribeAccessControlListAttribute(req); err != nil {
			if IsExceptedError(err, SlbAclNotExists) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("While DeleteAccessControlList，DescribeAccessControlListAttribute %s got an error: %#v.", d.Id(), err))
		}

		return resource.RetryableError(fmt.Errorf("DeleteAccessControlList %s timeout.", d.Id()))
	})
}
