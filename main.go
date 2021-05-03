package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	pinCode, state, district, email, password string

	age int

	rootCmd = &cobra.Command{
		Use:   "covaccine-notifier [FLAGS]",
		Short: "CoWIN Vaccine availability notifier India",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(args)
		},
	}
)

const (
	pinCodeEnv       = "PIN_CODE"
	stateNameEnv     = "STATE_NAME"
	districtNameEnv  = "DISTRICT_NAME"
	ageEnv           = "AGE"
	emailIDEnv       = "EMAIL_ID"
	emailPasswordEnv = "EMAIL_PASSOWORD"

	searchInterval = time.Minute * 3
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&pinCode, "pincode", "c", os.Getenv(pinCodeEnv), "Search by pin code")
	rootCmd.PersistentFlags().StringVarP(&state, "state", "s", os.Getenv(stateNameEnv), "Search by state name")
	rootCmd.PersistentFlags().StringVarP(&district, "district", "d", os.Getenv(districtNameEnv), "Search by district name")
	rootCmd.PersistentFlags().IntVarP(&age, "age", "a", getIntEnv(ageEnv), "Search appointment for age")
	rootCmd.PersistentFlags().StringVarP(&email, "email", "e", os.Getenv(emailIDEnv), "Email address to send notifications")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", os.Getenv(emailPasswordEnv), "Email ID password for auth")
}

// Execute executes the main command
func Execute() error {
	return rootCmd.Execute()
}

func checkFlags() error {
	if len(pinCode) == 0 &&
		len(state) == 0 &&
		len(district) == 0 {
		return errors.New("Please pass one of the pinCode or state & district name combination options")
	}
	if len(pinCode) == 0 && (len(state) == 0 || len(district) == 0) {
		return errors.New("Missing state or district name option")
	}
	if age == 0 {
		return errors.New("Missing age option")
	}
	if len(email) == 0 || len(password) == 0 {
		return errors.New("Missing email creds")
	}
	return nil
}

func main() {
	Execute()
}

func getIntEnv(envVar string) int {
	v := os.Getenv(envVar)
	if len(v) == 0 {
		return 0
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func Run(args []string) error {
	if err := checkFlags(); err != nil {
		return err
	}
	if err := checkSlots(); err != nil {
		return err
	}
	ticker := time.NewTicker(searchInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := checkSlots(); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkSlots() error {
	// Search for slots
	var err error
	if len(pinCode) != 0 {
		err = searchByPincode(pinCode)
	}
	err = searchByStateDistrict(age, state, district)
	return err
}
