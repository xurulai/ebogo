SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `xx_goods_query` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `goods_id` BIGINT NOT NULL UNIQUE,
  `category_id` BIGINT NOT NULL,
  `brand_name` VARCHAR(64) NOT NULL,
  `code` VARCHAR(64) NOT NULL UNIQUE,
  `status` TINYINT NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `market_price` BIGINT NOT NULL,
  `price` BIGINT NOT NULL,
  `brief` TEXT,
  `created_at` DATETIME NULL,
  `updated_at` DATETIME NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `xx_room_goods` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `room_id` BIGINT NOT NULL,
  `goods_id` BIGINT NOT NULL,
  `weight` INT NOT NULL DEFAULT 0,
  `is_current` TINYINT NOT NULL DEFAULT 0,
  `created_at` DATETIME NULL,
  `updated_at` DATETIME NULL,
  INDEX `idx_room` (`room_id`),
  INDEX `idx_goods` (`goods_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `xx_stock` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `goods_id` BIGINT NOT NULL UNIQUE,
  `stocknum` BIGINT NOT NULL DEFAULT 0,
  `lock` BIGINT NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `xx_stock_record` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `order_id` BIGINT NOT NULL,
  `goods_id` BIGINT NOT NULL,
  `num` BIGINT NOT NULL,
  `status` INT NOT NULL DEFAULT 1,
  INDEX `idx_order` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `xx_order` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `order_id` BIGINT NOT NULL UNIQUE,
  `user_id` BIGINT NOT NULL,
  `pay_amount` BIGINT NOT NULL DEFAULT 0,
  `status` VARCHAR(32) NOT NULL DEFAULT 'pending',
  `receive_address` VARCHAR(255) NOT NULL,
  `receive_name` VARCHAR(64) NOT NULL,
  `receive_phone` VARCHAR(32) NOT NULL,
  `created_at` DATETIME NULL,
  `updated_at` DATETIME NULL,
  INDEX `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `xx_order_detail` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `order_id` BIGINT NOT NULL,
  `user_id` BIGINT NOT NULL,
  `goods_id` BIGINT NOT NULL,
  `num` BIGINT NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `status` VARCHAR(32) NOT NULL DEFAULT 'pending',
  `price` BIGINT NOT NULL,
  `brief` TEXT,
  `pay_amount` BIGINT NOT NULL,
  `created_at` DATETIME NULL,
  `updated_at` DATETIME NULL,
  INDEX `idx_order` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Seed data
INSERT INTO `xx_goods_query`
  (`goods_id`, `category_id`, `brand_name`, `code`, `status`, `title`, `market_price`, `price`, `brief`, `created_at`, `updated_at`)
VALUES
  (1, 1, 'DemoBrand', 'SKU-0001', 1, '演示商品', 9999, 8999, '演示用商品', NOW(), NOW())
ON DUPLICATE KEY UPDATE `title`=VALUES(`title`), `price`=VALUES(`price`), `updated_at`=NOW();

INSERT INTO `xx_room_goods` (`room_id`, `goods_id`, `weight`, `is_current`, `created_at`, `updated_at`)
VALUES (1, 1, 0, 1, NOW(), NOW());

INSERT INTO `xx_stock` (`goods_id`, `stocknum`, `lock`)
VALUES (1, 100, 0)
ON DUPLICATE KEY UPDATE `stocknum`=VALUES(`stocknum`), `lock`=VALUES(`lock`);





