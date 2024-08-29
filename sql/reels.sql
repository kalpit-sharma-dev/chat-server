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


//insert for reels

INSERT INTO reels (user_id, video_url, created_at) 
VALUES 
(1, 'https://file-examples.com/storage/fe45dfa76e66c6232a111c9/2017/04/file_example_MP4_480_1_5MG.mp4', NOW()),
(20, 'https://samplelib.com/lib/preview/mp4/sample-5s.mp4', NOW())


