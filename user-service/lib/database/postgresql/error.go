package postgresql

// https://www.postgresql.org/docs/11/errcodes-appendix.html
var (
	// Class 23 — Integrity Constraint Violation
	IntegrityConstraintViolation = "23000"
	RestrictViolation            = "23001"
	NotNullViolation             = "23502"
	ForeignKeyViolation          = "23503"
	UniqueViolation              = "23505"
	CheckViolation               = "23514"
	ExclusionViolation           = "23P01"

	// Class 40 — Transaction Rollback
	TransactionRollback                     = "40000"
	TransactionIntegrityConstraintViolation = "40002"
	DeadlockDetected                        = "40P01"

	// Class 42 — Syntax Error or Access Rule Violation
	SyntaxError           = "42601"
	InsufficientPrivilege = "42501"
)
