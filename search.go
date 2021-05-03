package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/pkg/errors"
)

// https://apisetu.gov.in/public/api/cowin
const (
	baseURL                     = "https://cdn-api.co-vin.in/api"
	calendarByPinURLFormat      = "/v2/appointment/sessions/public/calendarByPin?pincode=%s&date=%s"
	calendarByDistrictURLFormat = "/v2/appointment/sessions/public/calendarByDistrict?district_id=%d&date=%s"
	listStatesURLFormat         = "/v2/admin/location/states"
	listDistrictsURLFormat      = "/v2/admin/location/districts/%d"
)

type StateList struct {
	States []struct {
		StateID    int    `json:"state_id"`
		StateName  string `json:"state_name"`
		StateNameL string `json:"state_name_l"`
	} `json:"states"`
	TTL int `json:"ttl"`
}

type DistrictList struct {
	Districts []struct {
		StateID       int    `json:"state_id"`
		DistrictID    int    `json:"district_id"`
		DistrictName  string `json:"district_name"`
		DistrictNameL string `json:"district_name_l"`
	} `json:"districts"`
	TTL int `json:"ttl"`
}

type Appointments struct {
	Centers []struct {
		CenterID      int     `json:"center_id"`
		Name          string  `json:"name"`
		NameL         string  `json:"name_l"`
		StateName     string  `json:"state_name"`
		StateNameL    string  `json:"state_name_l"`
		DistrictName  string  `json:"district_name"`
		DistrictNameL string  `json:"district_name_l"`
		BlockName     string  `json:"block_name"`
		BlockNameL    string  `json:"block_name_l"`
		Pincode       int     `json:"pincode"`
		Lat           float64 `json:"lat"`
		Long          float64 `json:"long"`
		From          string  `json:"from"`
		To            string  `json:"to"`
		FeeType       string  `json:"fee_type"`
		VaccineFees   []struct {
			Vaccine string `json:"vaccine"`
			Fee     string `json:"fee"`
		} `json:"vaccine_fees"`
		Sessions []struct {
			SessionID         string   `json:"session_id"`
			Date              string   `json:"date"`
			AvailableCapacity int      `json:"available_capacity"`
			MinAgeLimit       int      `json:"min_age_limit"`
			Vaccine           string   `json:"vaccine"`
			Slots             []string `json:"slots"`
		} `json:"sessions"`
	} `json:"centers"`
}

func timeNow() string {
	return time.Now().Format("02-01-2006")
}

func queryServer(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "hi_IN")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Request failed with statusCode: %d", resp.StatusCode))
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func searchByPincode(pinCode string) error {
	response, err := queryServer(fmt.Sprintf(calendarByPinURLFormat, pinCode, timeNow()))
	if err != nil {
		return errors.Wrap(err, "Failed to fetch appointment sessions")
	}
	return getAvailableSessions(response, age)
}

func getStateIDByName(state string) (int, error) {
	response, err := queryServer(listStatesURLFormat)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to list states")
	}
	states := StateList{}
	if err := json.Unmarshal(response, &states); err != nil {
		return 0, err
	}
	for _, s := range states.States {
		if strings.ToLower(s.StateName) == strings.ToLower(state) {
			return s.StateID, nil
		}
	}
	return 0, errors.New("Invalid state name passed")
}

func getDistrictIDByName(stateID int, district string) (int, error) {
	response, err := queryServer(fmt.Sprintf(listDistrictsURLFormat, stateID))
	if err != nil {
		return 0, errors.Wrap(err, "Failed to list states")
	}
	dl := DistrictList{}
	if err := json.Unmarshal(response, &dl); err != nil {
		return 0, err
	}
	for _, d := range dl.Districts {
		if strings.ToLower(d.DistrictName) == strings.ToLower(district) {
			return d.DistrictID, nil
		}
	}
	return 0, errors.New("Invalid district name passed")
}

func searchByStateDistrict(age int, state, district string) error {
	stateID, err := getStateIDByName(state)
	if err != nil {
		return err
	}
	districtID, err := getDistrictIDByName(stateID, district)
	if err != nil {
		return err
	}
	response, err := queryServer(fmt.Sprintf(calendarByDistrictURLFormat, districtID, timeNow()))
	if err != nil {
		return errors.Wrap(err, "Failed to fetch appointment sessions")
	}
	return getAvailableSessions(response, age)
}

func getAvailableSessions(response []byte, age int) error {
	appnts := Appointments{}
	err := json.Unmarshal(response, &appnts)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 1, 8, 1, '\t', 0)
	for _, center := range appnts.Centers {
		for _, s := range center.Sessions {
			if s.MinAgeLimit <= age && s.AvailableCapacity != 0 {
				fmt.Fprintln(w, fmt.Sprintf("Center\t%s", center.Name))
				fmt.Fprintln(w, fmt.Sprintf("State\t%s", center.StateName))
				fmt.Fprintln(w, fmt.Sprintf("District\t%s", center.DistrictName))
				fmt.Fprintln(w, fmt.Sprintf("PinCode\t%d", center.Pincode))
				fmt.Fprintln(w, fmt.Sprintf("Fee\t%s", center.FeeType))
				if len(center.VaccineFees) != 0 {
					fmt.Fprintln(w, fmt.Sprintf("Vaccine\t"))
				}
				for _, v := range center.VaccineFees {
					fmt.Fprintln(w, fmt.Sprintf("\tName\t%s", v.Vaccine))
					fmt.Fprintln(w, fmt.Sprintf("\tFees\t%s", v.Fee))
				}
				fmt.Fprintln(w, fmt.Sprintf("Sessions\t"))
				fmt.Fprintln(w, fmt.Sprintf("\tDate\t%s", s.Date))
				fmt.Fprintln(w, fmt.Sprintf("\tAvailableCapacity\t%d", s.AvailableCapacity))
				fmt.Fprintln(w, fmt.Sprintf("\tMinAgeLimit\t%d", s.MinAgeLimit))
				fmt.Fprintln(w, fmt.Sprintf("\tVaccine\t%s", s.Vaccine))
				fmt.Fprintln(w, fmt.Sprintf("\tSlots"))
				for _, slot := range s.Slots {
					fmt.Fprintln(w, fmt.Sprintf("\t\t%s", slot))
				}
				fmt.Fprintln(w, "-----------------------------")
			}
		}
	}
	if err := w.Flush(); err != nil {
		return err
	}
	if buf.Len() == 0 {
		log.Print("No slots available, rechecking")
		return nil
	}
	log.Print("Found available slots, sending email")
	return sendMail(email, password, buf.String())
}
