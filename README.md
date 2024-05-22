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

**TODO**

- [ ] benchmark and test for disaster recovery
- [ ] implement a leader process to have data sync between instances  
