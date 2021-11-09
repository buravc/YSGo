# YSGo
YS Go Task I was assigned to.

# Todo:

* read cron timer length from env var (DONE, "GodisBackupInterval" env variable is now being used with options being 1h, 15m and 10s)
* env variable processing needs to be done via regex (to allow more flex configurations for cron timer)
* write an http server
* complete unit tests
* prepare docker-compose and api doc file (swagger)
* build actions
* an isolation layer between httpServer and godisServer
* deploy to aws (or azure)
* make README have the usage of this lib. (endpoints, build etc.)
