-- +migrate Up
CREATE TABLE IF NOT EXISTS public.user_keys (
    user_uuid uuid NULL,
    public_key text NULL,
    private_key text NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL
);

CREATE INDEX idx_user_keys_deleted_at ON public.user_keys USING btree (deleted_at);

ALTER TABLE
    public.user_keys
ADD
    CONSTRAINT fk_users_tokens FOREIGN KEY (user_uuid) REFERENCES public.users("uuid");

-- +migrate Down
DROP TABLE IF EXISTS public.user_keys;