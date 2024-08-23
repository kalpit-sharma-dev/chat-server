CREATE TABLE reels (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    video_url VARCHAR(255),
    created_at DATETIME
);

CREATE TABLE likes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    reel_id INT
);

CREATE TABLE comments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    reel_id INT,
    content TEXT,
    created_at DATETIME
);
