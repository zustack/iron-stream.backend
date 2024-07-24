DROP TABLE IF EXISTS users;
CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  name VARCHAR(55) NOT NULL,
  surname VARCHAR(55) NOT NULL,
  username VARCHAR(55) NOT NULL UNIQUE,
  email VARCHAR(55) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  is_admin BOOL,
  is_active BOOL DEFAULT TRUE,
  email_token INT,
  verified   BOOL DEFAULT FALSE, 
  courses    TEXT DEFAULT '',  
  pc         VARCHAR(255) DEFAULT '',  
  os VARCHAR(20) DEFAULT '',  
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
