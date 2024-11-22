CREATE TABLE IF NOT EXISTS guests (
  id bigserial PRIMARY KEY,
  ip inet NOT NULL,
  message text NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
)
