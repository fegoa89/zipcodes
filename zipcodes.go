// zipcodes is a package that uses the GeoNames Postal Code dataset from http://www.geonames.org
// in order to perform zipcode lookup operations
package zipcodes

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// ZipCodeLocation struct represents each line of the dataset
type ZipCodeLocation struct {
	ZipCode   string
	PlaceName string
	AdminName string
	Lat       float64
	Lon       float64
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
		}
	}

	if err := scanner.Err(); err != nil {
		return Zipcodes{}, fmt.Errorf("zipcodes: error while opening file %v", err)
	}
	return zipcodeMap, nil
}
