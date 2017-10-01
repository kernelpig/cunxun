-- 评论列表视图
CREATE VIEW `articlelistview` AS (
  select a.*, ifnull(u.nickname, '') as nickname, ifnull(count(c.relate_id),0) as comment_count from article a left join comment c on a.id = c.relate_id left join user u on a.creater_uid = u.id group by a.id
);