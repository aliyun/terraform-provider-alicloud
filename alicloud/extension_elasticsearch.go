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

const ClientNodeDisk = "20"
const ClientNodeDiskType = "cloud_efficiency"

const WaitInstanceActiveTimeout = 3600

var DataNodeSpec = "elasticsearch.n4.small"
var DataNodeAmount = "2"
var DataNodeDisk = "20"
var DataNodeDiskType = "cloud_ssd"

var DataNodeSpecForUpdate = "elasticsearch.sn2ne.large"
var DataNodeAmountForUpdate = "3"
var DataNodeDiskForUpdate = "30"

var MasterNodeSpecForUpdate = "elasticsearch.sn2ne.xlarge"
