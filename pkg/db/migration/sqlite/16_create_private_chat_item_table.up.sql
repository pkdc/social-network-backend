CREATE TABLE private_chat_item (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  source_id INTEGER NOT NULL,
  target_id INTEGER NOT NULL,
  chat_noti INTEGER NOT NULL,
  last_msg_at DATETIME NOT NULL,
  FOREIGN KEY (source_id) REFERENCES user(id),
  FOREIGN KEY (target_id) REFERENCES user(id)
);