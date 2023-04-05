DROP TABLE IF EXISTS `t_task`;
CREATE TABLE `t_task`
(
    `id`      bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `num`     bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '开始block number',
    `name`    varchar(256) NOT NULL DEFAULT '' COMMENT 'task name'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;