module keeper/services/session_manager

go 1.23.4

toolchain go1.24.4

replace keeper/services/lock_manager => ../lock_manager

require (
	google.golang.org/grpc v1.73.0
	keeper/services/lock_manager v0.0.0-00010101000000-000000000000
)

require (
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
