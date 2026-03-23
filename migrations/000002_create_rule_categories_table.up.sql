CREATE TABLE IF NOT EXISTS rule_categories (
                                               id SERIAL PRIMARY KEY,
                                               name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_categories_name ON rule_categories(name);

INSERT INTO rule_categories (name) VALUES
                                       ('Security Rules'),
                                       ('Performance Rules'),
                                       ('Code Quality Rules'),
                                       ('Documentation Rules');