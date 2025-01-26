-- +migrate Up
CREATE TABLE IF NOT EXISTS public.users (
    "uuid" uuid DEFAULT gen_random_uuid() NOT NULL,
    trace_id UUID NOT NULL,
    username text NOT NULL,
    "password" text NOT NULL,
    email text NOT NULL,
    phone_number text NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    is_email_verified bool NOT NULL DEFAULT false,
    is_phone_number_verified bool NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,
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