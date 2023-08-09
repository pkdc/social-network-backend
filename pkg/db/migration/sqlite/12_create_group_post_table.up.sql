CREATE TABLE group_post (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  author INTEGER NOT NULL,
  group_id INTEGER NOT NULL,
  message_ TEXT NOT NULL,
  image_ TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  FOREIGN KEY (author) REFERENCES user(id),
  FOREIGN KEY (group_id) REFERENCES group_(id)
);
