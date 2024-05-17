
# Heimdall (the api gateway)

Heimdall is a lightweight API Gateway built for anyone who wants to use a simple gateway.

Heimdall core functionalities : 
- load balancing
- support for HTTP and GRPC
- request body simple checking
- self-healing by circuit-breaking and health-checking



## Deployment

Change the config.yaml file to your prefernece. For more info about fields of config, check the Documentation section.

### Add your HTTP APIs
 The only thing to do is to update the config.yaml.

### Add your GRPC APIs
 You need to do four more steps. 

-  Add your proto files to the Proto folder
-  Build your Protos by buildProto.sh located in the Proto folder.
 ```bash
  $ ./Proto/buildProto.sh
```    
- Select a name for your GPRC service and assign it in internal/proxy/grpc/grpcInfo.go file. 
```go
package proxyGrpc

type HeimdallGrpcService string

const (
   HeimdallGrpcService_YOUR_SERVICE_NAME HeimdallGrpcService = "YOUR_SERVICE_NAME"
   // Note that this name should be the same as the name you have set in config.yaml file for this API.
)
```
- Render a GRPC client to connect to your service by modifying the method "establishConnection" located in internal/proxy/grpc/grpcProxy.go file.

```go
func (g *grpcProxy) establishConnection(identifier HeimdallGrpcService, host string, mux *runtime.ServeMux) error {

  	//You just need to add the below Case to your code!

   switch identifier {
   case HeimdallGrpcService_YOUR_SERVICE_NAME:
      err = golang.RegisterYourServiceHandlerFromEndpoint(ctx, mux, host, opts)
   }
   return err
}
```

### depoly

First change the config.yaml file as described earlier.

then select one ot these two approaches:

* Run Golang

```bash
  go run ./main.go
```    

* Run by docker

    Build the docker image by running the dockerFile and start the service.

```bash
  $ docker build -f ./Dockerfile -t heimdal:1.0.6 .

  $ docker-compose up
```
## Documentation

The only part you must know to start using the
field name | optional |  describe
---------- | --------- | ----------
heimdall_port | false |     port which API Gateway will listen to. By default  80  
apis_config| false | list of Apis Policy
match_policy | false | describe each API
apis_config.match_policy.connection_type | true |    HTTP or GRPC API. select between "http" and "grpc". (default is http) 
apis_config.match_policy.name | false | name of this app. For GRPC APIs you have to set the same name as the name you have selected in internal/proxy/grpcProxy.go 
apis_config.match_policy.path | false | checks if the request does belong to this API. You can use /*any for accepting any sub-path.
apis_config.match_policy.per_method | false |  config for each separate HTTP method.
apis_config.match_policy.per_method.type | false |  specified HTTP method type. Multiple methods could be selected at once. Select among "post", "get", "put" and ... 
apis_config.match_policy.per_method.request_body_check | true | policy for checking each request received in Heimdall. request will be sent to specified API only if it passes this checks
apis_config.match_policy.per_method.request_body_check.mandatory_fields | true | list of fields that must be presented in request body
apis_config.match_policy.per_method.request_body_check.mandatory_fields.field_name | false | name of body field
apis_config.match_policy.per_method.request_body_check.mandatory_fields.type | false | type of field. select among ("string", "float64", "bool", "map", "slice")
apis_config.load_balance | false | policy for load balancing and specifies the hosts. By default round-robin is enabled 
apis_config.load_balance.type | true |   load balancing logic. select between "weighted_round_robin" or "round_robin". by default "round_robin" is enabled 
apis_config.load_balance.host_units | false |   hosts ip and theire weight if "weighted_round_robin" is selected. 
apis_config.load_balance.host_units.host | false |   host IP of your service.
apis_config.load_balance.host_units.load_balance_weight | true |   If "weighted_round_robin" is selected, set the weight of this host by this field.
apis_config.health_check_config | true |   health check policy for this API. 
apis_config.health_check_config.path | false | This path will be called for all the hosts you have mentioned at apis_config.load_balance.host_units for this api_config. 
apis_config.health_check_config.failure_threshold | true | consecutive failures threshold to count this host as down. By default, the failure_threshold is set to 3
apis_config.health_check_config.interval | true |   interval of calling the health URL of the host. By default, the interval is set to 1s.
apis_config.circuit_breaker_config | true | specifies the circuit break policy for host 
apis_config.circuit_breaker_config.quarantine_duration | true |   duration that the host will be removed from available hosts and won't receive new inputs.
apis_config.circuit_breaker_config.failure_tolerance_count | true |  count of consecutive failures to count this host as down 

## Usage/Examples

The sample config:

```yaml
heimdall_port: 80
apis_config:
  - match_policy:
      connection_type: "http"
      name: "your_http_service_name"
      path: "/path/ro/your/serivce/v1/method"
      per_method:
        - type: "GET"
        - type: "HEAD"
        - type: "POST"
          request_body_check:
            mandatory_fields:
              - field_name: "field1"
                type: "string"
              - field_name: "field2"
                type: "bool"
              - field_name: "field3"
                type: "float64"
              - field_name: "field4"
                type: "map"
              - field_name: "field5"
                type: "slice"
    load_balance:
      type: "weighted_round_robin"
      host_units:
        - host: "http://192.168.1.118:8080"
          load_balance_weight: 4
        - host: "http://192.168.1.118:4040"
          load_balance_weight: 2
        - host: "http://192.168.1.118:4043"
          load_balance_weight: 2
        - host: "http://192.168.1.118:80"
          load_balance_weight: 2
    health_check_config:
      path: "/health"
      failure_threshold: 3
      interval: "10s"
    circuit_breaker_config:
      quarantine_duration: "30s"
      failure_tolerance_count: 2


  - match_policy:
      connection_type: "grpc"
      name: "YOUR_SERVICE_NAME"
      path: "/PROTO_GENERATED_ADDRESS/*any" # like /backend.messaging.v1.MessagingService/*any
      per_method:
        - type: "POST"
    load_balance:
      type: "round_robin"
      host_units:
        - host: "192.168.1.118:50503"
        - host: "192.168.1.118:50501"
    circuit_breaker_config:
      examine_window: "1m"
      quarantine_duration: "5s"
      failure_tolerance_count: 2

```

