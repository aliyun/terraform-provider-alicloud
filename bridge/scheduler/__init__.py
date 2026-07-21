"""Public contracts for the import-safe Bridge scheduler core."""

from .model import (
    AdaptiveSchedule, CapabilityValidationContext, CheckpointUpgradePolicy, CommandRunner, DailySchedule,
    HandlerRunner, HeadlessRunner, IntervalSchedule, JobPurpose, JobResult,
    JobResultStatus, MisfirePolicy, ReplayPolicy, RunnerDefinition, RunnerKind,
    ScheduleDefinition, ScheduleKind, ScheduledJobDefinition, TaskObservation,
    definition_digest, definition_snapshot, validate_job_definition, validate_registry,
)
from .planner import TriggerPlanner
from .engine import (
    DurableResultPublisher, ExecutionDisposition, ExecutionOutcome, JobRegistration,
    ScheduledJobControlPlane, ScheduledJobState, ScheduledJobStatus, ScheduledSlot,
    SchedulerEngine, build_registrations, plan_due_slots,
)
from .runtime import (
    JobRunner, PermanentJobError, RetryableJobError, ScannerInvocation,
    ScannerRuntime, ScannerStopped,
)
from .control_plane_client import (
    HttpScheduledJobControlPlane, ScheduledJobControlPlaneError,
    ScheduledJobControlPlaneProtocolError, ScheduledJobControlPlaneRejected,
    ScheduledJobControlPlaneUnavailable,
)
from .composition import (
    EmptyResultPublisher, SCHEDULER_HOST_ID, SCHEDULER_WORKER_KEY,
    SchedulerComposition, SchedulerCompositionError,
)
from .migration import (
    JOB_ENV_PREFIX, SCHEDULER_ENABLE_ENV, SchedulerMigrationError,
    job_env_name, requested_new_jobs, scheduler_enabled, uses_new_engine,
)

__all__ = [
    "AdaptiveSchedule", "CapabilityValidationContext", "CheckpointUpgradePolicy", "CommandRunner", "DailySchedule",
    "HandlerRunner", "HeadlessRunner", "IntervalSchedule", "JobPurpose", "JobResult",
    "JobResultStatus", "MisfirePolicy", "ReplayPolicy", "RunnerDefinition", "RunnerKind",
    "ScheduleDefinition", "ScheduleKind", "ScheduledJobDefinition", "TaskObservation",
    "TriggerPlanner", "DurableResultPublisher", "ExecutionDisposition", "ExecutionOutcome",
    "JobRegistration", "ScheduledJobControlPlane", "ScheduledJobState", "ScheduledJobStatus",
    "ScheduledSlot", "SchedulerEngine", "build_registrations", "plan_due_slots", "JobRunner",
    "PermanentJobError", "RetryableJobError", "ScannerInvocation", "ScannerRuntime",
    "ScannerStopped", "definition_digest", "definition_snapshot", "validate_job_definition",
    "validate_registry", "HttpScheduledJobControlPlane", "ScheduledJobControlPlaneError",
    "ScheduledJobControlPlaneProtocolError", "ScheduledJobControlPlaneRejected",
    "ScheduledJobControlPlaneUnavailable",
    "EmptyResultPublisher", "SCHEDULER_HOST_ID", "SCHEDULER_WORKER_KEY",
    "SchedulerComposition", "SchedulerCompositionError", "JOB_ENV_PREFIX",
    "SCHEDULER_ENABLE_ENV", "SchedulerMigrationError", "job_env_name",
    "requested_new_jobs", "scheduler_enabled", "uses_new_engine",
]
