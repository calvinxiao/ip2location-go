// This ip2location package provides a fast lookup of country, region, city, latitude, longitude, ZIP code, time zone,
// ISP, domain name, connection type, IDD code, area code, weather station code, station name, MCC, MNC,
// mobile brand, elevation, and usage type from IP address by using IP2Location database.
package ip2location

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
	"net"
	"os"
	"strconv"

	"lukechampine.com/uint128"
)

type DBReader interface {
	io.ReadCloser
	io.ReaderAt
}

type ip2locationmeta struct {
	databasetype      uint8
	databasecolumn    uint8
	databaseday       uint8
	databasemonth     uint8
	databaseyear      uint8
	ipv4databasecount uint32
	ipv4databaseaddr  uint32
	ipv6databasecount uint32
	ipv6databaseaddr  uint32
	ipv4indexbaseaddr uint32
	ipv6indexbaseaddr uint32
	ipv4columnsize    uint32
	ipv6columnsize    uint32
}

// The IP2Locationrecord struct stores all of the available
// geolocation info found in the IP2Location database.
type IP2Locationrecord struct {
	Country_short      string
	Country_long       string
	Region             string
	City               string
	Isp                string
	Latitude           float32
	Longitude          float32
	Domain             string
	Zipcode            string
	Timezone           string
	Netspeed           string
	Iddcode            string
	Areacode           string
	Weatherstationcode string
	Weatherstationname string
	Mcc                string
	Mnc                string
	Mobilebrand        string
	Elevation          float32
	Usagetype          string
}

type DB struct {
	f    DBReader
	meta ip2locationmeta

	country_position_offset            uint32
	region_position_offset             uint32
	city_position_offset               uint32
	isp_position_offset                uint32
	domain_position_offset             uint32
	zipcode_position_offset            uint32
	latitude_position_offset           uint32
	longitude_position_offset          uint32
	timezone_position_offset           uint32
	netspeed_position_offset           uint32
	iddcode_position_offset            uint32
	areacode_position_offset           uint32
	weatherstationcode_position_offset uint32
	weatherstationname_position_offset uint32
	mcc_position_offset                uint32
	mnc_position_offset                uint32
	mobilebrand_position_offset        uint32
	elevation_position_offset          uint32
	usagetype_position_offset          uint32

	country_enabled            bool
	region_enabled             bool
	city_enabled               bool
	isp_enabled                bool
	domain_enabled             bool
	zipcode_enabled            bool
	latitude_enabled           bool
	longitude_enabled          bool
	timezone_enabled           bool
	netspeed_enabled           bool
	iddcode_enabled            bool
	areacode_enabled           bool
	weatherstationcode_enabled bool
	weatherstationname_enabled bool
	mcc_enabled                bool
	mnc_enabled                bool
	mobilebrand_enabled        bool
	elevation_enabled          bool
	usagetype_enabled          bool

	metaok bool
}

var defaultDB = &DB{}

var country_position = [25]uint8{0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
var region_position = [25]uint8{0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}
var city_position = [25]uint8{0, 0, 0, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}
var isp_position = [25]uint8{0, 0, 3, 0, 5, 0, 7, 5, 7, 0, 8, 0, 9, 0, 9, 0, 9, 0, 9, 7, 9, 0, 9, 7, 9}
var latitude_position = [25]uint8{0, 0, 0, 0, 0, 5, 5, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5}
var longitude_position = [25]uint8{0, 0, 0, 0, 0, 6, 6, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6}
var domain_position = [25]uint8{0, 0, 0, 0, 0, 0, 0, 6, 8, 0, 9, 0, 10, 0, 10, 0, 10, 0, 10, 8, 10, 0, 10, 8, 10}
var zipcode_position = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 0, 7, 7, 7, 0, 7, 0, 7, 7, 7, 0, 7}
var timezone_position = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 8, 7, 8, 8, 8, 7, 8, 0, 8, 8, 8, 0, 8}
var netspeed_position = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 11, 0, 11, 8, 11, 0, 11, 0, 11, 0, 11}
var iddcode_position = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 12, 0, 12, 0, 12, 9, 12, 0, 12}
var areacode_position = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 13, 0, 13, 0, 13, 10, 13, 0, 13}
var weatherstationcode_position = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 14, 0, 14, 0, 14, 0, 14}
var weatherstationname_position = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 15, 0, 15, 0, 15, 0, 15}
var mcc_position = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 16, 0, 16, 9, 16}
var mnc_position = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 17, 0, 17, 10, 17}
var mobilebrand_position = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 18, 0, 18, 11, 18}
var elevation_position = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 19, 0, 19}
var usagetype_position = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 12, 20}

const api_version string = "8.4.0"

var max_ipv4_range = uint128.From64(4294967295)
var max_ipv6_range = uint128.From64(0)
var from_v4mapped = uint128.From64(281470681743360)
var to_v4mapped = uint128.From64(281474976710655)
var from_6to4 = uint128.From64(0)
var to_6to4 = uint128.From64(0)
var from_teredo = uint128.From64(0)
var to_teredo = uint128.From64(0)
var last_32bits = uint128.From64(4294967295)

const countryshort uint32 = 0x00001
const countrylong uint32 = 0x00002
const region uint32 = 0x00004
const city uint32 = 0x00008
const isp uint32 = 0x00010
const latitude uint32 = 0x00020
const longitude uint32 = 0x00040
const domain uint32 = 0x00080
const zipcode uint32 = 0x00100
const timezone uint32 = 0x00200
const netspeed uint32 = 0x00400
const iddcode uint32 = 0x00800
const areacode uint32 = 0x01000
const weatherstationcode uint32 = 0x02000
const weatherstationname uint32 = 0x04000
const mcc uint32 = 0x08000
const mnc uint32 = 0x10000
const mobilebrand uint32 = 0x20000
const elevation uint32 = 0x40000
const usagetype uint32 = 0x80000

const all uint32 = countryshort | countrylong | region | city | isp | latitude | longitude | domain | zipcode | timezone | netspeed | iddcode | areacode | weatherstationcode | weatherstationname | mcc | mnc | mobilebrand | elevation | usagetype

const invalid_address string = "Invalid IP address."
const missing_file string = "Invalid database file."
const not_supported string = "This parameter is unavailable for selected data file. Please upgrade the data file."

