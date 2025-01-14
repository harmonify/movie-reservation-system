# Microservices

## Stateless Microservices

Stateless microservices do not maintain any internal state between requests. Each request is processed independently, without relying on data from previous interactions. Ideal for tasks like authentication, processing requests, or data transformation.

- Notification service (Email & SMS)
- User service

## Persistence Microservices

Persistent microservices maintain a state, often backed by a database or other storage. They manage and store data required for operations over time. Often tied to domain data models. Responsible for CRUD operations.

- Customer service
- Movie service
- Theater service
- Ticket service

## Aggregation Microservices

These microservices aggregate data or responses from multiple other services, combining them to fulfill a single request. Often act as a facade for multiple underlying services.

- Order service
- Report service
- Payment service
