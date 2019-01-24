package alicloud

type ESVersion string

const (
	ESVersion553WithXPack = ESVersion("5.5.3_with_X-Pack")
	ESVersion632WithXPck  = ESVersion("6.3.2_with_X-Pack")
)

type ElasticsearchStatus string

const (
	ElasticsearchStatusActive     = ElasticsearchStatus("active")
	ElasticsearchStatusActivating = ElasticsearchStatus("activating")
)

const MasterNodeDisk = "20"
const MasterNodeDiskType = "cloud_ssd"
const MasterNodeAmount = "3"

const WaitInstanceActiveTimeout = 3600

