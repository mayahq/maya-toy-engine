# InfluxDB components for Cascades FBP

## Usage

```
$ ./components/influxdb/write
Usage of ./components/influxdb/write:
  -debug=false: Enable debug mode
  -json=false: Print component documentation in JSON
  -port.err="": Component's output port endpoint
  -port.in="": Component's input port endpoint
  -port.options="": Component's input port endpoint
```

## Example

```
# Send options IP
'host=127.0.0.1:8086,user=root,password=root,db=demo' -> OPTIONS MyInfluxDB(influxdb/write)

# Connect errors port to console
MyInfluxDB ERR -> IN Errors(core/console)

# Send test data as IPs
'{"name":"tests","columns":["id","name","value"],"points":[[1,"Alex",true]]}' -> IN MyInfluxDB
'{"name":"tests","columns":["id","name","value"],"points":[[2,"Bob",false]]}' -> IN MyInfluxDB
'{"name":"tests","columns":["id","name","value"],"points":[[3,"Dick",true]]}' -> IN MyInfluxDB
```

