CREATE TABLE `user` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT "id",
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "创建时间",
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "更新时间",
  `phone` VARCHAR(64) NOT NULL COMMENT "手机号码",
  `nickname` VARCHAR(64) NOT NULL DEFAULT "" COMMENT "用户昵称",
  `hashed_password` VARCHAR(64) NOT NULL COMMENT "密码",
  `password_level` INT NOT NULL DEFAULT 0 COMMENT "密码级别",
  `register_source` VARCHAR(256) NOT NULL DEFAULT "" COMMENT "注册来源",
  `avatar` VARCHAR(256) NOT NULL DEFAULT "" COMMENT "用户头像",
  `role` INT NOT NULL DEFAULT 0 COMMENT "用户角色, 普通用户/管理员/超级管理员",
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_phone` (`phone`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;