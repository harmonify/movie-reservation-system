# Tx001 - Seat reservation feature

Regular users should be able to reserve seats for a showtime.

## Technical Requirements

- [ ] `POST /reservations`
  - [ ] Payment process is sync but reservation process will be done async. user -> reservation-service -> kafka -> reservation-process
- [ ] `GET /reservations`
  - [ ] for history

## Implementation

### Changes
