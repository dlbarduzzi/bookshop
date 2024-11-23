CREATE TABLE IF NOT EXISTS books (
  id bigserial PRIMARY KEY,
  title text NOT NULL,
  authors text[] NOT NULL,
  published_date DATE NOT NULL,
  page_count integer NOT NULL,
  categories text[] NOT NULL,
  version integer NOT NULL DEFAULT 1,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
)
