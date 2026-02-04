-- Sample data for Hello World example

INSERT INTO users (id, name, email, created_at) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'Alice Smith', 'alice@example.com', NOW()),
    ('550e8400-e29b-41d4-a716-446655440002', 'Bob Johnson', 'bob@example.com', NOW()),
    ('550e8400-e29b-41d4-a716-446655440003', 'Charlie Brown', 'charlie@example.com', NOW());

-- Verify
SELECT * FROM users ORDER BY created_at;
