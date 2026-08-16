-- 商城后台数据库初始化脚本 (PostgreSQL)
-- 优化说明：
-- 1. 去掉冗余自增 id 列，统一使用业务 *_id（雪花算法）作为主键
-- 2. is_default 由 varchar(1) 改为 boolean
-- 3. weight 由 varchar(255) 改为 weight_mg bigint，单位：毫克（最小重量单位，milligrams）
-- 4. spu.description / sku.sku_desc 由 varchar(255) 改为 text
-- 5. sku.price 保持 price_cent bigint，单位：分（最小货币单位，cents）
-- 6. 补全所有带 update_time 列的表的自动更新触发器
-- 7. 补全关联字段索引，避免 JOIN 全表扫描
-- 8. 补充表与关键字段 COMMENT
-- 9. 每张表 DROP 加 CASCADE，兼容外键重建（顺序不再受限）
-- 10. 核心关联加 FOREIGN KEY，并配 ON DELETE RESTRICT / CASCADE 策略
-- 11. 补全 user_role 的 update_time 触发器
-- 12. spu_image_list 索引由 image_name 改为 spu_id

-- 创建 schema（如不存在）
CREATE SCHEMA IF NOT EXISTS app;

-- =========================================================
-- Table structure for table attr
-- =========================================================
DROP TABLE IF EXISTS app.attr CASCADE;
CREATE TABLE app.attr (
  attr_id bigint NOT NULL,
  attr_name varchar(255) NOT NULL,
  category_id bigint NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (attr_id)
);
CREATE INDEX idx_attr_category_id ON app.attr (category_id);
COMMENT ON TABLE app.attr IS '平台属性（三级分类下的规格属性）';
-- 注：attr 初始数据由代码 seed 生成（雪花ID），见 seed 包。

-- =========================================================
-- Table structure for table attr_value
-- =========================================================
DROP TABLE IF EXISTS app.attr_value CASCADE;
CREATE TABLE app.attr_value (
  attr_value_id bigint NOT NULL,
  value_name varchar(255) NOT NULL,
  attr_id bigint NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (attr_value_id)
);
CREATE INDEX idx_attr_value_attr_id ON app.attr_value (attr_id);
COMMENT ON TABLE app.attr_value IS '平台属性值';
-- 注：attr_value 初始数据由代码 seed 生成（雪花ID），见 seed 包。

-- =========================================================
-- Table structure for table category1
-- =========================================================
DROP TABLE IF EXISTS app.category1 CASCADE;
CREATE TABLE app.category1 (
  category1_id bigint NOT NULL,
  name varchar(255) NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (category1_id),
  UNIQUE (name)
);
COMMENT ON TABLE app.category1 IS '一级分类';
--
-- Dumping data for table category1
--
INSERT INTO app.category1 (category1_id, name) VALUES (1,'图书、音像、电子书刊'),(2,'手机'),(3,'家用电器'),(4,'数码'),(5,'家居家装'),(6,'电脑办公'),(7,'厨具'),(8,'个护化妆'),(9,'服饰内衣'),(10,'钟表'),(11,'鞋靴'),(12,'母婴'),(13,'礼品箱包'),(14,'食品饮料、保健食品'),(15,'珠宝'),(16,'汽车用品'),(17,'运动健康');

-- =========================================================
-- Table structure for table category2
-- =========================================================
DROP TABLE IF EXISTS app.category2 CASCADE;
CREATE TABLE app.category2 (
  category2_id bigint NOT NULL,
  name varchar(255) NOT NULL,
  category1_id bigint NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (category2_id)
);
CREATE INDEX idx_category2_category1_id ON app.category2 (category1_id);
COMMENT ON TABLE app.category2 IS '二级分类';
--
-- Dumping data for table category2
--
INSERT INTO app.category2 (category2_id, name, category1_id) VALUES (1,'电子书刊',1),(2,'音像',1),(3,'英文原版',1),(4,'文艺',1),(5,'少儿',1),(6,'人文社科',1),(7,'经管励志',1),(8,'生活',1),(9,'科技',1),(10,'教育',1),(11,'港台图书',1),(12,'其他',1),(13,'手机通讯',2),(14,'运营商',2),(15,'手机配件',2),(16,'大 家 电',3),(17,'厨卫大电',3),(18,'厨房小电',3),(19,'生活电器',3),(20,'个护健康',3),(21,'五金家装',3),(22,'摄影摄像',4),(23,'数码配件',4),(24,'智能设备',4),(25,'影音娱乐',4),(26,'电子教育',4),(27,'虚拟商品',4),(28,'家纺',5),(29,'灯具',5),(30,'生活日用',5),(31,'家装软饰',5),(32,'宠物生活',5),(33,'电脑整机',6),(34,'电脑配件',6),(35,'外设产品',6),(36,'游戏设备',6),(37,'网络产品',6),(38,'办公设备',6),(39,'文具/耗材',6),(40,'服务产品',6),(41,'烹饪锅具',7),(42,'刀剪菜板',7),(43,'厨房配件',7),(44,'水具酒具',7),(45,'餐具',7),(46,'酒店用品',7),(47,'茶具/咖啡具',7),(48,'清洁用品',8),(49,'面部护肤',8),(50,'身体护理',8),(51,'口腔护理',8),(52,'女性护理',8),(53,'洗发护发',8),(54,'香水彩妆',8),(55,'女装',9),(56,'男装',9),(57,'内衣',9),(58,'洗衣服务',9),(59,'服饰配件',9),(60,'钟表',10),(61,'流行男鞋',11),(62,'时尚女鞋',11),(63,'奶粉',12),(64,'营养辅食',12),(65,'尿裤湿巾',12),(66,'喂养用品',12),(67,'洗护用品',12),(68,'童车童床',12),(69,'寝居服饰',12),(70,'妈妈专区',12),(71,'童装童鞋',12),(72,'安全座椅',12),(73,'潮流女包',13),(74,'精品男包',13),(75,'功能箱包',13),(76,'礼品',13),(77,'奢侈品',13),(78,'婚庆',13),(79,'进口食品',14),(80,'地方特产',14),(81,'休闲食品',14),(82,'粮油调味',14),(83,'饮料冲调',14),(84,'食品礼券',14),(85,'茗茶',14),(86,'时尚饰品',15),(87,'黄金',15),(88,'K金饰品',15),(89,'金银投资',15),(90,'银饰',15),(91,'钻石',15),(92,'翡翠玉石',15),(93,'水晶玛瑙',15),(94,'彩宝',15),(95,'铂金',15),(96,'木手串/把件',15),(97,'珍珠',15),(98,'维修保养',16),(99,'车载电器',16),(100,'美容清洗',16),(101,'汽车装饰',16),(102,'安全自驾',16),(103,'汽车服务',16),(104,'赛事改装',16),(105,'运动鞋包',17),(106,'运动服饰',17),(107,'骑行运动',17),(108,'垂钓用品',17),(109,'游泳用品',17),(110,'户外鞋服',17),(111,'户外装备',17),(112,'健身训练',17),(113,'体育用品',17);

