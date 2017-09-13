-- 评论列表视图
CREATE VIEW `articledetailview` AS (
  SELECT `article`.*, `user`.nickname FROM `article`, `user` WHERE `article`.creater_uid = `user`.id
);