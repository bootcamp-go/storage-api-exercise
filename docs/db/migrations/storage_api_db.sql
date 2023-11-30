-- DDL: Data Definition Language
DROP DATABASE IF EXISTS `storage_api_db`;

CREATE DATABASE `storage_api_db`;

USE `storage_api_db`;

CREATE TABLE `warehouses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `address` varchar(150) NOT NULL,
  `telephone` varchar(150) NOT NULL,
  `capacity` int NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `products` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `quantity` int NOT NULL,
  `code_value` varchar(255) NOT NULL,
  `is_published` boolean NOT NULL,
  `expiration` date NOT NULL,
  `price` decimal(10, 2) NOT NULL,
  `warehouse_id` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_products_warehouse_id` (`warehouse_id`),
  CONSTRAINT `fk_products_warehouse_id` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);