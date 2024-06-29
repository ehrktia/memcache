#!/usr/bin/bash
for i in {1..100}
do
    key="key-"$i
    value="value-"$i
    data="{
    \"key\": \"${key}\", 
    \"value\": \"${value}\" 
}"
# echo $data
    curl --location 'http://localhost:8080/save' \
    --data "${data}"
    sleep 1
done

# curl --location 'http://localhost:8080/save' \
# --header 'Content-Type: application/json' \
# --data '{
    # "key": "some-key",
    # "value": "some-value"
# }'
