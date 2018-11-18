-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE sources(
  id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
  _created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  _updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  name TEXT NOT NULL,
  type VARCHAR(64) NOT NULL DEFAULT 'prometheus',
  config jsonb
);

CREATE TRIGGER set_timestamp
  BEFORE UPDATE ON sources
  FOR EACH ROW
  EXECUTE PROCEDURE trigger_set_timestamp();

-- +migrate Down
-- SQL in section 'Down' is executed when this migration is applied
DROP TABLE sources;