Creative Information Manager
===

* Receive get request with a creative url and return the creative id of it


Runtime Environment
---

* golang (see ways of installation as follows)

  * centOS: `yum install golang`

  * ubuntu: `apt-get install golang`

  * macOS: `brew install golang`


Dependecy Installation
---

    make deps



Example
---

** Attention: creative url should be escaped/encoded, for example, in golang, you should use QueryEscape of net/url package

    # get a creative id of http://cdn.image2.cloudmobi.net/static/image/1000/1000/1501680592.jpg
    curl "http://127.0.0.1:12121/get_creative_id?creative_url=http%3A%2F%2Fcdn.image2.cloudmobi.net%2Fstatic%2Fimage%2F1000%2F1000%2F1501680592.jpg"