// get IP type and calculate IP number; calculates index too if exists
func (d *DB) checkip(ip string) (iptype uint32, ipnum uint128.Uint128, ipindex uint32) {
	iptype = 0
	ipnum = uint128.From64(0)
	ipnumtmp := uint128.From64(0)
	ipindex = 0
	ipaddress := net.ParseIP(ip)

	if ipaddress != nil {
		v4 := ipaddress.To4()

		if v4 != nil {
			iptype = 4
			ipnum = uint128.From64(uint64(binary.BigEndian.Uint32(v4)))
		} else {
			v6 := ipaddress.To16()

			if v6 != nil {
				iptype = 6
				ipnum.PutBytes(v6)

				if ipnum.Cmp(from_v4mapped) >= 0 && ipnum.Cmp(to_v4mapped) <= 0 {
					// ipv4-mapped ipv6 should treat as ipv4 and read ipv4 data section
					iptype = 4
					ipnum = ipnum.Sub(from_v4mapped)
				} else if ipnum.Cmp(from_6to4) >= 0 && ipnum.Cmp(to_6to4) <= 0 {
					// 6to4 so need to remap to ipv4
					iptype = 4
					ipnum = ipnum.Rsh(80)
					ipnum = ipnum.And(last_32bits)
				} else if ipnum.Cmp(from_teredo) >= 0 && ipnum.Cmp(to_teredo) <= 0 {
					// Teredo so need to remap to ipv4
					iptype = 4
					ipnum = uint128.Uint128{^ipnum.Lo, ^ipnum.Hi}
					ipnum = ipnum.And(last_32bits)
				}
			}
		}
	}
	if iptype == 4 {
		if d.meta.ipv4indexbaseaddr > 0 {
			ipnumtmp = ipnum.Rsh(16)
			ipnumtmp = ipnumtmp.Lsh(3)
			ipindex = uint32(ipnumtmp.Add(uint128.From64(uint64(d.meta.ipv4indexbaseaddr))).Lo)
		}
	} else if iptype == 6 {
		if d.meta.ipv6indexbaseaddr > 0 {
			ipnumtmp = ipnum.Rsh(112)
			ipnumtmp = ipnumtmp.Lsh(3)
			ipindex = uint32(ipnumtmp.Add(uint128.From64(uint64(d.meta.ipv6indexbaseaddr))).Lo)
		}
	}
	return
}

// read byte
func (d *DB) readuint8(pos int64) (uint8, error) {
	var retval uint8
	data := make([]byte, 1)
	_, err := d.f.ReadAt(data, pos-1)
	if err != nil {
		return 0, err
	}
	retval = data[0]
	return retval, nil
}

// read unsigned 32-bit integer from slices
func (d *DB) readuint32_row(row []byte, pos uint32) uint32 {
	var retval uint32
	data := row[pos : pos+4]
	retval = binary.LittleEndian.Uint32(data)
	return retval
}

// read unsigned 32-bit integer
func (d *DB) readuint32(pos uint32) (uint32, error) {
	pos2 := int64(pos)
	var retval uint32
	var data [4]byte
	_, err := d.f.ReadAt(data[:], pos2-1)
	if err != nil {
		return 0, err
	}

	retval = binary.LittleEndian.Uint32(data[:])
	return retval, nil
}

// read unsigned 128-bit integer
func (d *DB) readuint128(pos uint32) (uint128.Uint128, error) {
	pos2 := int64(pos)
	retval := uint128.From64(0)
	data := make([]byte, 16)
	_, err := d.f.ReadAt(data, pos2-1)
	if err != nil {
		return uint128.From64(0), err
	}

	// little endian to big endian
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	retval.PutBytes(data)
	return retval, nil
}

// read string
func (d *DB) readstr(pos uint32) (string, error) {
	pos2 := int64(pos)
	var retval string
	lenbyte := make([]byte, 1)
	_, err := d.f.ReadAt(lenbyte, pos2)
	if err != nil {
		return "", err
	}
	strlen := lenbyte[0]
	data := make([]byte, strlen)
	_, err = d.f.ReadAt(data, pos2+1)
	if err != nil {
		return "", err
	}
	retval = string(data[:strlen])
	return retval, nil
}

// read float from slices
func (d *DB) readfloat_row(row []byte, pos uint32) float32 {
	var retval float32
	data := row[pos : pos+4]
	bits := binary.LittleEndian.Uint32(data)
	retval = math.Float32frombits(bits)
	return retval
}

// read float
func (d *DB) readfloat(pos uint32) (float32, error) {
	pos2 := int64(pos)
	var retval float32
	data := make([]byte, 4)
	_, err := d.f.ReadAt(data, pos2-1)
	if err != nil {
		return 0, err
	}
	buf := bytes.NewReader(data)
	err = binary.Read(buf, binary.LittleEndian, &retval)
	if err != nil {
		fmt.Printf("binary read failed: %v", err)
	}
	return retval, nil
}

func fatal(db *DB, err error) (*DB, error) {
	_ = db.f.Close()
	return nil, err
}

// Open takes the path to the IP2Location BIN database file. It will read all the metadata required to
// be able to extract the embedded geolocation data, and return the underlining DB object.
func OpenDB(dbpath string) (*DB, error) {
	f, err := os.Open(dbpath)
	if err != nil {
		return nil, err
	}

	return OpenDBWithReader(f)
}

