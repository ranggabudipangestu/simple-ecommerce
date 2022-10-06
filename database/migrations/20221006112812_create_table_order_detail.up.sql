CREATE TABLE transaction_detail  (
  id int(11) NOT NULL AUTO_INCREMENT,
  transactionId int(11) NOT NULL DEFAULT 0,
  productId int(11) NOT NULL,
  qty int(11) NOT NULL,
  price double NOT NULL,
  total double NOT NULL,
  PRIMARY KEY (id)
) ENGINE = InnoDB;