#!/usr/bin/env python3
"""Regression tests for the duplicate required_providers guard in
scripts/local-example-check.sh.

The check used to write a `terraform { required_providers { alicloud = ... } }`
file into every example module unconditionally. When a documentation example
already declared its own `required_providers` block, Terraform refused to load
the module with "Duplicate required providers configuration" -- any two
required_providers blocks conflict, regardless of the providers each lists. The
fix only injects the block when the example declares none, and the new
--validate-only mode surfaces real duplicates without init, credentials or
cloud resources.
"""
import os
from pathlib import Path
import shutil
import subprocess
import tempfile
import textwrap
import unittest

SCRIPTS = Path(__file__).resolve().parents[1]
SCRIPT = SCRIPTS / 'local-example-check.sh'


def extract_function(name):
    """Return the bash text of a top-level function defined in the script."""
    lines = SCRIPT.read_text().splitlines()
    start = next((i for i, line in enumerate(lines)
                 if line.startswith(name + '() {')), None)
    if start is None:
        raise AssertionError('function %s() not found in %s' % (name, SCRIPT))
    for j in range(start + 1, len(lines)):
        if lines[j] == '}':
            return '\n'.join(lines[start:j + 1]) + '\n'
    raise AssertionError('closing brace of %s() not found' % name)


class RequiredProvidersDetectionTest(unittest.TestCase):
    """Pure HCL detection of an existing required_providers block.

    No terraform required: the helper is sourced straight from the script so
    the test exercises the real implementation, not a reimplementation.
    """

    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.func = extract_function('example_has_required_providers_block')

    def detects(self, hcl):
        result = subprocess.run(
            ['bash', '-c',
             self.func + '\nprintf "%s\\n" "$1" | '
             'example_has_required_providers_block && echo HAS || echo NO',
             '_', hcl],
            text=True, capture_output=True)
        self.assertEqual(result.returncode, 0, result.stderr)
        return result.stdout.strip()

    def test_resource_only_example_has_no_block(self):
        self.assertEqual(self.detects('resource "alicloud_vpc" "demo" {}\n'), 'NO')

    def test_full_alicloud_declaration_has_block(self):
        hcl = textwrap.dedent('''\
            terraform {
              required_providers {
                alicloud = { source = "registry.terraform.io/aliyun/alicloud" }
              }
            }
        ''')
        self.assertEqual(self.detects(hcl), 'HAS')

    def test_commented_out_block_is_not_detected(self):
        hcl = ('# terraform { required_providers { alicloud = {} } }\n'
               'resource "alicloud_vpc" "demo" {}\n')
        self.assertEqual(self.detects(hcl), 'NO')

    def test_slash_comment_block_is_not_detected(self):
        hcl = ('// terraform { required_providers { alicloud = {} } }\n'
               'resource "alicloud_vpc" "demo" {}\n')
        self.assertEqual(self.detects(hcl), 'NO')

    def test_multi_provider_single_block_has_block(self):
        hcl = textwrap.dedent('''\
            terraform {
              required_providers {
                alicloud = { source = "registry.terraform.io/aliyun/alicloud" }
                random   = { source = "hashicorp/random" }
              }
            }
        ''')
        self.assertEqual(self.detects(hcl), 'HAS')

    def test_non_alicloud_required_providers_still_treated_as_block(self):
        # A required_providers block that does not declare alicloud must still
        # be respected: injecting a second block would itself trigger
        # "Duplicate required providers configuration".
        hcl = textwrap.dedent('''\
            terraform {
              required_providers { aws = { source = "hashicorp/aws" } }
            }
        ''')
        self.assertEqual(self.detects(hcl), 'HAS')

    def test_required_version_only_has_no_block(self):
        self.assertEqual(self.detects('terraform { required_version = ">= 1.0" }\n'), 'NO')

    def test_source_url_with_double_slash_keeps_block(self):
        hcl = textwrap.dedent('''\
            terraform {
              required_providers {
                alicloud = { source = "https://registry.terraform.io/aliyun/alicloud" }
              }
            }
        ''')
        self.assertEqual(self.detects(hcl), 'HAS')


