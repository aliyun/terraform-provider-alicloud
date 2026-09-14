// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudEnsBucketAcls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEnsBucketAclsRead,
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"bucket_acls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bucket_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bucket_acl": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudEnsBucketAclsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := EnsServiceV2{client}

	bucketName := d.Get("bucket_name").(string)

	objectRaw, err := ensServiceV2.DescribeEnsBucketAcl(bucketName)
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_bucket_acls DescribeEnsBucketAcl Failed!!! %s", err)
			d.SetId(fmt.Sprintf("%s", bucketName))
			return nil
		}
		return WrapError(err)
	}

	acl := fmt.Sprint(objectRaw["BucketAcl"])
	d.SetId(fmt.Sprintf("%s", bucketName))
	d.Set("bucket_name", bucketName)
	d.Set("bucket_acl", acl)

	var ids []string
	var acls []map[string]interface{}
	ids = append(ids, bucketName)
	acls = append(acls, map[string]interface{}{
		"id":          bucketName,
		"bucket_name": bucketName,
		"bucket_acl":  acl,
	})
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("bucket_acls", acls); err != nil {
		return WrapError(err)
	}

	return nil
}
