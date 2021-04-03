package ip2location

import "fmt"
import "testing"


var ip = "8.8.8.8"
var ips = []string {
	"36.78.234.239",
	"249.176.59.9",
	"201.29.23.88",
	"8.126.72.242",
	"103.13.78.145",
	"17.113.247.153",
	"27.198.19.231",
	"117.106.191.195",
	"102.90.60.88",
	"185.12.218.105",
}

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

func BenchmarkGetAllOld(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		for _, _ip := range ips {
			db.Get_all(_ip)
		}
	}

	db.Close()
}

func BenchmarkGetAll2(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		for _, _ip := range ips {
			db.Get_all2(_ip, 2)
		}
	}

	db.Close()

}

func BenchmarkGetAll4(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		for _, _ip := range ips {
			db.Get_all2(_ip, 4)
		}
	}

	db.Close()

}

func BenchmarkGetAll8(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		for _, _ip := range ips {
			db.Get_all2(_ip, 8)
		}
	}

	db.Close()

}

func BenchmarkGetAll16(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		for _, _ip := range ips {
			db.Get_all2(_ip, 16)
		}
	}

	db.Close()

}

func BenchmarkGetAll32(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		for _, _ip := range ips {
			db.Get_all2(_ip, 32)
		}
	}

	db.Close()

}

func BenchmarkGetAll64(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		for _, _ip := range ips {
			db.Get_all2(_ip, 64)
		}
	}

	db.Close()

}

func BenchmarkGetAll128(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		for _, _ip := range ips {
			db.Get_all2(_ip, 128)
		}
	}

	db.Close()

}

func BenchmarkGetAll256(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		for _, _ip := range ips {
			db.Get_all2(_ip, 256)
		}
	}

	db.Close()

}

func BenchmarkGetAll512(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		for _, _ip := range ips {
			db.Get_all2(_ip, 512)
		}
	}

	db.Close()

}

func BenchmarkGetAll1024(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		for _, _ip := range ips {
			db.Get_all2(_ip, 1024)
		}
	}

	db.Close()

}

func BenchmarkGetAll2048(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		for _, _ip := range ips {
			db.Get_all2(_ip, 2048)
		}
	}

	db.Close()

}


func BenchmarkGetAll4096(b *testing.B) {
	db, err := OpenDB("./IP2LOCATION-LITE-DB11.BIN")

	if err != nil {
		b.Errorf("Open database fail.")
	}

	for i := 0; i < b.N; i++ {
		for _, _ip := range ips {
			db.Get_all2(_ip, 4096)
		}
	}

	db.Close()

}
