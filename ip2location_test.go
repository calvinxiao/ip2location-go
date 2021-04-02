package ip2location

import "fmt"
import "testing"


var ip = "8.8.8.8"

func TestSample(t *testing.T) {

	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		t.Errorf("Open database fail.")
	}

	results, err := db.Get_all(ip)
	
	if err != nil {
		t.Errorf("Get_all fail.")
	}

	fmt.Printf("country_short: %s\n", results.Country_short)
	fmt.Printf("country_long: %s\n", results.Country_long)
	fmt.Printf("region: %s\n", results.Region)
	fmt.Printf("city: %s\n", results.City)
	fmt.Printf("isp: %s\n", results.Isp)
	fmt.Printf("latitude: %f\n", results.Latitude)
	fmt.Printf("longitude: %f\n", results.Longitude)
	fmt.Printf("domain: %s\n", results.Domain)
	fmt.Printf("zipcode: %s\n", results.Zipcode)
	fmt.Printf("timezone: %s\n", results.Timezone)
	fmt.Printf("netspeed: %s\n", results.Netspeed)
	fmt.Printf("iddcode: %s\n", results.Iddcode)
	fmt.Printf("areacode: %s\n", results.Areacode)
	fmt.Printf("weatherstationcode: %s\n", results.Weatherstationcode)
	fmt.Printf("weatherstationname: %s\n", results.Weatherstationname)
	fmt.Printf("mcc: %s\n", results.Mcc)
	fmt.Printf("mnc: %s\n", results.Mnc)
	fmt.Printf("mobilebrand: %s\n", results.Mobilebrand)
	fmt.Printf("elevation: %f\n", results.Elevation)
	fmt.Printf("usagetype: %s\n", results.Usagetype)
	fmt.Printf("api version: %s\n", Api_version())
	
	db.Close()
}

func TestSampleNew(t *testing.T) {


	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		t.Errorf("Open database fail.")
	}

	results, err := db.Get_all2(ip, 128)
	
	if err != nil {
		t.Errorf("Get_all fail.")
	}

	fmt.Printf("country_short: %s\n", results.Country_short)
	fmt.Printf("country_long: %s\n", results.Country_long)
	fmt.Printf("region: %s\n", results.Region)
	fmt.Printf("city: %s\n", results.City)
	fmt.Printf("isp: %s\n", results.Isp)
	fmt.Printf("latitude: %f\n", results.Latitude)
	fmt.Printf("longitude: %f\n", results.Longitude)
	fmt.Printf("domain: %s\n", results.Domain)
	fmt.Printf("zipcode: %s\n", results.Zipcode)
	fmt.Printf("timezone: %s\n", results.Timezone)
	fmt.Printf("netspeed: %s\n", results.Netspeed)
	fmt.Printf("iddcode: %s\n", results.Iddcode)
	fmt.Printf("areacode: %s\n", results.Areacode)
	fmt.Printf("weatherstationcode: %s\n", results.Weatherstationcode)
	fmt.Printf("weatherstationname: %s\n", results.Weatherstationname)
	fmt.Printf("mcc: %s\n", results.Mcc)
	fmt.Printf("mnc: %s\n", results.Mnc)
	fmt.Printf("mobilebrand: %s\n", results.Mobilebrand)
	fmt.Printf("elevation: %f\n", results.Elevation)
	fmt.Printf("usagetype: %s\n", results.Usagetype)
	fmt.Printf("api version: %s\n", Api_version())
	
	db.Close()
}

func BenchmarkGetAll1(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		db.Get_all(ip)
	}

	db.Close()
}

func BenchmarkGetAll2(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		db.Get_all2(ip, 2)
	}

	db.Close()

}

func BenchmarkGetAll4(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		db.Get_all2(ip, 4)
	}

	db.Close()

}

func BenchmarkGetAll8(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		db.Get_all2(ip, 8)
	}

	db.Close()

}

func BenchmarkGetAll16(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		db.Get_all2(ip, 16)
	}

	db.Close()

}

func BenchmarkGetAll32(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		db.Get_all2(ip, 32)
	}

	db.Close()

}

func BenchmarkGetAll64(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		db.Get_all2(ip, 64)
	}

	db.Close()

}

func BenchmarkGetAll128(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		db.Get_all2(ip, 128)
	}

	db.Close()

}

func BenchmarkGetAll256(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		db.Get_all2(ip, 256)
	}

	db.Close()

}

func BenchmarkGetAll512(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		db.Get_all2(ip, 512)
	}

	db.Close()

}

func BenchmarkGetAll1024(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		db.Get_all2(ip, 1024)
	}

	db.Close()

}

func BenchmarkGetAll2048(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		db.Get_all2(ip, 2048)
	}

	db.Close()

}


func BenchmarkGetAll4096(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		db.Get_all2(ip, 4096)
	}

	db.Close()

}

func BenchmarkByteStack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var by [4]byte
		by[0] = 1
	}
}

func BenchmarkByteHeap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var by = make([]byte, 4)
		by[0] = 1
	}
}