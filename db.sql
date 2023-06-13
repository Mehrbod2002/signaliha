CREATE TABLE messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    message_id INT,
    coin VARCHAR(50),
    base_currency VARCHAR(50),
    platform VARCHAR(50),
    leverage VARCHAR(50),
    side VARCHAR(50),
    entries VARCHAR(255),
    margin VARCHAR(255),
    sl VARCHAR(255),
    timestamp INT,
    exit BOOL,
    risk BOOL
);

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    token VARCHAR(255) NOT NULL,
    limit INT DEFAULT -1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE history (
  id INT AUTO_INCREMENT PRIMARY KEY,
  token VARCHAR(255),
  time DATETIME,
  result VARCHAR(255),
  request VARCHAR(255)
);
