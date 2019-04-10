package zipcodes

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	zipcodesDataset, err := New("datasets/valid_dataset.txt")
	if err != nil {
		t.Errorf("Unexpected error while initializing struct %v", err)
	}
	if (reflect.TypeOf(*zipcodesDataset) != reflect.TypeOf(Zipcodes{})) {
		t.Errorf("Unexpected response type. Got %v, want %v", reflect.TypeOf(*zipcodesDataset), reflect.TypeOf(Zipcodes{}))
	}
}

func TestLoadDataset(t *testing.T) {
	// Wrong file format cases
	cases := []struct {
		Dataset       string
		ExpectedError string
	}{
		{
			"datasets/wrong_length_dataset.txt",
			"zipcodes: file line does not have 12 fields",
		},
		{
			"datasets/wrong_lat_dataset.txt",
			"zipcodes: error while converting WRONG to Latitude",
		},
		{
			"datasets/wrong_lon_dataset.txt",
			"zipcodes: error while converting WRONG to Longitude",
		},
	}

	for _, c := range cases {
		_, err := LoadDataset(c.Dataset)
		if err.Error() != c.ExpectedError {
			t.Errorf("Unexpected error. Got %s, want %s", err, c.ExpectedError)
		}
	}

	// Valid file format cases
	dataset, err := LoadDataset("datasets/valid_dataset.txt")
	if err != nil {
		t.Errorf("Unexpected error while initializing struct %v", err)
	}
	if (reflect.TypeOf(dataset) != reflect.TypeOf(Zipcodes{})) {
		t.Errorf("Unexpected response type. Got %v, want %v", reflect.TypeOf(dataset), reflect.TypeOf(Zipcodes{}))
	}
}

func TestLookup(t *testing.T) {
	zipcodesDataset, err := New("datasets/valid_dataset.txt")
	if err != nil {
		t.Errorf("Unexpected error while initializing struct %v", err)
	}

	// Looking for a zipcode that exists
	existingZipCode := "01945"
	foundedZC, err := zipcodesDataset.Lookup(existingZipCode)
	if err != nil {
		t.Errorf("Unexpected error while looking for zipcode %s", existingZipCode)
	}
	expectedZipCode := ZipCodeLocation{
		ZipCode:   "01945",
		PlaceName: "Guteborn",
		AdminName: "Brandenburg",
		Lat:       51.4167,
		Lon:       13.9333,
	}

	if reflect.DeepEqual(foundedZC, &expectedZipCode) != true {
		t.Errorf("Unexpected response when calling Lookup")
	}
	// Looking for a zipcode that does not exists
	missingZipCode := "XYZ"
	_, errZC := zipcodesDataset.Lookup(missingZipCode)
	if errZC.Error() != "zipcodes: zipcode XYZ not found !" {
		t.Errorf("Unexpected error while looking for zipcode %s", existingZipCode)
	}
}
