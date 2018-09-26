package core

import (
	"os"
	"strconv"
)

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

	for j := 2; j < c.Argc; j++ {
		var score float64
		if zobj == nil || zsetScore(zobj, c.Argv[j].Ptr.(string), &score) == C_ERR {
			addReplyError(c, "score get error ")
			return
		}
		var xy [2]float64
		if !decodeGeohash(score, xy) {
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

		buf := ""
		for i := 0; i < 11; i++ {
			count := 52 - (i+1)*5
			idx := (hash.bits >> (uint(count))) & 0x1f
			buf += string(geoAlphabet[idx])
		}
		addReplyStatus(c, buf)
	}
}

//获取经纬度
func GeoPosCommand(c *Client, s *Server) {

}

//获取两个位置的距离
func GeoDistCommand(c *Client, s *Server) {

}
