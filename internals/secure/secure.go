package secure

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	_error "payment-module/error"
	models "payment-module/internals/models"
	"payment-module/internals/responses"
	jwt "payment-module/jwtchecker"
	"payment-module/logger"

	"github.com/gofiber/fiber/v2"
)

func CheckPermission(c *fiber.Ctx, acceptedRoles []string) bool {
	logger.Info("Check permission start!")
	roles := jwt.GetRoleFromJwt(c)

	for _, v := range roles {
		if contains(acceptedRoles, v) {
			return true
		}
	}

	return false
}

func ValidateChecksum(secretKey string, message models.Message) (bool, *_error.SystemError) {
	dataChecksum := fmt.Sprintf("%s|%d|%s", message.Message, message.Time, secretKey)
	logger.Debugf("dataChecksum :%s", dataChecksum)
	hasher := sha256.New()
	hasher.Write([]byte(dataChecksum))
	sha1_hash := hex.EncodeToString(hasher.Sum(nil))

	logger.Debugf("checksum-in-request:%s checksum-system-build:%s", message.Checksum, sha1_hash)
	result := sha1_hash == message.Checksum
	if !result {
		return false, _error.NewChecksumInvalid()
	}
	return true, nil
}

func CheckUserRolePermission(acceptedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		hasAccess := CheckPermission(c, acceptedRoles)
		if !hasAccess {
			return c.Status(http.StatusForbidden).JSON(responses.BaseResponse{Status: http.StatusForbidden, Message: "Forbidden", Data: ""})
		}
		return c.Next()
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
