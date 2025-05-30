package utility

import (
	"dex/app/utility/net"
	"dex/app/utility/sequence"
	"strconv"
	"strings"
)

var Sequence sequence.Sequence

func initSequence() (err error) {
	machineIP := net.GetDomainIP()
	machineSignature := strings.Replace(machineIP, ".", "", 4)
	machineID, err := strconv.ParseInt(machineSignature, 16, 64)
	machineID %= 32
	if err != nil {
		return
	}
	// Create snowflake sequence
	Sequence = sequence.NewSnowflake(uint64(10000), uint64(machineID), uint64(1))
	return
}
