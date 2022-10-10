CREATE TABLE product  (
  id int(11) NOT NULL AUTO_INCREMENT,
  title varchar(100) NOT NULL,
  description varchar(255) NULL DEFAULT NULL,
  brandId int(11) NOT NULL,
  price double(10, 2) NOT NULL DEFAULT 0,
  createdAt datetime(0) NOT NULL ON UPDATE CURRENT_TIMESTAMP(0),
  updatedAt datetime(0) NOT NULL ON UPDATE CURRENT_TIMESTAMP(0),
  PRIMARY KEY (id)
) ENGINE = InnoDB;