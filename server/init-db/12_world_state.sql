-- 12_world_state.sql — World Engine 世界状态持久化

CREATE TABLE IF NOT EXISTS world_state (
    id          SERIAL PRIMARY KEY,
    epoch       BIGINT NOT NULL DEFAULT 0,
    state_data  JSONB NOT NULL DEFAULT '{}',
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS region_resources (
    id              SERIAL PRIMARY KEY,
    region_id       VARCHAR(64) NOT NULL,
    resource_id     VARCHAR(64) NOT NULL,
    quantity        INT NOT NULL DEFAULT 0,
    max_quantity    INT NOT NULL DEFAULT 100,
    last_harvested  TIMESTAMPTZ,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(region_id, resource_id)
);

CREATE TABLE IF NOT EXISTS world_events (
    id              VARCHAR(64) PRIMARY KEY,
    event_type      VARCHAR(32) NOT NULL,
    region_id       VARCHAR(64) NOT NULL,
    intensity       DOUBLE PRECISION NOT NULL DEFAULT 1.0,
    start_time      TIMESTAMPTZ NOT NULL,
    end_time        TIMESTAMPTZ NOT NULL,
    status          VARCHAR(16) NOT NULL DEFAULT 'active'
);

CREATE TABLE IF NOT EXISTS event_cooldowns (
    id              SERIAL PRIMARY KEY,
    event_type      VARCHAR(32) NOT NULL,
    region_id       VARCHAR(64) NOT NULL,
    end_time        TIMESTAMPTZ NOT NULL,
    UNIQUE(event_type, region_id)
);

CREATE INDEX IF NOT EXISTS idx_region_resources_region ON region_resources(region_id);
CREATE INDEX IF NOT EXISTS idx_world_events_region ON world_events(region_id);
CREATE INDEX IF NOT EXISTS idx_world_events_status ON world_events(status);
CREATE INDEX IF NOT EXISTS idx_event_cooldowns_region ON event_cooldowns(region_id);