// OpenDBWithReader takes a DBReader to the IP2Location BIN database file. It will read all the metadata required to
// be able to extract the embedded geolocation data, and return the underlining DB object.
func OpenDBWithReader(reader DBReader) (*DB, error) {
	var db = &DB{}

	_max_ipv6_range := big.NewInt(0)
	_max_ipv6_range.SetString("340282366920938463463374607431768211455", 10)
	max_ipv6_range = uint128.FromBig(_max_ipv6_range)

	_from_6to4 := big.NewInt(0)
	_from_6to4.SetString("42545680458834377588178886921629466624", 10)
	from_6to4 = uint128.FromBig(_from_6to4)

	_to_6to4 := big.NewInt(0)
	_to_6to4.SetString("42550872755692912415807417417958686719", 10)
	to_6to4 = uint128.FromBig(_to_6to4)

	_from_teredo := big.NewInt(0)
	_from_teredo.SetString("42540488161975842760550356425300246528", 10)
	from_teredo = uint128.FromBig(_from_teredo)

	_to_teredo := big.NewInt(0)
	_to_teredo.SetString("42540488241204005274814694018844196863", 10)
	to_teredo = uint128.FromBig(_to_teredo)

	db.f = reader

	var err error
	db.meta.databasetype, err = db.readuint8(1)
	if err != nil {
		return fatal(db, err)
	}
	db.meta.databasecolumn, err = db.readuint8(2)
	if err != nil {
		return fatal(db, err)
	}
	db.meta.databaseyear, err = db.readuint8(3)
	if err != nil {
		return fatal(db, err)
	}
	db.meta.databasemonth, err = db.readuint8(4)
	if err != nil {
		return fatal(db, err)
	}
	db.meta.databaseday, err = db.readuint8(5)
	if err != nil {
		return fatal(db, err)
	}
	db.meta.ipv4databasecount, err = db.readuint32(6)
	if err != nil {
		return fatal(db, err)
	}
	db.meta.ipv4databaseaddr, err = db.readuint32(10)
	if err != nil {
		return fatal(db, err)
	}
	db.meta.ipv6databasecount, err = db.readuint32(14)
	if err != nil {
		return fatal(db, err)
	}
	db.meta.ipv6databaseaddr, err = db.readuint32(18)
	if err != nil {
		return fatal(db, err)
	}
	db.meta.ipv4indexbaseaddr, err = db.readuint32(22)
	if err != nil {
		return fatal(db, err)
	}
	db.meta.ipv6indexbaseaddr, err = db.readuint32(26)
	if err != nil {
		return fatal(db, err)
	}
	db.meta.ipv4columnsize = uint32(db.meta.databasecolumn << 2)              // 4 bytes each column
	db.meta.ipv6columnsize = uint32(16 + ((db.meta.databasecolumn - 1) << 2)) // 4 bytes each column, except IPFrom column which is 16 bytes

	dbt := db.meta.databasetype

	if country_position[dbt] != 0 {
		db.country_position_offset = uint32(country_position[dbt]-2) << 2
		db.country_enabled = true
	}
	if region_position[dbt] != 0 {
		db.region_position_offset = uint32(region_position[dbt]-2) << 2
		db.region_enabled = true
	}
	if city_position[dbt] != 0 {
		db.city_position_offset = uint32(city_position[dbt]-2) << 2
		db.city_enabled = true
	}
	if isp_position[dbt] != 0 {
		db.isp_position_offset = uint32(isp_position[dbt]-2) << 2
		db.isp_enabled = true
	}
	if domain_position[dbt] != 0 {
		db.domain_position_offset = uint32(domain_position[dbt]-2) << 2
		db.domain_enabled = true
	}
	if zipcode_position[dbt] != 0 {
		db.zipcode_position_offset = uint32(zipcode_position[dbt]-2) << 2
		db.zipcode_enabled = true
	}
	if latitude_position[dbt] != 0 {
		db.latitude_position_offset = uint32(latitude_position[dbt]-2) << 2
		db.latitude_enabled = true
	}
	if longitude_position[dbt] != 0 {
		db.longitude_position_offset = uint32(longitude_position[dbt]-2) << 2
		db.longitude_enabled = true
	}
	if timezone_position[dbt] != 0 {
		db.timezone_position_offset = uint32(timezone_position[dbt]-2) << 2
		db.timezone_enabled = true
	}
	if netspeed_position[dbt] != 0 {
		db.netspeed_position_offset = uint32(netspeed_position[dbt]-2) << 2
		db.netspeed_enabled = true
	}
	if iddcode_position[dbt] != 0 {
		db.iddcode_position_offset = uint32(iddcode_position[dbt]-2) << 2
		db.iddcode_enabled = true
	}
	if areacode_position[dbt] != 0 {
		db.areacode_position_offset = uint32(areacode_position[dbt]-2) << 2
		db.areacode_enabled = true
	}
	if weatherstationcode_position[dbt] != 0 {
		db.weatherstationcode_position_offset = uint32(weatherstationcode_position[dbt]-2) << 2
		db.weatherstationcode_enabled = true
	}
	if weatherstationname_position[dbt] != 0 {
		db.weatherstationname_position_offset = uint32(weatherstationname_position[dbt]-2) << 2
		db.weatherstationname_enabled = true
	}
	if mcc_position[dbt] != 0 {
		db.mcc_position_offset = uint32(mcc_position[dbt]-2) << 2
		db.mcc_enabled = true
	}
	if mnc_position[dbt] != 0 {
		db.mnc_position_offset = uint32(mnc_position[dbt]-2) << 2
		db.mnc_enabled = true
	}
	if mobilebrand_position[dbt] != 0 {
		db.mobilebrand_position_offset = uint32(mobilebrand_position[dbt]-2) << 2
		db.mobilebrand_enabled = true
	}
	if elevation_position[dbt] != 0 {
		db.elevation_position_offset = uint32(elevation_position[dbt]-2) << 2
		db.elevation_enabled = true
	}
	if usagetype_position[dbt] != 0 {
		db.usagetype_position_offset = uint32(usagetype_position[dbt]-2) << 2
		db.usagetype_enabled = true
	}

	db.metaok = true

	return db, nil
}

// Open takes the path to the IP2Location BIN database file. It will read all the metadata required to
// be able to extract the embedded geolocation data.
//
// Deprecated: No longer being updated.
func Open(dbpath string) {
	db, err := OpenDB(dbpath)
	if err != nil {
		return
	}

	defaultDB = db
}

// Close will close the file handle to the BIN file.
//
// Deprecated: No longer being updated.
func Close() {
	defaultDB.Close()
}

// Api_version returns the version of the component.
func Api_version() string {
	return api_version
}

// populate record with message
func loadmessage(mesg string) IP2Locationrecord {
	var x IP2Locationrecord

	x.Country_short = mesg
	x.Country_long = mesg
	x.Region = mesg
	x.City = mesg
	x.Isp = mesg
	x.Domain = mesg
	x.Zipcode = mesg
	x.Timezone = mesg
	x.Netspeed = mesg
	x.Iddcode = mesg
	x.Areacode = mesg
	x.Weatherstationcode = mesg
	x.Weatherstationname = mesg
	x.Mcc = mesg
	x.Mnc = mesg
	x.Mobilebrand = mesg
	x.Usagetype = mesg

	return x
}