-- =========================================================
-- Table structure for table category3
-- =========================================================
DROP TABLE IF EXISTS app.category3 CASCADE;
CREATE TABLE app.category3 (
  category3_id bigint NOT NULL,
  name varchar(255) NOT NULL,
  category2_id bigint NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (category3_id)
);
CREATE INDEX idx_category3_category2_id ON app.category3 (category2_id);
COMMENT ON TABLE app.category3 IS '三级分类';
--
-- Dumping data for table category3
--
INSERT INTO app.category3 (category3_id, name, category2_id) VALUES (61,'手机',13),(62,'游戏手机',13),(63,'老人机',13),(64,'对讲机',13),(65,'以旧换新',14),(66,'移动卡',14),(67,'联通卡',14),(68,'电信卡',14),(69,'宽带社区',14),(70,'手机配件',15),(71,'手机壳',15),(72,'贴膜',15),(73,'充电器/数据线',15),(74,'移动电源',15),(75,'蓝牙耳机',15),(76,'存储卡',15),(77,'有线电视',15),(78,'手机耳机',15),(79,'院校教育',10),(80,'玩具乐器',12),(81,'相机',22),(82,'数码相机',22),(83,'单反相机',22),(84,'单电/微单相机',22),(85,'拍立得',22),(86,'运动相机',22),(87,'智能设备',24),(88,'智能手表',24),(89,'智能手环',24),(90,'VR眼镜',24),(91,'运动跟踪器',24),(92,'健康监测',24),(93,'电子书',24),(94,'雷神智能',24),(95,'Apple Watch',24),(96,'极路由',24),(97,'数码配件',23),(98,'存储卡',23),(99,'三脚架/云台',23),(100,'相机包',23),(101,'滤镜',23),(102,'闪光灯/补光灯',23),(103,'机身附件',23),(104,'镜头',23),(105,'移动电源',23),(106,'蓝牙运动耳机',23),(107,'苹果认证',23),(108,'保护套',23),(109,'耳机耳麦',23),(110,'电源',23),(111,'数码照相机',23),(112,'苹果鼠标',23),(113,'智能手表',23),(114,'智能手环',23),(115,'智能眼镜',23),(116,'VR眼镜',23),(117,'其他智能',23),(118,'智能设备',23),(119,'智能手表',23),(120,'智能手环',23),(121,'健康设备',24),(122,'智能眼镜',24),(123,'AppleWatch',24),(124,'运动跟踪器',24),(125,'智能配件',24),(126,'智能路由器',24),(127,'智能音箱',24),(128,'智能机器人',24),(129,'智能摄像',24),(130,'智能车',24),(131,'其他智能',24),(132,'影音娱乐',25),(133,'耳机/耳麦',25),(134,'音箱/音响',25),(135,'麦克风',25),(136,'声卡',25),(137,'智能盒子',25),(138,'收音机',25),(139,'专业音频',25),(140,'电子教育',26),(141,'复读机',26),(142,'点读机',26),(143,'电子词典',26),(144,'早教机',26),(145,'智能手表',24),(146,'智能手环',24),(147,'智能眼镜',24),(148,'其他智能',24),(149,'智能配件',24),(169,'玩具乐器',12),(170,'相机',22),(171,'数码相机',22),(172,'单反相机',22),(173,'单电/微单相机',22),(174,'拍立得',22),(175,'运动相机',22),(176,'智能设备',24),(177,'智能手表',24),(178,'智能手环',24),(179,'VR眼镜',24),(180,'运动跟踪器',24),(181,'健康监测',24),(182,'电子书',24),(183,'雷神智能',24),(184,'Apple Watch',24),(185,'极路由',24),(186,'数码配件',23),(187,'存储卡',23),(188,'三脚架/云台',23),(189,'相机包',23),(190,'滤镜',23),(191,'闪光灯/补光灯',23),(192,'机身附件',23),(193,'镜头',23),(194,'移动电源',23),(195,'蓝牙运动耳机',23),(196,'苹果认证',23),(197,'保护套',23),(198,'耳机耳麦',23),(199,'电源',23),(200,'数码照相机',23),(201,'苹果鼠标',23),(202,'智能手表',23),(203,'智能手环',23),(204,'智能眼镜',23),(205,'VR眼镜',23),(206,'其他智能',23),(207,'智能设备',23),(208,'智能手表',23),(209,'智能手环',23),(210,'健康设备',24),(211,'智能眼镜',24),(212,'AppleWatch',24),(213,'运动跟踪器',24),(214,'智能配件',24),(215,'智能路由器',24),(216,'智能音箱',24),(217,'智能机器人',24),(218,'智能摄像',24),(219,'智能车',24),(220,'其他智能',24),(221,'影音娱乐',25),(222,'耳机/耳麦',25),(223,'音箱/音响',25),(224,'麦克风',25),(225,'声卡',25),(226,'智能盒子',25),(227,'收音机',25),(228,'专业音频',25),(229,'电子教育',26),(230,'复读机',26),(231,'点读机',26),(232,'电子词典',26),(233,'早教机',26),(234,'智能手表',24),(235,'智能手环',24),(236,'智能眼镜',24),(237,'其他智能',24),(238,'智能配件',24),(239,'大家电',16),(240,'平板电视',16),(241,'洗衣机',16),(242,'冰箱',16),(243,'空调',16),(244,'热水器',16),(245,'家庭影院',16),(246,'DVD/蓝光',16),(247,'电视盒子',16),(248,'迷你音响',16),(249,'烟机/灶具',17),(250,'热水器',17),(251,'消毒柜',17),(252,'洗碗机',17),(253,'集成灶',17),(254,'净水器',17),(255,'饮水机',17),(256,'电饭煲',18),(257,'电压力锅',18),(258,'豆浆机',18),(259,'咖啡机',18),(260,'微波炉',18),(261,'电烤箱',18),(262,'电水壶/热水瓶',18),(263,'榨汁机/原汁机',18),(264,'电饼铛',18),(265,'电磁炉',18),(266,'面包机',18),(267,'煮蛋器',18),(268,'电炖锅',18),(269,'电火锅',18),(270,'电蒸锅',18),(271,'养生壶/煎药壶',18),(272,'酸奶机',18),(273,'料理机',18),(274,'电烧烤炉',18),(275,'果干机',18),(276,'多用途锅',18),(277,'空气炸锅',18),(278,'破壁机',18),(279,'电风扇',19),(280,'冷风扇',19),(281,'吸尘器',19),(282,'电暖器',19),(283,'加湿器',19),(284,'空气净化器',19),(285,'饮水机',19),(286,'电蚊香',19),(287,'扫地机器人',19),(288,'挂烫机/熨斗',19),(289,'除湿机',19),(290,'干衣机',19),(291,'其他生活电器',19),(292,'收纳整理',30),(293,'雨伞雨具',30),(294,'缝纫/针织用品',30),(295,'洗晒/熨烫',30),(296,'净化除尘',30),(297,'衣架',30),(298,'浴室用品',30),(299,'香薰',30),(300,'浴室角落',30),(301,'净化',30),(302,'居家布艺',31),(303,'窗纱',31),(304,'窗帘/窗饰',31),(305,'坐垫/椅垫',31),(306,'沙发垫',31),(307,'桌布/罩件',31),(308,'靠垫/抱枕',31),(309,'床品套件',31),(310,'被子',31),(311,'枕芯',31),(312,'毛巾/浴巾',31),(313,'地毯',31),(314,'地垫',31),(315,'毛毯/线毯',31),(316,'蚊帐',31),(317,'凉席',31),(318,'床幔/床帘',31),(319,'床垫/床褥',31),(320,'被芯',31),(321,'枕套',31),(322,'蚊帐',31),(323,'毛巾',31),(324,'浴室用品',30),(325,'宠物生活',32),(326,'宠物食品',32),(327,'宠物零食',32),(328,'宠物玩具',32),(329,'宠物日用',32),(330,'宠物出行',32),(331,'宠物护理',32),(332,'宠物窝笼',32),(333,'宠物服饰',32),(334,'水族',32),(335,'仓鼠类',32),(336,'爬虫/虫宠',32),(337,'宠物生活',32),(338,'猫粮',32),(339,'狗粮',32),(340,'电脑整机',33),(341,'笔记本',33),(342,'游戏本',33),(343,'轻薄本',33),(344,'商用笔记本',33),(345,'台式机',33),(346,'一体机',33),(347,'服务器',33),(348,'工作站',33),(349,'平板电脑',33),(350,'笔记本电脑',33),(351,'游戏本',33),(352,'台式机',33),(353,'一体机',33),(354,'服务器',33),(355,'工作站',33),(356,'平板电脑',33),(357,'电脑配件',34),(358,'CPU',34),(359,'主板',34),(360,'显卡',34),(361,'内存',34),(362,'硬盘',34),(363,'机箱',34),(364,'电源',34),(365,'显示器',34),(366,'散热器',34),(367,'光驱',34),(368,'声卡',34),(369,'网卡',34),(370,'CPU风扇',34),(371,'SSD固态硬盘',34),(372,'外设产品',35),(373,'鼠标',35),(374,'键盘',35),(375,'U盘',35),(376,'移动硬盘',35),(377,'摄像头',35),(378,'网卡',35),(379,'网络仪表',35),(380,'手写板',35),(381,'游戏设备',36),(382,'游戏机',36),(383,'游戏耳机',36),(384,'游戏手柄',36),(385,'游戏软件',36),(386,'游戏周边',36),(387,'网络产品',37),(388,'路由器',37),(389,'网卡',37),(390,'交换机',37),(391,'网络存储',37),(392,'网络摄像头',37),(393,'办公设备',38),(394,'打印机',38),(395,'投影仪',38),(396,'扫描仪',38),(397,'碎纸机',38),(398,'考勤机',38),(399,'文具/耗材',39),(400,'墨盒',39),(401,'硒鼓',39),(402,'纸类',39),(403,'办公文具',39),(404,'财务用品',39),(405,'学生文具',39),(406,'服务产品',40),(407,'安装服务',40),(408,'维修服务',40),(409,'清洗服务',40),(410,'延保服务',40),(411,'烹饪锅具',41),(412,'炒锅',41),(413,'煎锅',41),(414,'汤锅',41),(415,'蒸锅',41),(416,'奶锅',41),(417,'压力锅',41),(418,'电饭煲',41),(419,'刀剪菜板',42),(420,'刀具',42),(421,'砧板',42),(422,'厨房配件',43),(423,'保鲜盒',43),(424,'收纳',43),(425,'水具酒具',44),(426,'水杯',44),(427,'保温杯',44),(428,'餐具',45),(429,'碗',45),(430,'盘',45),(431,'筷勺',45),(432,'刀叉',45),(433,'酒店用品',46),(434,'茶具/咖啡具',47),(435,'茶具',47),(436,'咖啡具',47),(437,'清洁用品',48),(438,'纸品',48),(439,'清洁工具',45),(440,'面部护肤',49),(441,'面膜',49),(442,'乳液/面霜',49),(443,'化妆水',49),(444,'精华',49),(445,'眼霜',49),(446,'防晒',49),(447,'身体护理',50),(448,'沐浴',50),(449,'身体乳',50),(450,'洗手液',50),(451,'口腔护理',51),(452,'牙刷',51),(453,'牙膏',51),(454,'女性护理',52),(455,'卫生巾',52),(456,'洗发护发',53),(457,'洗发水',53),(458,'护发',53),(459,'香水彩妆',54),(460,'口红',54),(461,'粉底',54),(462,'女装',55),(463,'连衣裙',55),(464,'半身裙',55),(465,'男装',56),(466,'夹克',56),(467,'内衣',57),(468,'文胸',57),(469,'洗衣服务',58),(470,'服饰配件',59),(471,'围巾',59),(472,'帽子',59),(473,'钟表',60),(474,'流行男鞋',61),(475,'时尚女鞋',62),(476,'奶粉',63),(477,'婴儿奶粉',63),(478,'营养辅食',64),(479,'果泥',64),(480,'尿裤湿巾',65),(481,'纸尿裤',65),(482,'湿巾',65),(483,'喂养用品',66),(484,'奶瓶',66),(485,'童车童床',68),(486,'婴儿车',68),(487,'寝居服饰',69),(488,'儿童睡袋',69),(489,'妈妈专区',70),(490,'孕妇装',70),(491,'童装童鞋',71),(492,'童装',71),(493,'安全座椅',72),(494,'潮流女包',73),(495,'精品男包',74),(496,'功能箱包',75),(497,'礼品',76),(498,'奢侈品',77),(499,'婚庆',78),(500,'进口食品',79),(501,'地方特产',80),(502,'休闲食品',81),(503,'饼干蛋糕',81),(504,'粮油调味',82),(505,'大米',82),(506,'饮料冲调',83),(507,'饮用水',83),(508,'食品礼券',84),(509,'茗茶',85),(510,'绿茶',85),(511,'时尚饰品',86),(512,'黄金',87),(513,'K金饰品',88),(514,'金银投资',89),(515,'银饰',90),(516,'钻石',91),(517,'翡翠玉石',92),(518,'水晶玛瑙',93),(519,'彩宝',94),(520,'铂金',95),(521,'木手串/把件',96),(522,'珍珠',97),(523,'维修保养',98),(524,'车载电器',99),(525,'美容清洗',100),(526,'汽车装饰',101),(527,'安全自驾',102),(528,'汽车服务',103),(529,'赛事改装',104),(530,'运动鞋包',105),(531,'运动服饰',106),(532,'骑行运动',107),(533,'垂钓用品',108),(534,'游泳用品',109),(535,'户外鞋服',110),(536,'户外装备',111),(537,'健身训练',112),(538,'体育用品',113);

