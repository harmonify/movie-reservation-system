-- +migrate Up
CREATE TABLE public.user_sessions (
    id bigserial NOT NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    user_uuid uuid NULL,
    refresh_token text NULL,
    is_revoked bool DEFAULT false NULL,
    expired_at timestamptz NULL,
    ip_address text NULL,
    user_agent text NULL,
    CONSTRAINT user_sessions_pkey PRIMARY KEY (id)
);

CREATE INDEX idx_user_sessions_deleted_at ON public.user_sessions USING btree (deleted_at);

CREATE INDEX idx_user_sessions_user_uuid ON public.user_sessions USING btree (user_uuid);

ALTER TABLE
    public.user_sessions
ADD
    CONSTRAINT fk_users_tokens FOREIGN KEY (user_uuid) REFERENCES public.users("uuid");

-- +migrate Down
DROP TABLE public.user_sessions;