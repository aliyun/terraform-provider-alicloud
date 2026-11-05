package slb

// ListenerAclAttachmentRPCPostForTest lets the external resource test package
// replace the transport entry without exporting production symbols.
var ListenerAclAttachmentRPCPostForTest = &listenerAclAttachmentRpcPost

// ListenerAclAttachmentConfirmTimeoutForTest lets tests shorten the bounded
// confirmation window without adding production configuration.
var ListenerAclAttachmentConfirmTimeoutForTest = &listenerAclAttachmentConfirmTimeout

// ListenerAclAttachmentPollIntervalForTest lets tests shorten the convergence
// poll interval.
var ListenerAclAttachmentPollIntervalForTest = &listenerAclAttachmentPollInterval

// ListenerAclIdsForTest exposes the ACL identity extraction to the external
// test package without exporting a production symbol.
var ListenerAclIdsForTest = listenerAclIds
