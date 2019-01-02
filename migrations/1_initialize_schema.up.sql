CREATE TABLE info_boxes
(
  id               int unsigned AUTO_INCREMENT,
  text             varchar(255),
  text_size        int,
  url              varchar(255),
  icon             varchar(255),
  icon_color       varchar(255),
  text_color       varchar(255),
  background_color varchar(255),
  intervals        int,

  user_id          varchar(255),
  PRIMARY KEY (id)
);
