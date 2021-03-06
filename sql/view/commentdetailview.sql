-- 评论列表视图
CREATE VIEW `commentdetailview` AS (
  select c.*, IFNULL(u.nickname, "") as nickname, IFNULL(u.avatar, "") as avatar from comment c left join user u on c.creater_uid = u.id order by c.updated_at asc
);