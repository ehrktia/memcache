### memcache

in memory cache via http api   


<video src="docs/demo.mp4" type="video/mp4" controls width="800"></video>


**start locally**   

```sh
COORDINATOR=true go run cmd/main.go
```




#### endpoints

`/save`

to save data in to cache  
data format

```json
{
    "key":"some-key",
    "value":"some-value"
}
```

req type - **POST**

`/get`

to retrieve data from cache

data format

```json
{
    "key":"some-key"
}
```

#### about

this is a concurrent implementation of cache via http.
Can be used across any service layer as a stage or landing for some data
which you require to consume / refer / lookup

#### distributed

**WIP**

Coordinator  
Will be used to manage multiple instance of cache and sync data between instances. Available in `coordinator`

how to start  
```sh
cd coordinator && zig build run  
```

this will start coordinator and emit 5 consecutive multicast messages with port number available for communication  

**container**   
To build image from root of the project   

```sh
podman build -t coordinator:latest --squash-all -f coordinator/Dockerfile ./coordinator/
```

To run the container 

```sh
podman run --name coordinator --publish-all -h coordinator -d localhost/coordinator:latest
```




testing  
- run the local version of coordinator and use the go run command to test
