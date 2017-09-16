-- 评论列表视图
CREATE VIEW `articledetailview` AS (
  select a.*, IFNULL(u.nickname, "") as nickname from article a left join user u on a.creater_uid = u.id
);