func handleError(rec IP2Locationrecord, err error) IP2Locationrecord {
	if err != nil {
		fmt.Print(err)
	}
	return rec
}

// Get_all will return all geolocation fields based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_all(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, all))
}

// Get_country_short will return the ISO-3166 country code based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_country_short(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, countryshort))
}

// Get_country_long will return the country name based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_country_long(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, countrylong))
}

// Get_region will return the region name based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_region(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, region))
}

// Get_city will return the city name based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_city(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, city))
}

// Get_isp will return the Internet Service Provider name based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_isp(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, isp))
}

// Get_latitude will return the latitude based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_latitude(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, latitude))
}

// Get_longitude will return the longitude based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_longitude(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, longitude))
}

// Get_domain will return the domain name based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_domain(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, domain))
}

// Get_zipcode will return the postal code based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_zipcode(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, zipcode))
}

// Get_timezone will return the time zone based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_timezone(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, timezone))
}

// Get_netspeed will return the Internet connection speed based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_netspeed(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, netspeed))
}

// Get_iddcode will return the International Direct Dialing code based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_iddcode(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, iddcode))
}

// Get_areacode will return the area code based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_areacode(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, areacode))
}

// Get_weatherstationcode will return the weather station code based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_weatherstationcode(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, weatherstationcode))
}

// Get_weatherstationname will return the weather station name based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_weatherstationname(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, weatherstationname))
}

// Get_mcc will return the mobile country code based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_mcc(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, mcc))
}

// Get_mnc will return the mobile network code based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_mnc(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, mnc))
}

// Get_mobilebrand will return the mobile carrier brand based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_mobilebrand(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, mobilebrand))
}

// Get_elevation will return the elevation in meters based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_elevation(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, elevation))
}

// Get_usagetype will return the usage type based on the queried IP address.
//
// Deprecated: No longer being updated.
func Get_usagetype(ipaddress string) IP2Locationrecord {
	return handleError(defaultDB.query(ipaddress, usagetype))
}

// Get_all will return all geolocation fields based on the queried IP address.
func (d *DB) Get_all(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, all)
}

func (d *DB) Get_all2(ipaddress string, threshold uint32) (IP2Locationrecord, error) {
	return d.query2(ipaddress, all, threshold)
}

// Get_country_short will return the ISO-3166 country code based on the queried IP address.
func (d *DB) Get_country_short(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, countryshort)
}

// Get_country_long will return the country name based on the queried IP address.
func (d *DB) Get_country_long(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, countrylong)
}

// Get_region will return the region name based on the queried IP address.
func (d *DB) Get_region(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, region)
}

// Get_city will return the city name based on the queried IP address.
func (d *DB) Get_city(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, city)
}

// Get_isp will return the Internet Service Provider name based on the queried IP address.
func (d *DB) Get_isp(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, isp)
}

// Get_latitude will return the latitude based on the queried IP address.
func (d *DB) Get_latitude(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, latitude)
}

// Get_longitude will return the longitude based on the queried IP address.
func (d *DB) Get_longitude(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, longitude)
}

// Get_domain will return the domain name based on the queried IP address.
func (d *DB) Get_domain(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, domain)
}

// Get_zipcode will return the postal code based on the queried IP address.
func (d *DB) Get_zipcode(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, zipcode)
}

// Get_timezone will return the time zone based on the queried IP address.
func (d *DB) Get_timezone(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, timezone)
}

// Get_netspeed will return the Internet connection speed based on the queried IP address.
func (d *DB) Get_netspeed(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, netspeed)
}

// Get_iddcode will return the International Direct Dialing code based on the queried IP address.
func (d *DB) Get_iddcode(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, iddcode)
}

// Get_areacode will return the area code based on the queried IP address.
func (d *DB) Get_areacode(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, areacode)
}

// Get_weatherstationcode will return the weather station code based on the queried IP address.
func (d *DB) Get_weatherstationcode(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, weatherstationcode)
}

// Get_weatherstationname will return the weather station name based on the queried IP address.
func (d *DB) Get_weatherstationname(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, weatherstationname)
}

// Get_mcc will return the mobile country code based on the queried IP address.
func (d *DB) Get_mcc(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, mcc)
}

// Get_mnc will return the mobile network code based on the queried IP address.
func (d *DB) Get_mnc(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, mnc)
}

// Get_mobilebrand will return the mobile carrier brand based on the queried IP address.
func (d *DB) Get_mobilebrand(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, mobilebrand)
}

// Get_elevation will return the elevation in meters based on the queried IP address.
func (d *DB) Get_elevation(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, elevation)
}

// Get_usagetype will return the usage type based on the queried IP address.
func (d *DB) Get_usagetype(ipaddress string) (IP2Locationrecord, error) {
	return d.query(ipaddress, usagetype)
}

