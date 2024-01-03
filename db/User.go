package db

import (
	"encoding/json"
	"errors"
	"regexp"
	"structs"
	"unicode"

	"go.etcd.io/bbolt"
)

var (
	ErrUserOpenDatabase = errors.New("can't open database")
	ErrUserExist        = errors.New("user is already exists")
	ErrUsereNotFound    = errors.New("user Doesn't exist")

	ErrUserNameIsShortOrEmpty = errors.New("username is short or empty")
	ErrUserEmailInvalid       = errors.New("Email is Invalid ")

	ErrUserPasswordIsSome         = errors.New("the same password")
	ErrUserPasswordIsShortOrEmpty = errors.New("Password is short or empty")

	ErrUserAddrCity   = errors.New("City is Invalid")
	ErrUserAddrStreet = errors.New("Street is Invalid")
	ErrUserAddrArea   = errors.New("Area is Invalid")

	ErrUservisaInvalidNumber         = errors.New("Visa Number is Invalid")
	ErrUservisaInvalidCvv            = errors.New("Visa Cvv is Invalid")
	ErrUservisaInvalidExpirationDate = errors.New("Visa ExpirationDate is Invalid")
	ErrUservisaNumberIsExit          = errors.New("visa Number Is Exit")
	ErrUservisaIsExit                = errors.New("visa  Already added")
	ErrUserPhoneNumberInvalid        = errors.New("is NOT a valid phone number")
	ErrUserFirstNameInvalid          = errors.New("First Name  is Invalid")
	ErrUserNameInvalid               = errors.New("Name  is Invalid")
	ErrUserNotStrongPassword         = errors.New("is NOT a strong password")

	ErrUserNameIslong          = errors.New("username Must be less than 20")
	ErrUserNameNotStartLetter  = errors.New("username Must start with a letter")
	ErrUserNameNotAlphanumeric = errors.New("username Must all characters are alphanumeric")
)

type user struct {
	dataBase *bbolt.DB
}

func (u *user) AddNew(users structs.User) error {
	return u.dataBase.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))

		if b == nil {
			return ErrUserOpenDatabase
		}
		checkIfuserExist := b.Get([]byte(users.Username))
		if checkIfuserExist != nil {
			return ErrUserExist
		}
		if err := isAcceptableUsername(users.Username); err != nil {
			return err
		}

		if err := isStrongPassword(users.Password); err != nil {
			return err
		}

		if !isValidEmail(users.UserEmail) {
			return ErrUserEmailInvalid
		}

		dataUser, err := json.Marshal(users)
		if err != nil {
			return err
		}
		b.Put([]byte(users.Username), dataUser)
		return nil
	})
}

func (u *user) GetUser(username string, users *structs.User) error {
	return u.dataBase.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		data := b.Get([]byte(username))
		if data == nil {
			return ErrUsereNotFound
		}
		err := json.Unmarshal(data, users)
		if err != nil {
			return err
		}

		return nil
	})
}
func (u *user) UpdataPassword(username string, oldPassword string, newPassword string) error {
	return u.dataBase.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		userdata := b.Get([]byte(username))
		if userdata == nil {
			return ErrUsereNotFound
		}
		if err := isStrongPassword(newPassword); err != nil {
			return err
		}

		usr := structs.User{}
		if err := json.Unmarshal(userdata, &usr); err != nil {
			return err
		}
		if oldPassword != usr.Password {
			return errors.New("woring password")
		}

		usr.Password = newPassword
		data, err := json.Marshal(usr)
		if err != nil {
			return err
		}
		b.Put([]byte(usr.Username), data)

		return nil
	})

}

func (u *user) UpdataAddr(username string, addr structs.Addr) error {
	return u.dataBase.Batch(func(tx *bbolt.Tx) error {
		if addr.City == "" {
			return ErrUserAddrCity
		}
		if addr.Street == "" {
			return ErrUserAddrStreet
		}
		if len(addr.Area) == 0 {
			return ErrUserAddrArea
		}

		b := tx.Bucket([]byte("users"))

		userdata := b.Get([]byte(username))
		if userdata == nil {
			return ErrUserExist
		}
		usr := structs.User{}
		if err := json.Unmarshal(userdata, &usr); err != nil {
			return err
		}

		usr.UserAddr = addr

		data, err := json.Marshal(usr)
		if err != nil {
			return err
		}
		b.Put([]byte(usr.Username), data)

		return nil
	})

}

