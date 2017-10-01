-- 评论列表视图
CREATE VIEW `articledetailview` AS (
  select a.*, IFNULL(u.nickname, "") as nickname, IFNULL(count(c.id), 0) as comment_count from article a left join user u on a.creater_uid = u.id left join comment c on a.id = c.relate_id group by a.id
);