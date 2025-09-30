-- 插入测试库存数据
INSERT INTO `xx_stock` (`goods_id`, `stocknum`, `lock`) VALUES
(1, 100, 0),
(2, 200, 0),
(3, 50, 0),
(4, 300, 0),
(5, 150, 0)
ON DUPLICATE KEY UPDATE 
  `stocknum` = VALUES(`stocknum`),
  `lock` = VALUES(`lock`);