-- =========================================================
-- Table structure for table menu
-- =========================================================
DROP TABLE IF EXISTS app.menu CASCADE;
CREATE TABLE app.menu (
  menu_id bigint NOT NULL,
  pid bigint NULL,
  name varchar(100) NOT NULL,
  code varchar(100) NULL,
  to_code varchar(100) NULL,
  type integer NOT NULL,
  status boolean NOT NULL DEFAULT true,
  level integer NOT NULL,

  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (menu_id)
);
COMMENT ON TABLE app.menu IS '系统菜单/权限（对齐 MySQL 源结构，不含 id 自增主键）';
COMMENT ON COLUMN app.menu.status IS '菜单状态（true=启用，false=禁用），默认启用';
COMMENT ON COLUMN app.menu.type IS '菜单类型（1=目录，2=菜单/按钮，0=其它）';
COMMENT ON COLUMN app.menu.level IS '菜单层级（1=一级，2=二级，3=三级，4=按钮）';
--
-- Dumping data for table menu
--
INSERT INTO app.menu (menu_id, pid, name, code, to_code, type, level) VALUES
(1,0,'全部数据','','',1,1),
(7,1,'权限管理','Acl','',1,2),
(8,7,'用户管理','User','',1,3),
(11,8,'添加用户','btn.User.add','',2,4),
(12,8,'删除用户','btn.User.remove','',2,4),
(13,8,'修改用户','btn.User.update','',2,4),

