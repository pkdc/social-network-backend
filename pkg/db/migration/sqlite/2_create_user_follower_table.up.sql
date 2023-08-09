CREATE TABLE user_follower (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  source_id INTEGER NOT NULL,
  target_id INTEGER NOT NULL,
  status_ INTEGER NOT NULL DEFAULT 0,
  FOREIGN KEY (source_id) REFERENCES user(id),
  FOREIGN KEY (target_id) REFERENCES user(id)
);