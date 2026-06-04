CREATE TABLE IF NOT EXISTS body_metrics (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    weight_kg       NUMERIC(5,2),
    body_fat_percent NUMERIC(4,1),
    chest_cm        NUMERIC(5,1),
    waist_cm        NUMERIC(5,1),
    hips_cm         NUMERIC(5,1),
    bicep_cm        NUMERIC(4,1),
    measured_at     DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_body_metrics_user_id ON body_metrics(user_id, measured_at);