(9,7,'角色管理','Role','',1,3),

(16,9,'添加角色','btn.Role.add','',2,4),
(17,9,'修改角色','btn.Role.update','',2,4),
(18,9,'删除角色','btn.Role.remove','',2,4),
(10,7,'菜单管理','Permission','',1,3),
(19,10,'添加','btn.Permission.add','',2,4),
(20,10,'修改','btn.Permission.update','',2,4),
(21,10,'删除','btn.Permission.remove','',2,4),
(22,1,'商品管理','Product','',1,2),
(23,22,'分类管理','Category','',1,3),
(44,23,'添加子分类','btn.Category.add','',2,4),
(24,22,'平台属性管理','Attr','',1,3),
(50,24,'添加属性','btn.Attr.add','',2,4),
(51,24,'更新属性','btn.Attr.update','',2,4),
(52,24,'删除属性','btn.Attr.remove','',2,4),
(25,22,'品牌管理','Trademark','',1,3),
(47,25,'添加品牌','btn.Trademark.add','',2,4),
(48,25,'修改品牌','btn.Trademark.update','',2,4),
(49,25,'删除品牌','btn.Trademark.remove','',2,4),
(13146,25,'查询品牌','btn.Trademark.search','',2,4),
(26,22,'SPU管理','Spu','',1,3),
(53,26,'添加SPU','btn.Spu.add','',2,4),
(54,26,'添加SKU','btn.Spu.addsku','',2,4),
(55,26,'更新Spu','btn.Spu.update','',2,4),
(56,26,'查看SKU列表','btn.Spu.skus','',2,4),
(57,26,'删除Spu','btn.Spu.remove','',2,4),
(27,22,'Sku管理','Sku','',1,3),
(58,27,'Sku上架','btn.Sku.up','',2,4),
(15805,27,'Sku下架','btn.Sku.down','',2,4),
(59,27,'更新SKU','btn.Sku.update','',2,4),
(60,27,'Sku详情','btn.Sku.detail','',2,4),
(61,27,'删除Sku','btn.Sku.remove','',2,4),
(28,1,'订单管理','Order','',1,2),
(29,28,'订单列表','OrderList','',1,3),
(41,29,'查看订单详情','btn.OrderList.detail','OrderShow',2,4),
(42,29,'退款','btn.OrderList.refund','Refund',2,4),
(62,28,'退款管理','Refund','',1,3),
(30,1,'客户管理','ClientUser','',1,2),
(31,30,'客户列表','UserList','',1,3),
(43,31,'锁定客户','btn.UserList.lock','',2,4),
(32,1,'优惠管理','Discount','',1,2),
(33,32,'优惠活动管理','Activity','',1,3),
(35,33,'添加活动','btn.Activity.add','ActivityAdd',2,4),
(36,33,'修改活动','btn.Activity.update','ActivityEdit',2,4),
(37,33,'活动规则','btn.Activity.rule','ActivityRule',2,4),
(14988,33,'活动监督','btn.Activity.supervise','',0,4),
(34,32,'优惠券管理','Coupon','',1,3),
(38,34,'添加优惠券','btn.Coupon.add','CouponAdd',2,4),
(39,34,'修改优惠券','btn.Coupon.update','CouponEdit',2,4),
(40,34,'活动规则','btn.Coupon.rule','CouponRule',2,4),
(15804,32,'优惠时间管理','Time','',0,3),
(100,1,'全部','btn.all','',2,2);

