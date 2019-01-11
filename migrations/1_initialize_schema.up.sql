CREATE TABLE users
(
  id                 VARCHAR(255),
  intervals          INT,
  last_triggered     TIMESTAMP,
  last_info_boxes_id INT UNSIGNED,

  PRIMARY KEY (id)
);

CREATE TABLE info_boxes
(
  id               INT UNSIGNED AUTO_INCREMENT,
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