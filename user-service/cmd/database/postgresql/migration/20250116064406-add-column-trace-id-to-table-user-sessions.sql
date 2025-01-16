
-- +migrate Up
ALTER TABLE public.user_sessions
ADD COLUMN trace_id UUID NOT NULL;

-- +migrate Down
ALTER TABLE public.user_sessions
DROP COLUMN IF EXISTS trace_id;
