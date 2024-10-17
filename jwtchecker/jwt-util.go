package jwtchecker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"payment-module/internals/constants"
	logger "payment-module/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _config Config

/*
REQUIRED(Any middleware must have this)

For every middleware we need a config.
In config we also need to define a function which allows us to skip the middleware if return true.
By convention it should be named as "Filter" but any other name will work too.
*/
type Config struct {
	// when returned true, our middleware is skipped
	Filter func(c *fiber.Ctx) bool
	// function to run when there is error decoding jwt
	Unauthorized fiber.Handler
	// function to decode our jwt token
	Decode func(c *fiber.Ctx) (*jwt.MapClaims, error)
	// set jwt secret
	Secret string
	// set jwt expiry in seconds
	Expiry int64
}

/*
Middleware specific

Our middleware's config default values if not passed
*/
var ConfigDefault = Config{
	Filter:       nil,
	Decode:       nil,
	Unauthorized: nil,
	Secret:       "secret",
	Expiry:       60,
}

/*
Middleware specific
Function for generating default config
*/
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values if not passed
	if cfg.Filter == nil {
		cfg.Filter = ConfigDefault.Filter
	}

	// Set default secret if not passed
	if cfg.Secret == "" {
		cfg.Secret = ConfigDefault.Secret
	}

	// Set default expiry if not passed
	if cfg.Expiry == 0 {
		cfg.Expiry = ConfigDefault.Expiry
	}

	// this is the main jwt decode function of our middleware
	if cfg.Decode == nil {
		// Set default Decode function if not passed
		cfg.Decode = func(c *fiber.Ctx) (*jwt.MapClaims, error) {
			logger.Info("Start decode authenticate")

			authHeader := c.Get(constants.AUTHORIZATION)
			logger.Infof("Start decode authenticate with authHeader:" + authHeader)
			if authHeader == "" || len(authHeader) < 8 {
				return nil, errors.New("AUTHORIZATION HEADER IS REQUIRED")
			}

			// we parse our jwt token and check for validity against our secret
			token, err := jwt.Parse(
				authHeader[7:],
				func(token *jwt.Token) (interface{}, error) {
					// IMPORTANT: Validating the algorithm per https://godoc.org/github.com/dgrijalva/jwt-go#example-Parse--Hmac
					if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
						return nil, fmt.Errorf(
							"EXPECTED TOKEN ALGORITHM '%v' BUT GOT '%v'",
							jwt.SigningMethodRS256.Name,
							token.Header)
					}
					untypedKeyId, found := token.Header["alg"]
					if !found {
						return nil, fmt.Errorf("NO KEY ID KEY '%v' FOUND IN TOKEN HEADER", "alg")
					}
					keyId, ok := untypedKeyId.(string)
					if !ok {
						return nil, fmt.Errorf("FOUND KEY ID, BUT VALUE WAS NOT A STRING")
					}

					// keyBase64, found := auth0_constants.RsaPublicKeyBase64[keyId]
					// if !found {
					// 		return nil, fmt.Errorf("No public RSA key found corresponding to key ID from token '%v'", keyId)
					// }
					//keyStr := pubKeyHeader + "\n" + keyBase64 + "\n" + pubKeyFooter

					// Since the token is RSA (which we validated at the start of this function), the return type of this function actually has to be rsa.PublicKey!
					pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.Secret))
					if err != nil {
						return nil, fmt.Errorf("AN ERROR OCCURRED PARSING THE PUBLIC KEY BASE64 FOR KEY ID '%v'; THIS IS A CODE BUG", keyId)
					}

					return pubKey, nil
					//return []byte(cfg.Secret), nil
				},
			)

			if err != nil {
				return nil, errors.New("ERROR PARSING TOKEN")
			}

			claims, ok := token.Claims.(jwt.MapClaims)

			// if !(ok && token.Valid) {
			if !(ok) {
				return nil, errors.New("INVALID TOKEN")
			}

			if expiresAt, ok := claims["exp"]; ok && int64(expiresAt.(float64)) < time.Now().UTC().Unix() {
				return nil, errors.New("jwt is expired")
			}
			logger.Info("End decode authenticate successful.")
			return &claims, nil
		}
	}

	// Set default Unauthorized if not passed
	if cfg.Unauthorized == nil {
		cfg.Unauthorized = func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}

	return cfg
}

/*
Middleware specific
Function to generate a jwt token
*/
func Encode(claims *jwt.MapClaims, secret string, expiryAfter int64) (string, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(secret))
	if err != nil {
		return "", err
	}

	// setting default expiryAfter
	if expiryAfter == 0 {
		expiryAfter = ConfigDefault.Expiry
	}

	// or you can use time.Now().Add(time.Second * time.Duration(expiryAfter)).UTC().Unix()
	(*claims)["exp"] = time.Now().UTC().Unix() + expiryAfter

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// our signed jwt token string
	signedToken, err := token.SignedString(privateKey)

	if err != nil {
		return "", errors.New("ERROR CREATING A TOKEN")
	}

	return signedToken, nil
}

