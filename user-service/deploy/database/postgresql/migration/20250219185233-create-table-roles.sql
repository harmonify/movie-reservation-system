-- +migrate Up
CREATE TABLE IF NOT EXISTS public.roles (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);

INSERT INTO
    public.roles (name)
VALUES
    ('admin');

INSERT INTO
    public.roles (name)
VALUES
    ('user');

-- +migrate Down
DROP TABLE IF EXISTS public.roles;