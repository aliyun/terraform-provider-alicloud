#!/usr/bin/env python3
"""Portable input and legacy redispatch reconstruction tests."""

import unittest

from bridge.task_input_contract import (
    PORTABLE_INPUT_CONTRACT,
    HEADLESS_POLICY_REVISION,
    RUNTIME_INTERACTIVE,
    RUNTIME_PERSISTENT,
    TaskInputContractError,
    canonical_input_digest,
    portable_replacement_for_redispatch,
    portable_task_payload,
    validate_portable_task_input,
)


class TaskInputContractTest(unittest.TestCase):
    def test_payload_is_symmetric_across_runtime_modes(self):
        for runtime in ("INTERACTIVE", "PERSISTENT"):
            with self.subTest(runtime=runtime):
                payload = portable_task_payload(
                    item_id="84407231",
                    project="1086837",
                    kind="ticket",
                    prompt="continue",
                    origin_runtime=runtime,
                    trigger="TEST",
                )
                self.assertEqual(
                    payload["inputContract"], PORTABLE_INPUT_CONTRACT)
                self.assertEqual(
                    payload["inputContext"]["originRuntime"], runtime)

    def test_legacy_canonical_aone_input_is_rehydrated(self):
        source = {
            "kind": "ticket",
            "itemId": "84407231",
            "project": "1086837",
            "trigger": "INTERACTIVE",
            "policyRevision": "terraform-rd-single-writer-v6",
        }
        timeline = {
            "task": {
                "id": 603,
                "taskKey": "aone:1086837:84407231",
                "aoneId": "84407231",
                "sourceType": "AONE",
                "taskType": "ticket",
                "desiredRevision": "interactive:old",
                "processingRevision": "interactive:old",
                "currentSessionId": 1195,
                "sourceRef": {
                    "aoneId": "84407231",
                    "projectId": "1086837",
                    "title": "legacy interactive task",
                },
            },
            "sessions": [{
                "id": 1195,
                "inputPayload": source,
            }],
            "effectiveInputPayload": source,
            "effectiveInputDigest": canonical_input_digest(source),
        }

        replacement = portable_replacement_for_redispatch(
            timeline,
            target_runtime="PERSISTENT",
            selected_session_id=1195,
        )

        self.assertIsNotNone(replacement)
        payload = replacement["payload"]
        self.assertEqual(payload["itemId"], "84407231")
        self.assertEqual(payload["project"], "1086837")
        self.assertEqual(payload["inputContract"], "PORTABLE_V1")
        self.assertTrue(payload["inputContext"]["rehydrated"])
        self.assertIn("读取当前 Aone 详情和最新评论", payload["prompt"])
        self.assertIn("禁止重复评论", payload["prompt"])
        self.assertNotEqual(
            replacement["desiredRevision"], "interactive:old")
        self.assertEqual(
            replacement["replacementDigest"],
            canonical_input_digest(payload))
        self.assertEqual(
            replacement["desiredRevision"],
            "portable:%s" % replacement["replacementDigest"])

    def test_non_aone_legacy_input_is_typed_blocked(self):
        source = {"itemId": "local-1", "kind": "adhoc"}
        timeline = {
            "task": {
                "id": 1,
                "taskKey": "local:1",
                "sourceType": "LOCAL",
                "taskType": "adhoc",
                "currentSessionId": 2,
            },
            "effectiveInputPayload": source,
            "effectiveInputDigest": canonical_input_digest(source),
        }
        with self.assertRaises(TaskInputContractError) as raised:
            portable_replacement_for_redispatch(
                timeline, target_runtime="INTERACTIVE",
                selected_session_id=2)
        self.assertEqual(raised.exception.code, "INPUT_NOT_REHYDRATABLE")


class PolicyRevisionValidationTest(unittest.TestCase):
    """A PERSISTENT executor refuses a payload whose policyRevision is missing
    or stale, so rehydrate rebuilds it instead of passing the gap to the
    executor (where it surfaced as stale_task_policy_revision: frozen=<missing>).
    """

    BASE = {"itemId": "1", "project": "2", "prompt": "p",
            "inputContract": PORTABLE_INPUT_CONTRACT}

    def test_persistent_requires_current_policy_revision(self):
        result = validate_portable_task_input(
            dict(self.BASE, policyRevision=HEADLESS_POLICY_REVISION),
            target_runtime=RUNTIME_PERSISTENT)
        self.assertEqual(result["policyRevision"], HEADLESS_POLICY_REVISION)

    def test_persistent_rejects_missing_policy_revision(self):
        with self.assertRaisesRegex(
                TaskInputContractError, "POLICY_REVISION_NOT_CURRENT"):
            validate_portable_task_input(
                self.BASE, target_runtime=RUNTIME_PERSISTENT)

    def test_persistent_rejects_stale_policy_revision(self):
        with self.assertRaisesRegex(
                TaskInputContractError, "POLICY_REVISION_NOT_CURRENT"):
            validate_portable_task_input(
                dict(self.BASE, policyRevision="terraform-rd-single-writer-v4"),
                target_runtime=RUNTIME_PERSISTENT)

    def test_interactive_does_not_require_policy_revision(self):
        result = validate_portable_task_input(
            self.BASE, target_runtime=RUNTIME_INTERACTIVE)
        self.assertNotIn("policyRevision", result)


class RehydratePolicyRevisionTest(unittest.TestCase):
    """A rehydrated payload that lost its policyRevision is rebuilt with v6."""

    def _timeline(self, policy_revision):
        payload = {
            "itemId": "83450138", "project": "1086837", "kind": "ticket",
            "prompt": "查看aone最新评论", "inputContract": PORTABLE_INPUT_CONTRACT,
            "inputContext": {"originRuntime": "INTERACTIVE", "rehydrated": True},
        }
        if policy_revision is not None:
            payload["policyRevision"] = policy_revision
        return {
            "task": {"id": 1052, "aoneId": "83450138",
                     "taskKey": "aone:1086837:83450138",
                     "sourceRef": {"aoneId": "83450138", "projectId": "1086837",
                                   "title": "t"},
                     "desiredRevision": "portable:old",
                     "processingRevision": "portable:old"},
            "effectiveInputPayload": payload,
        }

    def test_missing_policy_revision_is_rebuilt_with_current(self):
        rep = portable_replacement_for_redispatch(
            self._timeline(None), target_runtime=RUNTIME_PERSISTENT)
        self.assertIsNotNone(rep, "missing policyRevision must trigger rebuild")
        self.assertEqual(rep["payload"]["policyRevision"], HEADLESS_POLICY_REVISION)

    def test_stale_policy_revision_is_rebuilt_with_current(self):
        rep = portable_replacement_for_redispatch(
            self._timeline("terraform-rd-single-writer-v4"),
            target_runtime=RUNTIME_PERSISTENT)
        self.assertIsNotNone(rep)
        self.assertEqual(rep["payload"]["policyRevision"], HEADLESS_POLICY_REVISION)

    def test_current_policy_revision_is_not_rebuilt(self):
        """A payload already carrying the current revision is portable as-is."""
        rep = portable_replacement_for_redispatch(
            self._timeline(HEADLESS_POLICY_REVISION),
            target_runtime=RUNTIME_PERSISTENT)
        self.assertIsNone(rep)


if __name__ == "__main__":
    unittest.main()
