-- Sample data for Blog example

-- Users
INSERT INTO users (id, username, email, bio, created_at) VALUES
    ('00000000-0000-0000-0000-000000000001', 'alice', 'alice@blog.com', 'Tech blogger and developer', NOW() - INTERVAL '30 days'),
    ('00000000-0000-0000-0000-000000000002', 'bob', 'bob@blog.com', 'Photography enthusiast', NOW() - INTERVAL '20 days'),
    ('00000000-0000-0000-0000-000000000003', 'charlie', 'charlie@blog.com', NULL, NOW() - INTERVAL '10 days');

-- Posts
INSERT INTO posts (id, title, content, published, author_id, created_at) VALUES
    -- Alice's posts
    ('10000000-0000-0000-0000-000000000001', 
     'Getting Started with ChameleonDB', 
     'ChameleonDB makes database access type-safe and intuitive...', 
     true, 
     '00000000-0000-0000-0000-000000000001', 
     NOW() - INTERVAL '5 days'),
    
    ('10000000-0000-0000-0000-000000000002', 
     'Graph Navigation vs SQL Joins', 
     'Traditional ORMs force you to think in terms of JOINs...', 
     true, 
     '00000000-0000-0000-0000-000000000001', 
     NOW() - INTERVAL '3 days'),
    
    ('10000000-0000-0000-0000-000000000003', 
     'Draft: Upcoming Features', 
     'This post is not yet published...', 
     false, 
     '00000000-0000-0000-0000-000000000001', 
     NOW() - INTERVAL '1 day'),
    
    -- Bob's posts
    ('10000000-0000-0000-0000-000000000004', 
     'Best Cameras for 2026', 
     'After testing 15 different cameras...', 
     true, 
     '00000000-0000-0000-0000-000000000002', 
     NOW() - INTERVAL '4 days');

-- Comments
INSERT INTO comments (id, content, author_id, post_id, created_at) VALUES
    -- Comments on Alice's first post
    ('20000000-0000-0000-0000-000000000001', 
     'Great introduction! Looking forward to trying it out.', 
     '00000000-0000-0000-0000-000000000002', 
     '10000000-0000-0000-0000-000000000001', 
     NOW() - INTERVAL '4 days'),
    
    ('20000000-0000-0000-0000-000000000002', 
     'Thanks Bob! Let me know if you have questions.', 
     '00000000-0000-0000-0000-000000000001', 
     '10000000-0000-0000-0000-000000000001', 
     NOW() - INTERVAL '4 days'),
    
    -- Comments on Alice's second post
    ('20000000-0000-0000-0000-000000000003', 
     'This is exactly what I was looking for!', 
     '00000000-0000-0000-0000-000000000003', 
     '10000000-0000-0000-0000-000000000002', 
     NOW() - INTERVAL '2 days'),
    
    -- Comments on Bob's post
    ('20000000-0000-0000-0000-000000000004', 
     'Which one would you recommend for beginners?', 
     '00000000-0000-0000-0000-000000000001', 
     '10000000-0000-0000-0000-000000000004', 
     NOW() - INTERVAL '3 days'),
    
    ('20000000-0000-0000-0000-000000000005', 
     'For beginners I recommend the Canon EOS R10', 
     '00000000-0000-0000-0000-000000000002', 
     '10000000-0000-0000-0000-000000000004', 
     NOW() - INTERVAL '3 days');

-- Verify data
SELECT 
    u.username,
    COUNT(DISTINCT p.id) as post_count,
    COUNT(DISTINCT c.id) as comment_count
FROM users u
LEFT JOIN posts p ON p.author_id = u.id
LEFT JOIN comments c ON c.author_id = u.id
GROUP BY u.username
ORDER BY u.username;
