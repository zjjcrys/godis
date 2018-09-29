package core

import (
	"fmt"
	"os"
	"strconv"
)

const RADIUS_COORDS = (1 << 0) /* Search around coordinates. */
const RADIUS_MEMBER = (1 << 1) /* Search around member. */
const RADIUS_NOSTORE = (1 << 2)

const SORT_NONE = 0
const SORT_ASC = 1
const SORT_DESC = 2

// geoaddCommand 命令实现
func GeoAddCommand(c *Client, s *Server) {
	// check params numbers
	if (c.Argc-2)%3 != 0 {
		/* Need an odd number of arguments if we got this far... */
		addReplyError(c, "syntax error. Try GEOADD key [x1] [y1] [name1] "+
			"[x2] [y2] [name2] ... ")
	}

	elements := (c.Argc - 2) / 3 //坐标数
	argc := 2 + elements*2       /* ZADD key score ele ... */
	argv := make([]*GodisObject, argc)
	argv[0] = CreateObject(ObjectTypeString, "zadd")
	argv[1] = c.Argv[1]

	for i := 0; i < elements; i++ {
		var xy [2]float64
		var hash GeoHashBits
		//提取经纬度
		if lngObj, ok1 := c.Argv[i*3+2].Ptr.(string); ok1 {
			if latObj, ok2 := c.Argv[i*3+3].Ptr.(string); ok2 {
				var ok error
				xy[0], ok = strconv.ParseFloat(lngObj, 64)
				xy[1], ok = strconv.ParseFloat(latObj, 64)
				if ok != nil {
					addReplyError(c, "lng lat type error")
					os.Exit(0)
				}
			}
		}
		geohashEncodeWGS84(xy[0], xy[1], GEO_STEP_MAX, &hash)
		bits := geohashAlign52Bits(hash)
		score := CreateObject(ObjectTypeString, bits)

		val := c.Argv[2+i*3+2]
		argv[2+i*2] = score // 设置有序集合元素的分值和名字
		argv[3+i*2] = val
	}
	c.Argc = argc
	c.Argv = argv
	zaddCommand(c)

	addReplyStatus(c, "OK")
}

//获取特定位置的hash值
func GeoHashCommand(c *Client, s *Server) {
	geoAlphabet := "0123456789bcdefghjkmnpqrstuvwxyz"
	zobj := lookupKey(c.Db, c.Argv[1])
	if zobj != nil && zobj.ObjectType != OBJ_ZSET {
		return
	}
	buf := ""
	for j := 2; j < c.Argc; j++ {
		var score float64
		if zobj == nil || zsetScore(zobj, c.Argv[j].Ptr.(string), &score) == C_ERR {
			addReplyError(c, "score get error ")
			return
		}
		var xy [2]float64
		if !decodeGeohash(score, &xy) {
			addReplyError(c, "hash get error")
			continue
		}
		r := [2]GeoHashRange{}
		var hash GeoHashBits
		r[0].min = -180
		r[0].max = 180
		r[1].min = -90
		r[1].max = 90
		geohashEncode(&r[0], &r[1], xy[0], xy[1], 26, &hash)

		temp := ""
		for i := 0; i < 11; i++ {
			count := 52 - (i+1)*5
			idx := (hash.bits >> (uint(count))) & 0x1f
			temp += string(geoAlphabet[idx])
		}
		buf += temp
		buf += ";"
	}
	addReplyStatus(c, buf)
}

//获取经纬度
func GeoPosCommand(c *Client, s *Server) {
	zobj := lookupKey(c.Db, c.Argv[1])
	if zobj != nil && zobj.ObjectType != OBJ_ZSET {
		return
	}
	buf := "lng:"

	for j := 2; j < c.Argc; j++ {
		var score float64
		if zobj == nil || zsetScore(zobj, c.Argv[j].Ptr.(string), &score) == C_ERR {
			addReplyError(c, "score get error ")
			return
		}
		var xy [2]float64
		if !decodeGeohash(score, &xy) {
			addReplyError(c, "hash get error")
			continue
		}

		buf += fmt.Sprint(xy[0])
		buf += ",lat:"
		buf += fmt.Sprint(xy[1])
		buf += ";"
	}
	addReplyStatus(c, buf)
}

//获取两个位置的距离
func GeoDistCommand(c *Client, s *Server) {
	if c.Argc >= 5 {
		addReplyError(c, "params error")
		return
	}
	zobj := lookupKey(c.Db, c.Argv[1])
	if zobj != nil && zobj.ObjectType != OBJ_ZSET {
		return
	}

	var score1, score2 float64
	var xyxy1, xyxy2 [2]float64
	if zsetScore(zobj, c.Argv[2].Ptr.(string), &score1) == C_ERR ||
		zsetScore(zobj, c.Argv[3].Ptr.(string), &score2) == C_ERR {
		addReplyError(c, "score get error ")
		return
	}

	if !decodeGeohash(score1, &xyxy1) || !decodeGeohash(score2, &xyxy2) {
		addReplyError(c, "hash get error")
		return
	}

	buf := geohashGetDistance(xyxy1[0], xyxy1[1], xyxy2[0], xyxy2[1])
	addReplyStatus(c, fmt.Sprint(buf))
}

func GeoRadiusCommand(c *Client, s *Server) {
	georadiusGeneric(c, RADIUS_COORDS)
}

func GeoRadiusByMemberCommand(c *Client, s *Server) {
	georadiusGeneric(c, RADIUS_MEMBER)
}

func georadiusGeneric(c *Client, flags int) {
	storedist := 0
	zobj := lookupKey(c.Db, c.Argv[1])
	if zobj != nil && zobj.ObjectType != OBJ_ZSET {
		return
	}

	var xy [2]float64
	if !decodeGeohash(score, &xy) {
		addReplyError(c, "hash get error")
		continue
	}

	//从参数中获取半径和单位
	var radius_meters, conversion float64
	radius_meters = c.Argv

	//提取所有可选参数
	withdist := 0
	withhash := 0
	withcoords := 0
	sort := SORT_NONE
	var count int64
	count = 0

}

func membersOfAllNeighbors() {

}

func membersOfGeoHashBox(zobj *GodisObject, hash GeoHashBits, ga *geoArray, lon float64, lat float64, radius float64) int {
	var min, max GeoHashFix52Bits

	scoresOfGeoHashBox(hash, &min, &max)
	return geoGetPointsInRange(zobj, float64(min), float64(max), lon, lat, radius, ga)
}

func scoresOfGeoHashBox(hash GeoHashBits, min *GeoHashFix52Bits, max *GeoHashFix52Bits) {
	*min = geohashAlign52Bits(hash)
	hash.bits++
	*max = geohashAlign52Bits(hash)
}

func geoGetPointsInRange(zobj *GodisObject, min float64, max float64, lon float64, lat float64, radius float64, ga *geoArray) int {
	range :=zRangeSpec{min:min,max:max,minEx:0,maxEx:1}
	var origincount uint = ga.used
	var member string
	if zobj.ObjectType==OBJ_ZSET {
	
	} else {
		//ziplist
	}
	return ga.used - origincount;
}