-- =========================================================
-- Table structure for table role
-- =========================================================
DROP TABLE IF EXISTS app.role CASCADE;
CREATE TABLE app.role (
  role_id bigint NOT NULL,
  role_name varchar(100) NOT NULL,
  remark varchar(255) NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (role_id),
  UNIQUE (role_name)
);
COMMENT ON TABLE app.role IS '角色';
--
-- Dumping data for table role
--
INSERT INTO app.role (role_id, role_name, remark) VALUES (1,'管理员','超级管理员'),(2,'普通用户','普通用户'),(3,'游客','游客，只能查看');

-- =========================================================
-- Table structure for table role_menu
-- =========================================================
DROP TABLE IF EXISTS app.role_menu CASCADE;
CREATE TABLE app.role_menu (
  role_id bigint NOT NULL,
  menu_id bigint NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (role_id, menu_id)
);
CREATE INDEX idx_role_menu_role_id ON app.role_menu (role_id);
CREATE INDEX idx_role_menu_menu_id ON app.role_menu (menu_id);
COMMENT ON TABLE app.role_menu IS '角色-菜单关联';
--
-- Dumping data for table role_menu
--
INSERT INTO app.role_menu (role_id, menu_id) VALUES
-- 管理员(role 1)：对齐 MySQL 源 role_id=1 的完整菜单集合
(1,1),(1,7),(1,8),(1,11),(1,12),(1,13),(1,9),(1,16),(1,17),(1,18),(1,10),(1,19),(1,20),(1,21),(1,22),(1,23),(1,44),(1,24),(1,50),(1,51),(1,52),(1,25),(1,47),(1,48),(1,49),(1,13146),(1,26),(1,53),(1,54),(1,55),(1,56),(1,57),(1,27),(1,58),(1,15805),(1,59),(1,60),(1,61),(1,28),(1,29),(1,41),(1,42),(1,62),(1,30),(1,31),(1,43),(1,32),(1,33),(1,35),(1,36),(1,37),(1,14988),(1,34),(1,38),(1,39),(1,40),(1,15804),(1,100),
-- 普通用户(role 2)：商品管理相关模块 + 各级 CRUD 菜单
-- 注意：menu 表中不存在 menu_id 2~6，原数据引用了这些不存在的 id 会导致外键约束失败，
-- 此处将根菜单映射为真实存在的 menu_id=22（商品管理）。
(2,22),(2,23),(2,44),(2,24),(2,50),(2,51),(2,52),(2,25),(2,47),(2,48),(2,49),(2,13146),(2,26),(2,53),(2,54),(2,55),(2,56),(2,57),(2,27),(2,58),(2,15805),(2,59),(2,60),(2,61),
-- 游客(role 3)
-- 同上，menu_id=2 不存在，映射为 22（商品管理）。
(3,22),(3,18),(3,23),(3,26),(3,31),(3,36),(3,41),(3,44),(3,47),(3,52);

