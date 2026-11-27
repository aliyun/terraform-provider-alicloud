#!/usr/bin/env python3
"""Regression tests for the HCL-aware provider requirement preparation.

Positive cases cover: examples without any declaration (terraform.tf is
injected), examples that already declare alicloud (left untouched),
comments/strings/heredocs that merely look like blocks (not detected), and
blocks declaring other providers (completed in place, same block).

Negative cases cover real duplicate required_providers blocks: they must
stay untouched so `terraform validate` reports "Duplicate required
providers configuration" instead of the check silently ignoring the error.
"""

import importlib.util
import json
import os
import shutil
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

HELPER = Path(__file__).with_name("required_providers.py")

spec = importlib.util.spec_from_file_location("required_providers", HELPER)
rp = importlib.util.module_from_spec(spec)
spec.loader.exec_module(rp)


def has_terraform():
    return shutil.which("terraform") is not None


class ExampleDir(object):
    """A temporary example directory driven through the real CLI."""

    def __init__(self, main_tf):
        self.tmp = tempfile.TemporaryDirectory()
        self.dir = Path(self.tmp.name)
        (self.dir / "main.tf").write_text(main_tf, encoding="utf-8")

    def close(self):
        self.tmp.cleanup()

    def prepare(self):
        proc = subprocess.run(
            [sys.executable, str(HELPER), "--example-dir", str(self.dir)],
            stdout=subprocess.PIPE, stderr=subprocess.PIPE,
            universal_newlines=True,
        )
        if proc.returncode != 0:
            raise AssertionError(
                "helper failed (%d): %s" % (proc.returncode, proc.stderr)
            )
        return proc.stdout.strip()

    def main_tf(self):
        return (self.dir / "main.tf").read_text(encoding="utf-8")

    def terraform_tf(self):
        path = self.dir / "terraform.tf"
        return path.read_text(encoding="utf-8") if path.exists() else None

    def validate(self):
        env = dict(os.environ, TF_CLI_CONFIG_FILE=os.devnull)
        proc = subprocess.run(
            ["terraform", "-chdir=" + str(self.dir), "validate", "-json"],
            stdout=subprocess.PIPE, stderr=subprocess.PIPE,
            universal_newlines=True, env=env,
        )
        return json.loads(proc.stdout)


ALICLOUD_DECLARED = """terraform {
  required_providers {
    alicloud = {
      source  = "registry.terraform.io/aliyun/alicloud"
      version = "~> 1.220"
    }
  }
}

provider "alicloud" {
  alias  = "west"
  region = "cn-hangzhou-west"
}

resource "alicloud_vpc" "example" {
  vpc_name = "example"
}
"""

OTHER_PROVIDERS_ONLY = """terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = ">= 3.1"
    }
  }
}

resource "random_pet" "name" {
  length = 2
}
"""

TWO_TRUE_DUPLICATE_BLOCKS = """terraform {
  required_providers {
    alicloud = {
      source = "registry.terraform.io/aliyun/alicloud"
    }
  }
}

terraform {
  required_providers {
    alicloud = {
      source  = "registry.terraform.io/aliyun/alicloud"
      version = "~> 1.220"
    }
  }
}

resource "alicloud_vpc" "example" {
  vpc_name = "example"
}
"""

# QA-R2-1 (AI 260937918): the interpolation expression legally spans lines
# and sits before the required_providers block, inside the terraform block.
MULTILINE_INTERP_IN_TERRAFORM_BLOCK = """terraform {
  required_version = ">= 1.${
    0
  }"

  required_providers {
    alicloud = {
      source  = "registry.terraform.io/aliyun/alicloud"
      version = "~> 1.220"
    }
  }
}

resource "alicloud_vpc" "example" {
  vpc_name = "example"
}
"""

# QA-R3-1 (A-shape): a multi-line interpolation whose expression contains a
# bare object literal. The object's closing brace used to decrement the
# interpolation depth to zero, the following newline then ended the string
# scan early, the truncation corrupted the brace stack and the
# required_providers block went undetected (helper injected terraform.tf and
# terraform reported "Duplicate required providers configuration").
MULTILINE_INTERP_OBJECT_LITERAL_IN_TERRAFORM_BLOCK = """terraform {
  backend "oss" {
    config = "${
      merge(
        { public = "https://example.com" },
      )
    }"
  }

  required_providers {
    alicloud = {
      source  = "registry.terraform.io/aliyun/alicloud"
      version = "~> 1.220"
    }
  }
}

resource "alicloud_vpc" "example" {
  vpc_name = "example"
}
"""

