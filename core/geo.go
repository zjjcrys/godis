package core

import (
		"strconv"
	"os"
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
				xy[0],ok = strconv.ParseFloat(lngObj, 64)
				xy[1],ok = strconv.ParseFloat(latObj, 64)
				if ok!=nil {
					addReplyError(c,"lng lat type error")
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
	//geoalphabet :="0123456789bcdefghjkmnpqrstuvwxyz"
	//robj:=lookupKey(c.Db,c.Argv[1])

}

//获取经纬度
func GeoPosCommand(c *Client, s *Server) {

}

//获取两个位置的距离
func GeoDistCommand(c *Client, s *Server) {

}