-- =========================================================
-- Table structure for table sale_attr
-- =========================================================
DROP TABLE IF EXISTS app.sale_attr CASCADE;
CREATE TABLE app.sale_attr (
  sale_attr_id bigint NOT NULL,
  sale_attr_name varchar(255) NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (sale_attr_id)
);
COMMENT ON TABLE app.sale_attr IS '销售属性';
--
-- Dumping data for table sale_attr
--
INSERT INTO app.sale_attr (sale_attr_id, sale_attr_name) VALUES (1,'颜色'),(2,'版本'),(3,'尺码');

-- =========================================================
-- Table structure for table sale_attr_value
-- =========================================================
DROP TABLE IF EXISTS app.sale_attr_value CASCADE;
CREATE TABLE app.sale_attr_value (
  sale_attr_value_id bigint NOT NULL,
  spu_id bigint NOT NULL,
  sale_attr_id bigint NOT NULL,
  sale_attr_value_name varchar(255) NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (sale_attr_value_id)
);
CREATE INDEX idx_sale_attr_value_spu_id ON app.sale_attr_value (spu_id);
COMMENT ON TABLE app.sale_attr_value IS '销售属性值';
-- 注：sale_attr_value 初始数据由代码 seed 生成（雪花ID），见 seed 包。

-- =========================================================
-- Table structure for table sku
-- =========================================================
DROP TABLE IF EXISTS app.sku CASCADE;
CREATE TABLE app.sku (
  sku_id bigint NOT NULL,
  spu_id bigint NOT NULL,
  category3_id bigint NOT NULL,
  tm_id bigint NOT NULL,
  sku_name varchar(255) NOT NULL,
  weight_mg bigint NOT NULL,
  price_cent bigint NOT NULL,
  sku_desc text NOT NULL,
  is_sale boolean NOT NULL DEFAULT false,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (sku_id)
);
CREATE INDEX idx_sku_spu_id ON app.sku (spu_id);
CREATE INDEX idx_sku_category_3_id ON app.sku (category3_id);
CREATE INDEX idx_sku_tm_id ON app.sku (tm_id);
COMMENT ON TABLE app.sku IS '商品 SKU（库存量单位）';
COMMENT ON COLUMN app.sku.price_cent IS '价格，单位：分（最小货币单位，cents）';
COMMENT ON COLUMN app.sku.weight_mg IS '重量，单位：毫克（最小重量单位，milligrams）';
COMMENT ON COLUMN app.sku.is_sale IS '是否上架（true=上架，false=下架）';
-- 注：sku 初始数据由代码 seed 生成（雪花ID），见 seed 包。

-- =========================================================
-- Table structure for table sku_attr_value
-- =========================================================
DROP TABLE IF EXISTS app.sku_attr_value CASCADE;
CREATE TABLE app.sku_attr_value (
  sku_attr_value_id bigint NOT NULL,
  attr_id bigint NOT NULL,
  sku_id bigint NOT NULL,
  attr_value_id bigint NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (sku_attr_value_id)
);
CREATE INDEX idx_sku_attr_value_sku_id ON app.sku_attr_value (sku_id);
CREATE INDEX idx_sku_attr_value_attr_id ON app.sku_attr_value (attr_id);
CREATE INDEX idx_sku_attr_value_attr_value_id ON app.sku_attr_value (attr_value_id);
COMMENT ON TABLE app.sku_attr_value IS 'SKU 与平台属性值关联';
-- 注：sku_attr_value 初始数据由代码 seed 生成（雪花ID），见 seed 包。

-- =========================================================
-- Table structure for table sku_image
-- =========================================================
DROP TABLE IF EXISTS app.sku_image CASCADE;
CREATE TABLE app.sku_image (
  image_id bigint NOT NULL,
  sku_id bigint NOT NULL,
  image_url varchar(255) NOT NULL,
  spu_image_id bigint NULL,
  is_default boolean NOT NULL DEFAULT false,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (image_id)
);
CREATE INDEX idx_sku_image_sku_id ON app.sku_image (sku_id);
COMMENT ON TABLE app.sku_image IS 'SKU 图片';
COMMENT ON COLUMN app.sku_image.is_default IS '是否默认图片（true=默认，false=否）';
-- 注：sku_image 初始数据由代码 seed 生成（雪花ID），见 seed 包。

-- =========================================================
-- Table structure for table sku_sale_attr_value
-- =========================================================
DROP TABLE IF EXISTS app.sku_sale_attr_value CASCADE;
CREATE TABLE app.sku_sale_attr_value (
  sku_sale_attr_value_id bigint NOT NULL,
  sku_id bigint NOT NULL,
  sale_attr_value_id bigint NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (sku_sale_attr_value_id)
);
CREATE INDEX idx_sku_sale_attr_value_sku_id ON app.sku_sale_attr_value (sku_id);
CREATE INDEX idx_sku_sale_attr_value_sale_attr_value_id ON app.sku_sale_attr_value (sale_attr_value_id);
COMMENT ON TABLE app.sku_sale_attr_value IS 'SKU 与销售属性值关联';
-- 注：sku_sale_attr_value 初始数据由代码 seed 生成（雪花ID），见 seed 包。

-- =========================================================
-- Table structure for table spu
-- =========================================================
DROP TABLE IF EXISTS app.spu CASCADE;
CREATE TABLE app.spu (
  spu_id bigint NOT NULL,
  spu_name varchar(255) NOT NULL,
  description text NOT NULL,
  category3_id bigint NOT NULL,
  tm_id bigint NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (spu_id)
);
CREATE INDEX idx_spu_category3_id ON app.spu (category3_id);
CREATE INDEX idx_spu_tm_id ON app.spu (tm_id);
COMMENT ON TABLE app.spu IS '商品 SPU（标准化产品单元）';
--
-- Dumping data for table spu
-- 注：spu 为业务演示数据，使用雪花ID，由代码 seed（seed/demo.go）生成，此处不写死。

