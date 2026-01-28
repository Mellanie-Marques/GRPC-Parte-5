module github.com/Mellanie-Marques/microservices/order

go 1.24.0

toolchain go1.24.5

require (
	github.com/Mellanie-Marques/microservices-proto/golang/order v0.0.0-00010101000000-000000000000
	github.com/Mellanie-Marques/microservices-proto/golang/payment v0.0.0-00010101000000-000000000000
	github.com/Mellanie-Marques/microservices-proto/golang/shipping v0.0.0-00010101000000-000000000000
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	google.golang.org/grpc v1.78.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.31.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace (
	github.com/Mellanie-Marques/microservices-proto/golang/order => ../../microservices-proto/golang/order
	github.com/Mellanie-Marques/microservices-proto/golang/payment => ../../microservices-proto/golang/payment
	github.com/Mellanie-Marques/microservices-proto/golang/shipping => ../../microservices-proto/golang/shipping
)
