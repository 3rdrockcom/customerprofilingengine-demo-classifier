# # Customer Profiling Engine Demo

## Classifier Testing Tool
This project allows a user to feed transaction data into the program and generate a classification and other useful statisitics.

### Instructions
```bash
# Get package
go get github.com/epointpayment/customerprofilingengine-demo-classifier
cd $GOPATH/src/github.com/epointpayment/customerprofilingengine-demo-classifier

# Build exec
go build -o classifier

# Run
./classifier -file=samples/sample.weekly.csv
```