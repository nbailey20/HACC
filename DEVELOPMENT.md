# Homemade Authentication Credential Client - HACC 2.0

## Contributing - HACC CLI

* Update code in hacc/ folder as desired
* Execute main tests with `go test ./...`
* Locally test with `go run . [normal HACC cmd]` inside the hacc folder
* Once code working, build binary with `GOOS=windows GOARCH=amd64 go build -o ../releases/hacc-windows-amd64.exe .`
* Push to Github