/*
REQUIRED(Any middleware must have this)
Our main middleware function used to initialize our middleware.
By convention we name it "New" but any other name will work too.
*/
func New(config Config) fiber.Handler {
	// For setting default config
	_config = configDefault(config)

	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Filter returns true
		if _config.Filter != nil && _config.Filter(c) {
			logger.Debug("Midddle was skipped")
			return c.Next()
		}
		logger.Debug("Midddle was run")

		claims, err := _config.Decode(c)
		if err == nil {
			c.Locals("jwtClaims", *claims)
			return c.Next()
		}

		return _config.Unauthorized(c)
	}
}

func GetPrivateKeyFromFile(filePath string) (string, error) {
	// we need to pass in a secret otherwise default secret is used
	keyData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// return jwt.ParseRSAPublicKeyFromPEM(keyData)
	return string(keyData), nil
}

func GetPublicKeyFromFile(filePath string) (string, error) {
	// we need to pass in a secret otherwise default secret is used
	keyData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// return jwt.ParseRSAPublicKeyFromPEM(keyData)
	return string(keyData), nil
}

func GetRoleFromJwt(c *fiber.Ctx) []string {
	return ExtractClaimArrayString(c, constants.RJWT_KEY_ROLE)
}

func GetHospitalFromJwt(c *fiber.Ctx) string {
	return ExtractClaimString(c, constants.HOSPITAL_CODE_JWT_KEY)
}

func GetUserNameFromJwt(c *fiber.Ctx) string {
	return ExtractClaimString(c, constants.JWT_KEY_USER_NAME)
}

func GetUserIdFromJwt(c *fiber.Ctx) primitive.ObjectID {
	userIdString := ExtractClaimString(c, constants.JWT_KEY_USER_ID)
	userId, err := primitive.ObjectIDFromHex(userIdString)
	if err != nil {
		return primitive.NilObjectID
	}

	return userId
}

func GetPhoneFromJwt(c *fiber.Ctx) string {
	return ExtractClaimString(c, constants.JWT_KEY_PHONE)
}

func GetProvinceIdFromJwt(c *fiber.Ctx) float64 {
	return ExtractClaimNumber(c, constants.HOSPITAL_ID_JWT_KEY)
}

func ExtractClaimArrayString(c *fiber.Ctx, key string) []string {
	claims, err := extractClaims(c)
	if err != nil {
		return nil
	}
	aInterface := claims[key].([]interface{})
	aString := make([]string, len(aInterface))
	for i, v := range aInterface {
		aString[i] = v.(string)
	}
	return aString
}

func ExtractClaimString(c *fiber.Ctx, key string) string {
	claims, err := extractClaims(c)
	if err != nil {
		return ""
	}
	return claims[key].(string)
}

func ExtractClaimNumber(c *fiber.Ctx, key string) float64 {
	claims, err := extractClaims(c)
	if err != nil {
		return -1
	}
	return claims[key].(float64)
}

func extractClaims(c *fiber.Ctx) (jwt.MapClaims, error) {
	authHeader := c.Get(constants.AUTHORIZATION)
	logger.Infof("Start decode authenticate with authHeader:" + authHeader)
	if authHeader == "" || len(authHeader) < 8 {
		return nil, errors.New("AUTHORIZATION HEADER IS REQUIRED")
	}
	// we parse our jwt token and check for validity against our secret
	token, err := jwt.Parse(
		authHeader[7:],
		func(token *jwt.Token) (interface{}, error) {
			// IMPORTANT: Validating the algorithm per https://godoc.org/github.com/dgrijalva/jwt-go#example-Parse--Hmac
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf(
					"EXPECTED TOKEN ALGORITHM '%v' BUT GOT '%v'",
					jwt.SigningMethodRS256.Name,
					token.Header)
			}
			untypedKeyId, found := token.Header["alg"]
			if !found {
				return nil, fmt.Errorf("NO KEY ID KEY '%v' FOUND IN TOKEN HEADER", "alg")
			}
			keyId, ok := untypedKeyId.(string)
			if !ok {
				return nil, fmt.Errorf("FOUND KEY ID, BUT VALUE WAS NOT A STRING")
			}

			// keyBase64, found := auth0_constants.RsaPublicKeyBase64[keyId]
			// if !found {
			// 		return nil, fmt.Errorf("No public RSA key found corresponding to key ID from token '%v'", keyId)
			// }
			//keyStr := pubKeyHeader + "\n" + keyBase64 + "\n" + pubKeyFooter

			// Since the token is RSA (which we validated at the start of this function), the return type of this function actually has to be rsa.PublicKey!
			pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(_config.Secret))
			if err != nil {
				return nil, fmt.Errorf("AN ERROR OCCURRED PARSING THE PUBLIC KEY BASE64 FOR KEY ID '%v'; THIS IS A CODE BUG", keyId)
			}

			return pubKey, nil
		},
	)
	if err != nil {
		logger.Infof("extract Claims error: %s", err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		logger.Info("Invalid JWT Token")
		return nil, nil
	}
}
