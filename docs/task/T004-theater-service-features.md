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
- [x] Customer API
  - [x] `GET /v1/theaters` to get all theaters.
    - [x] Each theater should have the following details: Title, Address, Phone, Email, Website, Location (latitude, longitude).
    - [x] Client should be able to filter theaters by radius (distance from the user's location). The minimum radius filter should be 100 meters.
  - [x] ~~`GET /v1/theaters/:theaterId` to get theater details~~.
  - [x] ~~`GET /v1/theaters/:theaterId/showtimes` to get all active showtimes (ongoing + upcoming 7 days). Showtime should have the following details: Theater name, Movie title, Room number, Start time, End time, Available seats count~~ (This API is not required since movie-search-service has a similar endpoint)
  - [x] `GET /v1/showtimes/:showtimeId` to get showtime details.
    - [x] Showtime should have the following details: Theater name, Movie title, Room name, Start time, End time, Seats details (seat number, seat status)