class ValidateOnlyFixtureTest(unittest.TestCase):
    """End-to-end --validate-only run against isolated fixture examples.

    Requires terraform (available in CI); skipped when absent.
    """

    def setUp(self):
        self.terraform = shutil.which('terraform')
        if self.terraform is None:
            self.skipTest('terraform is required for the validate-only regression')
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name)
        (self.root / 'scripts').mkdir()
        shutil.copy2(SCRIPT, self.root / 'scripts' / 'local-example-check.sh')

    def write_fixture(self, name, hcl):
        path = self.root / (name + '.html.markdown')
        path.write_text('# %s\n\n```terraform\n%s```\n' % (name, hcl))
        return path

    def run_script(self, fixture):
        env = dict(os.environ)
        env.pop('ALICLOUD_ACCESS_KEY', None)
        env.pop('ALICLOUD_SECRET_KEY', None)
        return subprocess.run(
            ['bash', 'scripts/local-example-check.sh', '--validate-only',
             '--file', str(fixture)],
            cwd=self.root, env=env, text=True, capture_output=True)

    def test_full_declaration_does_not_produce_duplicate(self):
        # The original bug: the script injected a second required_providers
        # block here, producing "Duplicate required providers configuration".
        fixture = self.write_fixture('full_declaration', textwrap.dedent('''\
            terraform {
              required_providers {
                alicloud = {
                  source  = "registry.terraform.io/aliyun/alicloud"
                  version = "~> 1.220.0"
                }
              }
            }
            provider "alicloud" { region = "cn-hangzhou" }
            resource "alicloud_vpc" "demo" { cidr_block = "172.16.0.0/12" }
        '''))
        result = self.run_script(fixture)
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertIn('duplicate required_providers check: success', result.stdout)
        self.assertNotIn('Duplicate required providers configuration', result.stdout)

    def test_example_without_required_providers_is_injected(self):
        fixture = self.write_fixture(
            'no_declaration',
            'resource "alicloud_vpc" "demo" { cidr_block = "172.16.0.0/12" }\n')
        result = self.run_script(fixture)
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertIn('duplicate required_providers check: success', result.stdout)

    def test_alias_and_multiple_terraform_blocks_preserved(self):
        fixture = self.write_fixture('alias_multi', textwrap.dedent('''\
            terraform {
              required_providers {
                alicloud = { source = "registry.terraform.io/aliyun/alicloud" }
                random   = { source = "hashicorp/random" }
              }
            }
            terraform {
              required_version = ">= 1.5"
            }
            provider "alicloud" { region = "cn-hangzhou" }
            provider "alicloud" { alias = "foo"; region = "cn-beijing" }
            resource "alicloud_vpc" "demo" { cidr_block = "172.16.0.0/12" }
        '''))
        result = self.run_script(fixture)
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertIn('example already declares required_providers', result.stdout)
        self.assertNotIn('Duplicate required providers configuration', result.stdout)

    def test_real_duplicate_required_providers_still_errors(self):
        fixture = self.write_fixture('real_duplicate', textwrap.dedent('''\
            terraform {
              required_providers { alicloud = { source = "registry.terraform.io/aliyun/alicloud" } }
            }
            terraform {
              required_providers { alicloud = { source = "registry.terraform.io/aliyun/alicloud" } }
            }
            resource "alicloud_vpc" "demo" { cidr_block = "172.16.0.0/12" }
        '''))
        result = self.run_script(fixture)
        self.assertNotEqual(result.returncode, 0, result.stdout + result.stderr)
        self.assertIn('duplicate required_providers check: fail', result.stdout)
        self.assertIn('Duplicate required providers configuration', result.stdout)


if __name__ == '__main__':
    unittest.main()