func (u *user) AddVisa(username string, visa structs.Visa) error {
	return u.dataBase.Batch(func(tx *bbolt.Tx) error {
		// check if visa valid

		if len(visa.Number) != 16 {
			return ErrUservisaInvalidNumber
		}
		if len(visa.Cvv) != 3 {
			return ErrUservisaInvalidCvv
		}
		if len(visa.ExpirationData) != 5 {
			return ErrUservisaInvalidExpirationDate
		}
		for _, v := range visa.Number {
			if !unicode.IsNumber(v) {
				return ErrUservisaInvalidNumber
			}
		}
		for _, v := range visa.Cvv {
			if !unicode.IsNumber(v) {
				return ErrUservisaInvalidCvv
			}
		}
		if !unicode.IsNumber(rune(visa.ExpirationData[0])) || !unicode.IsNumber(rune(visa.ExpirationData[1])) ||
			visa.ExpirationData[2] != '/' || !unicode.IsNumber(rune(visa.ExpirationData[3])) || !unicode.IsNumber(rune(visa.ExpirationData[4])) {
			return ErrUservisaInvalidExpirationDate
		}

		b := tx.Bucket([]byte("users"))

		usr := structs.User{}
		if err := u.GetUser(username, &usr); err != nil {
			return err
		}

		if visa.Number == usr.UserVisa.Number {
			return ErrUservisaIsExit
		}
		usr.UserVisa = visa

		data, err := json.Marshal(usr)
		if err != nil {
			return err
		}
		b.Put([]byte(usr.Username), data)

		return nil
	})

}

func (u *user) DeleteVisa(username, first3Number, last3Number string) error {
	return u.dataBase.Batch(func(tx *bbolt.Tx) error {

		for _, v := range first3Number {
			if !unicode.IsNumber(v) {
				return errors.New("first 3 number is not number")
			}
		}

		for _, v := range last3Number {
			if !unicode.IsNumber(v) {
				return errors.New("last 3 number is not number")
			}
		}

		b := tx.Bucket([]byte("users"))

		usr := structs.User{}
		if err := u.GetUser(username, &usr); err != nil {
			return err
		}




		usr.UserVisa = structs.Visa{}

		data, err := json.Marshal(usr)
		if err != nil {
			return err
		}
		b.Put([]byte(usr.Username), data)

		return nil
	})

}

func (u *user) UpdataPhone(username string, phone string) error {
	return u.dataBase.Batch(func(tx *bbolt.Tx) error {

		re := regexp.MustCompile(`^(079|078|077)\d{7}$`)
		if !re.MatchString(phone) {
			return ErrUserPhoneNumberInvalid
		}
		b := tx.Bucket([]byte("users"))

		usr := structs.User{}
		if err := u.GetUser(username, &usr); err != nil {
			return err
		}

		// updata number phone to dataBase
		usr.Phone = phone

		data, err := json.Marshal(usr)
		if err != nil {
			return err
		}
		b.Put([]byte(usr.Username), data)

		return nil
	})

}

func (u *user) UpdateName(username string, Name string) error {
	return u.dataBase.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if len(Name) == 0 {
			return ErrUserNameInvalid
		}

		usr := structs.User{}
		if err := u.GetUser(username, &usr); err != nil {
			return err
		}

		usr.Name = Name

		data, err := json.Marshal(usr)
		if err != nil {
			return err
		}
		b.Put([]byte(usr.Username), data)

		return nil
	})

}

func isStrongPassword(password string) error {
	// Check minimum length
	if len(password) < 8 {
		return errors.New("pasword is short")
	}

	// Check for at least one uppercase letter
	hasUpperCase := false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpperCase = true
			break
		}
	}

	if !hasUpperCase {
		return errors.New("Password least one uppercase letter")
	}

	// Check for at least one lowercase letter
	hasLowerCase := false
	for _, char := range password {
		if unicode.IsLower(char) {
			hasLowerCase = true
			break
		}
	}

	if !hasLowerCase {
		return errors.New(" Password least one lowercase letter")
	}

	hasDigit := false
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasDigit = true
			break
		}
	}

	if !hasDigit {
		return errors.New("Password at least one digit")
	}

	// Check for at least one special character
	hasSpecialChar := false
	specialChars := "!@#$%^&*()-=_+[]{}|;:'\",.<>?/"
	for _, char := range password {
		if unicode.IsSymbol(char) || unicode.IsPunct(char) {
			for _, special := range specialChars {
				if char == special {
					hasSpecialChar = true
					break
				}
			}
		}
		if hasSpecialChar {
			break
		}
	}
	if !hasSpecialChar {
		return errors.New("Password least one special character")
	}
	return nil
}

func isAcceptableUsername(username string) error {
	// Check length
	if len(username) < 4 {
		return ErrUserNameIsShortOrEmpty
	}
	if len(username) > 20 {
		return ErrUserNameIslong
	}

	// Check if the first character is a letter
	if !unicode.IsLetter(rune(username[0])) {
		return ErrUserNameNotStartLetter
	}

	// Check if all characters are alphanumeric
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return ErrUserNameNotAlphanumeric

		}
	}

	return nil
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}
