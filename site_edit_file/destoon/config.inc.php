<?php
defined('IN_DESTOON') or exit('Access Denied');
/*
	[DESTOON B2B System] Copyright (c) 2008-2018 www.destoon.com
	This is NOT a freeware, use is subject to license.txt
*/
$CFG['database'] = 'mysqli';
$CFG['pconnect'] = '0';
$CFG['db_host'] = '127.0.0.1';
$CFG['db_name'] = '{{.NameDB}}';
$CFG['db_user'] = '{{.DbUser}}';
$CFG['db_pass'] = '{{.DbPass}}';
$CFG['db_charset'] = 'utf8';
$CFG['db_expires'] = '0';
$CFG['tb_pre'] = 'destoon_';
$CFG['charset'] = 'utf-8';
$CFG['url'] = 'http://{{.Domain}}/';
$CFG['com_domain'] = '';
$CFG['com_dir'] = '1';
$CFG['com_rewrite'] = '0';
$CFG['com_vip'] = 'VIP';
$CFG['file_mod'] = 0777;
$CFG['cache'] = 'file';
$CFG['cache_pre'] = 'cij_';
$CFG['cache_dir'] = '';
$CFG['cookie_domain'] = '';
$CFG['cookie_path'] = '/';
$CFG['cookie_pre'] = 'cxd_';
$CFG['session'] = 'file';
$CFG['cdn'] = '0';
$CFG['editor'] = 'fckeditor';
$CFG['timezone'] = 'Etc/GMT-8';
$CFG['timediff'] = '0';
$CFG['skin'] = 'default';
$CFG['template'] = 'default';
$CFG['template_mobile'] = 'mobile';
$CFG['template_refresh'] = '1';
$CFG['language'] = 'zh-cn';
$CFG['authadmin'] = 'session';
$CFG['authkey'] = 'mKaae5cstnWcHInc';
$CFG['static'] = '';
$CFG['cloud_uid'] = '';
$CFG['cloud_key'] = '';
$CFG['edittpl'] = '0';
$CFG['executesql'] = '0';
$CFG['founderid'] = '1';
?>