# T001 - Movie service features

Build movie service features.

Admins should be able to manage movies and showtimes.

## Technical Requirements

- [ ] [Movie service database design](#database-design)
  - [ ] Each movie should have a title, description, and poster image.
  - [ ] Movies should be categorized by genre.
  - [ ] Movies should have showtimes.
- [ ] `POST /movies`
  - [ ] Check admin role
- [ ] `POST /movies/:id/showtimes`
  - [ ] Check admin role
- [ ] `GET /movies`
  - [ ] Filter by keyword, date, genre, and actors. (Note: performance will be improved using cache at later ticket.)
- [ ] `GET /movies/:id`
  - [ ] Include future showtimes (now - 1 week later), grouped by date then theater id then room id (but show the time), and its available seats count (Note: performance will be improved using cache at later ticket)
- [ ] `GET /movies/:movieId/rooms/:roomId`
  - [ ] Get available seats. (Note: performance will be improved using cache at later ticket.)
- [ ] `PUT /movies/:id`
  - [ ] Check admin role
- [ ] `DELETE /movies/:id`
  - [ ] Check admin role

## Implementation

### Database design

```puml
entity movie {
    movie_id uuid primary key default uuid4
    title varchar(512) not null
    description text not null
    poster_image text not null # valid url
    duration int4 not null
    created_at datetime default now
}

entity genre {
    genre_id uuid primary key default uuid4
    title varchar(127)
}

entity movie_genre {
    genre_id uuid not null
    movie_id uuid not null
    constraint(genre_id, movie_id) unique
}

entity movie_showtime {
    movie_id uuid not null
    start_at datetime not null
    theater_id uuid not null
    room_id uuid not null
}
```

### API

#### `POST /movies`

#### `GET /movies`

#### `GET /movies/:id`

#### `PUT /movies/:id`

#### `DELETE /movies/:id`

### Changes
