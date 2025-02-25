# T005 - Order service features

Regular users should be able to reserve tickets for a showtime.

## Technical Requirements

- [ ] `POST /v1/orders`
  - [ ] Payment process is sync but reservation process will be done async. user -> order-service -> kafka -> order-processor
- [ ] `GET /v1/orders`
  - [ ] for history
- [ ] `GET /v1/orders/:orderId`
  - [ ] for history. Include payment information.

## Implementation

### order-processor

Handle new reservation and if error happened, refund payment.

### Changes
