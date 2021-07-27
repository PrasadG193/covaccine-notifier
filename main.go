package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/PrasadG193/covaccine-notifier/pkg/notify"
)

var (
	pinCode, state, district, vaccine, fee   string
	username, password, token, mattermostURL string
	age, interval, minCapacity, dose         int

	rootCmd = &cobra.Command{
		Use:   "covaccine-notifier [FLAGS]",
		Short: "CoWIN Vaccine availability notifier India",
	}

	telegramCmd = &cobra.Command{
		Use:   "telegram [FLAGS]",
		Short: "Notify slots availability using Telegram",
		RunE: func(cmd *cobra.Command, args []string) error {
			notifier, err := notify.NewTelegram(username, token)
			if err != nil {
				return err
			}
			return Run(args, notifier)
		},
	}

	mattermostCmd = &cobra.Command{
		Use:   "mattermost [FLAGS]",
		Short: "Notify slots availability using Mattermost",
		RunE: func(cmd *cobra.Command, args []string) error {
			notifier, err := notify.NewMattermost(mattermostURL, token, username)
			if err != nil {
				return err
			}
			return Run(args, notifier)
		},
	}

	emailCmd = &cobra.Command{
		Use:   "email [FLAGS]",
		Short: "Notify slots availability using Email",
		RunE: func(cmd *cobra.Command, args []string) error {
			notifier := notify.NewEmail(username, password)
			return Run(args, notifier)
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
	tgApiTokenEnv     = "TG_TOKEN"
	tgUsernameEnv     = "TG_USERNAME"
	mmURLEnv          = "MATTERMOST_URL"
	mmUserEnv         = "MATTERMOST_USERNAME"
	mmTokenEnv        = "MATTERMOST_TOKEN"
	minCapacityEnv    = "MIN_CAPACITY"
	doseEnv           = "DOSE"

	defaultSearchInterval = 60
	defaultMinCapacity    = 1

	covishield = "covishield"
	covaxin    = "covaxin"

	free = "free"
	paid = "paid"
)

func init() {
	rootCmd.PersistentFlags().IntVarP(&age, "age", "a", getIntEnv(ageEnv), "Search appointment for age (required)")
	rootCmd.MarkPersistentFlagRequired("age")
	rootCmd.PersistentFlags().StringVarP(&pinCode, "pincode", "c", os.Getenv(pinCodeEnv), "Search by pin code")
	rootCmd.PersistentFlags().StringVarP(&state, "state", "s", os.Getenv(stateNameEnv), "Search by state name")
	rootCmd.PersistentFlags().StringVarP(&district, "district", "d", os.Getenv(districtNameEnv), "Search by district name")
	rootCmd.PersistentFlags().IntVarP(&interval, "interval", "i", getIntEnv(searchIntervalEnv), fmt.Sprintf("Interval to repeat the search. Default: (%v) second", defaultSearchInterval))
	rootCmd.PersistentFlags().StringVarP(&vaccine, "vaccine", "v", os.Getenv(vaccineEnv), fmt.Sprintf("Vaccine preferences - covishield (or) covaxin. Default: No preference"))
	rootCmd.PersistentFlags().StringVarP(&fee, "fee", "f", os.Getenv(feeEnv), fmt.Sprintf("Fee preferences - free (or) paid. Default: No preference"))
	rootCmd.PersistentFlags().IntVarP(&minCapacity, "min-capacity", "m", getIntEnv(minCapacityEnv), fmt.Sprintf("Filter by minimum vaccination capacity. Default: (%v)", defaultMinCapacity))
	rootCmd.PersistentFlags().IntVarP(&dose, "dose", "o", getIntEnv(doseEnv), "Dose preference - 1 or 2. Default: 0 (both)")

	rootCmd.AddCommand(emailCmd, telegramCmd, mattermostCmd)

	emailCmd.PersistentFlags().StringVarP(&username, "username", "u", os.Getenv(emailIDEnv), "Email address to send notifications")
	emailCmd.MarkPersistentFlagRequired("username")
	emailCmd.PersistentFlags().StringVarP(&password, "password", "p", os.Getenv(emailPasswordEnv), "Email ID password for auth")
	emailCmd.MarkPersistentFlagRequired("password")

	telegramCmd.PersistentFlags().StringVarP(&username, "username", "u", os.Getenv(tgUsernameEnv), "telegram username")
	telegramCmd.MarkPersistentFlagRequired("username")
	telegramCmd.PersistentFlags().StringVarP(&token, "token", "t", os.Getenv(tgApiTokenEnv), "telegram bot API token")
	telegramCmd.MarkPersistentFlagRequired("token")

	mattermostCmd.PersistentFlags().StringVarP(&mattermostURL, "url", "l", os.Getenv(mmURLEnv), "mattermost server url")
	mattermostCmd.MarkFlagRequired("url")
	mattermostCmd.PersistentFlags().StringVarP(&username, "username", "u", os.Getenv(mmUserEnv), "mattermost username")
	mattermostCmd.MarkFlagRequired("username")
	mattermostCmd.PersistentFlags().StringVarP(&token, "token", "t", os.Getenv(mmTokenEnv), "mattermost bot API token")
	mattermostCmd.MarkFlagRequired("token")
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
	if dose < 0 || dose > 2 {
		return errors.New("Invalid dose preference, please use 1 or 2")
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

func Run(args []string, notifier notify.Notifier) error {
	if err := checkFlags(); err != nil {
		return err
	}
	if err := checkSlots(notifier); err != nil {
		return err
	}
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := checkSlots(notifier); err != nil {
				return err
			}
		}
	}
}

func checkSlots(notifier notify.Notifier) error {
	// Search for slots
	if len(pinCode) != 0 {
		return searchByPincode(notifier, pinCode)
	}
	return searchByStateDistrict(notifier, state, district)
}
