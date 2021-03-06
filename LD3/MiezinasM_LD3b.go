/*
============================================================================
Mantas Miežinas, IFF-2
LD3a
============================================================================
*/

package main

import "fmt"
import "os"
import "bufio"
import "io"
import "strings"
import "strconv"

//konstanta duomenu failo pavadinimui
const DataFile string = "MiezinasM_L2.txt"

/*
============================================================================

	Bazine struktura saugoti vienam modelio irasui

============================================================================
*/
type model struct {
	name     string
	quantity uint
	price    float32
}

/*
============================================================================

	Skaitliukas su metodais, kad butu lengviau dirbti

============================================================================
*/
type counter struct {
	price float32
	count uint
}

type counterSlice []counter

/*
============================================================================

	Esamo skaitliuko pozicijos radimas, jei tokio nera, grazinama -1

============================================================================
*/
func (slice *counterSlice) find(value float32) int {
	for pos, v := range *slice {
		if v.price == value {
			return pos
		}
	}
	return -1
}

/*
============================================================================

	Funkcija rasti pozicija, i kuria reikia iterpti skaitliuka,
	jei pozicija pabaigoje slice'o, tuomet graziname -1

============================================================================
*/
func (slice *counterSlice) getPos(value float32) int {
	for pos, v := range *slice {
		if v.price > value {
			return pos
		}
	}
	return -1
}

/*
============================================================================

	konteinerine struktura, skirta saugoti modeliams

============================================================================
*/
type manufacturer struct {
	name     string
	quantity uint
	models   []model
}

/*
============================================================================

	struktura, skirta bendrai atminciai, su pagalbiniais metodais:
	Add, Take ir Print

============================================================================
*/
type buffer struct {
	buff counterSlice
}

func (b *buffer) Add(c *counter) {
	i := (b.buff).find(c.price)
	switch i {
	case -1:
		{
			j := (b.buff).getPos(c.price)
			if j == -1 {
				b.buff = append(b.buff, *c)
			} else {
				b.buff = append(b.buff, counter{})
				copy(b.buff[j+1:], b.buff[j:])
				b.buff[j] = *c
			}
		}
	default:
		{
			b.buff[i].count += c.count
		}
	}
}

func (b *buffer) Take(c *counter) uint {
	if len(b.buff) == 0 {
		return 0
	}
	i := (b.buff).find(c.price)
	switch i {
	case -1:
		{
			return 0
		}
	default:
		{
			var j uint
			if b.buff[i].count >= c.count {
				j = c.count
			} else {
				j = b.buff[i].count
			}
			b.buff[i].count -= j
			if b.buff[i].count <= 0 {
				copy(b.buff[i:], b.buff[i+1:])
				b.buff[len(b.buff)-1] = counter{}
				b.buff = b.buff[:len(b.buff)-1]
			}

			return j
		}
	}
}

func (b *buffer) Print() {
	fmt.Println("Nesuvartota liko: \n")
	for _, v := range b.buff {
		fmt.Printf("%v\n", v)
	}
}

/*
============================================================================

	Pagalbine funkcija, tikrini klaidoms

============================================================================
*/

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*
============================================================================
ReadFile

	Pradiniu duomenu nuskaitymo funkcija, grazina gamintoju ir vartotoju
	sarasa.

============================================================================
*/

func ReadFile() ([]manufacturer, []counterSlice) {
	var Allmodels []manufacturer
	var users []counterSlice
	var readUsers bool = false
	var temp counterSlice

	file, err := os.Open(DataFile)
	check(err)
	read := bufio.NewReader(file)

	for j := 0; ; j++ {
	decide:
		line, _, err := read.ReadLine()
		if err == io.EOF {
			break
		}
		if !readUsers {
			if "UsersWant" == string(line) {
				readUsers = true
				goto decide
			}
			vars := strings.Split(string(line), " ")
			title := vars[0]
			count, _ := strconv.Atoi(vars[1])
			models := make([]model, count)
			for i := 0; i < count; i++ {
				line, _, err := read.ReadLine()
				check(err)
				vars := strings.Split(string(line), "\t")
				name := vars[0]
				quantity, err := strconv.Atoi(vars[1])
				check(err)
				price, err := strconv.ParseFloat(vars[2], 32)
				check(err)
				models[i] = model{name, uint(quantity), float32(price)}
			}
			Allmodels = append(Allmodels, manufacturer{title, uint(count), models})
		} else {
			if " " == string(line) {
				users = append(users, temp)
				temp = make(counterSlice, 0)
			} else {
				vars := strings.Split(string(line), "\t")
				price, err := strconv.ParseFloat(vars[0], 32)
				check(err)
				count, err := strconv.Atoi(vars[1])
				check(err)
				temp = append(temp, counter{float32(price), uint(count)})
			}
		}
	}
	file.Close()
	return Allmodels, users
}

