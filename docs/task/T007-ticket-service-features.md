# Tx002 - Ticket service features

Build ticket service features.

## Technical Requirements

- [ ] `POST /code`
  - [ ] User should be able to get ticket information using code they got from successful reservation process. QR Code?

## Implementation

### Possible solution

Get information about the movies, theater, movie showtime, and seat, using gRPC client from ticket service.
And then construct a PDF which then printer machine on theater could print them.

### Changes
