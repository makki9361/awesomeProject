CREATE TABLE IF NOT EXISTS rules (
                                     id SERIAL PRIMARY KEY,
                                     title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    category_id INTEGER NOT NULL REFERENCES rule_categories(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('draft', 'published', 'archived')),
    version INTEGER NOT NULL DEFAULT 1,
    created_by INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_rules_category ON rules(category_id);
CREATE INDEX idx_rules_status ON rules(status);
CREATE INDEX idx_rules_created_by ON rules(created_by);
CREATE INDEX idx_rules_created_at ON rules(created_at);
CREATE INDEX idx_rules_title ON rules(title);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_rules_updated_at
    BEFORE UPDATE ON rules
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

INSERT INTO rules (title, content, category_id, status, created_by) VALUES
    ('Example Rule', 'This is an example rule content', 1, 'published', 1);