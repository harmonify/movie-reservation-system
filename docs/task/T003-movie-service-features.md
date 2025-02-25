# T003 - Movie service features

Build movie service features.

Admins should be able to manage movies and showtimes.

## Technical Requirements

- [x] Movie service database design
  - [x] Each movie should have a title, description, and poster image.
  - [x] Movies should be categorized by genre.
  - [x] Movies should have showtimes.
- [x] Movie data seed
- [x] Movie Admin API (Check admin role)
  - [x] `GET /v1/admin/movies` Filter by keyword, date, genre, and actors (Note: performance will be improved using cache at later ticket.).
  - [x] `GET /v1/admin/movies/:movieId` Include future showtimes (now - 1 week later), grouped by date then theater id then room id (but show the time), and its available seats count (Note: performance will be improved using cache at later ticket).
  - [x] `POST /v1/admin/movies`
  - [x] `PUT /v1/admin/movies/:movieId`
  - [x] `DELETE /v1/admin/movies/:movieId`
- [x] Movie Customer API
  - [x] `GET /v1/movies` Filter by keyword, date, genre, and actors (Note: performance will be improved using cache at later ticket.).
