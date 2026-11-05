#!/usr/bin/perl
use strict; use warnings;

# Bulk rewrite for ECI test files migrated into alicloud/service/eci/
# (package eci_test). Maps root-package test helpers to acctest.* and
# leaf packages, modernizes provider factories to ProtoV5, and rewrites
# Service construction to eci.NewService.

my @files = @ARGV;
for my $f (@files) {
    open my $fh, '<', $f or die "open $f: $!";
    local $/; my $c = <$fh>; close $fh;

    # package
    $c =~ s/^package alicloud\b/package eci_test/m;

    # provider factory modernization (do before bare testAccProvider)
    $c =~ s/ProviderFactories:\s*testAccProviderFactory,/ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),/g;

    # Service construction: EciService{<expr>} (no nested braces) -> eci.NewService(<expr>)
    $c =~ s/\bEciService\{([^}]*)\}/eci.NewService($1)/g;

    # sweep helpers (local)
    $c =~ s/\bsharedClientForRegion\(/eciTestClientForRegion(/g;
    $c =~ s/\btestSweepPreCheckWithRegions\(/eciTestSweepPreCheckWithRegions(/g;

    # root test helpers -> acctest (lookbehind avoids already-qualified)
    $c =~ s/(?<![.\w])checkoutSupportedRegions\(/acctest.CheckoutSupportedRegions(/g;
    $c =~ s/(?<![.\w])testAccPreCheck\(/acctest.PreCheck(/g;
    $c =~ s/(?<![.\w])testAccProvider\b/acctest.Provider/g;
    $c =~ s/(?<![.\w])resourceAttrInit\(/acctest.ResourceAttrInit(/g;
    $c =~ s/(?<![.\w])resourceCheckInitWithDescribeMethod\(/acctest.ResourceCheckInitWithDescribeMethod(/g;
    $c =~ s/(?<![.\w])resourceAttrCheckInit\(/acctest.ResourceAttrCheckInit(/g;
    $c =~ s/(?<![.\w])resourceTestAccConfigFunc\(/acctest.ResourceTestAccConfig(/g;
    $c =~ s/(?<![.\w])defaultRegionToTest\b/acctest.DefaultRegionToTest/g;
    $c =~ s/(?<![.\w])CHECKSET\b/acctest.CHECKSET/g;

    # methods: lowercase -> exported (lookbehind requires leading dot)
    $c =~ s/\.resourceAttrMapUpdateSet\(\)/.ResourceAttrMapUpdateSet()/g;
    $c =~ s/\.checkResourceDestroy\(\)/.CheckResourceDestroy()/g;
    $c =~ s/\.dataSourceTestCheck\(/.DataSourceTestCheck(/g;

    # data source test types/fields
    $c =~ s/(?<![.\w])dataSourceAttr\b/acctest.DataSourceAttr/g;
    $c =~ s/(?<![.\w])dataSourceTestAccConfig\b/acctest.DataSourceTestAccConfig/g;
    $c =~ s/(?<![.\w])dataSourceTestCheckWithPreCheck\b/acctest.DataSourceTestCheckWithPreCheck/g;
    # field keys (colon directly follows; safe vs var decls using :=)
    $c =~ s/(?<![.\w])resourceId:/ResourceId:/g;
    $c =~ s/(?<![.\w])existMapFunc:/ExistMapFunc:/g;
    $c =~ s/(?<![.\w])fakeMapFunc:/FakeMapFunc:/g;
    $c =~ s/(?<![.\w])existConfig:/ExistConfig:/g;
    $c =~ s/(?<![.\w])fakeConfig:/FakeConfig:/g;

    # errs/logs/helper in sweep functions (same leaf map as impl)
    $c =~ s/(?<![.\w])WrapErrorf\(/errs.WrapErrorf(/g;
    $c =~ s/(?<![.\w])WrapError\(/errs.WrapError(/g;
    $c =~ s/(?<![.\w])DataDefaultErrorMsg\b/errs.DataDefaultErrorMsg/g;
    $c =~ s/(?<![.\w])DefaultErrorMsg\b/errs.DefaultErrorMsg/g;
    $c =~ s/(?<![.\w])AlibabaCloudSdkGoERROR\b/errs.AlibabaCloudSdkGoERROR/g;
    $c =~ s/(?<![.\w])FailedGetAttributeMsg\b/errs.FailedGetAttributeMsg/g;
    $c =~ s/(?<![.\w])addDebug\(/logs.AddDebug(/g;
    $c =~ s/(?<![.\w])incrementalWait\(/helper.IncrementalWait(/g;
    $c =~ s/(?<![.\w])NeedRetry\(/helper.NeedRetry(/g;

    open my $out, '>', $f or die "write $f: $!";
    print $out $c; close $out;
}
print "rewrote @files\n";
