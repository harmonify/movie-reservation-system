# Event-Driven Architecture

> P.S. Read [this document](./cdc-and-transactional-outbox.md) first to understand the context of this document.

This document outlines the pros and cons of different approaches when building this system with event-driven architecture mindset.

## Sending Account Verification Email

We will take an example of sending account verification email to the user when they registers in the system.

## Approach 1: user service --gRPC request--> notification service

User service sends a gRPC request to the notification service, the notification service receives the request, sends the account verification email to the user, and reply back to the user service.

### Pros

- Easy to implement
- Easy to test E2E
- Easy to debug

### Cons

- Message loss when system error (e.g. network error, host crashes, etc.) occurs.
- High temporal coupling:
  - At least one instance of notification service must be up and ready when the user service sends a request.
  - Necessity to implement client-side / 3rd-party service discovery: To be highly available and scalable.
  - Necessity to implement resilient and fault-tolerant strategies: There is always the possibility that the other service is unavailable or exhibits such high latency it is essentially unusable. We need to implement resilient patterns, such as circuit breaker to prevent flooding downstream services.

## Approach 2: user service --produce message--> message broker <-- subscribe & consume message --> notification service

User service publishes events to `public.user.registered.v1` topic on the message broker, then the notification service consumes the event and sends an account verification email to the user.

### Pros

- Message is not lost when system error occurs. The message is durably persisted in the message broker.
- Low temporal coupling:
  - If the notification service is unavailable when the user service publishes the event, it may consume and process the message later when it is up and ready.
- Scalable asynchronous communication:
  - The user service can publish the event and continue processing other requests without waiting for the notification service to respond.
    - The intermediary service can implement retry mechanisms or delayed retries to handle transient failures in the notification service, while still maintaining low synchronous response time to the user.
      - For example, when a user registers on the system through a client that makes HTTP request to the user service, we could return a HTTP response fast enough (<2 seconds) to the client because the user service does not concern itself to sends account verification email to the user synchronously.
        - The intermediary service could asynchronously (from the user perspective) makes gRPC calls to the notification service and handle transient failures with retry mechanism or other resilient patterns (could take up to 30 seconds). The notification service could asynchronously (from the user perspective) render and send the account verification email to the user through email provider (could take up to 10 seconds).
        - This is in contrast to the 1st approach, where a transient failure would require the user service to implement retry mechanism synchronously. Users could possibly wait too long that they lose interest in using our system. You can call me being speculative here, but when this scenario do happens in real-life, it is certainly undesirable because it would negatively impact business conversion.

### Cons

- Harder to implement, compared to the 1st approach. We need to consider the internal workings of the message broker, in this case Kafka. See [Apache Kafka technical challenges](../../reference/tools/kafka.md) for more info.
- Harder to test E2E, compared to the 1st approach.
- Harder to debug, compared to the 1st approach.
- High behavioral coupling & low cohesion:
  - Notification service now concerns itself with user domain business rules to send an account verification email to the user.
  - For example, to reduce the email spam rate (business rule), we want to send no more than 1 account verification email for the same user in 24 hours (technical implementation of the business rule).
  - Another example, in a system with different user roles, e.g. regular user, admin user, or superuser, we may not want to send an account verification email to the admin user because the admin user's account itself is manually created by a superuser. Business-logically (or by-product requirements) sending the account verification email to the admin user is redundant and not necessary.
  - I made those 2 examples only to convey how unpredictable real-life production systems will evolve.

## Approach 3: Hybrid approach: user service --produce message--> message broker <-- subscribe & consume message --> intermediary service --???--> notification service

User service publishes event to `public.user.registered.v1` topic like the 2nd approach. An intermediary service (could be a dedicated event processor, rule engine, or even user service itself) in the user domain consumes the event and then proceeds the flow (possibly applying business rules related to user domain) to send an account verification email to the user.

Possible flow of the intermediary service:

- Intermediary service subscribes & consumes events from `public.user.registered.v1` topic,
- Intermediary service applies business rules, and
- Intermediary service sends gRPC request to notification service,
  - Notification service receives the request,
  - Notification service sends the account verification email to the user, and
  - Notification service replies back to the intermediary service.

### Pros

- Low temporal coupling (this differs from 1st approach):
  - User service is now temporally decoupled from notification service.
  - Temporal coupling is now shifted between the intermediary service and notification service. This coupling is acceptable because if the notification service is unavailable when the user service sends a request, the message is not lost. The message is durably persisted in the message broker like the 2nd approach.
- Low behavioral coupling & high cohesion; better separation of concerns (this differs from 2nd approach):
  - Centralized business rules, which is flexible for an ever evolving business rules.
    - Easy to reason.
    - Easy to unit test business rules.
  - Notification service is now behaviorally decoupled from user service.
  - The intermediary service, which resides within the user domain, is responsible for applying user domain business rules regarding sending an account verification email to the user asynchronously.
- Scalable, for example:
  - When a user registers on the system through a client that makes HTTP request to the user service, we could return a HTTP response fast enough (<2 seconds) to the client because the user service does not concern itself to sends account verification email to the user synchronously
  - The intermediary service could asynchronously (from the user perspective) makes gRPC calls to the notification service and handle transient failures with retry mechanism or other resilient patterns if needed (could take up to 30 seconds).
  - The notification service could asynchronously (from the user perspective) render and send the account verification email to the user through email provider (could take up to 10 seconds).
  - This is in contrast to the 1st approach, where a transient failure would require the user service to implement retry mechanism synchronously. Users could possibly wait too long that they lose interest in using our system. You can call me being speculative here, but when this scenario do happens in real-life, it is certainly undesirable because it would negatively impact business conversion.

### Cons

- Somewhat harder to implement, compared to two previous approaches.
- Equally hard to test E2E, similar to the 2nd approach, due to the introduction of multiple components.
- Equally hard to debug, like the 2nd approach, because of the distributed nature and additional layers.
- Introduces additional latency and complexity:
  - The intermediary service adds another hop in the communication flow, increasing processing time and potential failure points.
- [More network calls / remote communication points = more point of errors](https://speakerdeck.com/ufried/getting-service-design-right):
  - This increases the number of potential failure points, as each additional remote call introduces the possibility of network issues or service outages.

## Summary

Event-driven architecture offers significant benefits for building scalable and resilient systems, but it also introduces complexity. Striking the right balance between simplicity and robustness is key. Each approach has its pros and cons, and the choice depends on the system's specific requirements and constraints. Here's a summary of what I learned and when to prefer each approach:

- **Approach 1**:  
  Use this approach if:
  - Real-time responses are critical, and the services involved are highly available.
  - You prioritize simplicity and ease of implementation over fault-tolerance.

- **Approach 2**:  
  Use this approach if:
  - You can tolerate eventual consistency.
  - Durability and resilience are critical.
    - You expect temporary unavailability of services.
  - You plan for scalable asynchronous communication.
  - You are comfortable with the operational overhead of managing a message broker and debugging distributed systems.

- **Approach 3**:  
  Use this approach if:
  - You can tolerate eventual consistency.
  - Durability and resilience are critical.
  - The business logic is complex and likely to evolve over time; maintainability is important.
    - You need a clear separation of concerns for better maintainability.
    - You are prepared to handle additional latency and complexity for greater flexibility.
