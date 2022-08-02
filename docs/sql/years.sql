CREATE TABLE `records` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `address` char(42) DEFAULT NULL,
  `pay_tx_hash` char(66) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `message` varchar(255) DEFAULT NULL,
  `signature` varchar(255) DEFAULT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0:未处理 1: 已经验证,待上链 2：已经上链 3：验证不通过',
  `block_number` bigint DEFAULT NULL,
  `c_time` bigint DEFAULT NULL,
  `update_time` bigint DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;