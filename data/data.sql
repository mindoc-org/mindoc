INSERT INTO md_options (option_title, option_name, option_value) SELECT '是否启用注册','ENABLED_REGISTER','false' WHERE NOT exists(SELECT * FROM md_options WHERE option_name = 'ENABLED_REGISTER');
INSERT INTO md_options (option_title, option_name, option_value) SELECT '是否启用文档历史','ENABLE_DOCUMENT_HISTORY','true' WHERE NOT exists(SELECT * FROM md_options WHERE option_name = 'ENABLE_DOCUMENT_HISTORY');
INSERT INTO md_options (option_title, option_name, option_value) SELECT '是否启用验证码','ENABLED_CAPTCHA','false' WHERE NOT exists(SELECT * FROM md_options WHERE option_name = 'ENABLED_CAPTCHA');
INSERT INTO md_options (option_title, option_name, option_value) SELECT '启用匿名访问','ENABLE_ANONYMOUS','true' WHERE NOT exists(SELECT * FROM md_options WHERE option_name = 'ENABLE_ANONYMOUS');
INSERT INTO md_options (option_title, option_name, option_value) SELECT '站点名称','SITE_NAME','MinDoc' WHERE NOT exists(SELECT * FROM md_options WHERE option_name = 'SITE_NAME');












