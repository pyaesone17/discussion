-- Sample seed data for testing

INSERT INTO users (username, email) VALUES
('alice', 'alice@example.com'),
('bob', 'bob@example.com'),
('charlie', 'charlie@example.com');

INSERT INTO topics (title, user_id) VALUES
('Welcome to the Discussion Forum', 1),
('Getting Started with Go', 2),
('Best Practices for REST APIs', 1);

INSERT INTO posts (topic_id, user_id, content) VALUES
(1, 1, 'Welcome everyone! This is our first discussion topic.'),
(1, 2, 'Thanks for creating this forum!'),
(2, 2, 'I''m learning Go and would love to share resources.'),
(2, 3, 'Check out the official Go documentation at golang.org'),
(3, 1, 'Let''s discuss API design patterns and best practices.'),
(3, 2, 'I recommend following RESTful principles and using proper HTTP methods.');
