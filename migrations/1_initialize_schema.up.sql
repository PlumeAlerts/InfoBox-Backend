CREATE TABLE users
(
  id               varchar(255),
  infobox_interval int,
  PRIMARY KEY (id)
);

CREATE TABLE info_boxes
(
  id               int unsigned AUTO_INCREMENT,
  title            varchar(255),
  text_size        int,
  url              varchar(255),
  icon             varchar(255),
  icon_color       varchar(255),
  text_color       varchar(255),
  background_color varchar(255),
  user_id          varchar(255),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users (id)
);
