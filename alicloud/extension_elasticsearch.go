package alicloud

type ESVersion string

const (
	ESVersion553WithXPack = ESVersion("5.5.3_with_X-Pack")
	ESVersion632WithXPack = ESVersion("6.3_with_X-Pack")
	ESVersion670WithXPack = ESVersion("6.7_with_X-Pack")
)

type ElasticsearchStatus string

const (
	ElasticsearchStatusActive     = ElasticsearchStatus("active")
	ElasticsearchStatusActivating = ElasticsearchStatus("activating")
)

const MasterNodeDisk = "20"
const MasterNodeDiskType = "cloud_ssd"
const MasterNodeAmount = "3"

const WaitInstanceActiveTimeout = 7200
