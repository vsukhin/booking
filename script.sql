CREATE TABLE `flights` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL DEFAULT '',
  `created_at` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `name` (`name`),
  KEY `created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `blocks` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `flight_id` int(11) NOT NULL,
  `rows` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`flight_id`) REFERENCES `flights`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `seat_numbers` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `block_id` int(11) NOT NULL,
  `type` int(11) NOT NULL,
  `number` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`block_id`) REFERENCES `blocks`(`id`),
  KEY `type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `seats` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `flight_id` int(11) NOT NULL,
  `index` int(11) NOT NULL,
  `type` int(11) NOT NULL,
  `row` int(11) NOT NULL, 
  `line` CHAR NOT NULL,
  `assigned` BOOLEAN NOT NULL DEFAULT FALSE,  
  `created_at` int(11) NOT NULL,
  `updated_at` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`flight_id`) REFERENCES `flights`(`id`),
  KEY `index` (`index`),
  KEY `type` (`type`),
  KEY `row` (`row`),  
  KEY `line` (`line`), 
  KEY `created_at` (`created_at`),
  KEY `updated_at` (`updated_at`)    
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
