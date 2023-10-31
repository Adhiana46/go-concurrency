package main

import (
	"time"

	"github.com/fatih/color"
)

type Barbershop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientsChan     chan string
	Open            bool
}

func (shop *Barbershop) addBarber(barberName string) {
	shop.NumberOfBarbers++

	go func() {
		isSleeping := false
		color.Yellow("%s goes to the waiting room to check for clients.", barberName)

		for {
			// if there are no clients, the buffer goes to sleep
			if len(shop.ClientsChan) == 0 {
				color.Yellow("There is nothing to do, so %s takes a nap.", barberName)
				isSleeping = true
			}

			client, shopOpen := <-shop.ClientsChan

			if !shopOpen {
				// shop is closed, so send the barber home and close this goroutine
				shop.sendBarberHome(barberName)
				return
			}

			// barber is sleeping, client wakes him up
			if isSleeping {
				color.Yellow("%s wakes %s up.", client, barberName)
				isSleeping = false
			}

			// cut hair
			shop.cutHair(barberName, client)
		}
	}()
}

func (shop *Barbershop) cutHair(barberName, clientName string) {
	color.Green("%s is cutting %s's hair.", barberName, clientName)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s is finished cutting %s's hair.", barberName, clientName)
}

func (shop *Barbershop) sendBarberHome(barberName string) {
	color.Cyan("%s is going home.", barberName)
	shop.BarbersDoneChan <- true
}

func (shop *Barbershop) closeShop() {
	color.Cyan("Closing shop for the day.")

	close(shop.ClientsChan)
	shop.Open = false

	for a := 1; a <= shop.NumberOfBarbers; a++ {
		<-shop.BarbersDoneChan
	}
	close(shop.BarbersDoneChan)

	color.Green("---------------------------------------------------------------------")
	color.Green("The barbershop is now closed for the day, and everyone has gone home.")
}

func (shop *Barbershop) addClient(clientName string) {
	color.Green("*** %s arrives!", clientName)

	if !shop.Open {
		color.Red("The shop is already closed, so %s leaves!", clientName)
		return
	}

	select {
	case shop.ClientsChan <- clientName:
		color.Yellow("%s takes a seat in the waiting room.", clientName)
	default:
		color.Red("The waiting room is full, so %s leaves.", clientName)
	}
}
