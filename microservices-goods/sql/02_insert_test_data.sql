-- 插入测试商品数据
INSERT INTO `xx_goods_query` (`goods_id`, `category_id`, `brand_name`, `code`, `status`, `title`, `market_price`, `price`, `brief`) VALUES
(1, 100, 'Apple', 'IPHONE15PRO', 1, 'iPhone 15 Pro 256GB 深空黑色', 999900, 899900, 'Apple iPhone 15 Pro，搭载A17 Pro芯片，钛金属设计，专业级摄像头系统'),
(2, 100, 'Apple', 'IPHONE15', 1, 'iPhone 15 128GB 粉色', 599900, 549900, 'Apple iPhone 15，搭载A16仿生芯片，4800万像素主摄，支持动态岛'),
(3, 101, 'Samsung', 'GALAXY-S24', 1, 'Samsung Galaxy S24 Ultra 512GB', 1199900, 1099900, 'Samsung Galaxy S24 Ultra，搭载Snapdragon 8 Gen 3，S Pen手写笔'),
(4, 102, 'Huawei', 'MATE60PRO', 1, '华为Mate 60 Pro 512GB 雅川青', 699900, 649900, '华为Mate 60 Pro，搭载麒麟9000S芯片，支持5G网络，徕卡影像'),
(5, 103, 'Xiaomi', 'MI14PRO', 1, '小米14 Pro 512GB 钛金属', 499900, 449900, '小米14 Pro，搭载骁龙8 Gen3芯片，徕卡专业摄影，120W快充')
ON DUPLICATE KEY UPDATE 
  `category_id` = VALUES(`category_id`),
  `brand_name` = VALUES(`brand_name`),
  `code` = VALUES(`code`),
  `status` = VALUES(`status`),
  `title` = VALUES(`title`),
  `market_price` = VALUES(`market_price`),
  `price` = VALUES(`price`),
  `brief` = VALUES(`brief`);

-- 插入测试直播间商品关联数据
INSERT INTO `xx_room_goods` (`room_id`, `goods_id`, `weight`, `is_current`) VALUES
(1, 1, 1, 1),  -- 直播间1，商品1，当前讲解
(1, 2, 2, 0),  -- 直播间1，商品2
(1, 3, 3, 0),  -- 直播间1，商品3
(2, 4, 1, 1),  -- 直播间2，商品4，当前讲解
(2, 5, 2, 0),  -- 直播间2，商品5
(3, 1, 1, 0),  -- 直播间3，商品1
(3, 3, 2, 0),  -- 直播间3，商品3
(3, 5, 3, 1)   -- 直播间3，商品5，当前讲解
ON DUPLICATE KEY UPDATE 
  `weight` = VALUES(`weight`),
  `is_current` = VALUES(`is_current`);