# QA-R3-1 (C-shape): the same layout with a for-object expression inside the
# multi-line interpolation; its braces must balance against the
# interpolation depth exactly like an object literal's.
MULTILINE_INTERP_FOR_OBJECT_IN_TERRAFORM_BLOCK = """terraform {
  backend "oss" {
    config = "${
      merge(
        { for k, v in var.tags : k => upper(v) },
      )
    }"
  }

  required_providers {
    alicloud = {
      source  = "registry.terraform.io/aliyun/alicloud"
      version = "~> 1.220"
    }
  }
}

resource "alicloud_vpc" "example" {
  vpc_name = "example"
}
"""

# AI 260951912: a template directive %{ ... } opens an expression context
# just like ${ ... } and may legally span lines.
MULTILINE_DIRECTIVE_IN_TERRAFORM_BLOCK = """terraform {
  backend "oss" {
    config = "%{
      if var.public
    }public%{ else
    }private%{ endif
    }"
  }

  required_providers {
    alicloud = {
      source  = "registry.terraform.io/aliyun/alicloud"
      version = "~> 1.220"
    }
  }
}

resource "alicloud_vpc" "example" {
  vpc_name = "example"
}
"""

# Same string shape outside the terraform block: legal cross-line
# interpolation in a locals block before the declaration.
MULTILINE_INTERP_IN_LOCALS = """locals {
  name = "prefix-${
    "dev"
  }-suffix"
}

terraform {
  required_providers {
    alicloud = {
      source  = "registry.terraform.io/aliyun/alicloud"
      version = "~> 1.220"
    }
  }
}

resource "alicloud_vpc" "example" {
  vpc_name = local.name
}
"""


