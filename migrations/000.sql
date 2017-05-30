CREATE TABLE installations (
  id              SERIAL PRIMARY KEY,
  username        VARCHAR(128) NOT NULL UNIQUE,
  installation_id INTEGER      NOT NULL UNIQUE,

  created_at      TIMESTAMP    NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMP    NOT NULL DEFAULT NOW()
);
