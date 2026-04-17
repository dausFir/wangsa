-- ============================================================
-- MIGRATION: Add audit fields to existing database
-- Run this ONCE on any existing database that was created
-- before v1.2.0. New databases use schema.sql directly.
-- ============================================================

-- users
ALTER TABLE users ADD COLUMN updated_by INTEGER REFERENCES users(id);
ALTER TABLE users ADD COLUMN version    INTEGER NOT NULL DEFAULT 1;
ALTER TABLE users ADD COLUMN is_deleted INTEGER NOT NULL DEFAULT 0 CHECK(is_deleted IN (0,1));

-- family_members
ALTER TABLE family_members ADD COLUMN updated_by INTEGER REFERENCES users(id);
ALTER TABLE family_members ADD COLUMN version    INTEGER NOT NULL DEFAULT 1;
ALTER TABLE family_members ADD COLUMN is_deleted INTEGER NOT NULL DEFAULT 0 CHECK(is_deleted IN (0,1));

-- marriages
ALTER TABLE marriages ADD COLUMN created_by INTEGER REFERENCES users(id);
ALTER TABLE marriages ADD COLUMN updated_by INTEGER REFERENCES users(id);
ALTER TABLE marriages ADD COLUMN updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE marriages ADD COLUMN version    INTEGER  NOT NULL DEFAULT 1;
ALTER TABLE marriages ADD COLUMN is_deleted INTEGER  NOT NULL DEFAULT 0 CHECK(is_deleted IN (0,1));

-- kas_transactions
ALTER TABLE kas_transactions ADD COLUMN updated_by INTEGER REFERENCES users(id);
ALTER TABLE kas_transactions ADD COLUMN version    INTEGER NOT NULL DEFAULT 1;
ALTER TABLE kas_transactions ADD COLUMN is_deleted INTEGER NOT NULL DEFAULT 0 CHECK(is_deleted IN (0,1));

-- addresses
ALTER TABLE addresses ADD COLUMN created_by INTEGER REFERENCES users(id);
ALTER TABLE addresses ADD COLUMN updated_by INTEGER REFERENCES users(id);
ALTER TABLE addresses ADD COLUMN updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE addresses ADD COLUMN version    INTEGER  NOT NULL DEFAULT 1;
ALTER TABLE addresses ADD COLUMN is_deleted INTEGER  NOT NULL DEFAULT 0 CHECK(is_deleted IN (0,1));

-- events
ALTER TABLE events ADD COLUMN updated_by INTEGER REFERENCES users(id);
ALTER TABLE events ADD COLUMN version    INTEGER NOT NULL DEFAULT 1;
ALTER TABLE events ADD COLUMN is_deleted INTEGER NOT NULL DEFAULT 0 CHECK(is_deleted IN (0,1));

-- audit_log (new table)
CREATE TABLE IF NOT EXISTS audit_log (
    id           INTEGER  PRIMARY KEY AUTOINCREMENT,
    table_name   TEXT     NOT NULL,
    record_id    INTEGER  NOT NULL,
    action       TEXT     NOT NULL CHECK(action IN ('create','update','delete','restore')),
    changed_by   INTEGER  REFERENCES users(id) ON DELETE SET NULL,
    changed_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    diff         TEXT
);
CREATE INDEX IF NOT EXISTS idx_audit_log_table_record ON audit_log(table_name, record_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_changed_by   ON audit_log(changed_by);
CREATE INDEX IF NOT EXISTS idx_audit_log_changed_at   ON audit_log(changed_at);

-- Partial indexes for is_deleted (run after columns exist)
CREATE INDEX IF NOT EXISTS idx_users_is_deleted           ON users(is_deleted);
CREATE INDEX IF NOT EXISTS idx_family_members_is_deleted  ON family_members(is_deleted);
CREATE INDEX IF NOT EXISTS idx_marriages_is_deleted       ON marriages(is_deleted);
CREATE INDEX IF NOT EXISTS idx_kas_transactions_is_deleted ON kas_transactions(is_deleted);
CREATE INDEX IF NOT EXISTS idx_addresses_is_deleted       ON addresses(is_deleted);
CREATE INDEX IF NOT EXISTS idx_events_is_deleted          ON events(is_deleted);
