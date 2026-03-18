CREATE TABLE IF NOT EXISTS words (
    id              SERIAL PRIMARY KEY,
    word            VARCHAR(100) NOT NULL UNIQUE,
    part_of_speech  VARCHAR(20) NOT NULL,
    definition      TEXT NOT NULL,
    example         TEXT NOT NULL,
    difficulty      VARCHAR(10) NOT NULL DEFAULT 'medium',
    assigned_date   DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_words_assigned_date ON words (assigned_date);
CREATE INDEX IF NOT EXISTS idx_words_difficulty ON words (difficulty);