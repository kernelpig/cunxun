CREATE TABLE `carpooling` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT "id",
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT "创建时间",
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT "更新时间",
  `from_city` VARCHAR(64) NOT NULL COMMENT "出发城市",
  `to_city` VARCHAR(64) NOT NULL COMMENT "到达城市",
  `depart_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT "发车时间",
  `people_count` INT NOT NULL COMMENT "座位数量",
  `contact` VARCHAR(64) NOT NULL COMMENT "联系方式",
  `status` INT NOT NULL COMMENT "状态",
  `remark` TEXT NOT NULL COMMENT "备注信息",
  `creater_uid` INT NOT NULL COMMENT "创建人id",
  `updater_uid` INT NOT NULL COMMENT "最近修改人id",
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;