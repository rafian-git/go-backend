package utility

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	random "math/rand"
	"net/mail"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/rafian-git/go-backend/pkg/apierror"

	"github.com/rafian-git/go-backend/pkg/log"

	"golang.org/x/crypto/bcrypt"
)

var logger = log.New().Named("backend-util")
var ctx = context.Background()

func GenerateOtp(max int) (string, error) {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), err
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GenerateToken(email string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
		msg := "error in generating token"
		logger.Error(ctx, fmt.Errorf("%s: %v", msg, err).Error())
		return "", apierror.New(apierror.NotFound, msg)
	}

	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil)), err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("password hashing failed: %v", err).Error())
		return "", apierror.New(apierror.NotFound, "password hashing failed !")
	}
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		logger.Error(ctx, fmt.Errorf("password didn't match: %v", err).Error())
		return false, apierror.New(apierror.NotFound, "password didn't match !")
	}
	return true, err
}

func IsEmail(emailOrPhone string) (string, error) {
	ifEmail, err := validateEmail(emailOrPhone)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("invalid email: %v", err).Error())
		return "", apierror.New(apierror.NotFound, fmt.Sprintf("invalid email ! - %s", emailOrPhone))
	}
	if ifEmail {
		return "email", err
	}
	return "phone", err
}

func validateEmail(email string) (bool, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("invalid email: %v", err).Error())
		return false, apierror.New(apierror.NotFound, fmt.Sprintf("invalid email ! - %s", email))
	}
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(email), nil
}

func EmailOrPhoneCheck(emailOrPhone string) (string, string, error) {
	var phone, email string = "", ""
	checkPhoneOrMail, err := IsEmail(emailOrPhone)
	if checkPhoneOrMail == "email" {
		email = emailOrPhone
		return email, phone, nil
	} else if checkPhoneOrMail == "phone" {
		phone = emailOrPhone
		return email, phone, nil
	} else {
		return email, phone, err
	}
}

func EmptyStringCheck(value string) bool {
	return value == ""
}

// FindPhoneOrEmailValue - takes in phone and mail 1 of which is an empty string.. returns the non-empty one
func FindPhoneOrEmailValue(phone string, email string) (string, string) {
	if EmptyStringCheck(phone) {
		return "email", email
	}
	return "phone", phone
}

func ValidateBDPhone(phone string) (bool, error) {
	valid := strings.HasPrefix(phone, "01")
	if !valid {
		return valid, fmt.Errorf("phone number must start with 01")
	}
	if len(phone) != 11 {
		return false, fmt.Errorf("phone number must be 11 digits")
	}
	valid, err := StringIsNumberCheck(phone)
	if err != nil {
		return valid, fmt.Errorf("phone number must be numeric")
	}
	return valid, nil
}

func ValidateIntlPhone(phone string) (bool, error) {
	if len(phone) != 10 {
		logger.Error(ctx, "phone number must be of 10 digits !")
		return false, apierror.New(apierror.NotFound, fmt.Sprintf("phone number must be of 10 digits ~ %s", phone))
	}
	valid, err := StringIsNumberCheck(phone)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("phone number must be numeric : %v", err).Error())
		return false, apierror.New(apierror.NotFound, fmt.Sprintf("phone number must be numeric ! ~ %s", phone))
	}
	return valid, nil
}

func StringIsNumberCheck(identifier string) (bool, error) {
	_, err := strconv.Atoi(identifier)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("type checking failed : %v", err).Error())
		return false, apierror.New(apierror.NotFound, "type checking failed !")
	}
	return true, nil
}

func ValidateIdentifier(identifier string, identifierType string) (bool, error) {
	if identifierType == "phone" {
		valid, err := ValidateIntlPhone(identifier)
		if err != nil {
			logger.Error(ctx, fmt.Errorf("error validating phone number: %v", err).Error())
			return valid, apierror.New(apierror.NotFound, fmt.Sprintf("error validating phone number - %s", identifier))
		}
		return valid, nil
	} else if identifierType == "email" {
		valid, err := validateEmail(identifier)
		if err != nil {
			logger.Error(ctx, fmt.Errorf("error validating email : %v", err).Error())
			return valid, apierror.New(apierror.NotFound, fmt.Sprintf("error validating email - %s", identifier))
		}
		return valid, nil
	}

	logger.Error(ctx, fmt.Errorf("invalid identifier or identifier type \nidentifier : %s \nidentifier_type : %s", identifier, identifierType).Error())
	return false, fmt.Errorf("invalid identifier or identifier type \nidentifier : %s \nidentifier_type : %s", identifier, identifierType)
}