class PreparationDecisionTest(unittest.TestCase):
    """Decision outcomes without requiring terraform on PATH."""

    def test_no_declaration_gets_terraform_tf_injected(self):
        example = ExampleDir('resource "alicloud_vpc" "example" {}\n')
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "injected")
        self.assertEqual(example.terraform_tf(), rp.INJECTED_TERRAFORM_TF)
        self.assertEqual(example.main_tf(), 'resource "alicloud_vpc" "example" {}\n')

    def test_line_commented_block_is_not_detected(self):
        main = '# terraform {\n#   required_providers {\n#     alicloud = {\n#       source = "registry.terraform.io/aliyun/alicloud"\n#     }\n#   }\n# }\n\nresource "alicloud_vpc" "example" {}\n'
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "injected")
        self.assertEqual(example.terraform_tf(), rp.INJECTED_TERRAFORM_TF)

    def test_block_commented_block_is_not_detected(self):
        main = '/*\nterraform {\n  required_providers {\n    alicloud = {\n      source = "registry.terraform.io/aliyun/alicloud"\n    }\n  }\n}\n*/\n\nresource "alicloud_vpc" "example" {}\n'
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "injected")
        self.assertEqual(example.terraform_tf(), rp.INJECTED_TERRAFORM_TF)

    def test_quoted_string_is_not_detected(self):
        main = 'variable "note" {\n  default = "terraform { required_providers { alicloud = { source = \\"registry.terraform.io/aliyun/alicloud\\" } } }"\n}\n'
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "injected")
        self.assertEqual(example.terraform_tf(), rp.INJECTED_TERRAFORM_TF)

    def test_heredoc_body_is_not_detected(self):
        main = 'variable "note" {\n  default = <<EOT\nterraform {\n  required_providers {\n    alicloud = {\n      source = "registry.terraform.io/aliyun/alicloud"\n    }\n  }\n}\nEOT\n}\n\nresource "alicloud_vpc" "example" {}\n'
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "injected")
        self.assertEqual(example.terraform_tf(), rp.INJECTED_TERRAFORM_TF)

    def test_declared_example_is_left_untouched(self):
        # The example's own source, version and provider alias must survive.
        example = ExampleDir(ALICLOUD_DECLARED)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "declared")
        self.assertIsNone(example.terraform_tf())
        self.assertEqual(example.main_tf(), ALICLOUD_DECLARED)

    def test_other_providers_are_completed_in_the_same_block(self):
        example = ExampleDir(OTHER_PROVIDERS_ONLY)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "completed")
        self.assertIsNone(example.terraform_tf())
        completed = example.main_tf()
        # One single required_providers block now, holding random AND
        # alicloud; the random entry is unchanged.
        self.assertEqual(completed.count("required_providers {"), 1)
        self.assertIn('random = {\n      source  = "hashicorp/random"\n      version = ">= 3.1"\n    }', completed)
        self.assertIn(
            'alicloud = {\n      source = "registry.terraform.io/aliyun/alicloud"\n    }',
            completed,
        )

    def test_quoted_provider_key_is_completed(self):
        main = 'terraform {\n  required_providers {\n    "random" = {\n      source = "hashicorp/random"\n    }\n  }\n}\n'
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "completed")
        self.assertEqual(example.main_tf().count("required_providers {"), 1)
        self.assertIn('alicloud = {', example.main_tf())

    def test_quoted_alicloud_declaration_is_declared(self):
        # The quoted key must be read from the original text (mask_hcl
        # blanks string bodies, the masked token is only quotes and
        # spaces); a declared alicloud requirement is recognized so the
        # example stays untouched instead of being "completed" into an
        # in-block duplicate.
        main = 'terraform {\n  required_providers {\n    "alicloud" = {\n      source = "registry.terraform.io/aliyun/alicloud"\n    }\n  }\n}\n'
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "declared")
        self.assertIsNone(example.terraform_tf())
        self.assertEqual(example.main_tf(), main)

    def test_completed_across_multiple_terraform_blocks(self):
        # Several terraform blocks with a single required_providers are
        # legal; the alicloud entry goes into the existing block.
        main = 'terraform {\n  required_providers {\n    random = {\n      source = "hashicorp/random"\n    }\n  }\n}\n\nterraform {\n  backend "local" {}\n}\n'
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "completed")
        self.assertEqual(example.main_tf().count("required_providers {"), 1)
        self.assertIn('alicloud = {', example.main_tf())

    def test_attribute_form_is_not_a_block(self):
        # required_providers = { ... } is an attribute, not a block; it is
        # not a valid declaration, so the check injects terraform.tf and
        # lets terraform surface the real configuration error.
        main = 'terraform {\n  required_providers = {\n    alicloud = {\n      source = "registry.terraform.io/aliyun/alicloud"\n    }\n  }\n}\n'
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "injected")
        self.assertEqual(example.terraform_tf(), rp.INJECTED_TERRAFORM_TF)

    def test_two_true_duplicate_blocks_are_left_untouched(self):
        original = TWO_TRUE_DUPLICATE_BLOCKS
        example = ExampleDir(original)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "duplicate")
        self.assertIsNone(example.terraform_tf())
        self.assertEqual(example.main_tf(), original)

    def test_two_blocks_with_second_alicloud_are_left_untouched(self):
        main = 'terraform {\n  required_providers {\n    random = {\n      source = "hashicorp/random"\n    }\n  }\n}\n\nterraform {\n  required_providers {\n    alicloud = {\n      source = "registry.terraform.io/aliyun/alicloud"\n    }\n  }\n}\n'
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "duplicate")
        self.assertIsNone(example.terraform_tf())
        self.assertEqual(example.main_tf(), main)

    def test_multiline_interpolation_in_terraform_block_is_declared(self):
        # QA-R2-1 (AI 260937918): a legal cross-line interpolation before
        # the required_providers block used to end the string scan at the
        # first newline, the unmasked interpolation body then corrupted
        # the brace stack, the block went undetected and the helper
        # wrongly injected terraform.tf (Duplicate required providers
        # configuration recurred). The example's own declaration must be
        # recognized and kept untouched.
        example = ExampleDir(MULTILINE_INTERP_IN_TERRAFORM_BLOCK)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "declared")
        self.assertIsNone(example.terraform_tf())
        self.assertEqual(example.main_tf(), MULTILINE_INTERP_IN_TERRAFORM_BLOCK)

    def test_multiline_interpolation_in_locals_is_declared(self):
        # The same cross-line interpolation outside the terraform block
        # (legal HCL in a locals block) must not hide the declaration
        # either.
        example = ExampleDir(MULTILINE_INTERP_IN_LOCALS)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "declared")
        self.assertIsNone(example.terraform_tf())
        self.assertEqual(example.main_tf(), MULTILINE_INTERP_IN_LOCALS)

    def test_multiline_interpolation_with_object_literal_is_declared(self):
        # QA-R3-1 (A-shape): a bare object literal inside a multi-line
        # interpolation used to end the interpolation depth early; the
        # block went undetected and the helper injected terraform.tf. The
        # example's own declaration must be recognized and kept untouched.
        example = ExampleDir(MULTILINE_INTERP_OBJECT_LITERAL_IN_TERRAFORM_BLOCK)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "declared")
        self.assertIsNone(example.terraform_tf())
        self.assertEqual(
            example.main_tf(), MULTILINE_INTERP_OBJECT_LITERAL_IN_TERRAFORM_BLOCK
        )

    def test_multiline_interpolation_with_for_object_is_declared(self):
        # QA-R3-1 (C-shape): a for-object expression's braces must balance
        # against the interpolation depth exactly like an object literal's.
        example = ExampleDir(MULTILINE_INTERP_FOR_OBJECT_IN_TERRAFORM_BLOCK)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "declared")
        self.assertIsNone(example.terraform_tf())
        self.assertEqual(
            example.main_tf(), MULTILINE_INTERP_FOR_OBJECT_IN_TERRAFORM_BLOCK
        )

    def test_multiline_directive_is_declared(self):
        # AI 260951912: a template directive %{ ... } opens an expression
        # context like ${ ... } and may span lines; it must not hide the
        # declaration either.
        example = ExampleDir(MULTILINE_DIRECTIVE_IN_TERRAFORM_BLOCK)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "declared")
        self.assertIsNone(example.terraform_tf())
        self.assertEqual(example.main_tf(), MULTILINE_DIRECTIVE_IN_TERRAFORM_BLOCK)

    def test_escaped_introducers_stay_literal_text(self):
        # $${ and %%{ are escaped template introducers (literal "${" and
        # "%{"): they open no expression context. An unterminated string
        # holding them still ends the scan at the next newline, and a
        # closed one ends at its own closing quote even with a stray
        # brace in between.
        text = '  note = "cost $${5 and %%{flag\n'
        start = text.index('"')
        end = rp._scan_quoted(text, start)
        self.assertEqual(text[end], "\n")
        text2 = '  note = "a $${ x } b"\n}\n'
        start2 = text2.index('"')
        end2 = rp._scan_quoted(text2, start2)
        self.assertEqual(text2[end2], '"')

    def test_multiline_interpolation_object_literal_other_providers_completed(self):
        # A-shape interpolation before an other-providers-only block: the
        # block must still be found and completed in place, no second
        # block, no terraform.tf.
        main = MULTILINE_INTERP_OBJECT_LITERAL_IN_TERRAFORM_BLOCK.replace(
            'alicloud = {\n      source  = "registry.terraform.io/aliyun/alicloud"\n      version = "~> 1.220"\n    }',
            'random = {\n      source = "hashicorp/random"\n    }',
        )
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "completed")
        self.assertIsNone(example.terraform_tf())
        completed = example.main_tf()
        self.assertEqual(completed.count("required_providers {"), 1)
        self.assertIn(
            'alicloud = {\n      source = "registry.terraform.io/aliyun/alicloud"\n    }',
            completed,
        )

    def test_object_literal_fake_block_in_interpolation_is_not_detected(self):
        # The object literal fix must not over-detect: a string nested in
        # the interpolation that merely looks like a required_providers
        # block stays string content, nothing is detected and terraform.tf
        # is injected.
        main = 'variable "note" {\n  default = "${\n    join(\n      " ",\n      ["terraform { required_providers { alicloud = {} } }"],\n    )\n  }"\n}\n\nresource "alicloud_vpc" "example" {}\n'
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "injected")
        self.assertEqual(example.terraform_tf(), rp.INJECTED_TERRAFORM_TF)

    def test_multiline_interpolation_with_other_providers_is_completed(self):
        # Cross-line interpolation before an other-providers-only block:
        # the block is still found and completed in place, no second
        # block, no terraform.tf.
        main = """terraform {
  required_version = ">= 1.${
    0
  }"

  required_providers {
    random = {
      source = "hashicorp/random"
    }
  }
}
"""
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "completed")
        self.assertIsNone(example.terraform_tf())
        completed = example.main_tf()
        self.assertEqual(completed.count("required_providers {"), 1)
        self.assertIn(
            'alicloud = {\n      source = "registry.terraform.io/aliyun/alicloud"\n    }',
            completed,
        )

    def test_multiline_interpolation_fake_block_is_not_detected(self):
        # A cross-line interpolation whose nested string merely looks
        # like a required_providers block stays string content: nothing
        # is detected and terraform.tf is injected.
        main = 'variable "note" {\n  default = "${\n    "terraform { required_providers { alicloud = {} } }"\n  }"\n}\n\nresource "alicloud_vpc" "example" {}\n'
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "injected")
        self.assertEqual(example.terraform_tf(), rp.INJECTED_TERRAFORM_TF)


