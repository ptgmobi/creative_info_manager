
GOBUILD = go build
GOTEST = go test
GOGET = go get -u

VARS=vars.mk
$(shell ./build_config ${VARS})
include ${VARS}

.PHONY: main clean test

main:
	${GOBUILD} -o bin/creative_info_manager src/main.go

deps:
	${GOGET} github.com/brg-liuwei/gotools
	${GOGET} github.com/go-sql-driver/mysql
	${GOGET} github.com/garyburd/redigo/redis

test:

clean:
	@rm bin/creative_info_manager