func RandomString(length int) string {
	random.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

// "hello""world" >> hello""world
func TrimDoubleQuote(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

const (
	// BootstrapServers : bootstrap servers list
	BootstrapServers string = "BOOTSTRAP_SERVERS"
	// Topic : topic
	Topic string = "TOPIC"
	// GroupID : consumer group
	GroupID string = "GROUP_ID"
	// DelayMs : between sent messages
	DelayMs string = "DELAY_MS"
	// Partition : partition from which to consume
	Partition string = "PARTITION"
)

// GetEnv : returns the environment variable value
func GetEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}

// trims the first character from a string
func TrimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

// trims the first character from a string if it is a plus sign
func TrimFirstPlusSign(s string) string {
	if len(s) > 0 && s[0:1] == "+" {
		s = TrimFirstRune(s)
	}
	return s
}

func EncodeSHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))

	bs := h.Sum(nil)

	return hex.EncodeToString(bs[:])
}

func StringToDate(date string) *time.Time {
	t, err := time.Parse("2006-01-02T15:04:05", date)

	if err != nil {
		return nil
	}
	return &t
}

func ParseStringDateToUnix(layout, date string) int64 {
	t, err := time.Parse(layout, date)

	if err != nil {
		logger.Error(ctx, err.Error())
		return 0
	}

	return t.Unix()
}

// ConvertToFloat: Function recieves a string value and convert it into float64
func ConvertToFloat(cell string) (float64, error) {
	isNegative := false

	if len(cell) > 2 && cell[0] == '(' && cell[len(cell)-1] == ')' {
		isNegative = true
		cell = cell[1 : len(cell)-1]
	}

	value, err := strconv.ParseFloat(strings.ReplaceAll(cell, ",", ""), 64)

	if isNegative {
		value *= -1
	}
	return value, err
}

func ConvertStringToHashedString(str string) (string, error) {
	// Create an MD5 hash object
	hash := md5.New()

	// Convert the concatenated string to bytes and hash it
	hash.Write([]byte(str))

	// Get the hashed result
	hashedBytes := hash.Sum(nil)

	// Convert the hashed result to a hexadecimal string
	hashedString := hex.EncodeToString(hashedBytes)
	return hashedString, nil
}

// GetNextGivenDayOfWeek to get the next occurrence of a specific day of the week - e.g. if targetDay is Sunday and if givenTime is on sunday then it will return next sunday (givenTime + 7 days)
func GetNextGivenDayOfWeek(givenTime time.Time, targetDay time.Weekday) time.Time {
	// Find the current weekday
	currentWeekday := givenTime.Weekday()

	// Calculate the number of days until the target day
	daysUntilTargetDay := (int(targetDay) - int(currentWeekday) + 7) % 7
	if daysUntilTargetDay == 0 {
		daysUntilTargetDay += 7
	}

	// Add the calculated number of days to the current time
	nextDay := givenTime.AddDate(0, 0, daysUntilTargetDay)

	return nextDay
}

// GetPrevGivenDayOfWeek to get the previous occurrence of a specific day of the week
func GetPrevGivenDayOfWeek(givenTime time.Time, targetDay time.Weekday) time.Time {
	// Get the weekday of the given time
	currentWeekday := givenTime.Weekday()

	// Calculate the number of days to subtract
	daysToSubtract := 0
	if currentWeekday > targetDay {
		daysToSubtract = int(currentWeekday) - int(targetDay)
	} else {
		daysToSubtract = 7 - int(targetDay) + int(currentWeekday)
	}

	// Subtract days and return previous occurrence
	return givenTime.AddDate(0, 0, -daysToSubtract)
}

func ConvertStringToWeekday(dayString string) time.Weekday {
	switch dayString {
	case "Sunday":
		return time.Sunday
	case "Monday":
		return time.Monday
	case "Tuesday":
		return time.Tuesday
	case "Wednesday":
		return time.Wednesday
	case "Thursday":
		return time.Thursday
	case "Friday":
		return time.Friday
	case "Saturday":
		return time.Saturday
	default:
		return time.Sunday // Default to Sunday if the input is not recognized
	}
}

