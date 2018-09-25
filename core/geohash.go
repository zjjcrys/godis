package core

const GEO_STEP_MAX = 26 /* 26*2 = 52 bits. */

/* Limits from EPSG:900913 / EPSG:3785 / OSGEO:41001 */
const GEO_LAT_MIN = -85.05112878
const GEO_LAT_MAX = 85.05112878
const GEO_LONG_MIN = -180
const GEO_LONG_MAX = 180

type GeoHashBits struct {
	bits uint64
	step uint8
}

type GeoHashRange struct {
	min float64
	max float64
}

func geohashEncodeWGS84(longitude float64, latitude float64, step uint8, hash *GeoHashBits) int {
	return geohashEncodeType(longitude, latitude, step, hash)
}

func geohashEncodeType(longitude float64, latitude float64, step uint8, hash *GeoHashBits) int {
	r := [2]GeoHashRange{}
	geohashGetCoordRange(&r[0], &r[1])
	return geohashEncode(&r[0], &r[1], longitude, latitude, step, hash)
}

/* These are constraints from EPSG:900913 / EPSG:3785 / OSGEO:41001 */
/* We can't geocode at the north/south pole. */
func geohashGetCoordRange(long_range *GeoHashRange, lat_range *GeoHashRange) {
	long_range.max = GEO_LONG_MAX
	long_range.min = GEO_LONG_MIN
	lat_range.max = GEO_LAT_MAX
	lat_range.min = GEO_LAT_MIN
}

func geohashEncode(long_range *GeoHashRange, lat_range *GeoHashRange, longitude float64, latitude float64, step uint8,
	hash *GeoHashBits) int {
	/* Check basic arguments sanity. */

	/* Return an error when trying to index outside the supported
	 * constraints. */
	if longitude > 180 || longitude < -180 ||
		latitude > 85.05112878 || latitude < -85.05112878 {
		return 0
	}

	hash.bits = 0
	hash.step = step

	if latitude < lat_range.min || latitude > lat_range.max ||
		longitude < long_range.min || longitude > long_range.max {
		return 0
	}

	var lat_offset float64
	var long_offset float64
	lat_offset =
		(latitude - lat_range.min) / (lat_range.max - lat_range.min)
	long_offset =
		(longitude - long_range.min) / (long_range.max - long_range.min)

	/* convert to fixed point based on the step size */
	mask := 1 << step
	lat_offset = lat_offset * float64(mask)
	long_offset = long_offset * float64(mask)
	hash.bits = interleave64(int32(lat_offset), int32(long_offset))
	return 1
}

/*
lat 放在偶数位，lng放在奇数位
*/
func interleave64(latOffset int32, lngOffset int32) uint64 {
	B := []uint64{0x5555555555555555, 0x3333333333333333,
		0x0F0F0F0F0F0F0F0F, 0x00FF00FF00FF00FF,
		0x0000FFFF0000FFFF}
	S := []uint8{1, 2, 4, 8, 16}
	x := uint64(latOffset)
	y := uint64(lngOffset)
	x = (x | (x << S[4])) & B[4]
	y = (y | (y << S[4])) & B[4]
	x = (x | (x << S[3])) & B[3]
	y = (y | (y << S[3])) & B[3]
	x = (x | (x << S[2])) & B[2]
	y = (y | (y << S[2])) & B[2]
	x = (x | (x << S[1])) & B[1]
	y = (y | (y << S[1])) & B[1]
	x = (x | (x << S[0])) & B[0]
	y = (y | (y << S[0])) & B[0]
	return x | (y << 1)
}

func geohashAlign52Bits(hash GeoHashBits) uint64 {
	bits := hash.bits
	bits <<= (52 - hash.step*2)
	return bits
}
