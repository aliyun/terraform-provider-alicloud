package datahub

const (
    projectsPath = "/projects"
    projectPath  = "/projects/%s"
    topicsPath   = "/projects/%s/topics"
    topicPath    = "/projects/%s/topics/%s"
    shardsPath   = "/projects/%s/topics/%s/shards"
    shardPath    = "/projects/%s/topics/%s/shards/%s"

    //connectorsPath        = "/projects/%s/topics/%s/connectors"
    connectorsPath        = "/projects/%s/topics/%s/connectors?mode=id"
    connectorPath         = "/projects/%s/topics/%s/connectors/%s"
    connectorDoneTimePath = "/projects/%s/topics/%s/connectors/%s?donetime"
    consumerGroupPath     = "/projects/%s/topics/%s/subscriptions/%s"

    subscriptionsPath = "/projects/%s/topics/%s/subscriptions"
    subscriptionPath  = "/projects/%s/topics/%s/subscriptions/%s"
    offsetsPath       = "/projects/%s/topics/%s/subscriptions/%s/offsets"
)
