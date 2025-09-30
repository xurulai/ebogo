-- 插入测试订单数据
INSERT INTO `xx_order` (`order_id`, `user_id`, `pay_amount`, `status`, `receive_address`, `receive_name`, `receive_phone`) VALUES
(1001, 1, 89900, 'pending', '北京市朝阳区某某街道123号', '张三', '13800138000'),
(1002, 1, 159900, 'paid', '上海市浦东新区某某路456号', '李四', '13900139000'),
(1003, 2, 299900, 'shipped', '广州市天河区某某大道789号', '王五', '13700137000');

-- 插入测试订单详情数据
INSERT INTO `xx_order_detail` (`order_id`, `user_id`, `goods_id`, `num`, `title`, `status`, `price`, `brief`, `pay_amount`) VALUES
(1001, 1, 1, 1, '商品A', 'pending', 89900, '这是商品A的简介', 89900),
(1002, 1, 2, 1, '商品B', 'paid', 159900, '这是商品B的简介', 159900),
(1003, 2, 3, 1, '商品C', 'shipped', 299900, '这是商品C的简介', 299900);




