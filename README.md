Creative Information Manager
===

* Receive http get request with a creative url and return the creative id of it


Runtime Environment
---

* golang (see ways of installation as follows)

  * centOS: `yum install golang`

  * ubuntu: `apt-get install golang`

  * macOS: `brew install golang`


Dependecy Installation
---

    make deps


Database Information
---

### Mysql

```sql
DROP TABLE IF EXISTS `creative_info`;

CREATE TABLE `creative_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `url` varchar(512) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '素材链接',
  `type` smallint(4) NOT NULL DEFAULT '0' COMMENT '素材类型 1:图片、2:视频',
  `size` bigint(11) NOT NULL DEFAULT '0' COMMENT '素材文件大小',
  `fail_times` smallint(4) NOT NULL DEFAULT '0' COMMENT '获取文件大小的失败次数',
  `create_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY uniq_url(`url`),
  KEY idx_size(`size`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='creative_info';
```

```
+-------------+--------------+------+-----+-------------------+----------------+
| Field       | Type         | Null | Key | Default           | Extra          |
+-------------+--------------+------+-----+-------------------+----------------+
| id          | int(11)      | NO   | PRI | NULL              | auto_increment |
| url         | varchar(512) | NO   | UNI |                   |                |
| type        | smallint(4)  | NO   |     | 0                 |                |
| size        | bigint(11)   | NO   | MUL | 0                 |                |
| fail_times  | smallint(4)  | NO   |     | 0                 |                |
| create_date | datetime     | NO   |     | CURRENT_TIMESTAMP |                |
+-------------+--------------+------+-----+-------------------+----------------+
```

### Redis

```
+---------------+--------------+
| Key           |  Value       | 
+---------------+--------------+
| url           |  id_size     | 
+---------------+--------------+
```


API 
===

* [get_creative_id](#get_creative_id) 
 
* [dump](#dump)


get_creative_id
--- 
 
* Description：get a creative id of given creative_url, which should better be escaped/encoded, for example, in golang, you should use QueryEscape of net/url package

* URL: http://54.255.165.222:12121/get_creative_id?creative_url={creative_url}
 
* Sample Response:
 
  ```
  {
      err_msg: "",
      creative_id: "img.4398",
      size: 107454
  }
  ``` 

[Back to TOC](#API) 


dump
--- 

* Description：get creative infos of given ids, which must be Comma separated, if there's no ids, returns 10 random creative urls

* URL: http://54.255.165.222:12345/dump?id={id1,id2,etc}
 
* Sample Response:
 
  ```
  {
      err_msg: "",
      creative_info: [
          {
              id: "17",
              url: "http://cdn.mvideo.cloudmobi.net/upload-files/view?path=7414443ebd13c933cb8e861ca37aabc2",
              type: 2,
              size: 5979835
          },
          {
              id: "69",
              url: "http://cdn.mvideo.cloudmobi.net/upload-files/view?path=5d4ca7f06d082e51f572f8c7b9d0a1bb",
              type: 2,
              size: 3661635
          },
          {
              id: "405",
              url: "http://cdn.mvideo.cloudmobi.net/upload-files/view?path=0ce22b68cc5fba6f5a68e25fd01b4bd7",
              type: 1,
              size: 23909
          },
          {
              id: "406",
              url: "http://cdn.mvideo.cloudmobi.net/upload-files/view?path=954e790039709ff0d9e738385c96e507",
              type: 1,
              size: 110706
          }
      ]
  }
  ``` 

[Back to TOC](#API) 
