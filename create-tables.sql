CREATE TABLE IF NOT EXISTS switches (
  id      INT AUTO_INCREMENT NOT NULL,
  app     VARCHAR(128) NOT NULL,
  tab     VARCHAR(128),
  `when`  DATETIME NOT NULL,
  PRIMARY KEY (`id`)
);
