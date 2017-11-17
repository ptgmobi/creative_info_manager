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

    # get a creative id
    curl "http://127.0.0.1:12121/get_creative_id?creative_url=http://cdn.image2.cloudmobi.net/static/image/1000/1000/1501680592.jpg"

