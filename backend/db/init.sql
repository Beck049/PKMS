-- Create database if not exists
CREATE DATABASE IF NOT EXISTS pkms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE pkms;

-- Create articles table
CREATE TABLE IF NOT EXISTS articles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL,
    create_date DATETIME NOT NULL,
    edit_date DATETIME NOT NULL,
    ref_count INT UNSIGNED DEFAULT 0,
    pin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_title (title),
    INDEX idx_path (path),
    INDEX idx_create_date (create_date),
    INDEX idx_edit_date (edit_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create tags table
CREATE TABLE IF NOT EXISTS tags (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create article_tags junction table
CREATE TABLE IF NOT EXISTS article_tags (
    article_id BIGINT UNSIGNED NOT NULL,
    tag_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (article_id, tag_id),
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    INDEX idx_article_id (article_id),
    INDEX idx_tag_id (tag_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create search_index table for Bleve
CREATE TABLE IF NOT EXISTS search_index (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    article_id BIGINT UNSIGNED NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    INDEX idx_article_id (article_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert initial tags from default.md
INSERT IGNORE INTO tags (id, name) VALUES 
    (1, '3C'),
    (2, 'Network'),
    (3, 'book'),
    (4, 'food'),
    (5, 'raw'),
    (6, 'sweet'),
    (7, 'snack'),
    (8, 'AI'),
    (9, 'Search'),
    (10, 'Introduction'),
    (11, 'Machine Learning'),
    (12, 'Deep Learning'),
    (13, 'Computer Vision'),
    (14, 'Natural Language Processing'),
    (15, 'Robotics');

-- Insert initial articles from default.md
INSERT IGNORE INTO articles (id, title, path, type, create_date, edit_date, ref_count, pin) VALUES 
    (1, 'WiFi', '3C/01.Wifi.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (2, 'NAS', '3C/02.NAS.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (3, 'USB', '3C/03.USB.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, true),
    (4, 'Story', 'Book/01.Story.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (5, 'Potato', 'Food/01.Potato.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (6, 'Fish', 'Food/02.Fish.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, true),
    (7, 'Pork', 'Food/03.Pork.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (8, 'Chocolate', 'Snack/01.Chocolate.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, true),
    (9, 'IceCream', 'Snack/02.IceCream.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, true),
    (10, 'Chips', 'Snack/03.Chips.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (11, 'Popcorn', 'Snack/04.Popcorn.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (12, 'Cookies', 'Snack/05.Cookies.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (13, 'Gum', 'Snack/06.Gum.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (14, 'Cheese', 'Snack/07.Cheese.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (15, 'Pocky', 'Snack/08.Pocky.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (16, 'Puffs', 'Snack/09.Puffs.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (17, 'Donuts', 'Snack/10.Donuts.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (18, 'Beef', 'Food/04.Beef.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (19, 'Chicken', 'Food/05.Chicken.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (20, 'Lamb', 'Food/06.Lamb.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (21, 'Duck', 'Food/07.Duck.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (22, 'Salmon', 'Food/08.Salmon.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (23, 'Deer', 'Food/09.Deer.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (24, 'Laptop', '3C/Computer/01.Laptop.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (25, 'CPU', '3C/Computer/02.CPU.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false),
    (26, 'RAM', '3C/Computer/03.RAM.md', 'markdown', '2024-03-14 00:00:00', '2024-03-14 00:00:00', 0, false);

-- Insert article_tags relationships from default.md
INSERT IGNORE INTO article_tags (article_id, tag_id) VALUES 
    (1, 1), -- WiFi -> 3C
    (1, 2), -- WiFi -> Network
    (2, 1), -- NAS -> 3C
    (2, 2), -- NAS -> Network
    (3, 1), -- USB -> 3C
    (4, 3), -- Story -> book
    (5, 4), -- Potato -> food
    (5, 6), -- Potato -> sweet
    (6, 4), -- Fish -> food
    (6, 5), -- Fish -> raw
    (7, 4), -- Pork -> food
    (7, 5), -- Pork -> raw
    (8, 4), -- Chocolate -> food
    (8, 6), -- Chocolate -> sweet
    (8, 7), -- Chocolate -> snack
    (9, 4), -- IceCream -> food
    (9, 6), -- IceCream -> sweet
    (9, 7), -- IceCream -> snack 
    (10, 4), -- Chips -> food
    (10, 6), -- Chips -> sweet
    (10, 7), -- Chips -> snack 
    (11, 4), -- Popcorn -> food
    (11, 6), -- Popcorn -> sweet
    (11, 7), -- Popcorn -> snack 
    (12, 4), -- Cookies -> food
    (12, 6), -- Cookies -> sweet
    (12, 7), -- Cookies -> snack 
    (13, 4), -- Gum -> food
    (13, 6), -- Gum -> sweet
    (13, 7), -- Gum -> snack 
    (14, 4), -- Cheese -> food
    (14, 6), -- Cheese -> sweet
    (14, 7), -- Cheese -> snack 
    (15, 4), -- Pocky -> food
    (15, 6), -- Pocky -> sweet
    (15, 7), -- Pocky -> snack 
    (16, 4), -- Puffs -> food
    (16, 6), -- Puffs -> sweet
    (16, 7), -- Puffs -> snack 
    (17, 4), -- Donuts -> food
    (17, 6), -- Donuts -> sweet
    (17, 7), -- Donuts -> snack 
    (18, 4), -- Beef -> food
    (18, 5), -- Beef -> raw
    (19, 4), -- Chicken -> food
    (19, 5), -- Chicken -> raw
    (20, 4), -- Lamb -> food
    (20, 5), -- Lamb -> raw
    (21, 4), -- Duck -> food
    (21, 5), -- Duck -> raw
    (22, 4), -- Salmon -> food
    (22, 5), -- Salmon -> raw
    (23, 4), -- Deer -> food
    (23, 5), -- Deer -> raw
    (24, 1), -- NAS -> 3C
    (25, 1), -- NAS -> 3C
    (26, 1); -- NAS -> 3C