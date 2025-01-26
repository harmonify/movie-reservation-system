-- +migrate Up
CREATE TABLE public.user_sessions (
    id bigserial NOT NULL,
    trace_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    user_uuid uuid NOT NULL,
    refresh_token text NOT NULL,
    is_revoked bool NOT NULL DEFAULT false,
    expired_at TIMESTAMP WITH TIME ZONE NOT NULL,
    ip_address text NOT NULL,
    user_agent text NOT NULL,
    CONSTRAINT user_sessions_pkey PRIMARY KEY (id),
    CONSTRAINT fk_users_tokens FOREIGN KEY (user_uuid) REFERENCES public.users("uuid")
);

CREATE INDEX idx_user_sessions_deleted_at ON public.user_sessions USING btree (deleted_at);

CREATE INDEX idx_user_sessions_user_uuid ON public.user_sessions USING btree (user_uuid);

-- +migrate Down
DROP TABLE public.user_sessions;