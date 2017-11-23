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
  `create_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY (`url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='creative_info';
```

```
+-------------+--------------+------+-----+-------------------+----------------+
| Field       | Type         | Null | Key | Default           | Extra          |
+-------------+--------------+------+-----+-------------------+----------------+
| id          | int(11)      | NO   | PRI | NULL              | auto_increment |
| url         | varchar(512) | NO   | UNI |                   |                |
| type        | smallint(4)  | NO   |     | 0                 |                |
| create_date | datetime     | NO   |     | CURRENT_TIMESTAMP |                |
+-------------+--------------+------+-----+-------------------+----------------+
```

### Redis

```
+---------------+--------------+-------+
| Key           |  Field       | Value |
+---------------+--------------+-------+
| creative_info |  url         | id    | 
+---------------+--------------+-------+
```


Example
---

*** Attention: creative url should be escaped/encoded, for example, in golang, you should use QueryEscape of net/url package

    # get a creative id of http://cdn.image2.cloudmobi.net/static/image/1000/1000/1501680592.jpg
    curl "http://127.0.0.1:12121/get_creative_id?creative_url=http%3A%2F%2Fcdn.image2.cloudmobi.net%2Fstatic%2Fimage%2F1000%2F1000%2F1501680592.jpg"
    
    #sample response:
    {
        err_msg: "",
        creative_id: "1151"
    }