class CompletionParsingTest(unittest.TestCase):
    """Direct parser checks for the in-place completion path."""

    def test_completion_positions(self):
        masked = rp.mask_hcl(OTHER_PROVIDERS_ONLY)
        blocks = rp.find_required_providers_blocks(OTHER_PROVIDERS_ONLY)
        self.assertEqual(len(blocks), 1)
        block = blocks[0]
        self.assertEqual(block["entries"], ["random"])
        # The closing brace offset must be the one of the required_providers
        # block, matched on the masked text.
        self.assertEqual(masked[block["close_pos"]], "}")
        self.assertEqual(masked.count("required_providers"), 1)

    def test_one_line_block_completion_keeps_code_prefix_out(self):
        # A one-liner has no indentation to inherit; the inserted entry
        # must not copy the code standing before the keyword on that line.
        main = 'terraform { required_providers { random = { source = "hashicorp/random" } } }\n'
        blocks = rp.find_required_providers_blocks(main)
        self.assertEqual(blocks[0]["entries"], ["random"])
        completed = rp.complete_block(main, blocks[0])
        self.assertNotIn("\nterraform {", completed)
        # The inserted entry must be newline-terminated: in a single-line
        # block the next argument follows on the same line, and HCL
        # requires a newline between arguments ("Missing newline after
        # argument" otherwise).
        self.assertIn('    source = "registry.terraform.io/aliyun/alicloud"\n  }\n', completed)
        reparsed = rp.find_required_providers_blocks(completed)
        self.assertEqual(len(reparsed), 1)
        self.assertEqual(sorted(reparsed[0]["entries"]), ["alicloud", "random"])

    def test_tab_indented_block_is_completed(self):
        main = 'terraform {\n\trequired_providers {\n\t\trandom = {\n\t\t\tsource = "hashicorp/random"\n\t\t}\n\t}\n}\n'
        blocks = rp.find_required_providers_blocks(main)
        completed = rp.complete_block(main, blocks[0])
        reparsed = rp.find_required_providers_blocks(completed)
        self.assertEqual(len(reparsed), 1)
        self.assertEqual(sorted(reparsed[0]["entries"]), ["alicloud", "random"])
        self.assertIn("\t  alicloud = {\n\t    source", completed)

    def test_first_entry_sharing_brace_line_is_pushed_to_next_line(self):
        # AI 260959162: a multi-line block whose first entry shares the
        # opening-brace line (required_providers { random = {, source on the
        # following lines) walks the defensive "entry += newline" branch: the
        # inserted alicloud entry must land on its own line and the code
        # sharing the brace line must be pushed after it, otherwise HCL sees
        # two arguments separated only by whitespace. Re-parsing the
        # completed text must yield both entries in one single block.
        main = 'terraform {\n  required_providers { random = {\n    source = "hashicorp/random"\n  } }\n}\n'
        blocks = rp.find_required_providers_blocks(main)
        self.assertEqual(blocks[0]["entries"], ["random"])
        completed = rp.complete_block(main, blocks[0])
        self.assertEqual(completed.count("required_providers {"), 1)
        # The inserted entry opens right after the brace on its own line and
        # is newline-terminated before the pre-existing code.
        self.assertIn('required_providers {\n    alicloud = {\n', completed)
        self.assertIn('    }\n random = {\n', completed)
        self.assertNotIn('required_providers { random', completed)
        # F-QA-R4-1: while the first entry shares the opening-brace line HCL
        # tolerates the entry's closing brace on the block's closing line;
        # pushing the entry to its own line exits that mode, so the block's
        # closing brace must move onto its own line ("} }").
        self.assertNotIn("} }", completed)
        self.assertIn('    source = "hashicorp/random"\n  }\n  }\n', completed)
        reparsed = rp.find_required_providers_blocks(completed)
        self.assertEqual(len(reparsed), 1)
        self.assertEqual(sorted(reparsed[0]["entries"]), ["alicloud", "random"])

    def test_completion_produces_single_block_with_both_providers(self):
        blocks = rp.find_required_providers_blocks(OTHER_PROVIDERS_ONLY)
        completed = rp.complete_block(OTHER_PROVIDERS_ONLY, blocks[0])
        self.assertEqual(completed.count("required_providers {"), 1)
        # Still parses as exactly one required_providers block with both
        # entries after the rewrite.
        reparsed = rp.find_required_providers_blocks(completed)
        self.assertEqual(len(reparsed), 1)
        self.assertEqual(sorted(reparsed[0]["entries"]), ["alicloud", "random"])
        self.assertIn(
            'alicloud = {\n      source = "registry.terraform.io/aliyun/alicloud"\n    }',
            completed,
        )

    def test_scan_quoted_spans_multiline_interpolation(self):
        # QA-R2-1 (AI 260937918): while an interpolation is open the scan
        # must not stop at the first newline; it ends at the string's real
        # closing quote, so the whole template body gets masked.
        text = '  name = "prefix-${\n    var.env\n  }-suffix"\n}\n'
        start = text.index('"')
        end = rp._scan_quoted(text, start)
        self.assertEqual(text[end], '"')
        self.assertEqual(text[start + 1:end], 'prefix-${\n    var.env\n  }-suffix')
        masked = rp.mask_hcl(text)
        # Newlines survive, the interpolation body (including its braces)
        # is blanked.
        self.assertEqual(masked.count("\n"), text.count("\n"))
        self.assertNotIn("var.env", masked)
        self.assertNotIn("prefix", masked)
        self.assertEqual(masked.replace(" ", "").replace("\n", ""), 'name=""}')

    def test_scan_quoted_counts_object_literal_braces(self):
        # QA-R3-1: the braces of a bare object literal inside an open
        # interpolation must feed the interpolation depth; otherwise the
        # object's closing brace ends the interpolation early and the next
        # newline truncates the scan before the real closing quote.
        text = '  cfg = "${\n    merge(\n      { public = "x" },\n    )\n  }"\n}\n'
        start = text.index('"')
        end = rp._scan_quoted(text, start)
        self.assertEqual(text[end], '"')
        self.assertEqual(
            text[start + 1:end], '${\n    merge(\n      { public = "x" },\n    )\n  }'
        )
        masked = rp.mask_hcl(text)
        self.assertEqual(masked.count("\n"), text.count("\n"))
        self.assertNotIn("merge", masked)
        self.assertNotIn("public", masked)
        self.assertEqual(masked.replace(" ", "").replace("\n", ""), 'cfg=""}')

    def test_scan_quoted_counts_for_object_braces(self):
        # QA-R3-1: a for-object expression balances like an object
        # literal; the scan must reach the string's real closing quote.
        text = '  tags = "${\n    merge(\n      { for k, v in var.tags : k => upper(v) },\n    )\n  }"\n}\n'
        start = text.index('"')
        end = rp._scan_quoted(text, start)
        self.assertEqual(text[end], '"')
        self.assertEqual(
            text[start + 1:end],
            '${\n    merge(\n      { for k, v in var.tags : k => upper(v) },\n    )\n  }',
        )
        masked = rp.mask_hcl(text)
        self.assertEqual(masked.count("\n"), text.count("\n"))
        self.assertNotIn("upper", masked)
        self.assertEqual(masked.replace(" ", "").replace("\n", ""), 'tags=""}')

    def test_scan_quoted_spans_multiline_directive(self):
        # AI 260951912: a template directive %{ ... } keeps the scan inside
        # the string exactly like an interpolation.
        text = '  cfg = "%{ if var.public\n}public%{ else\n}private%{ endif\n}"\n}\n'
        start = text.index('"')
        end = rp._scan_quoted(text, start)
        self.assertEqual(text[end], '"')
        self.assertEqual(
            text[start + 1:end],
            '%{ if var.public\n}public%{ else\n}private%{ endif\n}',
        )
        masked = rp.mask_hcl(text)
        self.assertEqual(masked.count("\n"), text.count("\n"))
        self.assertNotIn("endif", masked)
        self.assertEqual(masked.replace(" ", "").replace("\n", ""), 'cfg=""}')


