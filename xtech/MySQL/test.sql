USE test_db;

DROP TABLE IF EXISTS `records`;
DROP TABLE IF EXISTS `symbols`;

CREATE TABLE `symbols` (
	`id` INT PRIMARY KEY AUTO_INCREMENT,
	`symbol` VARCHAR(10) NOT NULL UNIQUE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `records` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `symbol` INT NOT NULL,
  `price_24h` float UNSIGNED DEFAULT NULL,
  `volume_24h` float UNSIGNED DEFAULT NULL,
  `last_trade_price` float UNSIGNED DEFAULT NULL,
  `stamp` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (symbol) REFERENCES symbols(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;