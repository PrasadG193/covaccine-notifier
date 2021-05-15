package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	pinCode, state, district, email, password, date, vaccine, fee string

	age, interval, minCapacity int

	rootCmd = &cobra.Command{
		Use:   "covaccine-notifier [FLAGS]",
		Short: "CoWIN Vaccine availability notifier India",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(args)
		},
	}
)

const (
	pinCodeEnv        = "PIN_CODE"
	stateNameEnv      = "STATE_NAME"
	districtNameEnv   = "DISTRICT_NAME"
	ageEnv            = "AGE"
	emailIDEnv        = "EMAIL_ID"
	emailPasswordEnv  = "EMAIL_PASSOWORD"
	searchIntervalEnv = "SEARCH_INTERVAL"
	vaccineEnv        = "VACCINE"
	feeEnv            = "FEE"
	minCapacityEnv    = "MIN_CAPACITY"

	defaultSearchInterval = 60
	defaultMinCapacity = 1

	covishield = "covishield"
	covaxin    = "covaxin"

	free = "free"
	paid = "paid"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&pinCode, "pincode", "c", os.Getenv(pinCodeEnv), "Search by pin code")
	rootCmd.PersistentFlags().StringVarP(&state, "state", "s", os.Getenv(stateNameEnv), "Search by state name")
	rootCmd.PersistentFlags().StringVarP(&district, "district", "d", os.Getenv(districtNameEnv), "Search by district name")
	rootCmd.PersistentFlags().IntVarP(&age, "age", "a", getIntEnv(ageEnv), "Search appointment for age")
	rootCmd.PersistentFlags().StringVarP(&email, "email", "e", os.Getenv(emailIDEnv), "Email address to send notifications")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", os.Getenv(emailPasswordEnv), "Email ID password for auth")
	rootCmd.PersistentFlags().IntVarP(&interval, "interval", "i", getIntEnv(searchIntervalEnv), fmt.Sprintf("Interval to repeat the search. Default: (%v) second", defaultSearchInterval))
	rootCmd.PersistentFlags().StringVarP(&vaccine, "vaccine", "v", os.Getenv(vaccineEnv), fmt.Sprintf("Vaccine preferences - covishield (or) covaxin. Default: No preference"))
	rootCmd.PersistentFlags().StringVarP(&fee, "fee", "f", os.Getenv(feeEnv), fmt.Sprintf("Fee preferences - free (or) paid. Default: No preference"))
    rootCmd.PersistentFlags().IntVarP(&minCapacity, "minCapacity", "m", getIntEnv(minCapacityEnv), fmt.Sprintf("Filter by minimum vaccination capacity. Default: (%v)", defaultMinCapacity))
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
	if interval == 0 {
		interval = defaultSearchInterval
	}
	if !(vaccine == "" || vaccine == covishield || vaccine == covaxin) {
		return errors.New("Invalid vaccine, please use covaxin or covishield")
	}
	if !(fee == "" || fee == free || fee == paid) {
		return errors.New("Invalid fee preference, please use free or paid")
	}
	if minCapacity == 0 {
	    minCapacity = defaultMinCapacity
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
	ticker := time.NewTicker(time.Second * time.Duration(interval))
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
	if len(pinCode) != 0 {
		return searchByPincode(pinCode)
	}
	return searchByStateDistrict(age, state, district)
}
