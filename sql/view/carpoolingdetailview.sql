-- 评论列表视图
CREATE VIEW `carpoolingdetailview` AS (
  select c.*, IFNULL(u.nickname, "") as nickname from carpooling c left join user u on c.creater_uid = u.id
);