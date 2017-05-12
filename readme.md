# Auth Service

Auth Service provides a simple gRPC service used to create and validate JWT tokens using
accounts created with the [account service](https://github.com/lileio/account_service)

``` protobuf
service AuthService {
  rpc Authenticate (AuthRequest) returns (AuthResponse) {}
  rpc Validate (ValidateRequest) returns (google.protobuf.Empty) {}
}
```

## Configuration

### SIGNING_TOKEN

`SIGNING_TOKEN` is the token used for the HMAC signing of the JWT token.

Choose a subtibly random string for your service and never expose it to other services.

### TOKEN_EXPIRY

`TOKEN_EXPIRY` sets how long the tokens are valid for before the user needs to re-authenticate.

The default is 48 hours, which is set as "48h". You can set any time which Go's [ParseDuration](https://golang.org/pkg/time/#ParseDuration) can parse.
