# current directory
CUR_DIR ?= $(shell pwd)

clean:
	@echo "Clean web generated code ...."
	@rm -rf rest/swagger-js
	@rm -rf rest/swagger-client*

build:
	@echo "Generate web code ...."
	@mkdir -p rest/swagger-js && git clone https://github.com/swagger-api/swagger-js.git $(CUR_DIR)/rest/swagger-js
	@cd $(CUR_DIR)/rest/swagger-js && npm install && npm run build:bundle
	@cd $(CUR_DIR)/rest/swagger-js/browser && mv index.js swagger-client.min.js && mv index.js.map swagger-client.min.js.map
	@cd $(CUR_DIR)/rest/swagger-js && npm install && npm run build:bundle:dev
	@cd $(CUR_DIR)/rest/swagger-js/browser && mv index.js swagger-client.js && mv index.js.map swagger-client.js.map
	@cd $(CUR_DIR)/rest/swagger-js && mv browser/* ../

browser: build
	@rm -rf rest/swagger-js
