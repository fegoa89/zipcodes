# zipcodes - Zip Code Lookups

A Zipcode lookup package that uses the GeoNames Postal Code dataset from http://www.geonames.org .
You can initialize it with a Postal Code dataset downloaded from http://download.geonames.org/export/zip .

## Install

Install with
```sh
go get github.com/fegoa89/zipcodes
```

### Initialize Struct
Initializes a zipcodes struct. It will throw an error if:
- The file does not exist / wrong format.
- Some of the lines contain less that 12 elements (in the readme.txt of each postal code dataset, they define up to 12 elements).
- Where latitude / longitude value are contains a wrong format (string that can not be converted to `float64`).

```golang
zipcodesDataset, err := zipcodes.New("path/to/my/dataset.txt")
```

### Lookup
Looks for a zipcode inside the map interface we loaded. If the object can not be found by the zipcode, it will return an error. 
When a object is found, returns its zipcode, place name, administrative name, latitude and longitude:

```golang
location, err := zipcodesDataset.Lookup("10395")
```

### DistanceInKm
Returns the line of sight distance between two zipcodes in kilometers:

```golang
location, err := zipcodesDataset.DistanceInKm("01945", "03058") // 49.87
```

### DistanceInMiles
Returns the line of sight distance between two zipcodes in miles:

```golang
location, err := zipcodesDataset.DistanceInMiles("01945", "03058") // 30.98
```

### DistanceInKmToZipCode
Calculates the distance between a zipcode and a give lat/lon in Kilometers:

```golang
location, err := zipcodesDataset.DistanceInKmToZipCode("01945", 51.4267, 13.9333) // 1.11
```

### DistanceInMilToZipCode
Calculates the distance between a zipcode and a give lat/lon in Miles:

```golang
location, err := zipcodesDataset.DistanceInMilToZipCode("01945", 51.4267, 13.9333) // 0.69
```

### GetZipcodesWithinKmRadius
Returns a list of zipcodes within the radius of this zipcode in Kilometers:

```golang
location, err := zipcodesDataset.GetZipcodesWithinKmRadius("01945", 50) // ["03058"]
```

### GetZipcodesWithinMlRadius
Returns a list of zipcodes within the radius of this zipcode in Miles:

```golang
location, err := zipcodesDataset.GetZipcodesWithinMlRadius("01945", 50) // ["03058"]
```
