-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE rule_sets(
  id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
  _created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  _updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  name TEXT NOT NULL,
  rules jsonb
);

CREATE TRIGGER set_timestamp
  BEFORE UPDATE ON rule_sets
  FOR EACH ROW
  EXECUTE PROCEDURE trigger_set_timestamp();

-- +migrate Down
-- SQL in section 'Down' is executed when this migration is applied
DROP TABLE rule_sets;