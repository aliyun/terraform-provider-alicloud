from __future__ import annotations

import unittest

from bridge.scheduler.role import require_scheduler_role


class SchedulerRoleTests(unittest.TestCase):
    def test_scheduler_role_is_allowed(self):
        require_scheduler_role({"JARVIS_BRIDGE_ROLE": "scheduler"})

    def test_worker_role_is_rejected_before_scheduler_startup(self):
        with self.assertRaisesRegex(RuntimeError, "JARVIS_BRIDGE_ROLE=scheduler"):
            require_scheduler_role({"JARVIS_BRIDGE_ROLE": "worker"})

    def test_unknown_role_is_rejected(self):
        with self.assertRaisesRegex(RuntimeError, "JARVIS_BRIDGE_ROLE=scheduler"):
            require_scheduler_role({"JARVIS_BRIDGE_ROLE": "unknown"})


if __name__ == "__main__":
    unittest.main()
