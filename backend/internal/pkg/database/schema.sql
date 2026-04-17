-- ============================================================
-- FAMILY HUB — PostgreSQL Schema
-- v1.2.0 — Full audit trail
-- ============================================================

-- ============================================================
-- AUTH MODULE
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id          BIGSERIAL    PRIMARY KEY,
    name        TEXT         NOT NULL,
    email       TEXT         NOT NULL UNIQUE,
    password    TEXT         NOT NULL,
    role        TEXT         NOT NULL DEFAULT 'member' CHECK(role IN ('super_admin','member')),
    avatar_url  TEXT,

    -- audit
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_by  BIGINT       REFERENCES users(id) ON DELETE SET NULL,
    version     INTEGER      NOT NULL DEFAULT 1,
    is_deleted  BOOLEAN      NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_users_email      ON users(email)      WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_users_is_deleted ON users(is_deleted);

-- ============================================================
-- SILSILAH MODULE
-- ============================================================
CREATE TABLE IF NOT EXISTS family_members (
    id            BIGSERIAL    PRIMARY KEY,
    user_id       BIGINT       REFERENCES users(id) ON DELETE SET NULL,
    full_name     TEXT         NOT NULL,
    nickname      TEXT,
    gender        TEXT         NOT NULL CHECK(gender IN ('male','female')),
    birth_date    DATE,
    birth_place   TEXT,
    death_date    DATE,
    photo_url     TEXT,
    parent_id     BIGINT       REFERENCES family_members(id) ON DELETE SET NULL,
    notes         TEXT,

    -- audit
    created_by    BIGINT       REFERENCES users(id) ON DELETE SET NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_by    BIGINT       REFERENCES users(id) ON DELETE SET NULL,
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    version       INTEGER      NOT NULL DEFAULT 1,
    is_deleted    BOOLEAN      NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_family_members_parent     ON family_members(parent_id)  WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_family_members_user       ON family_members(user_id)    WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_family_members_is_deleted ON family_members(is_deleted);

CREATE TABLE IF NOT EXISTS marriages (
    id            BIGSERIAL    PRIMARY KEY,
    husband_id    BIGINT       NOT NULL REFERENCES family_members(id) ON DELETE CASCADE,
    wife_id       BIGINT       NOT NULL REFERENCES family_members(id) ON DELETE CASCADE,
    marriage_date DATE,
    divorce_date  DATE,
    notes         TEXT,
    UNIQUE(husband_id, wife_id),

    -- audit
    created_by    BIGINT       REFERENCES users(id) ON DELETE SET NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_by    BIGINT       REFERENCES users(id) ON DELETE SET NULL,
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    version       INTEGER      NOT NULL DEFAULT 1,
    is_deleted    BOOLEAN      NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_marriages_husband    ON marriages(husband_id) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_marriages_wife       ON marriages(wife_id)    WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_marriages_is_deleted ON marriages(is_deleted);

-- ============================================================
-- KAS KELUARGA MODULE
-- ============================================================
CREATE TABLE IF NOT EXISTS kas_categories (
    id    BIGSERIAL PRIMARY KEY,
    name  TEXT      NOT NULL UNIQUE,
    icon  TEXT
);

CREATE TABLE IF NOT EXISTS kas_transactions (
    id           BIGSERIAL    PRIMARY KEY,
    category_id  BIGINT       REFERENCES kas_categories(id),
    type         TEXT         NOT NULL CHECK(type IN ('income','expense')),
    amount       NUMERIC(15,2) NOT NULL CHECK(amount > 0),
    description  TEXT,
    date         DATE         NOT NULL DEFAULT CURRENT_DATE,

    -- audit
    created_by   BIGINT       NOT NULL REFERENCES users(id),
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_by   BIGINT       REFERENCES users(id) ON DELETE SET NULL,
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    version      INTEGER      NOT NULL DEFAULT 1,
    is_deleted   BOOLEAN      NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_kas_transactions_date       ON kas_transactions(date)       WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_kas_transactions_type       ON kas_transactions(type)       WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_kas_transactions_is_deleted ON kas_transactions(is_deleted);

-- ============================================================
-- PETA DOMISILI MODULE
-- ============================================================
CREATE TABLE IF NOT EXISTS addresses (
    id               BIGSERIAL    PRIMARY KEY,
    family_member_id BIGINT       REFERENCES family_members(id) ON DELETE CASCADE,
    label            TEXT         NOT NULL DEFAULT 'Rumah',
    street           TEXT,
    city             TEXT         NOT NULL,
    province         TEXT,
    postal_code      TEXT,
    country          TEXT         NOT NULL DEFAULT 'Indonesia',
    latitude         DOUBLE PRECISION,
    longitude        DOUBLE PRECISION,
    is_current       BOOLEAN      NOT NULL DEFAULT TRUE,

    -- audit
    created_by       BIGINT       REFERENCES users(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_by       BIGINT       REFERENCES users(id) ON DELETE SET NULL,
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    version          INTEGER      NOT NULL DEFAULT 1,
    is_deleted       BOOLEAN      NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_addresses_member     ON addresses(family_member_id) WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_addresses_is_deleted ON addresses(is_deleted);

-- ============================================================
-- KALENDER ACARA MODULE
-- ============================================================
CREATE TABLE IF NOT EXISTS events (
    id           BIGSERIAL    PRIMARY KEY,
    title        TEXT         NOT NULL,
    description  TEXT,
    location     TEXT,
    start_at     TIMESTAMPTZ  NOT NULL,
    end_at       TIMESTAMPTZ,
    is_recurring BOOLEAN      NOT NULL DEFAULT FALSE,
    recur_rule   TEXT,
    color        TEXT         NOT NULL DEFAULT '#CC6649',

    -- audit
    created_by   BIGINT       NOT NULL REFERENCES users(id),
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_by   BIGINT       REFERENCES users(id) ON DELETE SET NULL,
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    version      INTEGER      NOT NULL DEFAULT 1,
    is_deleted   BOOLEAN      NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS event_attendees (
    event_id         BIGINT  NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    family_member_id BIGINT  NOT NULL REFERENCES family_members(id) ON DELETE CASCADE,
    rsvp             TEXT    DEFAULT 'pending' CHECK(rsvp IN ('pending','yes','no')),
    PRIMARY KEY(event_id, family_member_id)
);

CREATE INDEX IF NOT EXISTS idx_events_start      ON events(start_at)   WHERE is_deleted = FALSE;
CREATE INDEX IF NOT EXISTS idx_events_is_deleted ON events(is_deleted);

-- ============================================================
-- AUDIT LOG TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS audit_log (
    id           BIGSERIAL    PRIMARY KEY,
    table_name   TEXT         NOT NULL,
    record_id    BIGINT       NOT NULL,
    action       TEXT         NOT NULL CHECK(action IN ('create','update','delete','restore')),
    changed_by   BIGINT       REFERENCES users(id) ON DELETE SET NULL,
    changed_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    diff         JSONB
);

CREATE INDEX IF NOT EXISTS idx_audit_log_table_record ON audit_log(table_name, record_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_changed_by   ON audit_log(changed_by);
CREATE INDEX IF NOT EXISTS idx_audit_log_changed_at   ON audit_log(changed_at);


-- ============================================================
-- REFRESH TOKENS
-- Stored server-side so tokens can be revoked on logout.
-- One active token per user at a time (rotation invalidates old).
-- ============================================================
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id         BIGSERIAL    PRIMARY KEY,
    user_id    BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT         NOT NULL UNIQUE,  -- SHA-256 of the raw token
    expires_at TIMESTAMPTZ  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ                    -- NULL = still valid
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id    ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);

-- ============================================================
-- SEED DATA
-- ============================================================
INSERT INTO kas_categories(name, icon) VALUES
    ('Iuran Bulanan',   'wallet'),
    ('Konsumsi Acara',  'utensils'),
    ('Sumbangan Sosial','heart'),
    ('Hadiah',          'gift'),
    ('Lain-lain',       'more-horizontal')
ON CONFLICT (name) DO NOTHING;
