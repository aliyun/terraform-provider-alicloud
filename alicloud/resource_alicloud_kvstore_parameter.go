package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudKVStoreParameter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKVStoreParameterCreate,
		Read:   resourceAlicloudKVStoreParameterRead,
		Update: resourceAlicloudKVStoreParameterUpdate,
		Delete: resourceAlicloudKVStoreParameterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudKVStoreParameterCreate(d *schema.ResourceData, meta interface{}) error {
	instanceId := d.Get("instance_id").(string)
	name := d.Get("name").(string)
	d.SetId(fmt.Sprintf("%s%s%s", instanceId, COLON_SEPARATED, name))

	return resourceAlicloudKVStoreParameterUpdate(d, meta)
}

func resourceAlicloudKVStoreParameterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}
	parts := strings.SplitN(d.Id(), COLON_SEPARATED, 2)
	response, err := kvstoreService.DescribeRKVInstanceParameter(parts[0])
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	found := false
	d.Set("instance_id", parts[0])

	for _, i := range response.RunningParameters.Parameter {
		if i.ParameterName != "" {
			if i.ParameterName == parts[1] {
				found = true
				d.Set("name", i.ParameterName)
				d.Set("value", i.ParameterValue)
				break
			}
		}
	}

	for _, i := range response.ConfigParameters.Parameter {
		if i.ParameterName != "" {
			if i.ParameterName == parts[1] {
				found = true
				d.Set("name", i.ParameterName)
				d.Set("value", i.ParameterValue)
				break
			}
		}
	}

	if !found {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceAlicloudKVStoreParameterUpdate(d *schema.ResourceData, meta interface{}) error {
	parts := strings.SplitN(d.Id(), COLON_SEPARATED, 2)
	client := meta.(*connectivity.AliyunClient)
	update := false
	request := r_kvstore.CreateModifyInstanceConfigRequest()

	request.InstanceId = parts[0]

	if d.HasChange("value") {
		name := d.Get("name").(string)
		value := d.Get("value").(string)
		config := make(map[string]string)
		config[name] = value
		cfg, _ := json.Marshal(config)
		request.Config = string(cfg)
		update = true
	}

	if update {
		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			_, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
				return rkvClient.ModifyInstanceConfig(request)
			})
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("update parameter got an error: %#v", err))
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	return resourceAlicloudKVStoreParameterRead(d, meta)
}

func resourceAlicloudKVStoreParameterDelete(d *schema.ResourceData, meta interface{}) error {
	// In case of a delete we just ignore
	return nil
}
