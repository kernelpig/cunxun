-- 评论列表视图
CREATE VIEW `commentlistview` AS (
  SELECT `comment`.*, `user`.nickname FROM `comment`, `user` WHERE `comment`.creater_uid = `user`.id
);