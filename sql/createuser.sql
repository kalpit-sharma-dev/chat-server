-- Create the 'users' table
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    verification_code VARCHAR(50),
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    full_name VARCHAR(255),
    profile_picture_url VARCHAR(255),
    status_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Index for quick lookup by phone_number
CREATE INDEX idx_phone_number ON users (phone_number);
