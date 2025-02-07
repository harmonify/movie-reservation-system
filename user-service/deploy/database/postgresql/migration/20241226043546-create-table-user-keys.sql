-- +migrate Up
CREATE TABLE IF NOT EXISTS public.user_keys (
    user_uuid uuid NOT NULL,
    public_key text NOT NULL,
    private_key text NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE NULL,
    CONSTRAINT user_keys_pkey PRIMARY KEY (user_uuid),
    CONSTRAINT fk_users_tokens FOREIGN KEY (user_uuid) REFERENCES public.users("uuid")
);

CREATE INDEX idx_user_keys_deleted_at ON public.user_keys USING btree (deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS public.user_keys;