// main query
func (d *DB) query(ipaddress string, mode uint32) (IP2Locationrecord, error) {
	x := loadmessage(not_supported) // default message

	// read metadata
	if !d.metaok {
		x = loadmessage(missing_file)
		return x, nil
	}

	// check IP type and return IP number & index (if exists)
	iptype, ipno, ipindex := d.checkip(ipaddress)

	if iptype == 0 {
		x = loadmessage(invalid_address)
		return x, nil
	}

	var err error
	var colsize uint32
	var baseaddr uint32
	var low uint32
	var high uint32
	var mid uint32
	var rowoffset uint32
	var rowoffset2 uint32
	ipfrom := uint128.From64(0)
	ipto := uint128.From64(0)
	maxip := uint128.From64(0)

	if iptype == 4 {
		baseaddr = d.meta.ipv4databaseaddr
		high = d.meta.ipv4databasecount
		maxip = max_ipv4_range
		colsize = d.meta.ipv4columnsize
	} else {
		baseaddr = d.meta.ipv6databaseaddr
		high = d.meta.ipv6databasecount
		maxip = max_ipv6_range
		colsize = d.meta.ipv6columnsize
	}

	// reading index
	if ipindex > 0 {
		low, err = d.readuint32(ipindex)
		if err != nil {
			return x, err
		}
		high, err = d.readuint32(ipindex + 4)
		if err != nil {
			return x, err
		}

		fmt.Println(low, high)
	}

	if ipno.Cmp(maxip) >= 0 {
		ipno = ipno.Sub(uint128.From64(1))
	}

	for low <= high {
		mid = ((low + high) >> 1)
		rowoffset = baseaddr + (mid * colsize)
		rowoffset2 = rowoffset + colsize

		if iptype == 4 {
			ipfrom32, err := d.readuint32(rowoffset)
			if err != nil {
				return x, err
			}
			ipfrom = uint128.From64(uint64(ipfrom32))

			ipto32, err := d.readuint32(rowoffset2)
			if err != nil {
				return x, err
			}
			ipto = uint128.From64(uint64(ipto32))

		} else {
			ipfrom, err = d.readuint128(rowoffset)
			if err != nil {
				return x, err
			}

			ipto, err = d.readuint128(rowoffset2)
			if err != nil {
				return x, err
			}
		}

		if ipno.Cmp(ipfrom) >= 0 && ipno.Cmp(ipto) < 0 {
			var firstcol uint32 = 4 // 4 bytes for ip from
			if iptype == 6 {
				firstcol = 16 // 16 bytes for ipv6
				// rowoffset = rowoffset + 12 // coz below is assuming all columns are 4 bytes, so got 12 left to go to make 16 bytes total
			}

			row := make([]byte, colsize-firstcol) // exclude the ip from field
			_, err := d.f.ReadAt(row, int64(rowoffset+firstcol-1))
			if err != nil {
				return x, err
			}

			if mode&countryshort == 1 && d.country_enabled {
				// x.Country_short = readstr(readuint32(rowoffset + country_position_offset))
				if x.Country_short, err = d.readstr(d.readuint32_row(row, d.country_position_offset)); err != nil {
					return x, err
				}
			}

			if mode&countrylong != 0 && d.country_enabled {
				// x.Country_long = readstr(readuint32(rowoffset + country_position_offset) + 3)
				if x.Country_long, err = d.readstr(d.readuint32_row(row, d.country_position_offset) + 3); err != nil {
					return x, err
				}
			}

			if mode&region != 0 && d.region_enabled {
				// x.Region = readstr(readuint32(rowoffset + region_position_offset))
				if x.Region, err = d.readstr(d.readuint32_row(row, d.region_position_offset)); err != nil {
					return x, err
				}
			}

			if mode&city != 0 && d.city_enabled {
				// x.City = readstr(readuint32(rowoffset + city_position_offset))
				if x.City, err = d.readstr(d.readuint32_row(row, d.city_position_offset)); err != nil {
					return x, err
				}
			}

			if mode&isp != 0 && d.isp_enabled {
				// x.Isp = readstr(readuint32(rowoffset + isp_position_offset))
				if x.Isp, err = d.readstr(d.readuint32_row(row, d.isp_position_offset)); err != nil {
					return x, err
				}
			}

			if mode&latitude != 0 && d.latitude_enabled {
				// x.Latitude = readfloat(rowoffset + latitude_position_offset)
				x.Latitude = d.readfloat_row(row, d.latitude_position_offset)
			}

			if mode&longitude != 0 && d.longitude_enabled {
				// x.Longitude = readfloat(rowoffset + longitude_position_offset)
				x.Longitude = d.readfloat_row(row, d.longitude_position_offset)
			}

			if mode&domain != 0 && d.domain_enabled {
				// x.Domain = readstr(readuint32(rowoffset + domain_position_offset))
				if x.Domain, err = d.readstr(d.readuint32_row(row, d.domain_position_offset)); err != nil {
					return x, err
				}
			}

			if mode&zipcode != 0 && d.zipcode_enabled {
				// x.Zipcode = readstr(readuint32(rowoffset + zipcode_position_offset))
				if x.Zipcode, err = d.readstr(d.readuint32_row(row, d.zipcode_position_offset)); err != nil {
					return x, err
				}
			}

			if mode&timezone != 0 && d.timezone_enabled {
				// x.Timezone = readstr(readuint32(rowoffset + timezone_position_offset))
				if x.Timezone, err = d.readstr(d.readuint32_row(row, d.timezone_position_offset)); err != nil {
					return x, err
				}
			}

			if mode&netspeed != 0 && d.netspeed_enabled {
				// x.Netspeed = readstr(readuint32(rowoffset + netspeed_position_offset))
				if x.Netspeed, err = d.readstr(d.readuint32_row(row, d.netspeed_position_offset)); err != nil {
					return x, err
				}
			}

			if mode&iddcode != 0 && d.iddcode_enabled {
				// x.Iddcode = readstr(readuint32(rowoffset + iddcode_position_offset))
				if x.Iddcode, err = d.readstr(d.readuint32_row(row, d.iddcode_position_offset)); err != nil {
					return x, err
				}
			}

			if mode&areacode != 0 && d.areacode_enabled {
				// x.Areacode = readstr(readuint32(rowoffset + areacode_position_offset))
				if x.Areacode, err = d.readstr(d.readuint32_row(row, d.areacode_position_offset)); err != nil {
					return x, err
				}
			}

			if mode&weatherstationcode != 0 && d.weatherstationcode_enabled {
				// x.Weatherstationcode = readstr(readuint32(rowoffset + weatherstationcode_position_offset))
				if x.Weatherstationcode, err = d.readstr(d.readuint32_row(row, d.weatherstationcode_position_offset)); err != nil {
					return x, err
				}
			}

			if mode&weatherstationname != 0 && d.weatherstationname_enabled {
				// x.Weatherstationname = readstr(readuint32(rowoffset + weatherstationname_position_offset))
				if x.Weatherstationname, err = d.readstr(d.readuint32_row(row, d.weatherstationname_position_offset)); err != nil {
					return x, err
				}
			}

			if mode&mcc != 0 && d.mcc_enabled {
				// x.Mcc = readstr(readuint32(rowoffset + mcc_position_offset))
				if x.Mcc, err = d.readstr(d.readuint32_row(row, d.mcc_position_offset)); err != nil {
					return x, err
				}
			}

			if mode&mnc != 0 && d.mnc_enabled {
				// x.Mnc = readstr(readuint32(rowoffset + mnc_position_offset))
				if x.Mnc, err = d.readstr(d.readuint32_row(row, d.mnc_position_offset)); err != nil {
					return x, err
				}
			}

			if mode&mobilebrand != 0 && d.mobilebrand_enabled {
				// x.Mobilebrand = readstr(readuint32(rowoffset + mobilebrand_position_offset))
				if x.Mobilebrand, err = d.readstr(d.readuint32_row(row, d.mobilebrand_position_offset)); err != nil {
					return x, err
				}
			}

			if mode&elevation != 0 && d.elevation_enabled {
				// f, _ := strconv.ParseFloat(readstr(readuint32(rowoffset + elevation_position_offset)), 32)
				res, err := d.readstr(d.readuint32_row(row, d.elevation_position_offset))
				if err != nil {
					return x, err
				}

				f, _ := strconv.ParseFloat(res, 32)
				x.Elevation = float32(f)
			}

			if mode&usagetype != 0 && d.usagetype_enabled {
				// x.Usagetype = readstr(readuint32(rowoffset + usagetype_position_offset))
				if x.Usagetype, err = d.readstr(d.readuint32_row(row, d.usagetype_position_offset)); err != nil {
					return x, err
				}
			}

			return x, nil
		} else {
			if ipno.Cmp(ipfrom) < 0 {
				high = mid - 1
			} else {
				low = mid + 1
			}
		}
	}
	return x, nil
}

