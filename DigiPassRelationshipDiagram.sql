CREATE TABLE "accounts" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "name" varchar(255) NOT NULL,
  "email" varchar(255) UNIQUE NOT NULL,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "status" varchar(50) DEFAULT 'active'
);

CREATE TABLE "tenants" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "account_id" uuid NOT NULL,
  "name" varchar(255) NOT NULL,
  "slug" varchar(255) UNIQUE NOT NULL,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "status" varchar(50) DEFAULT 'active',
  "settings" jsonb
);

CREATE TABLE "clients" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "client_id" varchar(255) UNIQUE NOT NULL,
  "client_secret_hash" varchar(255) NOT NULL,
  "tenant_id" uuid NOT NULL,
  "name" varchar(255) NOT NULL,
  "description" text,
  "redirect_uris" text[] NOT NULL,
  "grant_types" varchar(50) NOT NULL,
  "response_types" varchar(50) NOT NULL,
  "scopes" varchar(100) NOT NULL,
  "is_confidential" boolean DEFAULT true,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "status" varchar(50) DEFAULT 'active'
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "tenant_id" uuid NOT NULL,
  "email" varchar(255) NOT NULL,
  "email_verified" boolean DEFAULT false,
  "password_hash" varchar(255),
  "given_name" varchar(255),
  "family_name" varchar(255),
  "picture_url" text,
  "locale" varchar(10),
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "last_login_at" timestamp,
  "status" varchar(50) DEFAULT 'active'
);

CREATE TABLE "authorization_codes" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "code" varchar(255) UNIQUE NOT NULL,
  "client_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "redirect_uri" text NOT NULL,
  "scopes" varchar(100) NOT NULL,
  "code_challenge" varchar(255),
  "code_challenge_method" varchar(10),
  "expires_at" timestamp NOT NULL,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "used_at" timestamp,
  "nonce" varchar(255)
);

CREATE TABLE "access_tokens" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "token_hash" varchar(255) UNIQUE NOT NULL,
  "client_id" uuid NOT NULL,
  "user_id" uuid,
  "scopes" varchar(100) NOT NULL,
  "expires_at" timestamp NOT NULL,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "revoked_at" timestamp,
  "session_id" uuid
);

CREATE TABLE "refresh_tokens" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "token_hash" varchar(255) UNIQUE NOT NULL,
  "access_token_id" uuid,
  "client_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "scopes" varchar(100) NOT NULL,
  "expires_at" timestamp NOT NULL,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "revoked_at" timestamp,
  "session_id" uuid
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user_id" uuid NOT NULL,
  "client_id" uuid,
  "user_agent" text,
  "ip_address" inet,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "last_activity_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "expires_at" timestamp NOT NULL,
  "revoked_at" timestamp
);

CREATE TABLE "user_consents" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user_id" uuid NOT NULL,
  "client_id" uuid NOT NULL,
  "scopes" varchar(100) NOT NULL,
  "granted_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "revoked_at" timestamp
);

CREATE TABLE "account_users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "account_id" uuid NOT NULL,
  "email" varchar(255) NOT NULL,
  "password_hash" varchar(255) NOT NULL,
  "role" varchar(50) NOT NULL,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" timestamp DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "audit_logs" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "tenant_id" uuid,
  "user_id" uuid,
  "account_user_id" uuid,
  "action" varchar(100) NOT NULL,
  "resource_type" varchar(50),
  "resource_id" uuid,
  "ip_address" inet,
  "user_agent" text,
  "metadata" jsonb,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "id_tokens" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "token_hash" varchar(255) UNIQUE NOT NULL,
  "authorization_code_id" uuid,
  "client_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "nonce" varchar(255),
  "acr" varchar(50),
  "amr" text,
  "azp" varchar(255),
  "at_hash" varchar(255),
  "c_hash" varchar(255),
  "issued_at" timestamp NOT NULL,
  "expires_at" timestamp NOT NULL,
  "created_at" timestamp DEFAULT (CURRENT_TIMESTAMP),
  "session_id" uuid
);

CREATE INDEX ON "accounts" ("email");

CREATE INDEX ON "tenants" ("account_id");

CREATE INDEX ON "tenants" ("slug");

CREATE UNIQUE INDEX ON "tenants" ("account_id", "name");

CREATE UNIQUE INDEX ON "clients" ("client_id");

CREATE INDEX ON "clients" ("tenant_id");

CREATE UNIQUE INDEX ON "users" ("tenant_id", "email");

CREATE INDEX ON "users" ("tenant_id");

CREATE UNIQUE INDEX ON "authorization_codes" ("code");

CREATE INDEX ON "authorization_codes" ("expires_at");

CREATE INDEX ON "authorization_codes" ("client_id");

CREATE INDEX ON "authorization_codes" ("user_id");

CREATE UNIQUE INDEX ON "access_tokens" ("token_hash");

CREATE INDEX ON "access_tokens" ("expires_at");

CREATE INDEX ON "access_tokens" ("user_id");

CREATE INDEX ON "access_tokens" ("client_id");

CREATE INDEX ON "access_tokens" ("session_id");

CREATE UNIQUE INDEX ON "refresh_tokens" ("token_hash");

CREATE INDEX ON "refresh_tokens" ("user_id");

CREATE INDEX ON "refresh_tokens" ("client_id");

