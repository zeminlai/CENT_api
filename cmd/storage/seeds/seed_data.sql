-- Clear existing data (if any)
TRUNCATE users, directories, notes CASCADE;

-- Reset sequences
ALTER SEQUENCE users_id_seq RESTART WITH 1;
ALTER SEQUENCE directories_id_seq RESTART WITH 1;
ALTER SEQUENCE notes_id_seq RESTART WITH 1;

-- Seed Users
INSERT INTO users (name, email, password, created_at, updated_at) VALUES
('John Doe', 'john@example.com', 'password123', NOW(), NOW()),
('Jane Smith', 'jane@example.com', 'password456', NOW(), NOW()),
('Bob Wilson', 'bob@example.com', 'password789', NOW(), NOW());

-- Seed Directories
INSERT INTO directories (name, user_id, parent_id, created_at, updated_at) VALUES
-- John's directories
('Work', 1, NULL, NOW(), NOW()),
('Personal', 1, NULL, NOW(), NOW()),
('Projects', 1, 1, NOW(), NOW()),  -- Sub-directory under Work

-- Jane's directories
('Study', 2, NULL, NOW(), NOW()),
('Research', 2, NULL, NOW(), NOW()),
('Literature Review', 2, 5, NOW(), NOW()),  -- Sub-directory under Research

-- Bob's directories
('Recipes', 3, NULL, NOW(), NOW()),
('Travel Plans', 3, NULL, NOW(), NOW());

-- Seed Notes
INSERT INTO notes (title, content, user_id, directory_id, created_at, updated_at) VALUES
-- John's notes
('Meeting Minutes', 'Discussion points from team meeting...', 1, 1, NOW(), NOW()),
('Project Ideas', 'List of potential new projects...', 1, 3, NOW(), NOW()),
('Shopping List', 'Groceries and supplies needed...', 1, 2, NOW(), NOW()),

-- Jane's notes
('Research Methodology', 'Steps for conducting the research...', 2, 5, NOW(), NOW()),
('Study Schedule', 'Weekly study plan and goals...', 2, 4, NOW(), NOW()),
('Literature Summary', 'Key findings from papers...', 2, 6, NOW(), NOW()),

-- Bob's notes
('Italian Pasta Recipe', 'Ingredients and cooking instructions...', 3, 7, NOW(), NOW()),
('Japan Trip Planning', 'Itinerary and accommodation details...', 3, 8, NOW(), NOW()),
('Meal Prep Ideas', 'Weekly meal planning suggestions...', 3, 7, NOW(), NOW()); 