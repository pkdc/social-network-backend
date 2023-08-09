CREATE TABLE group_message (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  source_id INTEGER NOT NULL,
  group_id INTEGER NOT NULL,
  message_ TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  FOREIGN KEY (source_id) REFERENCES user(id),
  FOREIGN KEY (group_id) REFERENCES group_(id)
);
