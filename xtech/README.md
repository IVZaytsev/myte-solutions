## Описание
Приложение для выдачи данных в указанном формате.
Данные сохраняются в БД каждые 30 секунд.
Также есть docker-compose.yml файл для быстрой установки.

## Зависимости
- GoLang ( + github.com/go-sql-driver/mysql )
- MySQL

## Запуск приложения
```
$ docker compose up
```
Приложение будет доступно по адресу: **[http://localhost:8080/](http://localhost:8080/)**
## База данных
```
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
```