package domain

import (
	"fmt"
)

type Trip struct {
	ID             string
	DriverId       string
	FromPoint      string
	ToPoint        string
	DepartureTime  string
	AvailableSeats string
}

func (t Trip) toString() string {

	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s",
		t.ID, t.DriverId, t.FromPoint, t.ToPoint, t.DepartureTime, t.AvailableSeats)
}

type ShareTrip struct {
	Trips []Trip
}

func NewShareTrip() *ShareTrip {
	return &ShareTrip{}
}

func (sht *ShareTrip) AddTrip(trip Trip) error {
	_, ok := sht.indexOf(trip.ID)
	if ok {
		return ErrAlreadyExists
	}
	sht.Trips = append(sht.Trips, trip)
	return nil
}

func (sht *ShareTrip) GetTrip() []Trip {
	return sht.Trips
}

func (sht *ShareTrip) indexOf(id string) (int, bool) {
	for i, trip := range sht.Trips {
		if trip.ID == id {
			return i, true
		}
	}
	return -1, false
}

//func (sht *ShareTrip) DeleteTrip(name string) {
//	for i := 0; i < len(sht.Trips); i++ {
//		if sht.Trips[i].Name == name {
//			res := sht.Trips[i].toString()
//			sht.Trips = slices.Delete(sht.Trips, i, i+1)
//			fmt.Printf("Trip '%s' was deleted:\n", res)
//			return
//		}
//	}
//	fmt.Println("There is no trip with this name")
//}
//
//func (sht *ShareTrip) UpdateTrip(trip Trip) error {
//	index, ok := sht.indexOf(trip.ID)
//	if !ok {
//		return ErrNotFound
//	}
//	sht.Trips[index] = trip
//	return nil
//}
//
//func (sht *ShareTrip) FindTrip(name string) {
//	res := make([]string, 0)
//	for i := 0; i < len(sht.Trips); i++ {
//		if strings.Contains(sht.Trips[i].Name, name) {
//			res = append(res, sht.Trips[i].toString())
//		}
//	}
//	if len(res) == 0 {
//		fmt.Printf("There is no trip containing text %s:\n", name)
//		return
//	}
//	fmt.Printf("These trips containing text %s: were found:\n", name)
//	fmt.Println(strings.Join(res, ",\n"))
//}
