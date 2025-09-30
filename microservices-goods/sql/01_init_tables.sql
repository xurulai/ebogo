-- 创建商品表
CREATE TABLE IF NOT EXISTS `xx_goods_query` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `goods_id` bigint NOT NULL COMMENT '商品ID',
  `category_id` bigint NOT NULL COMMENT '商品分类ID',
  `brand_name` varchar(255) NOT NULL COMMENT '品牌名称',
  `code` varchar(255) NOT NULL COMMENT '商品编码',
  `status` tinyint NOT NULL COMMENT '商品状态',
  `title` varchar(500) NOT NULL COMMENT '商品标题',
  `market_price` bigint NOT NULL COMMENT '市场价（单位：分）',
  `price` bigint NOT NULL COMMENT '销售价格（单位：分）',
  `brief` text COMMENT '商品简介',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_goods_id` (`goods_id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

-- 创建直播间商品关联表
CREATE TABLE IF NOT EXISTS `xx_room_goods` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `room_id` bigint NOT NULL COMMENT '直播间ID',
  `goods_id` bigint NOT NULL COMMENT '商品ID',
  `weight` int NOT NULL DEFAULT '0' COMMENT '权重（用于排序）',
  `is_current` tinyint NOT NULL DEFAULT '0' COMMENT '是否当前讲解商品：0-否，1-是',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_room_id` (`room_id`),
  KEY `idx_goods_id` (`goods_id`),
  KEY `idx_room_goods` (`room_id`, `goods_id`),
  KEY `idx_weight` (`weight`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='直播间商品关联表';