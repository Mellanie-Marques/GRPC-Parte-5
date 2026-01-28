CREATE DATABASE IF NOT EXISTS `order_db`;
CREATE DATABASE IF NOT EXISTS `payment`;

USE `order_db`;

-- Criar tabela de produtos
CREATE TABLE IF NOT EXISTS `products` (
  `id` bigint unsigned AUTO_INCREMENT,
  `product_code` varchar(100) NOT NULL UNIQUE,
  `name` varchar(255),
  `quantity` int DEFAULT 0,
  `created_at` datetime(3),
  `updated_at` datetime(3),
  `deleted_at` datetime(3),
  PRIMARY KEY (`id`),
  INDEX `idx_products_deleted_at` (`deleted_at`),
  INDEX `idx_products_product_code` (`product_code`)
);

-- Inserir produtos de teste
INSERT INTO products (product_code, name, quantity, created_at) VALUES
('prod1', 'Produto 1', 100, NOW()),
('prod2', 'Produto 2', 50, NOW()),
('prod3', 'Produto 3', 75, NOW()),
('prod4', 'Produto 4', 200, NOW()),
('ABC123', 'Item ABC', 150, NOW()),
('XYZ789', 'Item XYZ', 80, NOW());
