# Error

In general, there are two types of errors:

- business error, and
- system error.

## Business error

A business error occurs when the operation fails due to invalid business logic or rule violations. These errors are often predictable and result from conditions such as invalid input, failed validations, or conflicts with the system's business rules.

Examples of business errors:

- Submitting a form with missing required fields.
- Attempting to purchase an out-of-stock item.
- Violating a constraint, such as exceeding a credit limit.

### Handling business error

The most straightforward way to handle a business error is to do a backwards recovery, where the system compensates for the failed operation by reversing any changes made. This ensures the system remains in a consistent state. See [Saga pattern](../dev/saga.md).

We should provide detailed error messages to guide users either in correcting their input or retry again later.

We also need to log the error for audit or debugging purposes.

## System error

A system error occurs when the operation fails due to technical issues, such as unavailable resources, network problems, or unexpected failures in dependencies. These errors are typically less predictable and arise from infrastructure or external factors.

Examples of system errors:

- A database connection timeout.
- A dependency service is down or unreachable.
- Disk space exhaustion or server crash.

### Handling system error

Retry with exponential backoff is a common solution for transient system errors. It involves retrying the failed operation at increasing time intervals, which helps manage temporary issues without overwhelming the system. We also need to consider designing idempotent operations to safely retry without side effects.

Other solutions include:

- Implementing circuit breakers to prevent cascading failures.
- Monitoring and alerting to quickly detect and address recurring issues.
