# GeoIP2

The GeoIP2 Filter plugin looks up IP addresses of a record in a MaxMind GeoIP2 or GeoLite2 database <br /> and appends the location information, such as the country or the city, to the record. <br /> **For full documentation, refer to https://docs.fluentbit.io/manual/pipeline/filters/geoip2-filter**


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| database | Path to the MaxMind database file (.mmdb), e.g. /fluent-bit/geoip/GeoLite2-City.mmdb. The file can be mounted into the Fluent Bit pod with the volumes and volumeMounts of the FluentBit resource. | string |
| lookupKeys | Keys of the record whose value is an IP address to look up. | []string |
| records | Keys to append to the record, with the value looked up in the database. | [][GeoIP2Record](#geoip2record) |
# GeoIP2Record

GeoIP2Record appends Key to the record, with the value found at Value in the database for the IP address in LookupKey.


| Field | Description | Scheme |
| ----- | ----------- | ------ |
| key | Key appended to the record. | string |
| lookupKey | Lookup key holding the IP address to look up. It must be one of lookupKeys, otherwise the value is null. | string |
| value | Path of the value in the database, written as %{path}, e.g. %{country.names.en} or %{location.latitude}. The path is at most 63 characters long. | string |
