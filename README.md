## Build the Server ##

In case you don't have environment that support make command, 
please follow instruct [Pulling Pre-Build Image](#pulling-pre-build-image).

To build the server, first you must have docker up and running on your local machine,
then follow the below instruction.

- To build docker image server with default tag, run the following command
```bash
make build
```
- To build docker image server with specific tag, run the following command, replace `1.0.0` as your need
```bash
make TAG='1.0.0' build
```

Additional command:

- To clean up generated code and other file for server side, run the following command
```bash
make clean
```
- To generate code and file only, run the following command
```bash
make generate
``` 

## Pulling Pre-Build Image ##

To pre-build docker image you must first have your docker up and running on your local
machine, then run the following command:

```bash
docker pull cambodia-demo/envoy-grpc-sample
```

or with a specific tag

```bash
docker pull cambodia-demo/envoy-grpc-sample:1.0.0
```

## Running Server ##

Before you can running the server, first you must ensure that docker is running on your
local machine and you must already execute one of the following the step 
[Build the Server](#build-the-server) or [Pulling Pre-Build Image](#pulling-pre-build-image) above.
When you already done the above step then following the below instruction

If you build image locally on your machine, please update tag if you happen to build with custom tag:

```bash
docker run -d --name envoy-grpc-sample \
 -p 8080:8080 \
 -p 8090:8090 \
 -p 9901:9901 \
 envoy-grpc:1.0.0
```

The container listen and expose on 3 ports:

- envoy is listening on port 8080, the request either grpc or rest api can send to envoy on port 8080. Envoy will be responsible to convert rest into GRPC if it's rest api request. 
- the sample server is listening on port 8090 and it only accept GRPC request. Client can grpc request directly to the container via port 8090.
- envoy admin is listening on port 9901.

## Testing ##

#### Test Rest API with Tools ####

To test the result of the api you can use curl command if other tools that support http rest api request.

Post element data to the server:
```bash
curl \
  -d '{"name":"test 1", "age":18, "status": 4}' \
  http://localhost:8080/rest/v1.0/element
```

List element data from the server:
```bash
curl http://localhost:8080/rest/v1.0/element/list
```

or apply filter, example list only element where age is between 10 and 20:
```bash
curl http://localhost:8080/rest/v1.0/element/list?age=[10,20]
```

#### Test with Go client ####

The Go client is sent the grpc request direct thus envoy does not process any conversion but
forward the request to the upstream.

- Build the client command
```bash
make generate-client-go
```
- To list the element run the command, the example command below filter the element with age in range from 10 to 30
```bash
clients/go/sample -command list -filter '{"age":"[10,30]"}'
```
- To save the element run the command below
```bash
clients/go/sample -command save -data '{"name":"test 11", "age": 33, "status": 7}'
```

#### Test with Java client ####

You use IntelliJ IDE to open the project at `clients/java`. However the below instruct is provide for command only. 

- Build an executable jar file with make command
```bash
make generate-client-java CLIENT='java'
```
After the build success, you will see the executable jar file name `example-1.0-all.jar` under the folder `clients/java`

If you don't environment to support make command (Example MS Windows). You can run the gradle command below directly.
Make sure you are in the directory `clients/java/example` in your terminal/command prompt. 

For Linux user
```bash
./gradlew clean && ./gradlew build && ./gradlew shadowJar
```
For Window user
```batch
gradlew.bat clean && gradlew.bat build && gradlew.bat shadowJar
```

The executable file name `example-1.0-all.jar` is generated and available at the folder `clients/java/example/build/libs`

- Save element
```bash
java -jar clients/java/example-1.0-all.jar -c save -d '{"name":"PP 11", "age": 23, "status": 2}'
``` 

- List element
```bash
java -jar clients/java/example-1.0-all.jar -c list -f '{"age":"[30,38]"}'
```

Note: If you compile and build executable jar file without make command then make
sure to update the location of jar file `example-1.0-all.jar` accordingly. 

Additional Command:

Under the folder `clients/java`, there is also a makefile for managing java client code. It's provide the following command

- To clean the project run
```bash
make clean
```
- To build and generate an executable jar file
```bash
make build && make shadowJar
```

#### Test with Android client ####

To test with Android project, you're required to have Android Studio IDE. 
Open the project locate at `clients/android/example` with Android Studio.
The example is locate in the test file name `ExampleGrpcTest` at the package `com.example.grpcexample`.
You can execute the test file without required emulator or device. 

#### Test with iOS client ####

To test with iOS project, you're required MacOS and XCode IDE.
Open the project locate at `clients/ios/example` with XCode. Make sure to open
`.xcworkspace` instead of `.xcodeproj`

Before you can start running and testing grpc in swift code, first you need to
follow the instruction describe [here](https://github.com/grpc/grpc-swift#getting-the-plugins)
to build and install the plugin. On you done that, follow the below instruction:

- To generate swift protobuf and grpc file run the following command
```bash
make generate-client-ios
```
- Add the 2 generated file into project `sample.grpc.swift` and `sample.pb.swift`
- Install dependencies with pod
```bash
pod install
```
- The example code is locate at target `Example GrpcTests` in the file `Example_GrpcTests.swift`. Use Xcode IDE to run and test the code.

#### Test with Web client ####

- If you would like to compile swagger js client, run the following command
```bash
make generate-client-web-compile
```
- If you would like to only generate swagger definition, run the following command
```bash
make generate-client-web
```

Note: Normally you don't need to compile swagger js client all the time unless
there are new update available from Swagger JS community.

One you have done above step, use any web server to serve html document from the
directory `clients/web/rest`. In the browser turn on debug tools to see the console
log from javascript code then browser to your web server address, for example: `http://localhost:8081`

If you already have `npm` installed, you can follow below step:

- Install http-server, [https://www.npmjs.com/package/http-server](https://www.npmjs.com/package/http-server)
```sybase
npm install http-server -g
```
- Change directory to `clients/web/rest`
- Start the server
```bash
http-server -p 8081
```
- From your browser, open a new tab and enable developer tools
- Browser to the address http://localhost:8081