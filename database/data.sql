CREATE TABLE IF NOT EXISTS md_attachment
(
  `attachment_id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
  `book_id` INT DEFAULT '0' NOT NULL COMMENT '项目ID',
  `document_id` INT NULL COMMENT '文档ID',
  `file_name` VARCHAR(255)  NOT NULL COMMENT '文件名称',
  `file_path` VARCHAR(2000) NOT NULL COMMENT '文件本地',
  `file_size` FLOAT DEFAULT '0' NOT NULL  COMMENT '文件大小，单位字节',
  `http_path` VARCHAR(2000) NOT NULL  COMMENT '文件可访问的uri',
  `file_ext` VARCHAR(50) DEFAULT '' NOT NULL  COMMENT '文件扩展名',
  `create_time` DATETIME NOT NULL  COMMENT '创建时间',
  `create_at` INT DEFAULT '0' NOT NULL  COMMENT '创建人'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='附件表';

CREATE TABLE IF NOT EXISTS md_books
(
  `book_id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '项目主键',
  `book_name` VARCHAR(500) NOT NULL COMMENT '项目名称',
  `identify` VARCHAR(100) NOT NULL COMMENT '项目唯一标识',
  `order_index` INT DEFAULT '0' NOT NULL COMMENT '排序',
  `description` TEXT NULL COMMENT '项目描述',
  `label` VARCHAR(500) DEFAULT '' NOT NULL COMMENT '项目标签',
  `privately_owned` INT DEFAULT '0' NOT NULL COMMENT 'PrivatelyOwned 项目私有： 0 公开/ 1 私有',
  `private_token` VARCHAR(500) null COMMENT '当项目是私有时的访问Token',
  `status` INT DEFAULT '0' NOT NULL COMMENT '状态：0 正常/1 已删除',
  `editor` VARCHAR(50) DEFAULT '' NOT NULL COMMENT '编辑器类型：html 富文本/markdown ',
  `doc_count` INT DEFAULT '0' NOT NULL COMMENT '包含的文档数量',
  `comment_status` VARCHAR(20) DEFAULT 'open' NOT NULL COMMENT '评论开启状态',
  `comment_count` INT DEFAULT '0' NOT NULL COMMENT '评论数量',
  `cover` VARCHAR(1000) DEFAULT '' NOT NULL COMMENT '封面图片路径',
  `theme` VARCHAR(255) DEFAULT 'DEFAULT' NOT NULL COMMENT '主题风格',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `member_id` INT DEFAULT '0' NOT NULL COMMENT '创建人',
  `modify_time` datetime null COMMENT '修改时间',
  `modify_at` INT NULL COMMENT '修改人id',
  `version` BIGINT DEFAULT '0' NOT NULL COMMENT '当前版本',
  CONSTRAINT identify  UNIQUE (identify),
  KEY (`privately_owned`),
  KEY (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '项目表';


CREATE TABLE IF NOT EXISTS md_document_history
(
  `history_id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '历史记录表主键',
  `document_id` INT DEFAULT '0' NOT NULL COMMENT '文档ID',
  `document_name` VARCHAR(500) DEFAULT '' NOT NULL COMMENT '文档名称',
  `parent_id` INT DEFAULT '0' NOT NULL COMMENT '文档所属父文档ID',
  `markdown` LONGTEXT  NULL COMMENT 'markdown内容',
  `content` longtext NULL COMMENT 'html内容',
  `member_id` INT DEFAULT '0' NOT NULL COMMENT '创建的用户',
  `modify_time` datetime NOT NULL COMMENT '修改时间',
  `modify_at` INT DEFAULT '0' NOT NULL COMMENT '修改人',
  `version` BIGINT DEFAULT '0' NOT NULL COMMENT '当前版本',
  KEY (`document_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '文档历史表';

CREATE TABLE IF NOT EXISTS md_documents
(
  `document_id` INT AUTO_INCREMENT PRIMARY KEY  COMMENT '文档主键',
  `document_name` VARCHAR(500) DEFAULT '' NOT NULL COMMENT '文档名称',
  `identify` VARCHAR(100) DEFAULT 'null' null COMMENT '文档唯一标识',
  `book_id` INT DEFAULT '0' NOT NULL COMMENT '项目ID',
  `parent_id` INT DEFAULT '0' NOT NULL COMMENT '所属父文档ID',
  `order_sort` INT DEFAULT '0' NOT NULL COMMENT '排序',
  `markdown` LONGTEXT NULL COMMENT 'markdown内容',
  `release` LONGTEXT NULL COMMENT '当前发布的内容',
  `content` LONGTEXT NULL COMMENT 'html内容',
  `create_time` DATETIME NOT NULL COMMENT '创建时间',
  `member_id` INT DEFAULT '0' NOT NULL COMMENT '创建用户',
  `modify_time` DATETIME NOT NULL COMMENT '修改时间',
  `modify_at` INT DEFAULT '0' NOT NULL COMMENT '修改用户',
  `version` BIGINT DEFAULT '0' NOT NULL  COMMENT '当前版本',
  KEY (`member_id`),
  KEY (`identify`),
  KEY (`order_sort`),
  KEY (`parent_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '文档表';


CREATE TABLE IF NOT EXISTS md_logs
(
  `log_id` BIGINT AUTO_INCREMENT  PRIMARY KEY COMMENT '日志主键',
  `member_id` INT DEFAULT '0' NOT NULL COMMENT '产生日志的用户',
  `category` VARCHAR(255) DEFAULT 'operate' NOT NULL COMMENT '日志分类',
  `content` longtext NOT NULL COMMENT '日志内容',
  `original_data` longtext NOT NULL COMMENT '产生日志前的数据',
  `present_data` longtext NOT NULL COMMENT '产生日志后的数据',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `user_agent` VARCHAR(500) DEFAULT '' NOT NULL COMMENT '浏览器信息',
  `ip_address` VARCHAR(255) DEFAULT '' NOT NULL COMMENT '用户IP地址',
  KEY (`member_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '日志表';


CREATE TABLE IF NOT EXISTS md_member_token
(
  `token_id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
  `member_id` INT DEFAULT '0' NOT NULL COMMENT '用户ID',
  `token` VARCHAR(150) DEFAULT '' NOT NULL COMMENT 'token值',
  `email` VARCHAR(255) DEFAULT '' NOT NULL COMMENT '收件人邮箱',
  `is_valid` tinyINT(1) DEFAULT '0' NOT NULL COMMENT '是否已校验',
  `valid_time` datetime null COMMENT '校验时间',
  `send_time` datetime NOT NULL COMMENT '发送时间',
  KEY (`token`),
  KEY (`email`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '用户找回密码记录表';


CREATE TABLE IF NOT EXISTS md_members
(
  `member_id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
  `account` VARCHAR(100) DEFAULT '' NOT NULL COMMENT '账号',
  `password` VARCHAR(1000) DEFAULT '' NOT NULL COMMENT '密码',
  `description` VARCHAR(2000) DEFAULT '' NOT NULL COMMENT '描述',
  `email` VARCHAR(100) DEFAULT '' NOT NULL COMMENT '邮箱',
  `phone` VARCHAR(255) DEFAULT 'null' null COMMENT '手机号',
  `avatar` VARCHAR(1000) DEFAULT '' NOT NULL COMMENT '头像',
  `role` INT DEFAULT '1' NOT NULL COMMENT '用户角色：0 超级管理员 /1 管理员/ 2 普通用户 ',
  `status` INT DEFAULT '0' NOT NULL COMMENT '用户状态：0 正常/1 禁用',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `create_at` INT DEFAULT '0' NOT NULL COMMENT '创建人',
  `last_login_time` datetime null COMMENT '最后登录时间',
  `auth_method` VARCHAR(50) DEFAULT 'local' null COMMENT '认证方式: local 本地数据库 /ldap LDAP',
  CONSTRAINT account UNIQUE (account),
  CONSTRAINT `email`  UNIQUE (`email`),
  KEY (`role`),
  KEY (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '用户信息表';


CREATE TABLE IF NOT EXISTS md_options
(
  option_id INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
  option_title VARCHAR(500) DEFAULT '' NOT NULL COMMENT '选项名称',
  option_name VARCHAR(80) DEFAULT '' NOT NULL COMMENT '选项键',
  option_value longtext null COMMENT '选项值',
  remark longtext null COMMENT '备注',
  CONSTRAINT option_name UNIQUE (option_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '系统配置表';

CREATE TABLE IF NOT EXISTS md_relationship
(
  relationship_id INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
  member_id INT DEFAULT '0' NOT NULL COMMENT '用户ID',
  book_id INT DEFAULT '0' NOT NULL COMMENT '项目ID',
  role_id INT DEFAULT '0' NOT NULL COMMENT '角色：0 创始人(创始人不能被移除) / 1 管理员/2 编辑者/3 观察者',
  CONSTRAINT member_id UNIQUE (member_id, book_id),
  KEY (`role_id`),
  KEY (`member_id`),
  KEY (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '项目和用户关系表' ;

CREATE TABLE IF NOT EXISTS migrations
(
  id_migration INT(10) unsigned AUTO_INCREMENT comment 'surrogate key' PRIMARY KEY,
  name VARCHAR(255) null comment 'migration name, UNIQUE',
  created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL comment 'date migrated or rolled back',
  statements longtext null comment 'SQL statements for this migration',
  rollback_statements longtext null comment 'SQL statment for rolling back migration',
  status enum('update', 'rollback') null comment 'update indicates it is a normal migration while rollback means this migration is rolled back'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '迁移记录表';