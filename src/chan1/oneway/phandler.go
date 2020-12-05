package main

import (
	"fmt"
	"time"
)

type Person struct {
	Name    string
	Age     uint8
	Address Addr
}

type Addr struct {
	city     string
	district string
}

type PersonHandler interface {
	Batch(origs <-chan Person) <-chan Person
	Handle(orig *Person)
}

type PersonHandlerImpl struct{}

func (handler PersonHandlerImpl) Batch(origs <-chan Person) <-chan Person {
	dests := make(chan Person, 100)
	go func() {
		for p := range origs {
			handler.Handle(&p)
			dests <- p
		}
		fmt.Println("All the information has been handled.")
		close(dests)
	}()
	return dests
}

func (handler PersonHandlerImpl) Handle(orig *Person) {
	if orig.Address.district == "Haidian" {
		orig.Address.district = "Shijingshan"
	}
}

var personTotal = 200

var persons []Person = make([]Person, personTotal)

var personCount int

func init() {
	for i := 0; i < 200; i++ {
		name := fmt.Sprintf("%s%d", "P", i)
		p := Person{name, 32, Addr{"Beijing", "Haidian"}}
		persons[i] = p
	}
}

func main() {
	handler := getPersonHandler()
	origs := make(chan Person, 100)
	dests := handler.Batch(origs)
	fecthPerson(origs)
	sign := savePerson(dests)
	<-sign
}

func getPersonHandler() PersonHandler {
	return PersonHandlerImpl{}
}

func savePerson(dest <-chan Person) <-chan byte {
	sign := make(chan byte, 1)
	go func() {
		for {
			p, ok := <-dest
			if !ok {
				fmt.Println("All the information has been saved.")
				sign <- 0
				break
			}
			savePerson1(p)
		}
	}()
	return sign
}

func fecthPerson(origs chan<- Person) {
	origsCap := cap(origs)
	buffered := origsCap > 0
	goTicketTotal := origsCap / 2
	goTicket := initGoTicket(goTicketTotal)
	go func() {
		for {
			p, ok := fecthPerson1()
			if !ok {
				for {
					if !buffered || len(goTicket) == goTicketTotal {
						break
					}
					time.Sleep(time.Nanosecond)
				}
				fmt.Println("All the information has been fetched.")
				close(origs)
				break
			}
			if buffered {
				<-goTicket
				go func() {
					origs <- p
					goTicket <- 1
				}()
			} else {
				origs <- p
			}
		}
	}()
}

func initGoTicket(total int) chan byte {
	var goTicket chan byte
	if total == 0 {
		return goTicket
	}
	goTicket = make(chan byte, total)
	for i := 0; i < total; i++ {
		goTicket <- 1
	}
	return goTicket
}

func fecthPerson1() (Person, bool) {
	if personCount < personTotal {
		p := persons[personCount]
		personCount++
		return p, true
	}
	return Person{}, false
}

func savePerson1(p Person) bool {
	return true
}
