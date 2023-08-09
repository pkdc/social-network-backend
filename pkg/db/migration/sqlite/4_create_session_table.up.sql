CREATE TABLE session_table (
  session_token TEXT PRIMARY KEY,
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES user(id)
);