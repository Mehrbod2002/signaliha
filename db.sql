CREATE database arz;

CREATE USER 'signal'@'localhost' identified by 'DtrPuxeHW6wWQ#g^';

GRANT all privileges on *.* to 'signal'@'localhost' ;

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
    tp VARCHAR(255),
    sl VARCHAR(255),
    timestamp INT,
    `exit` BOOL,
    risk BOOL
);

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    `limit` INT DEFAULT -1,
    token VARCHAR(255) NOT NULL,
    last_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE history (
  id INT AUTO_INCREMENT PRIMARY KEY,
  token VARCHAR(255),
  time DATETIME,
  result VARCHAR(255),
  request VARCHAR(255)
);

CREATE TABLE admins (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL
);

INSERT into admins (username,password) Values ("admin","UEyqG!bAg$C");