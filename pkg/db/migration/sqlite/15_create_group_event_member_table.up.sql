CREATE TABLE group_event_member (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  event_id INTEGER NOT NULL,
  status_ INTEGER NOT NULL DEFAULT 0,
  FOREIGN KEY (user_id) REFERENCES user(id),
  FOREIGN KEY (event_id) REFERENCES group_event(id)
);
