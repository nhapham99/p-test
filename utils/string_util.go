package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/mail"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// Mang cac ky tu goc co dau
var SOURCE_CHARACTERS, LL_LENGTH = stringToRune(`ÀÁÂÃÈÉÊÌÍÒÓÔÕÙÚÝàáâãèéêìíòóôõùúýĂăĐđĨĩŨũƠơƯưẠạẢảẤấẦầẨẩẪẫẬậẮắẰằẲẳẴẵẶặẸẹẺẻẼẽẾếỀềỂểỄễỆệỈỉỊịỌọỎỏỐốỒồỔổỖỗỘộỚớỜờỞởỠỡỢợỤụỦủỨứỪừỬửỮữỰự`)

// Mang cac ky tu thay the khong dau
var DESTINATION_CHARACTERS, _ = stringToRune(`AAAAEEEIIOOOOUUYaaaaeeeiioooouuyAaDdIiUuOoUuAaAaAaAaAaAaAaAaAaAaAaAaEeEeEeEeEeEeEeEeIiIiOoOoOoOoOoOoOoOoOoOoOoOoUuUuUuUuUuUuUu`)

func stringToRune(s string) ([]string, int) {
	ll := utf8.RuneCountInString(s)
	var texts = make([]string, ll+1)
	var index = 0
	for _, runeValue := range s {
		texts[index] = string(runeValue)
		index++
	}
	return texts, ll
}

func binarySearch(sortedArray []string, key string, low int, high int) int {
	var middle int = (low + high) / 2
	if high < low {
		return -1
	}
	if key == sortedArray[middle] {
		return middle
	} else if key < sortedArray[middle] {
		return binarySearch(sortedArray, key, low, middle-1)
	} else {
		return binarySearch(sortedArray, key, middle+1, high)
	}
}

/** * Bo dau 1 ky tu * * @param ch * @return */
func removeAccentChar(ch string) string {
	var index int = binarySearch(SOURCE_CHARACTERS, ch, 0, LL_LENGTH)
	if index >= 0 {
		ch = DESTINATION_CHARACTERS[index]
	}
	return ch
}

/** * Bo dau 1 chuoi * * @param s * @return */
func RemoveAccent(s string) string {
	var buffer bytes.Buffer
	for _, runeValue := range s {
		buffer.WriteString(removeAccentChar(string(runeValue)))
	}
	return buffer.String()
}

/** * Bo dau 1 chuoi * * @param s * @return */
func GenerateCode(s string) string {
	var buffer bytes.Buffer
	for _, runeValue := range s {
		buffer.WriteString(removeAccentChar(string(runeValue)))
	}
	return strings.ReplaceAll(buffer.String(), " ", "-") + "-" + strconv.FormatInt(time.Now().Unix(), 10)
}

func GenOtp(max int) string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func IsValidPassword(s string) bool {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	UnUsed(hasSpecial)
	return hasUpper && hasLower && hasNumber
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
