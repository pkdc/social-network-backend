CREATE TABLE group_chat_item (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  group_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  chat_noti INTEGER NOT NULL,
  last_msg_at DATETIME NOT NULL,
  FOREIGN KEY (group_id) REFERENCES group_(id),
  FOREIGN KEY (user_id) REFERENCES user(id)
);