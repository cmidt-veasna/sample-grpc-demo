# current directory
CUR_DIR ?= $(shell pwd)

clean:
	@echo "Clean java generated code ...."
	@cd example && ./gradlew clean
	@rm -f $(CUR_DIR)/example/src/main/java/com/example/*.java
	@rm -f $(CUR_DIR)/example/src/main/grpc/com/example/*.java
	@rm -f $(CUR_DIR)/*.jar

build: clean
	@echo "Generate java code ...."
	@cd example && ./gradlew build && ./gradlew shadowJar && mv build/libs/example-1.0-all.jar ../ && ./gradlew clean