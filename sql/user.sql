CREATE TABLE `user` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT "id",
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "创建时间",
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "更新时间",
  `phone` VARCHAR(64) NOT NULL COMMENT "手机号码",
  `name` VARCHAR(64) NOT NULL COMMENT "用户登录名",
  `nickname` VARCHAR(64) NOT NULL COMMENT "用户昵称",
  `hashed_password` VARCHAR(64) NOT NULL COMMENT "密码",
  `password_level` INT NOT NULL COMMENT "密码级别",
  `avatar` VARCHAR(256) NOT NULL COMMENT "头像",
  `register_source` VARCHAR(256) NOT NULL COMMENT "注册来源",
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_phone` (`phone`),
  UNIQUE KEY `uk_name` (`name`)
);