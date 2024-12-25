-- +migrate Up
CREATE TABLE IF NOT EXISTS public.users (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    username text NULL,
    "password" text NULL,
    email text NULL,
    phone_number text NULL,
    first_name text NULL,
    last_name text NULL,
    is_email_verified bool DEFAULT false NULL,
    is_phone_number_verified bool DEFAULT false NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    CONSTRAINT uni_users_email UNIQUE (email),
    CONSTRAINT uni_users_phone_number UNIQUE (phone_number),
    CONSTRAINT uni_users_username UNIQUE (username),
    CONSTRAINT users_pkey PRIMARY KEY (uuid)
);

CREATE UNIQUE INDEX idx_user_email ON public.users USING btree (email);

CREATE UNIQUE INDEX idx_user_phone_number ON public.users USING btree (phone_number);

CREATE UNIQUE INDEX idx_user_username ON public.users USING btree (username);

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS public.users;