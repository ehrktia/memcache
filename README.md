### memcache

in memory cache via http api

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
cd coordinator && zig run src/main.zig  
```

this will start coordinator and emit 5 consecutive multicast messages with port number available for communication  

testing  
use the test implementation in go under `coordinator/test-udp` for testing this.  
