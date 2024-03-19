package config

import "os"

const (
	NOHOST           = "no host here"
	NOROUTER         = "no router here"
	HOOKHOST         = "the host: %v not in the hookwhite list"
	NOMESSAGETABLE   = "Please add the Proto Message Table"
	IPLIMITED        = "the IP is limited"
	BUCKETEMPTY      = "the token bucket is empty"
	NOPOOL           = "no grpc connect pool"
	WRONGPATHPATTERN = "wrong url pattern"
	WRONGPATH        = "%v is wrong url"
	GRPCPROXYEORROR  = "gRPC proxying error"
	GRPCPATHEORROR   = "path is wrong"
	SYSTEMERROR      = "sysem error"
	CONFIGFILEERROR  = "fatal error config file: %v"
	RELOADROUTER     = "DO not reload the router"
	IPADDRERROR      = "Bad data"
	NOPROTOMESSAGE   = "the proto message name : %v not in the prototable"
)

type MashType string

const (
	Http MashType = "Http"
	Grpc MashType = "Grpc"
)

type HttpType string

const (
	Full     HttpType = "full"
	Nohook   HttpType = "nohook"
	Onlyhook HttpType = "onlyhook"
)

type ParamType string

const (
	String ParamType = "string"
	Float  ParamType = "float"
	Int    ParamType = "int"
)

type BalanceType string

const (
	RoundRobin  BalanceType = "RoundRobin"
	WeightRobin BalanceType = "WeightRobin"
)

var (
	GateWay_Log   = os.Getenv("gateway_log")
	GateWay_Debug = os.Getenv("gateway_isdebug")
)
