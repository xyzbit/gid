CREATE DATABASE `gid` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

CREATE TABLE gid.schedule (
	id integer auto_increment NOT NULL,
	name VARCHAR(255) NOT NULL COMMENT '调度模块名称',
	last_id BIGINT DEFAULT 0 NOT NULL COMMENT '调度状态-最后生成id',
	CONSTRAINT schedule_PK PRIMARY KEY (id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_0900_ai_ci
COMMENT='gid 调度记录';
CREATE UNIQUE INDEX idx_name USING BTREE ON gid.schedule (name);
