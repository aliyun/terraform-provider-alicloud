#!/usr/bin/env python3
"""Focused tests for durable Task policy generations."""

import unittest

from bridge.task_policy import (
    HEADLESS_POLICY_REVISION,
    policy_desired_revision,
)


class TaskPolicyTest(unittest.TestCase):
    def test_current_revision_is_v6(self):
        self.assertEqual(
            HEADLESS_POLICY_REVISION, "terraform-rd-single-writer-v6")

    def test_policy_and_input_are_part_of_desired_revision(self):
        payload = {
            "itemId": "84846271",
            "kind": "ticket",
            "prompt": "continue on source ticket",
            "policyRevision": HEADLESS_POLICY_REVISION,
        }
        current = policy_desired_revision("comment:11", payload)
        old = policy_desired_revision(
            "comment:11", dict(payload, policyRevision="old"),
            policy_revision="terraform-rd-single-writer-v5")
        changed = policy_desired_revision(
            "comment:11", dict(payload, prompt="changed prompt"))

        self.assertIn("|policy:v6|input:", current)
        self.assertNotEqual(current, old)
        self.assertNotEqual(current, changed)
        self.assertLessEqual(len(current), 127)


if __name__ == "__main__":
    unittest.main()