-- =========================================================
-- Table structure for table spu_image_list
-- =========================================================
DROP TABLE IF EXISTS app.spu_image_list CASCADE;
CREATE TABLE app.spu_image_list (
  image_id bigint NOT NULL,
  image_name varchar(255) NULL,
  image_url varchar(255) NULL,
  spu_id bigint NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (image_id)
);
CREATE INDEX idx_spu_image_list_spu_id ON app.spu_image_list (spu_id);
COMMENT ON TABLE app.spu_image_list IS 'SPU 图片';
--
-- Dumping data for table spu_image_list
-- 注：spu_image_list 为业务演示数据，使用雪花ID，由代码 seed（seed/demo.go）生成，此处不写死。

-- =========================================================
-- Table structure for table spu_sale_attr
-- =========================================================
DROP TABLE IF EXISTS app.spu_sale_attr CASCADE;
CREATE TABLE app.spu_sale_attr (
  spu_sale_attr_id bigint NOT NULL,
  base_sale_attr_id bigint NOT NULL,
  sale_attr_name varchar(255) NULL,
  spu_id bigint NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (spu_sale_attr_id)
);
CREATE INDEX idx_spu_sale_attr_spu_id ON app.spu_sale_attr (spu_id);
COMMENT ON TABLE app.spu_sale_attr IS 'SPU 与销售属性关联';
-- 注：spu_sale_attr 初始数据由代码 seed 生成（雪花ID），见 seed 包。

-- =========================================================
-- Table structure for table trademark
-- =========================================================
DROP TABLE IF EXISTS app.trademark CASCADE;
CREATE TABLE app.trademark (
  tm_id bigint NOT NULL,
  tm_name varchar(255) NOT NULL,
  logo_url varchar(255) NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (tm_id),
  UNIQUE (tm_name)
);
COMMENT ON TABLE app.trademark IS '品牌（商标）';
--
-- Dumping data for table trademark
--
INSERT INTO app.trademark (tm_id, tm_name, logo_url) VALUES (1,'华为','/api/static/img/product/default/huawei.jpg'),(2,'小米','/api/static/img/product/default/xiaomi.png'),(3,'OPPO','/api/static/img/product/default/oppo.jpeg'),(4,'vivo','/api/static/img/product/default/vivo.jpeg'),(5,'雅诗兰黛','/api/static/img/product/default/ysld.jpeg'),(6,'Apple','/api/static/img/product/default/apple_logo.jpeg'),(7,'荣耀','/api/static/img/product/default/honor.png');

-- =========================================================
-- Table structure for table users
-- =========================================================
DROP TABLE IF EXISTS app.users CASCADE;
CREATE TABLE app.users (
  user_id bigint NOT NULL,
  username varchar(255) NOT NULL,
  password varchar(255) NOT NULL,
  name varchar(255) NULL,
  phone varchar(255) NULL,
  avatar varchar(255) NULL,
  status boolean NOT NULL DEFAULT true,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id),
  UNIQUE (username)
);
COMMENT ON TABLE app.users IS '系统用户';
COMMENT ON COLUMN app.users.status IS '账号状态（true=正常，false=禁用/锁定），默认正常';
-- 注：users / user_role 的初始数据改为由代码 seed 生成（user_id 使用雪花ID），
--     不再在此硬编码，避免与线上雪花算法生成的 ID 不一致。

-- =========================================================
-- Table structure for table user_role
-- =========================================================
DROP TABLE IF EXISTS app.user_role CASCADE;
CREATE TABLE app.user_role (
  user_id bigint NOT NULL,
  role_id bigint NOT NULL,
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id, role_id)
);
CREATE INDEX idx_user_role_user_id ON app.user_role (user_id);
CREATE INDEX idx_user_role_role_id ON app.user_role (role_id);
COMMENT ON TABLE app.user_role IS '用户-角色关联';
-- 注：user_role 的初始关联由代码 seed 生成（见 seed 包）。

-- =========================================================
-- update_time 自动更新触发器
-- =========================================================
CREATE OR REPLACE FUNCTION set_update_time()
RETURNS TRIGGER AS $$
BEGIN
  NEW.update_time = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_attr_update_time BEFORE UPDATE ON app.attr FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_attr_value_update_time BEFORE UPDATE ON app.attr_value FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_category1_update_time BEFORE UPDATE ON app.category1 FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_category2_update_time BEFORE UPDATE ON app.category2 FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_category3_update_time BEFORE UPDATE ON app.category3 FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_menu_update_time BEFORE UPDATE ON app.menu FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_role_update_time BEFORE UPDATE ON app.role FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_role_menu_update_time BEFORE UPDATE ON app.role_menu FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_sale_attr_update_time BEFORE UPDATE ON app.sale_attr FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_sale_attr_value_update_time BEFORE UPDATE ON app.sale_attr_value FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_sku_update_time BEFORE UPDATE ON app.sku FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_sku_attr_value_update_time BEFORE UPDATE ON app.sku_attr_value FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_sku_image_update_time BEFORE UPDATE ON app.sku_image FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_sku_sale_attr_value_update_time BEFORE UPDATE ON app.sku_sale_attr_value FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_spu_update_time BEFORE UPDATE ON app.spu FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_spu_image_list_update_time BEFORE UPDATE ON app.spu_image_list FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_spu_sale_attr_update_time BEFORE UPDATE ON app.spu_sale_attr FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_trademark_update_time BEFORE UPDATE ON app.trademark FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_users_update_time BEFORE UPDATE ON app.users FOR EACH ROW EXECUTE FUNCTION set_update_time();
