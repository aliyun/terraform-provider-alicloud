package alicloud

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/samber/lo"
)

func resourceAlicloudMongoDBPublicNetworkAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMongoDBPublicNetworkAddressCreate,
		Read:   resourceAlicloudMongoDBPublicNetworkAddressRead,
		Delete: resourceAlicloudMongoDBPublicNetworkAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"replica_sets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replica_set_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							// Though this should be always "Public".
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudMongoDBPublicNetworkAddressCreate(d *schema.ResourceData, meta interface{}) error {
	// only one public network address per instance.
	instanceId := d.Get("db_instance_id").(string)
	d.SetId(instanceId)

	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	if err := ddsService.AllocatePublicNetworkAddress(instanceId); err != nil {
		return WrapError(err)
	}

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, ddsService.RdsMongodbDBInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudMongoDBPublicNetworkAddressRead(d, meta)
}

func resourceAlicloudMongoDBPublicNetworkAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	object, err := ddsService.DescribeReplicaSetRole(d.Id())
	if err != nil {
		return WrapError(err)
	}

	replicaSets := transferToMongoReplicaSets(object, true)
	allPublicNetworkAddresses := lo.Filter(replicaSets, func(replica map[string]interface{}, idx int) bool {
		if networkType, ok := replica["network_type"]; ok && networkType.(string) == "Public" {
			return true
		}
		return false
	})

	d.Set("db_instance_id", d.Id())

	d.Set("replica_sets", allPublicNetworkAddresses)

	return nil
}

func resourceAlicloudMongoDBPublicNetworkAddressDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddsService := MongoDBService{client}

	if err := ddsService.ReleasePublicNetworkAddress(d.Id()); err != nil {
		return WrapError(err)
	}

	stateConf := BuildStateConf([]string{}, []string{"NotExist"}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, ddsService.RdsMongoDBPublicNetworkAddressStateRefreshFunc(d.Id()))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}

	return nil
}