func Contains(slice []string, element string) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

func GetNextWeekday(wd time.Weekday) time.Weekday {
	return time.Weekday((int(wd) + 1) % 7)
}

func GetPreviousWeekday(wd time.Weekday) time.Weekday {
	return time.Weekday((int(wd) - 1 + 7) % 7)
}

var daysOfWeek = map[string]time.Weekday{
	"Sunday":    time.Sunday,
	"Monday":    time.Monday,
	"Tuesday":   time.Tuesday,
	"Wednesday": time.Wednesday,
	"Thursday":  time.Thursday,
	"Friday":    time.Friday,
	"Saturday":  time.Saturday,
}

func ParseWeekday(v string) (time.Weekday, error) {
	if d, ok := daysOfWeek[v]; ok {
		return d, nil
	}
	return time.Sunday, fmt.Errorf("invalid weekday '%s'", v)
}

// GetNextAndPrevMartketDay getting next and previous market days from given day
func GetNextAndPrevMartketDay(marketDays []string, dayString string) (nextDay time.Weekday, prevDay time.Weekday, err error) {
	logger.Info(context.Background(), "backend ~ utility ~ getting next and previous market days")

	givenDay, err := ParseWeekday(dayString)
	if err != nil {
		logger.Error(context.Background(), "backend ~ utility ~ error while parsing given day to time.Weekday")
		return 0, 0, fmt.Errorf("error while parsing given day to time.Weekday")
	}
	prevWeekDay, nextWeekDay := givenDay, givenDay
	for {
		prevWeekDay = GetPreviousWeekday(prevWeekDay)
		if Contains(marketDays, prevWeekDay.String()) {
			prevDay = prevWeekDay
			break
		}
	}
	for {
		nextWeekDay = GetNextWeekday(nextWeekDay)
		if Contains(marketDays, nextWeekDay.String()) {
			nextDay = nextWeekDay
			break
		}
	}
	logger.Info(context.Background(), fmt.Sprintf("backend ~ utility ~ for given day : %s ~ previous market day : %s ~ next market day : %s", dayString, prevDay, nextDay))

	// for i := 0; i < len(marketDays); i++ {
	// 	if marketDays[i] == dayString {
	// 		//for next day
	// 	  if i == len(marketDays)-1 {
	// 		nextDay = marketDays[0]
	// 	  } else {
	// 		nextDay = marketDays[i+1]
	// 	  }
	// 	  //for previous day
	// 	  if i == 0 {
	// 		prevDay = marketDays[len(marketDays)-1]
	// 	  } else {
	// 		prevDay = marketDays[i-1]
	// 	  }

	// 	  break
	// 	}
	//   }
	return
}

func ReadJSONFromFile(filename string) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		logger.Error(context.Background(), err.Error())
		return nil, err
	}
	defer f.Close()

	var result interface{}
	if err := json.NewDecoder(f).Decode(&result); err != nil {
		logger.Error(context.Background(), err.Error())
		return nil, err
	}

	return result, nil
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func GetDateFromUnixTime(unixTime int64) time.Time {
	//logger.Info(ctx, fmt.Sprintf("getting date from unix time %d", unixTime))
	date := time.Unix(unixTime, 0).UTC()
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}

func GetDateWithLocationFromUnixTime(unixTime int64, location *time.Location) time.Time {
	//logger.Info(ctx, fmt.Sprintf("getting location specific date from unix UTC time %d", unixTime))
	date := time.Unix(unixTime, 0).UTC().In(location)
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

func DecodeToObj(msg []byte, dest interface{}) error {
	b := bytes.NewBuffer(msg)
	dec := gob.NewDecoder(b)
	err := dec.Decode(dest)
	if err != nil {
		return err
	}
	return nil
}

func GetTimeFromTimeStamp(timeStamp, layout, location string) (time.Time, error) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, err
	}

	parsedTime, err := time.ParseInLocation(layout, timeStamp, loc)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

func StripPrefixToFirstUnderscore(s string) string {
	index := strings.Index(s, "_")
	if index == -1 || index == len(s)-1 {
		return strings.ToLower(s) // no underscore or ends with underscore
	}
	return strings.ToLower(s[index+1:])
}
