CREATE TABLE group_event (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  author INTEGER NOT NULL,
  group_id INTEGER NOT NULL,
  title TEXT NOT NULL,
  description_ TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  date_ DATETIME NOT NULL,
  FOREIGN KEY (author) REFERENCES user(id),
  FOREIGN KEY (group_id) REFERENCES group_(id)
);