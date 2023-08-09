CREATE TABLE group_request (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  group_id INTEGER NOT NULL,
  status_ TEXT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES user(id),
  FOREIGN KEY (group_id) REFERENCES group_(id)
);