@unittest.skipUnless(has_terraform(), "terraform binary not available")
class TerraformValidationTest(unittest.TestCase):
    """End-to-end config validation of the prepared directories."""

    def assert_no_duplicate(self, diagnostics):
        summaries = [d["summary"] for d in diagnostics]
        self.assertNotIn("Duplicate required providers configuration", summaries)

    def test_declared_example_validates_without_duplicate_error(self):
        example = ExampleDir(ALICLOUD_DECLARED)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "declared")
        result = example.validate()
        self.assert_no_duplicate(result["diagnostics"])

    def test_injected_example_validates_without_duplicate_error(self):
        example = ExampleDir('resource "alicloud_vpc" "example" {}\n')
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "injected")
        result = example.validate()
        self.assert_no_duplicate(result["diagnostics"])

    def test_completed_example_validates_without_duplicate_error(self):
        example = ExampleDir(OTHER_PROVIDERS_ONLY)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "completed")
        result = example.validate()
        self.assert_no_duplicate(result["diagnostics"])
        # Both provider requirements are visible to terraform (the
        # providers themselves are simply not installed without init).
        details = " ".join(d.get("detail", "") for d in result["diagnostics"])
        self.assertIn("registry.terraform.io/aliyun/alicloud", details)
        self.assertIn("hashicorp/random", details)

    def test_one_line_block_completion_is_valid_hcl(self):
        # QA-1: completing a required_providers block written on a single
        # line (legal HCL as long as it holds a single argument) used to
        # emit the inserted entry and the existing argument on one line,
        # which terraform rejects with "Missing newline after argument";
        # the completed module must stay free of syntax errors.
        main = 'terraform {\n  required_providers { random = { source = "hashicorp/random" } }\n}\n'
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "completed")
        result = example.validate()
        summaries = [d["summary"] for d in result["diagnostics"]]
        self.assertNotIn("Missing newline after argument", summaries)
        self.assertNotIn("Unsupported block type", summaries)
        self.assert_no_duplicate(result["diagnostics"])
        details = " ".join(d.get("detail", "") for d in result["diagnostics"])
        self.assertIn("hashicorp/random", details)

    def test_first_entry_sharing_brace_line_completion_is_valid_hcl(self):
        # F-QA-R4-1 (end-to-end): while the first entry shares the
        # opening-brace line HCL parses the block in single-line mode and
        # tolerates the entry's closing brace on the block's closing line
        # ("} }"). The completion used to keep that tolerated tail while
        # pushing the entry to its own line, and terraform then reported
        # "Missing newline after argument". The completed module must stay
        # free of syntax errors, keep the single block and expose both
        # provider requirements.
        main = 'terraform {\n  required_providers { random = {\n    source = "hashicorp/random"\n  } }\n}\n'
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "completed")
        result = example.validate()
        summaries = [d["summary"] for d in result["diagnostics"]]
        self.assertNotIn("Missing newline after argument", summaries)
        self.assertNotIn("Invalid single-argument block definition", summaries)
        self.assert_no_duplicate(result["diagnostics"])
        details = " ".join(d.get("detail", "") for d in result["diagnostics"])
        self.assertIn("registry.terraform.io/aliyun/alicloud", details)
        self.assertIn("hashicorp/random", details)

    def test_two_true_duplicate_blocks_are_reported_by_terraform(self):
        # The mandated negative case: two real required_providers blocks
        # must fail configuration loading with the duplicate error.
        example = ExampleDir(TWO_TRUE_DUPLICATE_BLOCKS)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "duplicate")
        result = example.validate()
        self.assertFalse(result["valid"])
        summaries = [d["summary"] for d in result["diagnostics"]]
        self.assertIn("Duplicate required providers configuration", summaries)

    def test_string_fake_block_does_not_hide_real_declaration(self):
        # An example whose only "declaration" lives inside a string gets
        # terraform.tf injected; the resulting module still loads (no
        # duplicate) because the string is not a second block.
        main = 'variable "note" {\n  default = "terraform { required_providers { alicloud = { source = \\"registry.terraform.io/aliyun/alicloud\\" } } }"\n}\n'
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "injected")
        result = example.validate()
        self.assert_no_duplicate(result["diagnostics"])

    def test_multiline_interpolation_before_required_providers(self):
        # QA-R2-1 (AI 260937918): a legal cross-line interpolation before
        # the required_providers block made the helper inject terraform.tf
        # and terraform then reported "Duplicate required providers
        # configuration" (the original RL-11 failure). With the fix the
        # example's own declaration is kept and the module loads without
        # the duplicate error.
        example = ExampleDir(MULTILINE_INTERP_IN_TERRAFORM_BLOCK)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "declared")
        self.assertIsNone(example.terraform_tf())
        result = example.validate()
        self.assert_no_duplicate(result["diagnostics"])

    def test_multiline_interpolation_in_locals_validates(self):
        # The same cross-line interpolation shape in a locals block: the
        # declaration is kept and the prepared module still loads.
        example = ExampleDir(MULTILINE_INTERP_IN_LOCALS)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "declared")
        result = example.validate()
        self.assert_no_duplicate(result["diagnostics"])

    def test_multiline_interpolation_with_object_literal_validates(self):
        # QA-R3-1 (A-shape, end-to-end): the object literal used to hide
        # the declaration, the helper injected terraform.tf and terraform
        # reported "Duplicate required providers configuration" (RL-11
        # recurred). With the depth fix the declaration is kept and the
        # module loads without the duplicate error.
        example = ExampleDir(MULTILINE_INTERP_OBJECT_LITERAL_IN_TERRAFORM_BLOCK)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "declared")
        self.assertIsNone(example.terraform_tf())
        result = example.validate()
        self.assert_no_duplicate(result["diagnostics"])

    def test_multiline_interpolation_with_for_object_validates(self):
        # QA-R3-1 (C-shape, end-to-end): same guarantee for a for-object
        # expression inside the multi-line interpolation.
        example = ExampleDir(MULTILINE_INTERP_FOR_OBJECT_IN_TERRAFORM_BLOCK)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "declared")
        self.assertIsNone(example.terraform_tf())
        result = example.validate()
        self.assert_no_duplicate(result["diagnostics"])

    def test_multiline_directive_validates(self):
        # AI 260951912 (end-to-end): a cross-line template directive must
        # not hide the declaration; the prepared module loads without the
        # duplicate error.
        example = ExampleDir(MULTILINE_DIRECTIVE_IN_TERRAFORM_BLOCK)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "declared")
        self.assertIsNone(example.terraform_tf())
        result = example.validate()
        self.assert_no_duplicate(result["diagnostics"])

    def test_multiline_interpolation_object_literal_completion_validates(self):
        # QA-R3-1 (A-shape, completed path): an object-literal interpolation
        # before an other-providers-only block; the completed module still
        # loads with both provider requirements visible.
        main = MULTILINE_INTERP_OBJECT_LITERAL_IN_TERRAFORM_BLOCK.replace(
            'alicloud = {\n      source  = "registry.terraform.io/aliyun/alicloud"\n      version = "~> 1.220"\n    }',
            'random = {\n      source = "hashicorp/random"\n    }',
        )
        example = ExampleDir(main)
        self.addCleanup(example.close)
        self.assertEqual(example.prepare(), "completed")
        result = example.validate()
        self.assert_no_duplicate(result["diagnostics"])
        details = " ".join(d.get("detail", "") for d in result["diagnostics"])
        self.assertIn("registry.terraform.io/aliyun/alicloud", details)
        self.assertIn("hashicorp/random", details)


if __name__ == "__main__":
    unittest.main()
