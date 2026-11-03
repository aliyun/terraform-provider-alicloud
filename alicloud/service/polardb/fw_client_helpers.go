package polardb

import "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

// Client returns the provider's AliyunClient. The ACube-generated resource calls
// r.Client() for brevity, but fwadapt.ResourceBase only exposes Meta(). This thin
// adapter bridges the gap without modifying the generated file.
func (r *networkChannelResource) Client() *connectivity.AliyunClient {
	return r.Meta()
}
