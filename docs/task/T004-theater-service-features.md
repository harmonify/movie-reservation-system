# T001 - Movie service features

Build movie service features.

## Technical Requirements

- [ ] [Theater service database design](#database-design)
  - [ ] Each theater should have a title, location, seats for each room.
  - [ ] Location will be a single composite field for simplicity.
- [ ] Theater data seed
- [ ] `GET /theater`
- [ ] `GET /theater/:id`

## Implementation

### Database design

```puml
entity theater {
    theater_id uuid primary key default uuid4
    title varchar(512) not null
    location text not null
}

entity theater_room {
    room_id uuid primary key default uuid4
    theater_id uuid not null
    title varchar(127)
}

entity seat {
    room_id uuid not null
    row varchar(255) not null
    column int4 not null
    constraint(room_id, row, column) unique
}
```

### API

#### `GET /theater`

#### `GET /theater/:id`

### Changes
