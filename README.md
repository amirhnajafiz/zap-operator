# :zap: Zap Operator

```Zap``` is a blazing fast, structured, leveled logging in Go. Our service's logs are in JSON format.
In this project I create an operator to get these logs from OKD's pods ```stdout```, and publish
them over the following topics over a ```NATS``` cluster.

1. ```{service-name}.logs.debug```
2. ```{service-name}.logs.warning```
3. ```{service-name}.logs.info```
4. ```{service-name}.logs.error```
5. ```{service-name}.logs.unknown```

## :flashlight: setup

Use the ```amirhossein21/zapoperator:latest``` image to make a deployment of this operator on your cluster. Make sure to set a ```app```
label in the deployment that you want this operator to observe. After that, set the ```DEPLOYMENT``` environment variable
to the value that you set for ```app``` label.

Also make sure to set the following environmental variables:

- ```NAMESPACE``` : the namespace where you are deploying this operator
- ```NATS_HOST``` and ```NATS_TOPIC``` : information of your ```NATS``` cluster
