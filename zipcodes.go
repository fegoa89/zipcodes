// zipcodes is a package that uses the GeoNames Postal Code dataset from http://www.geonames.org
// in order to perform zipcode lookup operations
package zipcodes

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	earthRadiusKm = 6371
	earthRadiusMi = 3958
)

// ZipCodeLocation struct represents each line of the dataset
type ZipCodeLocation struct {
	ZipCode   string
	PlaceName string
	AdminName string
	Lat       float64
	Lon       float64
	StateCode string
}

// Zipcodes contains the whole list of structs representing
// the zipcode dataset
type Zipcodes struct {
	DatasetList map[string]ZipCodeLocation
}

// New loads the dataset that this packages uses and
// returns a struct that contains the dataset as a map interface
func New(datasetPath string) (*Zipcodes, error) {
	zipcodes, err := LoadDataset(datasetPath)
	if err != nil {
		return nil, err
	}
	return &zipcodes, nil
}

// Lookup looks for a zipcode inside the map interface
func (zc *Zipcodes) Lookup(zipCode string) (*ZipCodeLocation, error) {
	foundedZipcode := zc.DatasetList[zipCode]
	if (foundedZipcode == ZipCodeLocation{}) {
		return &ZipCodeLocation{}, fmt.Errorf("zipcodes: zipcode %s not found !", zipCode)
	}
	return &foundedZipcode, nil
}

// DistanceInKm returns the line of sight distance between two zipcodes in Kilometers
func (zc *Zipcodes) DistanceInKm(zipCodeA string, zipCodeB string) (float64, error) {
	return zc.CalculateDistance(zipCodeA, zipCodeB, earthRadiusKm)
}

// DistanceInMiles returns the line of sight distance between two zipcodes in Miles
func (zc *Zipcodes) DistanceInMiles(zipCodeA string, zipCodeB string) (float64, error) {
	return zc.CalculateDistance(zipCodeA, zipCodeB, earthRadiusMi)
}

// CalculateDistance returns the line of sight distance between two zipcodes in Kilometers
func (zc *Zipcodes) CalculateDistance(zipCodeA string, zipCodeB string, radius float64) (float64, error) {
	locationA, errLocA := zc.Lookup(zipCodeA)
	if errLocA != nil {
		return 0, errLocA
	}

	locationB, errLocB := zc.Lookup(zipCodeB)
	if errLocB != nil {
		return 0, errLocB
	}

	return DistanceBetweenPoints(locationA.Lat, locationA.Lon, locationB.Lat, locationB.Lon, radius), nil
}

// DistanceInKmToZipcode calculates the distance between a zipcode and a give lat/lon in Kilometers
func (zc *Zipcodes) DistanceInKmToZipCode(zipCode string, latitude, longitude float64) (float64, error) {
	location, errLoc := zc.Lookup(zipCode)
	if errLoc != nil {
		return 0, errLoc
	}

	return DistanceBetweenPoints(location.Lat, location.Lon, latitude, longitude, earthRadiusKm), nil
}

// DistanceInMilToZipcode calculates the distance between a zipcode and a give lat/lon in Miles
func (zc *Zipcodes) DistanceInMilToZipCode(zipCode string, latitude, longitude float64) (float64, error) {
	location, errLoc := zc.Lookup(zipCode)
	if errLoc != nil {
		return 0, errLoc
	}

	return DistanceBetweenPoints(location.Lat, location.Lon, latitude, longitude, earthRadiusMi), nil
}

// GetZipcodesWithinKmRadius get all zipcodes within the radius of this zipcode
func (zc *Zipcodes) GetZipcodesWithinKmRadius(zipCode string, radius float64) ([]string, error) {
	zipcodeList := []string{}
	location, errLoc := zc.Lookup(zipCode)
	if errLoc != nil {
		return zipcodeList, errLoc
	}

	return zc.FindZipcodesWithinRadius(location, radius, earthRadiusKm), nil
}

// GetZipcodesWithinMlRadius get all zipcodes within the radius of this zipcode
func (zc *Zipcodes) GetZipcodesWithinMlRadius(zipCode string, radius float64) ([]string, error) {
	zipcodeList := []string{}
	location, errLoc := zc.Lookup(zipCode)
	if errLoc != nil {
		return zipcodeList, errLoc
	}

	return zc.FindZipcodesWithinRadius(location, radius, earthRadiusMi), nil
}

// FindZipcodesWithinRadius finds zipcodes within a given radius
func (zc *Zipcodes) FindZipcodesWithinRadius(location *ZipCodeLocation, maxRadius float64, earthRadius float64) []string {
	zipcodeList := []string{}
	for _, elm := range zc.DatasetList {
		if elm.ZipCode != location.ZipCode {
			distance := DistanceBetweenPoints(location.Lat, location.Lon, elm.Lat, elm.Lon, earthRadius)
			if distance < maxRadius {
				zipcodeList = append(zipcodeList, elm.ZipCode)
			}
		}
	}

	return zipcodeList
}

func hsin(t float64) float64 {
	return math.Pow(math.Sin(t/2), 2)
}

// degreesToRadians converts degrees to radians
func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// DistanceBetweenPoints returns the distance between two lat/lon
// points using the Haversin distance formula.
func DistanceBetweenPoints(latitude1, longitude1, latitude2, longitude2 float64, radius float64) float64 {
	lat1 := degreesToRadians(latitude1)
	lon1 := degreesToRadians(longitude1)
	lat2 := degreesToRadians(latitude2)
	lon2 := degreesToRadians(longitude2)
	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := hsin(diffLat) + math.Cos(lat1)*math.Cos(lat2)*hsin(diffLon)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := c * radius

	return math.Round(distance*100) / 100
}

// LoadDataset reads and loads the dataset into a map interface
func LoadDataset(datasetPath string) (Zipcodes, error) {
	file, err := os.Open(datasetPath)
	if err != nil {
		log.Fatal(err)
		return Zipcodes{}, fmt.Errorf("zipcodes: error while opening file %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	zipcodeMap := Zipcodes{DatasetList: make(map[string]ZipCodeLocation)}
	for scanner.Scan() {
		splittedLine := strings.Split(scanner.Text(), "\t")
		if len(splittedLine) != 12 {
			return Zipcodes{}, fmt.Errorf("zipcodes: file line does not have 12 fields")
		}
		lat, errLat := strconv.ParseFloat(splittedLine[9], 64)
		if errLat != nil {
			return Zipcodes{}, fmt.Errorf("zipcodes: error while converting %s to Latitude", splittedLine[9])
		}
		lon, errLon := strconv.ParseFloat(splittedLine[10], 64)
		if errLon != nil {
			return Zipcodes{}, fmt.Errorf("zipcodes: error while converting %s to Longitude", splittedLine[10])
		}

		zipcodeMap.DatasetList[splittedLine[1]] = ZipCodeLocation{
			ZipCode:   splittedLine[1],
			PlaceName: splittedLine[2],
			AdminName: splittedLine[3],
			Lat:       lat,
			Lon:       lon,
			StateCode: splittedLine[4],
		}
	}

	if err := scanner.Err(); err != nil {
		return Zipcodes{}, fmt.Errorf("zipcodes: error while opening file %v", err)
	}
	return zipcodeMap, nil
}
