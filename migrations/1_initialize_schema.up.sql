CREATE TABLE users
(
  id                  VARCHAR(255),
  last_triggered      TIMESTAMP,
  annotation_interval INT DEFAULT 15,
  last_annotation_id  INT,

  PRIMARY KEY (id)
);

CREATE TABLE annotations
(
  id               SERIAL,
  text             VARCHAR(255) NOT NULL,
  text_size        INT          NOT NULL,
  url              VARCHAR(255),
  icon             VARCHAR(255),
  icon_color       VARCHAR(255) DEFAULT '#FFF',
  text_color       VARCHAR(255) DEFAULT '#FFF',
  background_color VARCHAR(255) DEFAULT '#FFF',

  user_id          VARCHAR(255),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users (id)
);