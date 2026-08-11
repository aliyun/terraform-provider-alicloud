#!/usr/bin/env python3
"""Tests for the fail-closed CloudSpec operation Git diff guard."""

import os
import shutil
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path
from unittest import mock

sys.path.insert(0, str(Path(__file__).resolve().parent.parent))

from bootstrap.cloudspec_operation_diff_guard import (
    OperationDiffGuardError,
    capture_repository_baseline,
    discover_cloudspec_components,
    inspect_repository,
    path_is_cloudspec_operation,
)


class _GitRepository:
    def __init__(self, root: Path):
        self.root = root
        self.git("init", "-q")
        self.git("config", "user.name", "Operation Guard Test")
        self.git("config", "user.email", "operation-guard@example.invalid")

    def git(self, *args: str) -> str:
        return subprocess.check_output(
            ["git", "-C", os.fspath(self.root), *args],
            text=True,
            stderr=subprocess.STDOUT,
        ).strip()

    def write(self, relative: str, content: str = "model") -> Path:
        target = self.root / relative
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(content, encoding="utf-8")
        return target

    def commit_all(self, message: str = "snapshot") -> str:
        self.git("add", "-A")
        self.git("commit", "-qm", message)
        return self.git("rev-parse", "HEAD")

    def seed_component(self, root: str = "service") -> None:
        prefix = root + "/" if root else ""
        self.write(prefix + "main.cspec", "namespace example")
        self.write(prefix + "operations/get.cspec", "operation get unique")
        self.write(prefix + "operations/update.cspec", "operation update unique")
        self.write(prefix + "operations/delete.cspec", "operation delete unique")
        self.write(prefix + "operations/copy-source.cspec",
                   "operation copy source exact\nline two\nline three")
        self.write(prefix + "resources/widget.cspec", "resource widget")