/*
============================================================================
PrintTable

	Atspausdina pradinius duomenis lentelemis
============================================================================
*/

func PrintTable(printOut []manufacturer, users []counterSlice) {
	fmt.Println("-----------------------------Pradiniai Duomenys------------------------------")
	for _, element := range printOut {
		fmt.Println(element.name)
		fmt.Println("-----------------------------------------------------------------------------")
		fmt.Printf("%20v %45v %10v\n", "Modelio Pavadinimas", "Kiekis", "Kaina")
		for i, modelis := range element.models {
			fmt.Printf("%2s %17s %45d %10.2f\n", strconv.Itoa(i+1)+")", modelis.name, modelis.quantity, modelis.price)
		}
		fmt.Println("-----------------------------------------------------------------------------")
	}
	fmt.Println("------------------------------------Users------------------------------------")
	for i, element := range users {
		fmt.Printf("%s %d\n", "User Nr.", i)
		fmt.Println("------------")
		fmt.Println("Price  Count")
		for _, user := range element {
			fmt.Printf("%5.2f  %5d\n", user.price, user.count)
		}
		fmt.Println("------------")
	}
}

/*
============================================================================
Use

	Vartotojo funkcija, elementai siunciami kanalu valdytojui,
	jei nepavyko isimti visu norimu elementu, tuomet gauname atgal
	counter struktura su tuo, kiek liko pasiimti dar, siunciame pakartotinai
============================================================================
*/
func Use(consumerMessages chan counter, user counterSlice, taken chan counter) {
	for _, el := range user {
		consumerMessages <- el
	consume:
		select {
		case c := <-taken:
			{
				consumerMessages <- c
				goto consume
			}
		default:
			break
		}
	}
}

/*
============================================================================
Make

	Gamintojo funkcija, elementai siunciami kanalu valdytojui
============================================================================
*/
func Make(producerMessages chan counter, models manufacturer) {
	for _, element := range models.models {
		producerMessages <- counter{element.price, element.quantity}
	}
}

/*
============================================================================
Valdytojas

	Priima counter strukturas is Make ir Use, per atitinkamus kanalus,
	tuomet arba ideda arba isema elementu is bendros atminties
============================================================================
*/
func Valdytojas(userChannels chan counter, producerChannels chan counter, done chan bool, buff *buffer, taken chan counter, howMuch int, removes int) {

	for i := 0; i < howMuch+removes; {
		select { //Alternatyvos pasirinkimas
		case counter := <-producerChannels: //Siuo atveju priimame duomenis is gamintoju ir dedame i bendraja atmint
			{
				buff.Add(&counter)
				i++
			}
		case counter := <-userChannels: //Siuo atveju priimame duomenis is vartotoju ir isimame imanoma kieki
			{ //elemento is bendros atminties, jei lieka neisimtu, siunciame counter objekta atgal (su likuciu)
				counter.count -= buff.Take(&counter)
				if counter.count == 0 {
					i++
				} else {
					taken <- counter
				}
			}

		default:
		}
	}

	done <- true //Siunciame signala, jog galime baigti main()

}

func main() {
	done := make(chan bool, 1)

	var howMuch int = 0
	var removes int = 0

	var buff buffer
	var AllModels []manufacturer
	var Users []counterSlice

	AllModels, Users = ReadFile()

	for _, el := range AllModels {
		howMuch += len(el.models)
	}

	for _, el := range Users {
		removes += len(el)
	}

	PrintTable(AllModels, Users)

	var producerMessages chan counter
	var consumerMessages chan counter
	var taken chan counter

	producerMessages = make(chan counter, len(AllModels))
	consumerMessages = make(chan counter, len(Users))
	taken = make(chan counter, len(Users))

	for _, element := range AllModels {
		go Make(producerMessages, element)
	}
	for _, element := range Users {
		go Use(consumerMessages, element, taken)
	}

	go Valdytojas(consumerMessages, producerMessages, done, &buff, taken, howMuch, removes)

	<-done //Valdytojas baige savo darba
	buff.Print()
}
