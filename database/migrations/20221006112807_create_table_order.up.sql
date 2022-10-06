CREATE TABLE transaction  (
  id int(11) NOT NULL AUTO_INCREMENT,
  transactionNumber varchar(50) NOT NULL DEFAULT '',
  deliveryAddress varchar(255) NULL DEFAULT NULL,
  totalTransaction double NOT NULL DEFAULT 0,
  createdAt datetime(0) NULL DEFAULT NULL,
  totalQty int(11) NULL DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE = InnoDB;