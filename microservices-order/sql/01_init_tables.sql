-- 创建订单表
CREATE TABLE IF NOT EXISTS `xx_order` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `order_id` bigint NOT NULL COMMENT '订单ID',
    `user_id` bigint NOT NULL COMMENT '用户ID',
    `pay_amount` bigint NOT NULL DEFAULT '0' COMMENT '支付金额（分）',
    `status` varchar(50) NOT NULL DEFAULT 'pending' COMMENT '订单状态',
    `receive_address` varchar(255) NOT NULL DEFAULT '' COMMENT '收货地址',
    `receive_name` varchar(100) NOT NULL DEFAULT '' COMMENT '收货人姓名',
    `receive_phone` varchar(20) NOT NULL DEFAULT '' COMMENT '收货人电话',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_order_id` (`order_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单表';

-- 创建订单详情表
CREATE TABLE IF NOT EXISTS `xx_order_detail` (
    `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `order_id` bigint NOT NULL COMMENT '订单ID',
    `user_id` bigint NOT NULL COMMENT '用户ID',
    `goods_id` bigint NOT NULL COMMENT '商品ID',
    `num` bigint NOT NULL COMMENT '商品数量',
    `title` varchar(255) NOT NULL COMMENT '商品标题',
    `status` varchar(50) NOT NULL DEFAULT 'pending' COMMENT '详情状态',
    `price` bigint NOT NULL COMMENT '商品价格（分）',
    `brief` text COMMENT '商品简介',
    `pay_amount` bigint NOT NULL COMMENT '支付金额（分）',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_order_id` (`order_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_goods_id` (`goods_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单详情表';




