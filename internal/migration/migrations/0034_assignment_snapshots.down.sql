-- Rollback snapshot metadata changes

DROP INDEX IF EXISTS idx_billing_assignments_user_status;

ALTER TABLE billing_operation_assignments
DROP COLUMN IF EXISTS snapshot_metadata;
