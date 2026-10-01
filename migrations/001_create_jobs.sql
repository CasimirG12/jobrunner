BEGIN TRANSACTION;

CREATE TYPE job_status AS ENUM (
  'pending',
  'running',
  'completed',
  'failed'
);


CREATE TABLE jobs (
  id UUID PRIMARY KEY,
  kind TEXT NOT NULL,
  status job_status NOT NULL,
  input TEXT NOT NULL,
  attempts INTEGER NOT NULL DEFAULT 0,
  result TEXT NOT NULL DEFAULT '',
  error_message TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  started_at TIMESTAMPTZ,
  ended_at TIMESTAMPTZ,
  CONSTRAINT jobs_kind_check CHECK (kind IN ('word_count', 'calculate_hash'))
);

CREATE INDEX jobs_status_created_at_idx
ON jobs(created_at)
WHERE status = 'pending';

COMMIT;
