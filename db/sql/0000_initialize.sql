-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE EXTENSION IF NOT EXISTS "hstore";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW._updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +migrate StatementEnd

-- +migrate Down
-- SQL in section 'Up' is executed when this migration is applied