CREATE INDEX ON "refresh_tokens" ("session_id");

CREATE INDEX ON "sessions" ("user_id");

CREATE INDEX ON "sessions" ("expires_at");

CREATE INDEX ON "sessions" ("client_id");

CREATE UNIQUE INDEX ON "user_consents" ("user_id", "client_id");

CREATE INDEX ON "user_consents" ("user_id");

CREATE INDEX ON "user_consents" ("client_id");

CREATE UNIQUE INDEX ON "account_users" ("account_id", "email");

CREATE INDEX ON "account_users" ("account_id");

CREATE INDEX ON "audit_logs" ("tenant_id", "created_at");

CREATE INDEX ON "audit_logs" ("user_id", "created_at");

CREATE INDEX ON "audit_logs" ("account_user_id");

CREATE UNIQUE INDEX ON "id_tokens" ("token_hash");

CREATE INDEX ON "id_tokens" ("user_id");

CREATE INDEX ON "id_tokens" ("client_id");

CREATE INDEX ON "id_tokens" ("authorization_code_id");

CREATE INDEX ON "id_tokens" ("session_id");

CREATE INDEX ON "id_tokens" ("expires_at");

COMMENT ON COLUMN "accounts"."status" IS 'active, suspended, deleted';

COMMENT ON COLUMN "tenants"."slug" IS 'tenant identifier in URLs';

COMMENT ON COLUMN "tenants"."settings" IS 'tenant-specific configuration';

COMMENT ON COLUMN "clients"."client_id" IS 'OAuth client_id';

COMMENT ON COLUMN "clients"."client_secret_hash" IS 'hashed secret';

COMMENT ON COLUMN "clients"."redirect_uris" IS 'allowed redirect URIs';

COMMENT ON COLUMN "clients"."grant_types" IS 'authorization_code, refresh_token, etc.';

COMMENT ON COLUMN "clients"."response_types" IS 'code, token, id_token';

COMMENT ON COLUMN "clients"."scopes" IS 'allowed scopes';

COMMENT ON COLUMN "users"."password_hash" IS 'null for social logins';

COMMENT ON COLUMN "authorization_codes"."code_challenge" IS 'PKCE code challenge';

COMMENT ON COLUMN "authorization_codes"."code_challenge_method" IS 'S256 or plain';

COMMENT ON COLUMN "authorization_codes"."used_at" IS 'prevent code reuse';

COMMENT ON COLUMN "authorization_codes"."nonce" IS 'OpenID Connect nonce';

COMMENT ON COLUMN "access_tokens"."token_hash" IS 'hash of actual token';

COMMENT ON COLUMN "access_tokens"."user_id" IS 'null for client_credentials';

COMMENT ON COLUMN "account_users"."role" IS 'owner, admin, member';

COMMENT ON COLUMN "audit_logs"."action" IS 'login, logout, token_issued, etc.';

COMMENT ON COLUMN "audit_logs"."resource_type" IS 'client, user, token';

COMMENT ON COLUMN "id_tokens"."token_hash" IS 'hash of actual JWT';

COMMENT ON COLUMN "id_tokens"."authorization_code_id" IS 'links to auth code that generated this';

COMMENT ON COLUMN "id_tokens"."nonce" IS 'matches nonce from auth request';

COMMENT ON COLUMN "id_tokens"."acr" IS 'Authentication Context Class Reference';

COMMENT ON COLUMN "id_tokens"."amr" IS 'Authentication Methods References';

COMMENT ON COLUMN "id_tokens"."azp" IS 'Authorized party - client_id';

COMMENT ON COLUMN "id_tokens"."at_hash" IS 'Access Token hash for validation';

COMMENT ON COLUMN "id_tokens"."c_hash" IS 'Code hash for validation';

ALTER TABLE "tenants" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "clients" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "authorization_codes" ADD FOREIGN KEY ("client_id") REFERENCES "clients" ("id");

ALTER TABLE "authorization_codes" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "access_tokens" ADD FOREIGN KEY ("client_id") REFERENCES "clients" ("id");

ALTER TABLE "access_tokens" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "access_tokens" ADD FOREIGN KEY ("session_id") REFERENCES "sessions" ("id");

ALTER TABLE "refresh_tokens" ADD FOREIGN KEY ("access_token_id") REFERENCES "access_tokens" ("id");

ALTER TABLE "refresh_tokens" ADD FOREIGN KEY ("client_id") REFERENCES "clients" ("id");

ALTER TABLE "refresh_tokens" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "refresh_tokens" ADD FOREIGN KEY ("session_id") REFERENCES "sessions" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("client_id") REFERENCES "clients" ("id");

ALTER TABLE "user_consents" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_consents" ADD FOREIGN KEY ("client_id") REFERENCES "clients" ("id");

ALTER TABLE "account_users" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "audit_logs" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "audit_logs" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "audit_logs" ADD FOREIGN KEY ("account_user_id") REFERENCES "account_users" ("id");

ALTER TABLE "id_tokens" ADD FOREIGN KEY ("authorization_code_id") REFERENCES "authorization_codes" ("id");

ALTER TABLE "id_tokens" ADD FOREIGN KEY ("client_id") REFERENCES "clients" ("id");

ALTER TABLE "id_tokens" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "id_tokens" ADD FOREIGN KEY ("session_id") REFERENCES "sessions" ("id");