class CloudspecOperationDiffGuardTest(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory()
        self.root = Path(self.temporary.name)
        self.repo = _GitRepository(self.root)
        self.repo.seed_component()
        self.repo.commit_all("initial component")

    def tearDown(self):
        self.temporary.cleanup()

    def test_discovery_and_write_path_classification_use_actual_git_repo(self):
        self.repo.seed_component("组件 空间")
        self.repo.git("add", "组件 空间/main.cspec")

        self.assertEqual(
            discover_cloudspec_components(self.root),
            ("service", "组件 空间"),
        )
        self.assertTrue(path_is_cloudspec_operation(
            self.root, "service/operations/new op.cspec"))
        self.assertTrue(path_is_cloudspec_operation(
            self.root, self.root / "组件 空间/operations/查询.cspec"))
        self.assertFalse(path_is_cloudspec_operation(
            self.root, "service/resources/widget.cspec"))
        self.assertFalse(path_is_cloudspec_operation(
            self.root, "unrelated/operations/get.cspec"))
        with self.assertRaises(OperationDiffGuardError):
            path_is_cloudspec_operation(self.root, self.root.parent / "outside.cspec")

    def test_committed_staged_unstaged_and_untracked_changes_are_combined(self):
        baseline = capture_repository_baseline(self.root)

        self.repo.write("service/operations/get.cspec", "committed edit")
        self.repo.commit_all("committed operation edit")
        self.repo.write("service/operations/update.cspec", "unstaged edit")
        self.repo.write("service/operations/staged new.cspec", "staged add")
        self.repo.git("add", "service/operations/staged new.cspec")
        self.repo.write("service/operations/未跟踪 空格.cspec", "untracked add")

        result = inspect_repository(self.root, baseline)

        self.assertFalse(result.blocked, result.reason)
        self.assertEqual(result.current_branch, baseline.branch)
        self.assertEqual(result.baseline_sha, baseline.head)
        self.assertEqual(set(result.paths), {
            "service/operations/get.cspec",
            "service/operations/update.cspec",
            "service/operations/staged new.cspec",
            "service/operations/未跟踪 空格.cspec",
        })
        by_path = {
            change.new_path or change.old_path: change.status
            for change in result.changes
        }
        self.assertEqual(by_path["service/operations/get.cspec"], "M")
        self.assertEqual(by_path["service/operations/update.cspec"], "M")
        self.assertEqual(by_path["service/operations/staged new.cspec"], "A")
        self.assertEqual(by_path["service/operations/未跟踪 空格.cspec"], "A")

    def test_delete_rename_and_copy_preserve_old_and_new_paths(self):
        baseline = capture_repository_baseline(self.root)
        self.repo.git("rm", "-q", "service/operations/delete.cspec")
        self.repo.git(
            "mv", "service/operations/get.cspec",
            "service/operations/renamed 获取.cspec",
        )
        shutil.copyfile(
            self.root / "service/operations/copy-source.cspec",
            self.root / "service/operations/copied 副本.cspec",
        )
        self.repo.git("add", "service/operations/copied 副本.cspec")

        result = inspect_repository(self.root, baseline)

        self.assertFalse(result.blocked, result.reason)
        self.assertIn(
            ("R", "service/operations/get.cspec",
             "service/operations/renamed 获取.cspec"),
            {(change.status, change.old_path, change.new_path)
             for change in result.changes},
        )
        self.assertIn(
            ("C", "service/operations/copy-source.cspec",
             "service/operations/copied 副本.cspec"),
            {(change.status, change.old_path, change.new_path)
             for change in result.changes},
        )
        self.assertIn(
            ("D", "service/operations/delete.cspec", None),
            {(change.status, change.old_path, change.new_path)
             for change in result.changes},
        )
        self.assertIn("service/operations/get.cspec", result.paths)
        self.assertIn("service/operations/renamed 获取.cspec", result.paths)
        self.assertIn("service/operations/copy-source.cspec", result.paths)
        self.assertIn("service/operations/copied 副本.cspec", result.paths)

    def test_root_component_and_filename_with_newline_are_nul_safe(self):
        self.repo.seed_component("")
        self.repo.commit_all("root component")
        baseline = capture_repository_baseline(self.root)
        unusual = "operations/line one\n第二行.cspec"
        self.repo.write(unusual, "untracked")

        result = inspect_repository(self.root, baseline)

        self.assertFalse(result.blocked, result.reason)
        self.assertIn("", result.components)
        self.assertEqual(result.paths, (unusual,))

    def test_non_cloudspec_operations_and_other_component_files_do_not_match(self):
        baseline = capture_repository_baseline(self.root)
        self.repo.write("scripts/operations/job.py", "not CloudSpec")
        self.repo.write("service/resources/widget.cspec", "resource changed")
        self.repo.write("service/README.md", "docs")

        result = inspect_repository(self.root, baseline)

        self.assertFalse(result.blocked, result.reason)
        self.assertFalse(result.has_operation_changes)
        self.assertEqual(result.paths, ())
        self.assertEqual(result.changes, ())

    def test_non_operation_type_change_does_not_block_repository_inspection(self):
        baseline = capture_repository_baseline(self.root)
        resource = self.root / "service/resources/widget.cspec"
        resource.unlink()
        resource.symlink_to("../main.cspec")

        result = inspect_repository(self.root, baseline)

        self.assertFalse(result.blocked, result.reason)
        self.assertEqual(result.paths, ())
        self.assertEqual(result.changes, ())

    def test_deleted_main_still_uses_baseline_component(self):
        baseline = capture_repository_baseline(self.root)
        self.repo.git("rm", "-q", "service/main.cspec")
        self.repo.git("rm", "-q", "service/operations/get.cspec")

        result = inspect_repository(self.root, baseline)

        self.assertFalse(result.blocked, result.reason)
        self.assertEqual(result.paths, ("service/operations/get.cspec",))

    def test_wholly_untracked_component_is_detected(self):
        baseline = capture_repository_baseline(self.root)
        self.repo.write("新组件/main.cspec", "namespace new")
        self.repo.write("新组件/operations/create item.cspec", "operation create")

        result = inspect_repository(self.root, baseline)

        self.assertFalse(result.blocked, result.reason)
        self.assertEqual(result.paths, ("新组件/operations/create item.cspec",))

    def test_invalid_repo_and_baseline_are_blocked_without_exception(self):
        invalid_repo = inspect_repository(self.root / "missing", "not-a-sha")
        self.assertTrue(invalid_repo.blocked)
        self.assertTrue(invalid_repo.reason)

        invalid_type = inspect_repository(None, None)  # type: ignore[arg-type]
        self.assertTrue(invalid_type.blocked)
        self.assertTrue(invalid_type.reason)

        invalid_baseline = inspect_repository(self.root, "not-a-sha")
        self.assertTrue(invalid_baseline.blocked)
        self.assertIn("baseline SHA", invalid_baseline.reason)
        self.assertEqual(invalid_baseline.current_head,
                         self.repo.git("rev-parse", "HEAD"))

    def test_git_timeout_is_returned_as_blocked_without_exception(self):
        baseline = capture_repository_baseline(self.root)
        with mock.patch(
                "bootstrap.cloudspec_operation_diff_guard.subprocess.run",
                side_effect=subprocess.TimeoutExpired(["git", "rev-parse"], 5)):
            result = inspect_repository(self.root, baseline)

        self.assertTrue(result.blocked)
        self.assertIn("timed out after 5.0 seconds", result.reason)
        self.assertEqual(result.paths, ())

    def test_branch_change_and_detached_head_are_blocked(self):
        baseline = capture_repository_baseline(self.root)
        self.repo.git("checkout", "-qb", "other-branch")

        changed_branch = inspect_repository(self.root, baseline)
        self.assertTrue(changed_branch.blocked)
        self.assertIn("branch changed", changed_branch.reason)

        self.repo.git("checkout", "-q", "--detach")
        detached = inspect_repository(self.root, baseline.head)
        self.assertTrue(detached.blocked)
        self.assertIn("detached", detached.reason)

    def test_non_ancestor_baseline_is_blocked(self):
        baseline = capture_repository_baseline(self.root)
        self.repo.git("checkout", "-q", "--orphan", "diverged")
        self.repo.git("commit", "-qm", "unrelated root")

        result = inspect_repository(self.root, baseline.head)

        self.assertTrue(result.blocked)
        self.assertIn("merge-base", result.reason)


if __name__ == "__main__":
    unittest.main()
