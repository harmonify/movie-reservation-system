
-- +migrate Up
ALTER TABLE public.users
ADD COLUMN trace_id UUID NOT NULL;

-- +migrate Down
ALTER TABLE public.users
DROP COLUMN IF EXISTS trace_id;
