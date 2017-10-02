-- 评论列表视图
CREATE VIEW `columndetailview` AS (
  select a.*, ifnull(u.nickname, '') as nickname, ifnull(count(a.id),0) as column_count from `column` a left join user u on a.creater_uid = u.id group by a.id
);