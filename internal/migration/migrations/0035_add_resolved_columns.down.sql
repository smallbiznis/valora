-- Rollback: Remove resolved_at and resolved_by columns

DROP INDEX IF EXISTS idx_billing_assignments_resolved_at;

ALTER TABLE billing_operation_assignments
  DROP COLUMN IF EXISTS resolved_at,
  DROP COLUMN IF EXISTS resolved_by;
