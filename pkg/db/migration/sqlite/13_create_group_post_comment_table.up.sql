CREATE TABLE group_post_comment (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  author INTEGER NOT NULL,
  group_post_id INTEGER NOT NULL,
  message_ TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  FOREIGN KEY (author) REFERENCES user(id),
  FOREIGN KEY (group_post_id) REFERENCES group_post(id)
);