func (d *DB) query2(ipaddress string, mode uint32, threshold uint32) (IP2Locationrecord, error) {
	x := loadmessage(not_supported) // default message

	// read metadata
	if !d.metaok {
		x = loadmessage(missing_file)
		return x, nil
	}

	// check IP type and return IP number & index (if exists)
	iptype, ipno, ipindex := d.checkip(ipaddress)

	if iptype == 0 {
		x = loadmessage(invalid_address)
		return x, nil
	}

	var err error
	var colsize uint32
	var baseaddr uint32
	var low uint32
	var high uint32
	var mid uint32
	var rowoffset uint32
	var rowoffset2 uint32
	ipfrom := uint128.From64(0)
	ipto := uint128.From64(0)
	maxip := uint128.From64(0)

	if iptype == 4 {
		baseaddr = d.meta.ipv4databaseaddr
		high = d.meta.ipv4databasecount
		maxip = max_ipv4_range
		colsize = d.meta.ipv4columnsize
	} else {
		baseaddr = d.meta.ipv6databaseaddr
		high = d.meta.ipv6databasecount
		maxip = max_ipv6_range
		colsize = d.meta.ipv6columnsize
	}

	// reading index
	if ipindex > 0 {
		low, err = d.readuint32(ipindex)
		if err != nil {
			return x, err
		}
		high, err = d.readuint32(ipindex + 4)
		if err != nil {
			return x, err
		}
	}

	if ipno.Cmp(maxip) >= 0 {
		ipno = ipno.Sub(uint128.From64(1))
	}

	for low <= high {
		if high-low > threshold {
			mid = ((low + high) >> 1)
			rowoffset = baseaddr + (mid * colsize)
			rowoffset2 = rowoffset + colsize

			if iptype == 4 {
				ipfrom32, err := d.readuint32(rowoffset)
				if err != nil {
					return x, err
				}
				ipfrom = uint128.From64(uint64(ipfrom32))

				ipto32, err := d.readuint32(rowoffset2)
				if err != nil {
					return x, err
				}
				ipto = uint128.From64(uint64(ipto32))

			} else {
				ipfrom, err = d.readuint128(rowoffset)
				if err != nil {
					return x, err
				}

				ipto, err = d.readuint128(rowoffset2)
				if err != nil {
					return x, err
				}
			}

			cmpFrom := ipno.Cmp(ipfrom)

			if cmpFrom >= 0 && ipno.Cmp(ipto) < 0 {
				var firstcol uint32 = 4 // 4 bytes for ip from
				if iptype == 6 {
					firstcol = 16 // 16 bytes for ipv6
					// rowoffset = rowoffset + 12 // coz below is assuming all columns are 4 bytes, so got 12 left to go to make 16 bytes total
				}

				row := make([]byte, colsize-firstcol) // exclude the ip from field
				_, err := d.f.ReadAt(row, int64(rowoffset+firstcol-1))
				if err != nil {
					return x, err
				}

				if mode&countryshort == 1 && d.country_enabled {
					// x.Country_short = readstr(readuint32(rowoffset + country_position_offset))
					if x.Country_short, err = d.readstr(d.readuint32_row(row, d.country_position_offset)); err != nil {
						return x, err
					}
				}

				if mode&countrylong != 0 && d.country_enabled {
					// x.Country_long = readstr(readuint32(rowoffset + country_position_offset) + 3)
					if x.Country_long, err = d.readstr(d.readuint32_row(row, d.country_position_offset) + 3); err != nil {
						return x, err
					}
				}

				if mode&region != 0 && d.region_enabled {
					// x.Region = readstr(readuint32(rowoffset + region_position_offset))
					if x.Region, err = d.readstr(d.readuint32_row(row, d.region_position_offset)); err != nil {
						return x, err
					}
				}

				if mode&city != 0 && d.city_enabled {
					// x.City = readstr(readuint32(rowoffset + city_position_offset))
					if x.City, err = d.readstr(d.readuint32_row(row, d.city_position_offset)); err != nil {
						return x, err
					}
				}

				if mode&isp != 0 && d.isp_enabled {
					// x.Isp = readstr(readuint32(rowoffset + isp_position_offset))
					if x.Isp, err = d.readstr(d.readuint32_row(row, d.isp_position_offset)); err != nil {
						return x, err
					}
				}

				if mode&latitude != 0 && d.latitude_enabled {
					// x.Latitude = readfloat(rowoffset + latitude_position_offset)
					x.Latitude = d.readfloat_row(row, d.latitude_position_offset)
				}

				if mode&longitude != 0 && d.longitude_enabled {
					// x.Longitude = readfloat(rowoffset + longitude_position_offset)
					x.Longitude = d.readfloat_row(row, d.longitude_position_offset)
				}

				if mode&domain != 0 && d.domain_enabled {
					// x.Domain = readstr(readuint32(rowoffset + domain_position_offset))
					if x.Domain, err = d.readstr(d.readuint32_row(row, d.domain_position_offset)); err != nil {
						return x, err
					}
				}

				if mode&zipcode != 0 && d.zipcode_enabled {
					// x.Zipcode = readstr(readuint32(rowoffset + zipcode_position_offset))
					if x.Zipcode, err = d.readstr(d.readuint32_row(row, d.zipcode_position_offset)); err != nil {
						return x, err
					}
				}

				if mode&timezone != 0 && d.timezone_enabled {
					// x.Timezone = readstr(readuint32(rowoffset + timezone_position_offset))
					if x.Timezone, err = d.readstr(d.readuint32_row(row, d.timezone_position_offset)); err != nil {
						return x, err
					}
				}

				if mode&netspeed != 0 && d.netspeed_enabled {
					// x.Netspeed = readstr(readuint32(rowoffset + netspeed_position_offset))
					if x.Netspeed, err = d.readstr(d.readuint32_row(row, d.netspeed_position_offset)); err != nil {
						return x, err
					}
				}

				if mode&iddcode != 0 && d.iddcode_enabled {
					// x.Iddcode = readstr(readuint32(rowoffset + iddcode_position_offset))
					if x.Iddcode, err = d.readstr(d.readuint32_row(row, d.iddcode_position_offset)); err != nil {
						return x, err
					}
				}

				if mode&areacode != 0 && d.areacode_enabled {
					// x.Areacode = readstr(readuint32(rowoffset + areacode_position_offset))
					if x.Areacode, err = d.readstr(d.readuint32_row(row, d.areacode_position_offset)); err != nil {
						return x, err
					}
				}

				if mode&weatherstationcode != 0 && d.weatherstationcode_enabled {
					// x.Weatherstationcode = readstr(readuint32(rowoffset + weatherstationcode_position_offset))
					if x.Weatherstationcode, err = d.readstr(d.readuint32_row(row, d.weatherstationcode_position_offset)); err != nil {
						return x, err
					}
				}

				if mode&weatherstationname != 0 && d.weatherstationname_enabled {
					// x.Weatherstationname = readstr(readuint32(rowoffset + weatherstationname_position_offset))
					if x.Weatherstationname, err = d.readstr(d.readuint32_row(row, d.weatherstationname_position_offset)); err != nil {
						return x, err
					}
				}

				if mode&mcc != 0 && d.mcc_enabled {
					// x.Mcc = readstr(readuint32(rowoffset + mcc_position_offset))
					if x.Mcc, err = d.readstr(d.readuint32_row(row, d.mcc_position_offset)); err != nil {
						return x, err
					}
				}

				if mode&mnc != 0 && d.mnc_enabled {
					// x.Mnc = readstr(readuint32(rowoffset + mnc_position_offset))
					if x.Mnc, err = d.readstr(d.readuint32_row(row, d.mnc_position_offset)); err != nil {
						return x, err
					}
				}

				if mode&mobilebrand != 0 && d.mobilebrand_enabled {
					// x.Mobilebrand = readstr(readuint32(rowoffset + mobilebrand_position_offset))
					if x.Mobilebrand, err = d.readstr(d.readuint32_row(row, d.mobilebrand_position_offset)); err != nil {
						return x, err
					}
				}

				if mode&elevation != 0 && d.elevation_enabled {
					// f, _ := strconv.ParseFloat(readstr(readuint32(rowoffset + elevation_position_offset)), 32)
					res, err := d.readstr(d.readuint32_row(row, d.elevation_position_offset))
					if err != nil {
						return x, err
					}

					f, _ := strconv.ParseFloat(res, 32)
					x.Elevation = float32(f)
				}

				if mode&usagetype != 0 && d.usagetype_enabled {
					// x.Usagetype = readstr(readuint32(rowoffset + usagetype_position_offset))
					if x.Usagetype, err = d.readstr(d.readuint32_row(row, d.usagetype_position_offset)); err != nil {
						return x, err
					}
				}

				return x, nil
			} else {
				if cmpFrom < 0 {
					high = mid - 1
				} else {
					low = mid + 1
				}
			}
		} else {
			for mid = low; mid <= high; mid++ {
				rowoffset = baseaddr + (mid * colsize)
				rowoffset2 = rowoffset + colsize

				if iptype == 4 {
					ipfrom32, err := d.readuint32(rowoffset)
					if err != nil {
						return x, err
					}
					ipfrom = uint128.From64(uint64(ipfrom32))

					ipto32, err := d.readuint32(rowoffset2)
					if err != nil {
						return x, err
					}
					ipto = uint128.From64(uint64(ipto32))

				} else {
					ipfrom, err = d.readuint128(rowoffset)
					if err != nil {
						return x, err
					}

					ipto, err = d.readuint128(rowoffset2)
					if err != nil {
						return x, err
					}
				}
				if ipno.Cmp(ipfrom) >= 0 && ipno.Cmp(ipto) < 0 {
					var firstcol uint32 = 4 // 4 bytes for ip from
					if iptype == 6 {
						firstcol = 16 // 16 bytes for ipv6
						// rowoffset = rowoffset + 12 // coz below is assuming all columns are 4 bytes, so got 12 left to go to make 16 bytes total
					}

					row := make([]byte, colsize-firstcol) // exclude the ip from field
					_, err := d.f.ReadAt(row, int64(rowoffset+firstcol-1))
					if err != nil {
						return x, err
					}

					if mode&countryshort == 1 && d.country_enabled {
						// x.Country_short = readstr(readuint32(rowoffset + country_position_offset))
						if x.Country_short, err = d.readstr(d.readuint32_row(row, d.country_position_offset)); err != nil {
							return x, err
						}
					}

					if mode&countrylong != 0 && d.country_enabled {
						// x.Country_long = readstr(readuint32(rowoffset + country_position_offset) + 3)
						if x.Country_long, err = d.readstr(d.readuint32_row(row, d.country_position_offset) + 3); err != nil {
							return x, err
						}
					}

					if mode&region != 0 && d.region_enabled {
						// x.Region = readstr(readuint32(rowoffset + region_position_offset))
						if x.Region, err = d.readstr(d.readuint32_row(row, d.region_position_offset)); err != nil {
							return x, err
						}
					}

					if mode&city != 0 && d.city_enabled {
						// x.City = readstr(readuint32(rowoffset + city_position_offset))
						if x.City, err = d.readstr(d.readuint32_row(row, d.city_position_offset)); err != nil {
							return x, err
						}
					}

					if mode&isp != 0 && d.isp_enabled {
						// x.Isp = readstr(readuint32(rowoffset + isp_position_offset))
						if x.Isp, err = d.readstr(d.readuint32_row(row, d.isp_position_offset)); err != nil {
							return x, err
						}
					}

					if mode&latitude != 0 && d.latitude_enabled {
						// x.Latitude = readfloat(rowoffset + latitude_position_offset)
						x.Latitude = d.readfloat_row(row, d.latitude_position_offset)
					}

					if mode&longitude != 0 && d.longitude_enabled {
						// x.Longitude = readfloat(rowoffset + longitude_position_offset)
						x.Longitude = d.readfloat_row(row, d.longitude_position_offset)
					}

					if mode&domain != 0 && d.domain_enabled {
						// x.Domain = readstr(readuint32(rowoffset + domain_position_offset))
						if x.Domain, err = d.readstr(d.readuint32_row(row, d.domain_position_offset)); err != nil {
							return x, err
						}
					}

					if mode&zipcode != 0 && d.zipcode_enabled {
						// x.Zipcode = readstr(readuint32(rowoffset + zipcode_position_offset))
						if x.Zipcode, err = d.readstr(d.readuint32_row(row, d.zipcode_position_offset)); err != nil {
							return x, err
						}
					}

					if mode&timezone != 0 && d.timezone_enabled {
						// x.Timezone = readstr(readuint32(rowoffset + timezone_position_offset))
						if x.Timezone, err = d.readstr(d.readuint32_row(row, d.timezone_position_offset)); err != nil {
							return x, err
						}
					}

					if mode&netspeed != 0 && d.netspeed_enabled {
						// x.Netspeed = readstr(readuint32(rowoffset + netspeed_position_offset))
						if x.Netspeed, err = d.readstr(d.readuint32_row(row, d.netspeed_position_offset)); err != nil {
							return x, err
						}
					}

					if mode&iddcode != 0 && d.iddcode_enabled {
						// x.Iddcode = readstr(readuint32(rowoffset + iddcode_position_offset))
						if x.Iddcode, err = d.readstr(d.readuint32_row(row, d.iddcode_position_offset)); err != nil {
							return x, err
						}
					}

					if mode&areacode != 0 && d.areacode_enabled {
						// x.Areacode = readstr(readuint32(rowoffset + areacode_position_offset))
						if x.Areacode, err = d.readstr(d.readuint32_row(row, d.areacode_position_offset)); err != nil {
							return x, err
						}
					}

					if mode&weatherstationcode != 0 && d.weatherstationcode_enabled {
						// x.Weatherstationcode = readstr(readuint32(rowoffset + weatherstationcode_position_offset))
						if x.Weatherstationcode, err = d.readstr(d.readuint32_row(row, d.weatherstationcode_position_offset)); err != nil {
							return x, err
						}
					}

					if mode&weatherstationname != 0 && d.weatherstationname_enabled {
						// x.Weatherstationname = readstr(readuint32(rowoffset + weatherstationname_position_offset))
						if x.Weatherstationname, err = d.readstr(d.readuint32_row(row, d.weatherstationname_position_offset)); err != nil {
							return x, err
						}
					}

					if mode&mcc != 0 && d.mcc_enabled {
						// x.Mcc = readstr(readuint32(rowoffset + mcc_position_offset))
						if x.Mcc, err = d.readstr(d.readuint32_row(row, d.mcc_position_offset)); err != nil {
							return x, err
						}
					}

					if mode&mnc != 0 && d.mnc_enabled {
						// x.Mnc = readstr(readuint32(rowoffset + mnc_position_offset))
						if x.Mnc, err = d.readstr(d.readuint32_row(row, d.mnc_position_offset)); err != nil {
							return x, err
						}
					}

					if mode&mobilebrand != 0 && d.mobilebrand_enabled {
						// x.Mobilebrand = readstr(readuint32(rowoffset + mobilebrand_position_offset))
						if x.Mobilebrand, err = d.readstr(d.readuint32_row(row, d.mobilebrand_position_offset)); err != nil {
							return x, err
						}
					}

					if mode&elevation != 0 && d.elevation_enabled {
						// f, _ := strconv.ParseFloat(readstr(readuint32(rowoffset + elevation_position_offset)), 32)
						res, err := d.readstr(d.readuint32_row(row, d.elevation_position_offset))
						if err != nil {
							return x, err
						}

						f, _ := strconv.ParseFloat(res, 32)
						x.Elevation = float32(f)
					}

					if mode&usagetype != 0 && d.usagetype_enabled {
						// x.Usagetype = readstr(readuint32(rowoffset + usagetype_position_offset))
						if x.Usagetype, err = d.readstr(d.readuint32_row(row, d.usagetype_position_offset)); err != nil {
							return x, err
						}
					}

					return x, nil
				}
			}
			break
		}
	}
	return x, nil
}

