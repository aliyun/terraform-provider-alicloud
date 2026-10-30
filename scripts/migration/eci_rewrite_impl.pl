#!/usr/bin/perl
use strict;
use warnings;

# Bulk rewrite for ECI implementation files migrated into alicloud/service/eci/.
# - package alicloud -> eci
# - factory funcs -> Resource*/DataSource* (exported, SDKv2 no New prefix)
# - EciService type -> Service; EciService{client} -> NewService(client)
# - root-package helpers -> leaf packages (errs/logs/helper/convert/validators/tags)
# Negative lookbehind (?<![.\w]) prevents matching already-qualified or substrings.

my @files = @ARGV;
for my $f (@files) {
    open my $fh, '<', $f or die "open $f: $!";
    local $/; my $c = <$fh>; close $fh;

    # package
    $c =~ s/^package alicloud\b/package eci/m;

    # factory function definitions (exact func decls)
    $c =~ s/\bfunc resourceAlicloudEciContainerGroup\(\)/func ResourceEciContainerGroup()/g;
    $c =~ s/\bfunc resourceAlicloudEciImageCache\(\)/func ResourceEciImageCache()/g;
    $c =~ s/\bfunc resourceAlicloudEciVirtualNode\(\)/func ResourceEciVirtualNode()/g;
    $c =~ s/\bfunc dataSourceAlicloudEciContainerGroups\(\)/func DataSourceEciContainerGroups()/g;
    $c =~ s/\bfunc dataSourceAlicloudEciImageCaches\(\)/func DataSourceEciImageCaches()/g;
    $c =~ s/\bfunc dataSourceAlicloudEciVirtualNodes\(\)/func DataSourceEciVirtualNodes()/g;
    $c =~ s/\bfunc dataSourceAlicloudEciZones\(\)/func DataSourceEciZones()/g;

    # Service type and construction
    $c =~ s/\bEciService\{client\}/NewService(client)/g;
    $c =~ s/\btype EciService struct\b/type Service struct/g;
    $c =~ s/\bfunc \(s \*EciService\)/func (s *Service)/g;

    # errs package (longer/const identifiers first, then funcs)
    $c =~ s/(?<![.\w])DataDefaultErrorMsg\b/errs.DataDefaultErrorMsg/g;
    $c =~ s/(?<![.\w])DefaultErrorMsg\b/errs.DefaultErrorMsg/g;
    $c =~ s/(?<![.\w])AlibabaCloudSdkGoERROR\b/errs.AlibabaCloudSdkGoERROR/g;
    $c =~ s/(?<![.\w])ProviderERROR\b/errs.ProviderERROR/g;
    $c =~ s/(?<![.\w])FailedToReachTargetStatus\b/errs.FailedToReachTargetStatus/g;
    $c =~ s/(?<![.\w])NotFoundWithResponse\b/errs.NotFoundWithResponse/g;
    $c =~ s/(?<![.\w])FailedGetAttributeMsg\b/errs.FailedGetAttributeMsg/g;
    $c =~ s/(?<![.\w])NotFoundMsg\b/errs.NotFoundMsg/g;
    $c =~ s/(?<![.\w])IdMsg\b/errs.IdMsg/g;
    $c =~ s/(?<![.\w])WrapErrorf\(/errs.WrapErrorf(/g;
    $c =~ s/(?<![.\w])WrapError\(/errs.WrapError(/g;
    $c =~ s/(?<![.\w])NotFoundError\(/errs.NotFoundError(/g;
    $c =~ s/(?<![.\w])NotFoundErr\(/errs.NotFoundErr(/g;
    $c =~ s/(?<![.\w])IsExpectedErrors\(/errs.IsExpectedErrors(/g;
    $c =~ s/(?<![.\w])Error\(/errs.Error(/g;

    # logs
    $c =~ s/(?<![.\w])addDebug\(/logs.AddDebug(/g;

    # helper (constants and funcs)
    $c =~ s/(?<![.\w])PageSizeLarge\b/helper.PageSizeLarge/g;
    $c =~ s/(?<![.\w])NeedRetry\(/helper.NeedRetry(/g;
    $c =~ s/(?<![.\w])incrementalWait\(/helper.IncrementalWait(/g;
    $c =~ s/(?<![.\w])BuildStateConf\(/helper.BuildStateConf(/g;
    $c =~ s/(?<![.\w])buildClientToken\(/helper.BuildClientToken(/g;
    $c =~ s/(?<![.\w])dataResourceIdHash\(/helper.DataResourceIdHash(/g;
    $c =~ s/(?<![.\w])writeToFile\(/helper.WriteToFile(/g;
    $c =~ s/(?<![.\w])Trim\(/helper.Trim(/g;

    # convert
    $c =~ s/(?<![.\w])expandStringList\(/convert.ExpandStringList(/g;
    $c =~ s/(?<![.\w])formatInt\(/convert.FormatInt(/g;

    # validators (lookbehind avoids validation.StringInSlice)
    $c =~ s/(?<![.\w])StringInSlice\(/validators.StringInSlice(/g;

    # tags
    $c =~ s/(?<![.\w])tagsSchema\(/tags.TagsSchema(/g;
    $c =~ s/(?<![.\w])tagsToMap\(/tags.TagsToMap(/g;
    $c =~ s/(?<![.\w])ignoredTags\(/tags.IgnoredTags(/g;

    open my $out, '>', $f or die "write $f: $!";
    print $out $c; close $out;
}
print "rewrote @files\n";
