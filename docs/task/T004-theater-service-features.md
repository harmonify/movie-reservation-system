# T004 - Theater service features

Build theater service features.

## Technical Requirements

- [x] Theater service database design
  - [x] Each theater should have a title, location, seats for each room.
  - [x] Location will be a single composite field for simplicity.
- [ ] Theater data seed
- [x] Theater Admin API
  - [x] `GET /v1/admin/theaters`
  - [x] `GET /v1/admin/theaters/:theaterId`
  - [x] `POST /v1/admin/theaters`
  - [x] `PUT /v1/admin/theaters/:theaterId`
  - [x] `DELETE /v1/admin/theaters/:theaterId`
- [x] Showtime Admin API
  - [x] `GET /v1/admin/showtimes`
  - [x] `GET /v1/admin/showtimes/:showtimeId`
  - [x] `POST /v1/admin/showtimes`
  - [x] `PUT /v1/admin/showtimes/:showtimeId`
  - [x] `DELETE /v1/admin/showtimes/:showtimeId`
- [ ] Customer API
  - [ ] `GET /v1/showtimes/:showtimeId` to get available seats (Note: performance will be improved using cache at later ticket).