func (d *DB) Close() {
	_ = d.f.Close()
}

// Printrecord is used to output the geolocation data for debugging purposes.
func Printrecord(x IP2Locationrecord) {
	fmt.Printf("country_short: %s\n", x.Country_short)
	fmt.Printf("country_long: %s\n", x.Country_long)
	fmt.Printf("region: %s\n", x.Region)
	fmt.Printf("city: %s\n", x.City)
	fmt.Printf("isp: %s\n", x.Isp)
	fmt.Printf("latitude: %f\n", x.Latitude)
	fmt.Printf("longitude: %f\n", x.Longitude)
	fmt.Printf("domain: %s\n", x.Domain)
	fmt.Printf("zipcode: %s\n", x.Zipcode)
	fmt.Printf("timezone: %s\n", x.Timezone)
	fmt.Printf("netspeed: %s\n", x.Netspeed)
	fmt.Printf("iddcode: %s\n", x.Iddcode)
	fmt.Printf("areacode: %s\n", x.Areacode)
	fmt.Printf("weatherstationcode: %s\n", x.Weatherstationcode)
	fmt.Printf("weatherstationname: %s\n", x.Weatherstationname)
	fmt.Printf("mcc: %s\n", x.Mcc)
	fmt.Printf("mnc: %s\n", x.Mnc)
	fmt.Printf("mobilebrand: %s\n", x.Mobilebrand)
	fmt.Printf("elevation: %f\n", x.Elevation)
	fmt.Printf("usagetype: %s\n", x.Usagetype)
}
