This is demo app using BitGo API

### to say hello bitgo 
curl -X GET http://localhost:8081

### to initialise bitgo api
#### you need to get a BitGo API token
curl -X POST http://localhost:8081

### to get tbtc address
curl -X GET http://localhost:8081/address

### to send tbtc to destination
curl -X POST http://localhost:8081/send/2NCykwAEJnbktzbyxguAZo7qHCawGqduJYs
