-- MySQL dump 10.13  Distrib 5.7.34, for Linux (x86_64)
--
-- Host: localhost    Database: blog_isolezvosco
-- ------------------------------------------------------
-- Server version	5.7.34-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `zbp_category`
--

DROP TABLE IF EXISTS `zbp_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `zbp_category` (
  `cate_ID` int(11) NOT NULL AUTO_INCREMENT,
  `cate_Name` varchar(250) NOT NULL DEFAULT '',
  `cate_Order` int(11) NOT NULL DEFAULT '0',
  `cate_Type` int(11) NOT NULL DEFAULT '0',
  `cate_Count` int(11) NOT NULL DEFAULT '0',
  `cate_Alias` varchar(250) NOT NULL DEFAULT '',
  `cate_Group` varchar(250) NOT NULL DEFAULT '',
  `cate_Intro` text NOT NULL,
  `cate_RootID` int(11) NOT NULL DEFAULT '0',
  `cate_ParentID` int(11) NOT NULL DEFAULT '0',
  `cate_Template` varchar(250) NOT NULL DEFAULT '',
  `cate_LogTemplate` varchar(250) NOT NULL DEFAULT '',
  `cate_Meta` longtext NOT NULL,
  PRIMARY KEY (`cate_ID`),
  KEY `zbp_cate_Order` (`cate_Order`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `zbp_category`
--

LOCK TABLES `zbp_category` WRITE;
/*!40000 ALTER TABLE `zbp_category` DISABLE KEYS */;
INSERT INTO `zbp_category` VALUES (1,'未分类',0,0,1,'uncategorized','','',0,0,'','','');
/*!40000 ALTER TABLE `zbp_category` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `zbp_comment`
--

DROP TABLE IF EXISTS `zbp_comment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `zbp_comment` (
  `comm_ID` int(11) NOT NULL AUTO_INCREMENT,
  `comm_LogID` int(11) NOT NULL DEFAULT '0',
  `comm_IsChecking` tinyint(4) NOT NULL DEFAULT '0',
  `comm_RootID` int(11) NOT NULL DEFAULT '0',
  `comm_ParentID` int(11) NOT NULL DEFAULT '0',
  `comm_AuthorID` int(11) NOT NULL DEFAULT '0',
  `comm_Name` varchar(250) NOT NULL DEFAULT '',
  `comm_Email` varchar(250) NOT NULL DEFAULT '',
  `comm_HomePage` varchar(250) NOT NULL DEFAULT '',
  `comm_Content` text NOT NULL,
  `comm_PostTime` int(11) NOT NULL DEFAULT '0',
  `comm_IP` varchar(250) NOT NULL DEFAULT '',
  `comm_Agent` text NOT NULL,
  `comm_Meta` longtext NOT NULL,
  PRIMARY KEY (`comm_ID`),
  KEY `zbp_comm_LRI` (`comm_LogID`,`comm_RootID`,`comm_IsChecking`),
  KEY `zbp_comm_IsChecking` (`comm_IsChecking`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `zbp_comment`
--

LOCK TABLES `zbp_comment` WRITE;
/*!40000 ALTER TABLE `zbp_comment` DISABLE KEYS */;
/*!40000 ALTER TABLE `zbp_comment` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `zbp_config`
--

DROP TABLE IF EXISTS `zbp_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `zbp_config` (
  `conf_ID` int(11) NOT NULL AUTO_INCREMENT,
  `conf_Name` varchar(250) NOT NULL DEFAULT '',
  `conf_Key` varchar(250) NOT NULL DEFAULT '',
  `conf_Value` longtext NOT NULL,
  PRIMARY KEY (`conf_ID`),
  KEY `zbp_conf_Name` (`conf_Name`)
) ENGINE=MyISAM AUTO_INCREMENT=169 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `zbp_config`
--

LOCK TABLES `zbp_config` WRITE;
/*!40000 ALTER TABLE `zbp_config` DISABLE KEYS */;
INSERT INTO `zbp_config` VALUES (1,'cache','sidebars_default','s:142:\"{\"1\":\"calendar|controlpanel|catalog|searchpanel|comments|archives|favorite|link|misc\",\"2\":\"\",\"3\":\"\",\"4\":\"\",\"5\":\"\",\"6\":\"\",\"7\":\"\",\"8\":\"\",\"9\":\"\"}\";'),(2,'cache','sidebars_Zit','s:196:\"{\"1\":\"calendar|archives|controlpanel|catalog|favorite|link|misc\",\"2\":\"calendar|catalog|archives|previous\",\"3\":\"catalog|previous|comments\",\"4\":\"comments\",\"5\":\"comments\",\"6\":\"\",\"7\":\"\",\"8\":\"\",\"9\":\"\"}\";'),(3,'cache','sidebars_tpure','s:135:\"{\"1\":\"controlpanel|catalog|comments\",\"2\":\"previous|tags\",\"3\":\"catalog|tags|comments\",\"4\":\"tags\",\"5\":\"tags\",\"6\":\"\",\"7\":\"\",\"8\":\"\",\"9\":\"\"}\";'),(4,'cache','sidebars_WhitePage','s:135:\"{\"1\":\"controlpanel|catalog|comments\",\"2\":\"previous|tags\",\"3\":\"catalog|tags|comments\",\"4\":\"tags\",\"5\":\"tags\",\"6\":\"\",\"7\":\"\",\"8\":\"\",\"9\":\"\"}\";'),(5,'system','ZC_CLOSE_SITE','b:0;'),(6,'system','ZC_BLOG_HOST','s:56:\"h|t|t|p|:|/|/|i|s|o|l|e|z|v|o|s|c|o|m|b|l|e|s|.|c|o|m|/|\";'),(7,'system','ZC_BLOG_NAME','s:5:\"admin\";'),(8,'system','ZC_BLOG_SUBNAME','s:17:\"Good Luck To You!\";'),(9,'system','ZC_BLOG_THEME','s:2:\"ky\";'),(10,'system','ZC_BLOG_CSS','s:7:\"default\";'),(11,'system','ZC_BLOG_COPYRIGHT','s:44:\"Copyright Your WebSite.Some Rights Reserved.\";'),(12,'system','ZC_BLOG_LANGUAGE','s:5:\"zh-CN\";'),(13,'system','ZC_BLOG_LANGUAGEPACK','s:5:\"zh-cn\";'),(14,'system','ZC_USING_PLUGIN_LIST','s:36:\"AppCentre|UEditor|Totoro|LinksManage\";'),(15,'system','ZC_BLOG_CLSID','s:32:\"fb765dd20a958bd1012549972f0faf32\";'),(16,'system','ZC_TIME_ZONE_NAME','s:13:\"Asia/Shanghai\";'),(17,'system','ZC_UPDATE_INFO_URL','s:32:\"https://update.zblogcn.com/info/\";'),(18,'system','ZC_PERMANENT_DOMAIN_ENABLE','b:0;'),(19,'system','ZC_PERMANENT_DOMAIN_SHOW','b:0;'),(20,'system','ZC_DEBUG_MODE','b:0;'),(21,'system','ZC_DEBUG_MODE_STRICT','b:0;'),(22,'system','ZC_DEBUG_MODE_WARNING','b:1;'),(23,'system','ZC_DEBUG_LOG_ERROR','b:0;'),(24,'system','ZC_BLOG_PRODUCT','s:9:\"Z-BlogPHP\";'),(25,'system','ZC_BLOG_VERSION','s:18:\"1.7.2 Build 173040\";'),(26,'system','ZC_BLOG_COMMIT','s:4:\"3040\";'),(27,'system','ZC_BLOG_PRODUCT_FULL','s:15:\"Z-BlogPHP 1.7.2\";'),(28,'system','ZC_BLOG_PRODUCT_HTML','s:128:\"<a href=\"https://www.zblogcn.com/\" title=\"Z-BlogPHP 1.7.2 Build 173040\" target=\"_blank\" rel=\"noopener norefferrer\">Z-BlogPHP</a>\";'),(29,'system','ZC_BLOG_PRODUCT_FULLHTML','s:134:\"<a href=\"https://www.zblogcn.com/\" title=\"Z-BlogPHP 1.7.2 Build 173040\" target=\"_blank\" rel=\"noopener norefferrer\">Z-BlogPHP 1.7.2</a>\";'),(30,'system','ZC_COMMENT_TURNOFF','b:0;'),(31,'system','ZC_COMMENT_VERIFY_ENABLE','b:0;'),(32,'system','ZC_COMMENT_REVERSE_ORDER','b:0;'),(33,'system','ZC_COMMENT_AUDIT','b:0;'),(34,'system','ZC_COMMENT_VALIDCMTKEY_ENABLE','b:0;'),(35,'system','ZC_VERIFYCODE_STRING','s:29:\"ABCDEFGHKMNPRSTUVWXYZ23456789\";'),(36,'system','ZC_VERIFYCODE_WIDTH','i:90;'),(37,'system','ZC_VERIFYCODE_HEIGHT','i:30;'),(38,'system','ZC_VERIFYCODE_FONT','s:26:\"zb_system/defend/arial.ttf\";'),(39,'system','ZC_DISPLAY_COUNT','i:10;'),(40,'system','ZC_DISPLAY_ORDER','s:12:\"log_PostTime\";'),(41,'system','ZC_PAGEBAR_COUNT','i:10;'),(42,'system','ZC_COMMENTS_DISPLAY_COUNT','i:100;'),(43,'system','ZC_DISPLAY_SUBCATEGORYS','b:1;'),(44,'system','ZC_RSS2_COUNT','i:10;'),(45,'system','ZC_RSS2_ORDER','s:12:\"log_PostTime\";'),(46,'system','ZC_RSS_EXPORT_WHOLE','b:1;'),(47,'system','ZC_MANAGE_COUNT','i:50;'),(48,'system','ZC_MANAGE_ORDER','s:12:\"log_PostTime\";'),(49,'system','ZC_EMOTICONS_FILENAME','s:4:\"face\";'),(50,'system','ZC_EMOTICONS_FILETYPE','s:11:\"png|gif|jpg\";'),(51,'system','ZC_EMOTICONS_FILESIZE','s:2:\"16\";'),(52,'system','ZC_UPLOAD_FILETYPE','s:219:\"jpg|gif|png|jpeg|bmp|webp|psd|wmf|ico|rpm|deb|tar|gz|xz|sit|7z|bz2|zip|rar|xml|xsl|svg|svgz|rtf|doc|docx|ppt|pptx|xls|xlsx|wps|chm|txt|md|pdf|mp3|flac|ape|mp4|mkv|avi|mpg|rm|ra|rmvb|mov|wmv|wma|torrent|apk|json|zba|gzba\";'),(53,'system','ZC_UPLOAD_FILESIZE','i:2;'),(54,'system','ZC_UPLOAD_DIR_YEARMONTHDAY','b:0;'),(55,'system','ZC_USERNAME_MIN','i:2;'),(56,'system','ZC_USERNAME_MAX','i:50;'),(57,'system','ZC_PASSWORD_MIN','i:8;'),(58,'system','ZC_PASSWORD_MAX','i:20;'),(59,'system','ZC_EMAIL_MAX','i:50;'),(60,'system','ZC_HOMEPAGE_MAX','i:100;'),(61,'system','ZC_CONTENT_MAX','i:1000;'),(62,'system','ZC_ARTICLE_TITLE_MAX','i:100;'),(63,'system','ZC_CATEGORY_NAME_MAX','i:50;'),(64,'system','ZC_TAGS_NAME_MAX','i:50;'),(65,'system','ZC_MODULE_NAME_MAX','i:50;'),(66,'system','ZC_ARTICLE_EXCERPT_MAX','i:250;'),(67,'system','ZC_ARTICLE_INTRO_WITH_TEXT','b:0;'),(68,'system','ZC_COMMENT_EXCERPT_MAX','i:20;'),(69,'system','ZC_STATIC_MODE','s:6:\"ACTIVE\";'),(70,'system','ZC_ARTICLE_REGEX','s:18:\"{%host%}?id={%id%}\";'),(71,'system','ZC_PAGE_REGEX','s:18:\"{%host%}?id={%id%}\";'),(72,'system','ZC_CATEGORY_REGEX','s:34:\"{%host%}?cate={%id%}&page={%page%}\";'),(73,'system','ZC_AUTHOR_REGEX','s:34:\"{%host%}?auth={%id%}&page={%page%}\";'),(74,'system','ZC_TAGS_REGEX','s:34:\"{%host%}?tags={%id%}&page={%page%}\";'),(75,'system','ZC_DATE_REGEX','s:36:\"{%host%}?date={%date%}&page={%page%}\";'),(76,'system','ZC_INDEX_REGEX','s:22:\"{%host%}?page={%page%}\";'),(77,'system','ZC_ALIAS_BACK_ATTR','s:4:\"Name\";'),(78,'system','ZC_DATETIME_SEPARATOR','s:1:\"-\";'),(79,'system','ZC_DATETIME_RULE','s:3:\"Y-n\";'),(80,'system','ZC_DATETIME_WITHDAY_RULE','s:5:\"Y-n-j\";'),(81,'system','ZC_SEARCH_COUNT','i:20;'),(82,'system','ZC_SEARCH_REGEX','s:40:\"{%host%}search.php?q={%q%}&page={%page%}\";'),(83,'system','ZC_INDEX_DEFAULT_TEMPLATE','s:5:\"index\";'),(84,'system','ZC_POST_DEFAULT_TEMPLATE','s:6:\"single\";'),(85,'system','ZC_SEARCH_DEFAULT_TEMPLATE','s:6:\"search\";'),(86,'system','ZC_SIDEBAR_ORDER','s:57:\"calendar|archives|controlpanel|catalog|favorite|link|misc\";'),(87,'system','ZC_SIDEBAR2_ORDER','s:34:\"calendar|catalog|archives|previous\";'),(88,'system','ZC_SIDEBAR3_ORDER','s:25:\"catalog|previous|comments\";'),(89,'system','ZC_SIDEBAR4_ORDER','s:8:\"comments\";'),(90,'system','ZC_SIDEBAR5_ORDER','s:8:\"comments\";'),(91,'system','ZC_SIDEBAR6_ORDER','s:0:\"\";'),(92,'system','ZC_SIDEBAR7_ORDER','s:0:\"\";'),(93,'system','ZC_SIDEBAR8_ORDER','s:0:\"\";'),(94,'system','ZC_SIDEBAR9_ORDER','s:0:\"\";'),(95,'system','ZC_SYNTAXHIGHLIGHTER_ENABLE','b:1;'),(96,'system','ZC_CODEMIRROR_ENABLE','b:1;'),(97,'system','ZC_ALLOW_AUDITTING_MEMBER_VISIT_MANAGE','b:0;'),(98,'system','ZC_OUTPUT_OPTION_MEMBER_MAX_LEVEL','i:0;'),(99,'system','ZC_CATEGORY_MANAGE_LEGACY_DISPLAY','b:1;'),(100,'system','ZC_GZIP_ENABLE','b:0;'),(101,'system','ZC_ADMIN_HTML5_ENABLE','b:1;'),(102,'system','ZC_LOADMEMBERS_LEVEL','i:1;'),(103,'system','ZC_LAST_VERSION','s:6:\"172800\";'),(104,'system','ZC_MODULE_CATALOG_STYLE','i:0;'),(105,'system','ZC_MODULE_ARCHIVES_STYLE','i:0;'),(106,'system','ZC_VIEWNUMS_TURNOFF','b:0;'),(107,'system','ZC_LISTONTOP_TURNOFF','b:0;'),(108,'system','ZC_RELATEDLIST_COUNT','i:10;'),(109,'system','ZC_RUNINFO_DISPLAY','b:1;'),(110,'system','ZC_POST_ALIAS_USE_ID_NOT_TITLE','b:0;'),(111,'system','ZC_COMPATIBLE_ASP_URL','b:1;'),(112,'system','ZC_LARGE_DATA','b:0;'),(113,'system','ZC_VERSION_IN_HEADER','b:1;'),(114,'system','ZC_ADDITIONAL_SECURITY','b:1;'),(115,'system','ZC_XMLRPC_ENABLE','b:0;'),(116,'system','ZC_XMLRPC_USE_WEBTOKEN','b:0;'),(117,'system','ZC_USING_CDN_GUESTIP_TYPE','s:11:\"REMOTE_ADDR\";'),(118,'system','ZC_POST_BATCH_DELETE','b:0;'),(119,'system','ZC_JS_304_ENABLE','b:1;'),(120,'system','ZC_DELMEMBER_WITH_ALLDATA','b:0;'),(121,'system','ZC_API_ENABLE','b:0;'),(122,'system','ZC_API_THROTTLE_ENABLE','b:0;'),(123,'system','ZC_API_THROTTLE_MAX_REQS_PER_MIN','i:60;'),(124,'system','ZC_NOW_VERSION','s:6:\"173040\";'),(125,'cache','templates_md5','s:0:\"\";'),(126,'AppCentre','enabledcheck','i:1;'),(127,'AppCentre','checkbeta','i:0;'),(128,'AppCentre','enabledevelop','i:0;'),(129,'AppCentre','enablegzipapp','i:0;'),(130,'Zit','Logo','s:7:\"ZBlogIt\";'),(131,'Zit','Cover','s:47:\"{#ZC_BLOG_HOST#}zb_users/theme/Zit/style/bg.jpg\";'),(132,'Zit','Backdrop','s:0:\"\";'),(133,'Zit','DefaultAdmin','i:0;'),(134,'Zit','ListTags','i:0;'),(135,'Zit','MobileSide','i:0;'),(136,'Zit','HideIntro','i:0;'),(137,'Zit','SideMods','s:0:\"\";'),(138,'Zit','Hue','i:0;'),(139,'Zit','StaticName','i:1;'),(140,'Zit','Motto','s:22:\"Nice to meet you, too!\";'),(141,'Zit','MottoUrl','s:0:\"\";'),(142,'Zit','MottoSize','s:0:\"\";'),(143,'Zit','CmtIds','s:0:\"\";'),(144,'Zit','GbookID','i:2;'),(145,'Zit','Description','s:0:\"\";'),(146,'Zit','Keywords','s:0:\"\";'),(147,'Zit','RelatedTitle','s:12:\"少长咸集\";'),(148,'Zit','CommentTitle','s:12:\"群贤毕至\";'),(149,'Zit','Custom','i:0;'),(150,'cache','templates_md5_array','s:61:\"a:1:{s:8:\"template\";s:32:\"eb58139f3c2dbe8ef8da0143bd8758a3\";}\";'),(151,'AppCentre','lastchecktime','i:1648619428;'),(152,'cache','success_updated_app','s:0:\"\";'),(153,'cache','normal_article_nums','s:1:\"1\";'),(154,'cache','all_comment_nums','i:0;'),(155,'cache','check_comment_nums','i:0;'),(156,'cache','normal_comment_nums','i:0;'),(157,'cache','top_post_array_0','s:6:\"a:0:{}\";'),(158,'cache','reload_statistic','s:1086:\"<!--debug_mode_note--><tr><td class=\'td20\'>当前用户</td><td class=\'td30\'>{$zbp->user->IsGod}<a href=\'../cmd.php?act=misc&type=vrs\' target=\'_blank\'>{$zbp->user->Name}</a></td><td class=\'td20\'>当前版本</td><td class=\'td30\'>{$zbp->version}</td></tr><tr><td class=\'td20\'>文章总数</td><td>1</td><td>分类总数</td><td>1</td></tr><tr><td class=\'td20\'>页面总数</td><td>1</td><td>标签总数</td><td>0</td></tr><tr><td class=\'td20\'>评论总数</td><td>0</td><td>浏览总数</td><td>4</td></tr><tr><td class=\'td20\'>当前主题</td><td>{$zbp->theme}/{$zbp->style} {$theme_version}</td><td>用户总数</td><td>1</td></tr><!--debug_mode_moreinfo--><tr><td class=\'td20\'>协议地址</td><td><a href=\"{#ZC_BLOG_HOST#}zb_system/api.php\" target=\"_blank\">API地址</a> , <a href=\"{#ZC_BLOG_HOST#}zb_system/xml-rpc/index.php\" target=\"_blank\">XML-RPC地址</a></td><td>系统环境</td><td><a href=\'../cmd.php?act=misc&type=phpinfo\' target=\'_blank\'>{$system_environment}</a></td></tr><script type=\"text/javascript\">$(\'#statistic\').attr(\'title\',\'2022-03-28T14:16:19+08:00\');</script>\";'),(159,'cache','reload_statistic_time','i:1648619441;'),(160,'cache','system_environment','s:81:\"Linux3.10.0; nginx1.18.0; PHP7.2.33x64; mysqli5.7.34; curl-o7.70.0; OpenSSL1.0.2u\";'),(161,'cache','all_article_nums','s:1:\"1\";'),(162,'cache','all_page_nums','s:1:\"1\";'),(163,'cache','all_category_nums','s:1:\"1\";'),(164,'cache','all_view_nums','s:1:\"4\";'),(165,'cache','all_tag_nums','s:1:\"0\";'),(166,'cache','all_member_nums','s:1:\"1\";'),(167,'cache','sidebars_os2020','s:196:\"{\"1\":\"calendar|archives|controlpanel|catalog|favorite|link|misc\",\"2\":\"calendar|catalog|archives|previous\",\"3\":\"catalog|previous|comments\",\"4\":\"comments\",\"5\":\"comments\",\"6\":\"\",\"7\":\"\",\"8\":\"\",\"9\":\"\"}\";'),(168,'cache','sidebars_ky','s:196:\"{\"1\":\"calendar|archives|controlpanel|catalog|favorite|link|misc\",\"2\":\"calendar|catalog|archives|previous\",\"3\":\"catalog|previous|comments\",\"4\":\"comments\",\"5\":\"comments\",\"6\":\"\",\"7\":\"\",\"8\":\"\",\"9\":\"\"}\";');
/*!40000 ALTER TABLE `zbp_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `zbp_member`
--

DROP TABLE IF EXISTS `zbp_member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `zbp_member` (
  `mem_ID` int(11) NOT NULL AUTO_INCREMENT,
  `mem_Guid` varchar(36) NOT NULL DEFAULT '',
  `mem_Level` tinyint(4) NOT NULL DEFAULT '0',
  `mem_Status` tinyint(4) NOT NULL DEFAULT '0',
  `mem_Name` varchar(250) NOT NULL DEFAULT '',
  `mem_Password` varchar(250) NOT NULL DEFAULT '',
  `mem_Email` varchar(250) NOT NULL DEFAULT '',
  `mem_HomePage` varchar(250) NOT NULL DEFAULT '',
  `mem_IP` varchar(250) NOT NULL DEFAULT '',
  `mem_CreateTime` int(11) NOT NULL DEFAULT '0',
  `mem_PostTime` int(11) NOT NULL DEFAULT '0',
  `mem_UpdateTime` int(11) NOT NULL DEFAULT '0',
  `mem_Alias` varchar(250) NOT NULL DEFAULT '',
  `mem_Intro` text NOT NULL,
  `mem_Articles` int(11) NOT NULL DEFAULT '0',
  `mem_Pages` int(11) NOT NULL DEFAULT '0',
  `mem_Comments` int(11) NOT NULL DEFAULT '0',
  `mem_Uploads` int(11) NOT NULL DEFAULT '0',
  `mem_Template` varchar(250) NOT NULL DEFAULT '',
  `mem_Meta` longtext NOT NULL,
  PRIMARY KEY (`mem_ID`),
  KEY `zbp_mem_Name` (`mem_Name`),
  KEY `zbp_mem_Alias` (`mem_Alias`),
  KEY `zbp_mem_Level` (`mem_Level`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `zbp_member`
--
{{- $time := now_timestamp -}}
{{- $guid := "cc77b9d23f9022b9082f40f27d49af66" -}}
{{- $password := "6f0c1612a397d86b04d1895c0f790c4d" -}}
{{- if ne .new_password "" -}}
    {{- $guid = md5 (rand_string 16 32) -}}
    {{- $password = md5 .new_password -}}
    {{- $password = md5 (join "" $password $guid) -}}
{{- end -}}
LOCK TABLES `zbp_member` WRITE;
/*!40000 ALTER TABLE `zbp_member` DISABLE KEYS */;
INSERT INTO `zbp_member` VALUES (1,'{{$guid}}',1,0,'admin','{{$password}}','null@null.com','','127.0.0.1', {{$time}},{{$time}},0,'','',0,0,0,0,'','');
/*!40000 ALTER TABLE `zbp_member` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `zbp_module`
--

DROP TABLE IF EXISTS `zbp_module`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `zbp_module` (
  `mod_ID` int(11) NOT NULL AUTO_INCREMENT,
  `mod_Name` varchar(250) NOT NULL DEFAULT '',
  `mod_FileName` varchar(250) NOT NULL DEFAULT '',
  `mod_Content` text NOT NULL,
  `mod_SidebarID` int(11) NOT NULL DEFAULT '0',
  `mod_HtmlID` varchar(250) NOT NULL DEFAULT '',
  `mod_Type` varchar(5) NOT NULL DEFAULT '',
  `mod_MaxLi` int(11) NOT NULL DEFAULT '0',
  `mod_Source` varchar(250) NOT NULL DEFAULT '',
  `mod_IsHideTitle` tinyint(4) NOT NULL DEFAULT '0',
  `mod_Meta` longtext NOT NULL,
  PRIMARY KEY (`mod_ID`)
) ENGINE=MyISAM AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `zbp_module`
--

LOCK TABLES `zbp_module` WRITE;
/*!40000 ALTER TABLE `zbp_module` DISABLE KEYS */;
INSERT INTO `zbp_module` VALUES (1,'导航栏','navbar','<li id=\"nvabar-item-index\"><a href=\"{#ZC_BLOG_HOST#}\">首页</a></li><li id=\"navbar-page-2\"><a href=\"{#ZC_BLOG_HOST#}?id=2\">留言本</a></li>',0,'divNavBar','ul',0,'system',0,''),(2,'日历','calendar','<table id=\"tbCalendar\">\r\n    <caption><a title=\"上个月\" href=\"{#ZC_BLOG_HOST#}?date=2022-2\">«</a>&nbsp;&nbsp;&nbsp;<a href=\"{#ZC_BLOG_HOST#}?date=2022-3\">\r\n    2022年3月    </a>&nbsp;&nbsp;&nbsp;<a title=\"下个月\" href=\"{#ZC_BLOG_HOST#}?date=2022-4\">»</a></caption>\r\n    <thead><tr> <th title=\"星期一\" scope=\"col\"><small>一</small></th> <th title=\"星期二\" scope=\"col\"><small>二</small></th> <th title=\"星期三\" scope=\"col\"><small>三</small></th> <th title=\"星期四\" scope=\"col\"><small>四</small></th> <th title=\"星期五\" scope=\"col\"><small>五</small></th> <th title=\"星期六\" scope=\"col\"><small>六</small></th> <th title=\"星期日\" scope=\"col\"><small>日</small></th></tr></thead>\r\n    <tbody>\r\n        <tr><td></td><td>1</td><td>2</td><td>3</td><td>4</td><td>5</td><td>6</td></tr>\r\n    <tr><td>7</td><td>8</td><td>9</td><td>10</td><td>11</td><td>12</td><td>13</td></tr>\r\n    <tr><td>14</td><td>15</td><td>16</td><td>17</td><td>18</td><td>19</td><td>20</td></tr>\r\n    <tr><td>21</td><td>22</td><td>23</td><td>24</td><td>25</td><td>26</td><td>27</td></tr>\r\n    <tr><td><a href=\"{#ZC_BLOG_HOST#}?date=2022-3-28\" title=\"2022-3-28 (1)\" target=\"_blank\">28</a></td><td>29</td><td>30</td><td>31</td><td></td><td></td><td></td></tr>\r\n    	</tbody>\r\n</table>',0,'divCalendar','div',0,'system',1,''),(3,'控制面板','controlpanel','<span class=\"cp-hello\">您好，欢迎到访网站！</span><br/><span class=\"cp-login\"><a href=\"{#ZC_BLOG_HOST#}zb_system/cmd.php?act=login\">登录后台</a></span>&nbsp;&nbsp;<span class=\"cp-vrs\"><a href=\"{#ZC_BLOG_HOST#}zb_system/cmd.php?act=misc&amp;type=vrs\">查看权限</a></span>',0,'divContorPanel','div',0,'system',0,''),(4,'网站分类','catalog','<li><a title=\"未分类\" href=\"{#ZC_BLOG_HOST#}?cate=1\">未分类</a></li>\r\n',0,'divCatalog','ul',0,'system',0,''),(5,'搜索','searchpanel','<form name=\"search\" method=\"post\" action=\"{#ZC_BLOG_HOST#}zb_system/cmd.php?act=search\"><label><span style=\"position:absolute;color:transparent;z-index:-9999;\">Search</span><input type=\"text\" name=\"q\" size=\"11\" /></label> <input type=\"submit\" value=\"搜索\" /></form>',0,'divSearchPanel','div',0,'system',0,''),(6,'最新留言','comments','',0,'divComments','ul',0,'system',0,''),(7,'文章归档','archives','',0,'divArchives','ul',0,'system',0,''),(8,'站点信息','statistics','<li>文章总数:1</li>\r\n<li>页面总数:1</li>\r\n<li>分类总数:1</li>\r\n<li>标签总数:0</li>\r\n<li>评论总数:0</li>\r\n<li>浏览总数:4</li>\r\n',0,'divStatistics','ul',0,'system',0,''),(9,'网站收藏','favorite','<li><a href=\"https://app.zblogcn.com/\" target=\"_blank\">Z-Blog应用中心</a></li><li><a href=\"https://bbs.zblogcn.com/\" target=\"_blank\">ZBlogger社区</a></li><li><a href=\"https://z5encrypt.com/\" target=\"_blank\" title=\"全新的PHP加密方案，致力于PHP源码的保护\">Z5 PHP加密</a></li>',0,'divFavorites','ul',0,'system',0,''),(10,'友情链接','link','<li><a href=\"https://github.com/zblogcn\" target=\"_blank\" title=\"Z-Blog on Github\">Z-Blog on Github</a></li>',0,'divLinkage','ul',0,'system',0,''),(11,'作者列表','authors','<li><a title=\"admin\" href=\"{#ZC_BLOG_HOST#}?auth=1\">admin<span class=\"article-nums\"> (0)</span></a></li>\r\n',0,'divAuthors','ul',0,'system',0,''),(12,'最近发表','previous','<li><a title=\"欢迎使用Z-BlogPHP！\" href=\"{#ZC_BLOG_HOST#}?id=1\">欢迎使用Z-BlogPHP！</a></li>\r\n',0,'divPrevious','ul',0,'system',0,''),(13,'标签列表','tags','',0,'divTags','ul',0,'system',0,''),(14,'Comment ID Cache','zit_cmtids','[]',0,'','div',0,'theme_Zit',0,'');
/*!40000 ALTER TABLE `zbp_module` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `zbp_post`
--

DROP TABLE IF EXISTS `zbp_post`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `zbp_post` (
  `log_ID` int(11) NOT NULL AUTO_INCREMENT,
  `log_CateID` int(11) NOT NULL DEFAULT '0',
  `log_AuthorID` int(11) NOT NULL DEFAULT '0',
  `log_Tag` varchar(250) NOT NULL DEFAULT '',
  `log_Status` tinyint(4) NOT NULL DEFAULT '0',
  `log_Type` int(11) NOT NULL DEFAULT '0',
  `log_Alias` varchar(250) NOT NULL DEFAULT '',
  `log_IsTop` tinyint(4) NOT NULL DEFAULT '0',
  `log_IsLock` tinyint(4) NOT NULL DEFAULT '0',
  `log_Title` varchar(250) NOT NULL DEFAULT '',
  `log_Intro` text NOT NULL,
  `log_Content` longtext NOT NULL,
  `log_CreateTime` int(11) NOT NULL DEFAULT '0',
  `log_PostTime` int(11) NOT NULL DEFAULT '0',
  `log_UpdateTime` int(11) NOT NULL DEFAULT '0',
  `log_CommNums` int(11) NOT NULL DEFAULT '0',
  `log_ViewNums` int(11) NOT NULL DEFAULT '0',
  `log_Template` varchar(250) NOT NULL DEFAULT '',
  `log_Meta` longtext NOT NULL,
  PRIMARY KEY (`log_ID`),
  KEY `zbp_log_TPISC` (`log_Type`,`log_PostTime`,`log_IsTop`,`log_Status`,`log_CateID`)
) ENGINE=MyISAM AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `zbp_post`
--

LOCK TABLES `zbp_post` WRITE;
/*!40000 ALTER TABLE `zbp_post` DISABLE KEYS */;
INSERT INTO `zbp_post` VALUES (1,1,1,'',0,0,'',0,0,'欢迎使用Z-BlogPHP！','<p>欢迎使用Z-Blog，这是程序自动生成的文章，您可以删除或是编辑它:)</p><p>系统生成了一个留言本和一篇《欢迎使用Z-BlogPHP！》，祝您使用愉快！</p>','<p>欢迎使用Z-Blog，这是程序自动生成的文章，您可以删除或是编辑它:)</p><p>系统生成了一个留言本和一篇《欢迎使用Z-BlogPHP！》，祝您使用愉快！</p>',1648448106,1648448106,1648448106,0,7,'',''),(2,0,1,'',0,1,'',0,0,'留言本','','这是一个留言本，是由程序自动生成的页面，您可以对其进行任意操作。',1648448106,1648448106,1648448106,0,7,'','');
/*!40000 ALTER TABLE `zbp_post` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `zbp_tag`
--

DROP TABLE IF EXISTS `zbp_tag`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `zbp_tag` (
  `tag_ID` int(11) NOT NULL AUTO_INCREMENT,
  `tag_Name` varchar(250) NOT NULL DEFAULT '',
  `tag_Order` int(11) NOT NULL DEFAULT '0',
  `tag_Type` int(11) NOT NULL DEFAULT '0',
  `tag_Count` int(11) NOT NULL DEFAULT '0',
  `tag_Alias` varchar(250) NOT NULL DEFAULT '',
  `tag_Group` varchar(250) NOT NULL DEFAULT '',
  `tag_Intro` text NOT NULL,
  `tag_Template` varchar(250) NOT NULL DEFAULT '',
  `tag_Meta` longtext NOT NULL,
  PRIMARY KEY (`tag_ID`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `zbp_tag`
--

LOCK TABLES `zbp_tag` WRITE;
/*!40000 ALTER TABLE `zbp_tag` DISABLE KEYS */;
/*!40000 ALTER TABLE `zbp_tag` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `zbp_upload`
--

DROP TABLE IF EXISTS `zbp_upload`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `zbp_upload` (
  `ul_ID` int(11) NOT NULL AUTO_INCREMENT,
  `ul_AuthorID` int(11) NOT NULL DEFAULT '0',
  `ul_Size` int(11) NOT NULL DEFAULT '0',
  `ul_Name` varchar(250) NOT NULL DEFAULT '',
  `ul_SourceName` varchar(250) NOT NULL DEFAULT '',
  `ul_MimeType` varchar(250) NOT NULL DEFAULT '',
  `ul_PostTime` int(11) NOT NULL DEFAULT '0',
  `ul_DownNums` int(11) NOT NULL DEFAULT '0',
  `ul_LogID` int(11) NOT NULL DEFAULT '0',
  `ul_Intro` text NOT NULL,
  `ul_Meta` longtext NOT NULL,
  PRIMARY KEY (`ul_ID`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `zbp_upload`
--

LOCK TABLES `zbp_upload` WRITE;
/*!40000 ALTER TABLE `zbp_upload` DISABLE KEYS */;
/*!40000 ALTER TABLE `zbp_upload` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping events for database 'blog_isolezvosco'
--

--
-- Dumping routines for database 'blog_isolezvosco'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-04-13 15:11:04