CREATE TRIGGER trg_user_role_update_time BEFORE UPDATE ON app.user_role FOR EACH ROW EXECUTE FUNCTION set_update_time();

-- =========================================================
-- 外键约束（核心关联 + 合适的 ON DELETE 策略）
-- 说明：
--   1. 写在全部 CREATE TABLE / 触发器之后，避免建表时引用尚未存在的父表。
--   2. ON DELETE RESTRICT：禁止删除仍被子表引用的父记录，强制先清理子数据，
--      防止"手滑删大类目连带清光商品"等事故。
--   3. ON DELETE CASCADE：删除父记录时自动连带删除子表引用行，
--      对应代码中 DeleteSpu 等手工"先子后父"的清理逻辑（保留代码也无妨）。
--   4. app.sale_attr 为字典表（颜色/版本/尺码），几乎不删，故 sale_attr_value
--      未对其加外键，以保持字典灵活性。
-- =========================================================

-- 类目三级链
ALTER TABLE app.category2
  ADD CONSTRAINT fk_category2_cat1 FOREIGN KEY (category1_id) REFERENCES app.category1 (category1_id) ON DELETE RESTRICT;
ALTER TABLE app.category3
  ADD CONSTRAINT fk_category3_cat2 FOREIGN KEY (category2_id) REFERENCES app.category2 (category2_id) ON DELETE RESTRICT;

-- 平台属性
ALTER TABLE app.attr
  ADD CONSTRAINT fk_attr_cat3 FOREIGN KEY (category_id) REFERENCES app.category3 (category3_id) ON DELETE RESTRICT;
ALTER TABLE app.attr_value
  ADD CONSTRAINT fk_attr_value_attr FOREIGN KEY (attr_id) REFERENCES app.attr (attr_id) ON DELETE CASCADE;

-- SPU / SKU 主链
ALTER TABLE app.spu
  ADD CONSTRAINT fk_spu_cat3 FOREIGN KEY (category3_id) REFERENCES app.category3 (category3_id) ON DELETE RESTRICT,
  ADD CONSTRAINT fk_spu_tm FOREIGN KEY (tm_id) REFERENCES app.trademark (tm_id) ON DELETE RESTRICT;
ALTER TABLE app.sku
  ADD CONSTRAINT fk_sku_spu FOREIGN KEY (spu_id) REFERENCES app.spu (spu_id) ON DELETE CASCADE,
  ADD CONSTRAINT fk_sku_cat3 FOREIGN KEY (category3_id) REFERENCES app.category3 (category3_id) ON DELETE RESTRICT,
  ADD CONSTRAINT fk_sku_tm FOREIGN KEY (tm_id) REFERENCES app.trademark (tm_id) ON DELETE RESTRICT;

-- SPU 从表
ALTER TABLE app.spu_image_list
  ADD CONSTRAINT fk_spu_img_spu FOREIGN KEY (spu_id) REFERENCES app.spu (spu_id) ON DELETE CASCADE;
ALTER TABLE app.spu_sale_attr
  ADD CONSTRAINT fk_spu_sale_attr_spu FOREIGN KEY (spu_id) REFERENCES app.spu (spu_id) ON DELETE CASCADE,
  ADD CONSTRAINT fk_spu_sale_attr_base FOREIGN KEY (base_sale_attr_id) REFERENCES app.sale_attr (sale_attr_id) ON DELETE RESTRICT;

-- SKU 从表
ALTER TABLE app.sku_image
  ADD CONSTRAINT fk_sku_img_sku FOREIGN KEY (sku_id) REFERENCES app.sku (sku_id) ON DELETE CASCADE,
  ADD CONSTRAINT fk_sku_img_spi FOREIGN KEY (spu_image_id) REFERENCES app.spu_image_list (image_id) ON DELETE SET NULL;
ALTER TABLE app.sku_attr_value
  ADD CONSTRAINT fk_sku_attr_val_sku FOREIGN KEY (sku_id) REFERENCES app.sku (sku_id) ON DELETE CASCADE,
  ADD CONSTRAINT fk_sku_attr_val_attr FOREIGN KEY (attr_id) REFERENCES app.attr (attr_id) ON DELETE RESTRICT,
  ADD CONSTRAINT fk_sku_attr_val_attrval FOREIGN KEY (attr_value_id) REFERENCES app.attr_value (attr_value_id) ON DELETE RESTRICT;
ALTER TABLE app.sku_sale_attr_value
  ADD CONSTRAINT fk_sku_sale_av_sku FOREIGN KEY (sku_id) REFERENCES app.sku (sku_id) ON DELETE CASCADE,
  ADD CONSTRAINT fk_sku_sale_av_sav FOREIGN KEY (sale_attr_value_id) REFERENCES app.sale_attr_value (sale_attr_value_id) ON DELETE CASCADE;

-- RBAC：角色-菜单、用户-角色
ALTER TABLE app.role_menu
  ADD CONSTRAINT fk_role_menu_role FOREIGN KEY (role_id) REFERENCES app.role (role_id) ON DELETE CASCADE,
  ADD CONSTRAINT fk_role_menu_menu FOREIGN KEY (menu_id) REFERENCES app.menu (menu_id) ON DELETE CASCADE;
ALTER TABLE app.user_role
  ADD CONSTRAINT fk_user_role_user FOREIGN KEY (user_id) REFERENCES app.users (user_id) ON DELETE CASCADE,
  ADD CONSTRAINT fk_user_role_role FOREIGN KEY (role_id) REFERENCES app.role (role_id) ON DELETE CASCADE;

-- =========================================================
-- 设置默认 search_path 为 app
-- =========================================================
ALTER DATABASE "vue_admin" SET search_path TO app, public;
