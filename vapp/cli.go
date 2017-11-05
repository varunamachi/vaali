package vapp

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
	cli "gopkg.in/urfave/cli.v1"
)

//AskSecret - asks password from user, does not echo charectors
func askSecret() (secret string, err error) {
	var pbyte []byte
	pbyte, err = terminal.ReadPassword(int(syscall.Stdin))
	if err == nil {
		secret = string(pbyte)
		fmt.Println()
	}
	return secret, err
}

//ArgGetter - this struct and its method are helpers to combine getting args
//from commandline arguments or from reading from console. Also handles errors
//when required arguments are not provided
type ArgGetter struct {
	Ctx *cli.Context
	Err error
}

func readInput(text *string) (err error) {
	scanner := bufio.NewScanner(os.Stdin)
	*text = ""
	for scanner.Scan() {
		*text = scanner.Text()
		break
	}
	err = scanner.Err()
	return err
}

//GetString - gives a string argument either from commandline or from blocking
//user input, this method doesnt complain even if the arg-value is empty
func (retriever *ArgGetter) GetString(key string) (val string) {
	if retriever.Err != nil {
		return val
	}
	val = retriever.Ctx.String(key)
	if !retriever.Ctx.IsSet(key) && len(val) == 0 {
		fmt.Print(key + ": ")
		err := readInput(&val)
		if err != nil {
			val = ""
		}
	}
	return val
}

//GetRequiredString - gives a string argument either from commandline or from
//blocking user input, this method sets the error if required arg-val is empty
func (retriever *ArgGetter) GetRequiredString(key string) (val string) {
	if retriever.Err != nil {
		return val
	}
	val = retriever.Ctx.String(key)
	if !retriever.Ctx.IsSet(key) && len(val) == 0 {
		fmt.Print(key + "*: ")
		err := readInput(&val)
		if err != nil || len(val) == 0 {
			val = ""
			retriever.Err = fmt.Errorf("Required argument %s not provided", key)
		}
	}
	return val
}

//GetRequiredSecret - gives a string argument either from commandline or from
//blocking user input, this method sets the error if required arg-val is empty
func (retriever *ArgGetter) GetRequiredSecret(key string) (val string) {
	if retriever.Err != nil {
		return val
	}
	val = retriever.Ctx.String(key)
	if !retriever.Ctx.IsSet(key) && len(val) == 0 {
		fmt.Print(key + "*: ")
		var err error
		val, err = askSecret()
		if err != nil || len(val) == 0 {
			val = ""
			retriever.Err = fmt.Errorf("Required argument %s not provided", key)
		}
	}
	return val
}

//GetRequiredInt - gives a Integer argument either from commandline or from
//blocking user input, this method sets the error if required arg-val is empty
func (retriever *ArgGetter) GetRequiredInt(key string) (val int) {
	if retriever.Err != nil {
		return val
	}
	val = retriever.Ctx.Int(key)
	if !retriever.Ctx.IsSet(key) && val == 0 {
		fmt.Print(key + ": ")
		var strval string
		err := readInput(&strval)
		if err != nil || len(strval) == 0 {
			val = 0
			retriever.Err = fmt.Errorf("Required argument %s not provided", key)
		} else {
			val, err = strconv.Atoi(strval)
			if err != nil {
				retriever.Err = fmt.Errorf("Invalid value for %s given", key)
				val = 0
			}
		}
	}
	return val
}

//GetInt - gives a Integer argument either from commandline or from blocking
//user input, this method doesnt complain even if the arg-value is empty
func (retriever *ArgGetter) GetInt(key string) (val int) {
	if retriever.Err != nil {
		return val
	}
	val = retriever.Ctx.Int(key)
	if !retriever.Ctx.IsSet(key) && val == 0 {
		fmt.Print(key + ": ")
		var strval string
		err := readInput(&strval)
		if err != nil || len(strval) == 0 {
			val = 0
		} else {
			val, err = strconv.Atoi(strval)
			if err != nil {
				val = 0
			}
		}
	}
	return val
}

//GetRequiredBool - gives a Boolean argument either from commandline or from
//blocking user input, this method sets the error if required arg-val is empty
func (retriever *ArgGetter) GetRequiredBool(key string) (val bool) {
	if retriever.Err != nil {
		return val
	}
	val = retriever.Ctx.Bool(key)
	// if !retriever.Ctx.IsSet(key) {
	// 	fmt.Print(key + ": ")
	// 	var strval string
	// 	err := readInput(&strval)
	// 	trimmed := strings.TrimSpace(strval)
	// 	if err != nil || len(trimmed) == 0 {
	// 		val = false
	// 		retriever.Err = fmt.Errorf("Required argument %s not provided", key)
	// 	} else {
	// 		val = strings.ToUpper(trimmed) == "TRUE" || trimmed == "1"
	// 		if err != nil {
	// 			retriever.Err = fmt.Errorf("Invalid value for %s given", key)
	// 			val = false
	// 		}
	// 	}
	// }
	return val
}

//GetBool - gives a Boolean argument either from commandline or from blocking
//user input, this method doesnt complain even if the arg-value is empty
func (retriever *ArgGetter) GetBool(key string) (val bool) {
	if retriever.Err != nil {
		return val
	}
	val = retriever.Ctx.Bool(key)
	// if !retriever.Ctx.IsSet(key) {
	// 	fmt.Print(key + ": ")
	// 	var strval string
	// 	err := readInput(&strval)
	// 	if err != nil || len(strval) == 0 {
	// 		val = false
	// 	} else {
	// 		trimmed := strings.TrimSpace(strval)
	// 		val = strings.ToUpper(trimmed) == "TRUE" || trimmed == "1"
	// 		if err != nil {
	// 			val = false
	// 		}
	// 	}
	// }
	return val
}
