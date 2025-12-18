# GoPLC User Task Files

This directory contains user-written Go task files that are auto-discovered by the GoPLC runtime.

**Setup:** Story 1.5 (Task Discovery) + Story 1.8 (Example Task)

**Task Interface:**
Tasks must implement the Task interface (defined in Story 1.5):
- `Name() string` - Task identifier
- `Config() TaskConfig` - Scan rate, enabled status
- `Execute(ctx TaskContext) error` - Control logic

**Example:** An example task will be added in Story 1.8.