# CCONV Cringle Tech Test

## Development set up

Assuming you have a valid go environment set up at `${GOPATH}` then
```
mkdir -p ${GOPATH}/src/github.com/robotlovesyou
cd ${GOPATH}/src/github.com/robotlovesyou
git clone git@github.com:robotlovesyou/cringletest.git
cd cringletest
```

## Testing
To run the tests
```
go test -v ./...
```

### Testing CurrencyLayer and Sendgrid connections (optional)

To prevent overuse of the currencylayer and sendgrid apis their tests will be skipped unless certain
environment variables are present. The easiest way to make these environment variables available is via a .env file which the project will
automatically load

```
CURRENCYLAYER_API_KEY=[currencylayer api key]
SENDGRID_API_KEY=[sendgrid api key]
SENDGRID_FROM_ADDRESS=[some email address]
SENDGRID_TEST_TO_ADDRESS=[some other email address]
```

## Building and running

To build and install
```
go install ./cconv
```

### Configuration
The cconv binary relies upon configuration settings being available from the environment (for a real world app I would use something more robust but this is a tech test!). For simplicity it will load a .env file from the current working directory. It requires
```
CURRENCYLAYER_API_KEY=[currencylayer api key]
SENDGRID_API_KEY=[sendgrid api key]
SENDGRID_FROM_ADDRESS=[some email address]
```

### Running

Assuming ${GOPATH}/bin is in your $PATH then the cli can be used by running
```
cconv
```

The cli includes instructions explaining how it should be utilised for the three requested modes of operation so I won't repeat them here