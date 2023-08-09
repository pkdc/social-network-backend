CREATE TABLE user_message (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  source_id INTEGER NOT NULL,
  target_id INTEGER NOT NULL,
  message_ TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  FOREIGN KEY (source_id) REFERENCES user(id),
  FOREIGN KEY (target_id) REFERENCES user(id)
);
