CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS user_account (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "created_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "deactivated_at" TIMESTAMP(0) WITH TIME ZONE,
  "first_name" TEXT,
  "last_name" TEXT,
  "birth_date" TIMESTAMP(0) WITH TIME ZONE,
  "email" citext UNIQUE NOT NULL,
  "password" TEXT NOT NULL,
  "roles" TEXT [ ] NOT NULL DEFAULT '{}',
  "profil_name" TEXT NOT NULL,
  "short_id" TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS user_session (
  "id" uuid PRIMARY KEY,
  "created_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "deactivated_at" TIMESTAMP(0) WITH TIME ZONE NOT NULL,
  "ip" TEXT NOT NULL,
  "name" TEXT NOT NULL,
  "location" TEXT NOT NULL,
  "user_id" uuid NOT NULL REFERENCES user_account ON DELETE CASCADE
);