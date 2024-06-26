
BINDATA_CMD = httpsrv-bindata
BINDATA_ARGS_WEBUI = -src webui/ -dst pkg/webui/ -inc hchart.js,chartjs/chart.js,sea.min.js

all: bindata_build
	@echo ""
	@echo "build complete"
	@echo ""

bindata_build:
	# statik -src=webui -include=hchart.js -dest=pkg -ns webui -p webui -f
	# go get -d github.com/rakyll/statik
	# go get -d github.com/hooto/httpsrv-bindata
	# go install github.com/hooto/httpsrv-bindata
	# $(BINDATA_CMD) $(BINDATA_ARGS_WEBUI)
	mkdir -p ./.build_temp/webui/chartjs
	cp -rp webui/hchart.js .build_temp/webui/
	cp -rp webui/sea.min.js .build_temp/webui/
	cp -rp webui/chartjs/chart.js .build_temp/webui/chartjs\/chart.js
	statik -src .build_temp/webui/ -dest pkg/ -p webui -ns webui -f
	rm -fr .build_temp

bindata_clean:
	rm -f pkg/webui/statik.go

install:
	install -m 755 ${EXE_CLI} ${APP_PATH}

clean: bindata_clean
	@echo ""
	@echo "clean complete"
	@echo ""
