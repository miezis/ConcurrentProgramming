# coding=utf-8
"""
============================================================================
Mantas Miežinas, IFF-2
CHANa
============================================================================
"""
from pycsp.parallel import *
from os import fstat
import sys

# Skaitliuko struktūra (klasė)
class MyCounter:
    def __init__(self, price, count):
        self.price = price
        self.count = count

# Bazinė struktūra (klasė), saugoti vienam modelio įrašui
class Model:
    def __init__(self, name, quantity, price):
        self.name = name
        self.quantity = quantity
        self.price = price

# Konteinerinė struktūra (klasė), skirta saugoti modeliams
class Manufacturer:
    def __init__(self, name, quantity, models):
        self.name = name
        self.quantity = quantity
        self.models = models

# Klasė, skirta bendrai atminčiai, su pagalbiniais metodais
class Buffer:
    def __init__(self):
        self.buffer = []

    # Metodas, randantis esamo elemento indeksą,
    # jei tokio nėra - gražina -1
    def find(self, x):
        for index in xrange(len(self.buffer)):
            if self.buffer[index].price == x:
                return index
        return -1

    # Metodas, randantis indeksą, kuriuo reik įterpti naują elementą,
    # gražina -1, jei reikia dėti į galą
    def get_pos(self, x):
        for index in xrange(len(self.buffer)):
            if self.buffer[index].price > x:
                return index
        return -1

    # Metodas, pridedantis elementą į bendrą atmintį,
    # gražina indeksą, į kur idėjo elementą
    def add(self, price, count):
        index = self.find(price)
        if index != -1: # jei jau yra elementas, su tokia kaina,
            self.buffer[index].count += count # tai pridedame prie esamo
        else:
            index = self.get_pos(price) # randame, kur įterpti
            if index == -1: # jei -1, tai pridėti į galą
                self.buffer.append(MyCounter(price, count))
            else:
                self.buffer.insert(index, MyCounter(price, count))

        # gražiname indeksą, kur idėjome
        if index != -1:
            return index
        else:
            return len(self.buffer) - 1

    # Metodas, atimantis elementus iš bendros atminties,
    # gražina kiekį, kiek pavyko išimti bei indeksą iš
    # kurio elemento paimė
    def take(self, price, count):
        if len(self.buffer) == 0: # jei tuščia
            return 0, 0

        index = self.find(price)

        if index == -1: # jei nėra tokio elemento
            return 0, 0
        else:
            if self.buffer[index].count >= count: # jei esamo elemento kiekis didesnis ar lygus tam, kiek išimti
                taken = count
            else: # jei mažesnis
                taken = self.buffer[index].count

            # atminusuojame, kiek išimame
            self.buffer[index].count -= taken

            # jei po išimimo 0, tai šaliname elementą
            if self.buffer[index].count <= 0:
                del self.buffer[index]

            return taken, index

    # Metodas, kuris atspausdina bendros atminties turinį
    def print_list(self):
        print "Nesuvartota liko:\n"
        for index in xrange(len(self.buffer)):
            print str(self.buffer[index].price) + "\t" + str(self.buffer[index].count)

# Pradinių duomenų skaitymo funkcija
def read_file(manufacturers, counters):
    fin = open("MiezinasM_CHAN.txt", "r")
    read_users = False
    temp = []

    while fin.tell() != fstat(fin.fileno()).st_size:
        line = fin.readline()

        if line == "UsersWant\n":
            read_users = True
            line = fin.readline()

        if not read_users:
            vars = line.split()
            title = vars[0]
            count = int(vars[1])
            models = []
            for i in range(0, count):
                line = fin.readline()
                vars = line.split()
                name = vars[0]
                quantity = int(vars[1])
                price = float(vars[2])
                models.append(Model(name, quantity, price))
            manufacturers.append(Manufacturer(title, count, models))
        else:
            if " \n" == line:
                counters.append(temp)
                temp = []
            else:
                vars = line.split()
                price = float(vars[0])
                count = int(vars[1])
                temp.append(MyCounter(price, count))

    fin.close()

# Atspausdina pradinius duomenis lentelėmis
def print_table(manufacturers, counters):
    print "-----------------------------Pradiniai Duomenys------------------------------"
    for element in manufacturers:
        print element.name + "\n-----------------------------------------------------------------------------"
        print "%19s %46s %10s" % ("Modelio Pavadinimas", "Kiekis", "Kaina")
        print "-----------------------------------------------------------------------------"
        for model in element.models:
            print "%19s %46d %10.2f" % (model.name, model.quantity, model.price)
        print "-----------------------------------------------------------------------------"
    print "---------------------------------Vartotojai----------------------------------"
    for i in xrange(len(counters)):
        print "%s %d" % ("Vartotojai Nr.", i + 1)
        print "-----------------------------------------------------------------------------"
        print "%5s %70s" % ("Kaina", "Kiekis")
        print "-----------------------------------------------------------------------------"
        for j in xrange(len(counters[i])):
            print "%5.2f %70d" % (counters[i][j].price, counters[i][j].count)
        print "-----------------------------------------------------------------------------"

# Gamintojo procesas, elementai siunčiami kanalu manager procesui, baigus darbą
# pranešame išsiųsdami -1
@process
def producer(chan_out, models):
    for element in models.models:
        output = MyCounter(element.price, element.quantity)
        chan_out(output)
    chan_out(-1)

# Vartotojo procesas, elementai siunčiami kanalu manager procesui, taken kanalu
# grįžta paimtų elementų kiekis, jei jis mažesnis, nei skaitliuko elementų kiekis,
# tai siunčiame pakartotinai, tik atminusuodami count (t.y. kiek jau pavyko išimti),
# baigus darbą pranešame išsiųsdami -1
@process
def consumer(chan_out, taken, counters):
    for counter in counters:
        count = 0
        while count < counter.count:
            chan_out(MyCounter(counter.price, counter.count - count))
            add = taken()
            if add != -1:
                count += add
            else:
                break
    chan_out(-1)

# Valdytojo procesas, priima MyCounter klasės objektus iš producer ir consumer procesų,
# FairSelect išrenką kanalą iš kurio imsime duomenis ir alternatyvų veiksmą (add arba take),
# t.y. prideda arba atima duomenis iš bendros atminties
@process
def manager(logger, prodList, consList, takeList, buffer):
    consumersCount = [len(consList)]
    producersCount = [len(prodList)]

    # alternatyva elemento pridėjimui
    @choice
    def add(prod_cnt, channel_input=None):
        counter = channel_input
        if not isinstance(counter, MyCounter):
            prod_cnt[0] -= 1
        else:
            index = buffer.add(counter.price, counter.count)
            logger("%10s %10.2f %10d %10d\n" % ("Pridejo", counter.price, counter.count, index + 1))

    # alternatyva elemento išimimui
    @choice
    def take(con_cnt, prod_cnt, taken_chan, channel_input=None):
        counter = channel_input
        if not isinstance(counter, MyCounter):
            con_cnt[0] -= 1
        else:
            taken, index = buffer.take(counter.price, counter.count)
            if taken >= 0 and prod_cnt[0] > 0:
                taken_chan(taken)
                if taken > 0:
                    logger("%10s %10.2f %10d %10d\n" % ("Iseme", counter.price, taken, index + 1))
            else:
                taken_chan(-1)


    # kol yra vartotojų arba gamintojų
    while consumersCount[0] > 0 or producersCount[0] > 0:
        FairSelect( # alternatyvos pasirinkimas
            InputGuard(prodList[0], action=add(producersCount)),
            InputGuard(prodList[1], action=add(producersCount)),
            InputGuard(prodList[2], action=add(producersCount)),
            InputGuard(prodList[3], action=add(producersCount)),
            InputGuard(prodList[4], action=add(producersCount)),
            InputGuard(consList[0], action=take(consumersCount, producersCount, takeList[0])),
            InputGuard(consList[1], action=take(consumersCount, producersCount, takeList[1])),
            InputGuard(consList[2], action=take(consumersCount, producersCount, takeList[2])),
            InputGuard(consList[3], action=take(consumersCount, producersCount, takeList[3]))
        )

    # išsiunčiame, jog baigiame darbą (stebėtojui)
    logger(-1)
    # uždarome kanalus
    poison(*prodList)
    poison(*consList)
    poison(*takeList)

# Stebėtojo procesas, gauna duomenis logger kanalu iš manager proceso ir prideda informaciją į
# bendrą log'ą, kai gauname -1 - tai reiškia, kad manager baigia savo darbą ir galime
# spausdinti informaciją į ekraną, galiausiai uždarome logger kanalą
@process
def observer(logger):
    income = ""
    info = ""
    while income != -1:
        info += income
        income = logger()
    print "------------------------------Atlikti Veiksmai-------------------------------"
    print "%10s %10s %10s %10s" % ("Veiksmas", "Reiksme", "Kiekis", "Pozicija")
    print "-----------------------------------------------------------------------------"
    print info
    print "-----------------------------------------------------------------------------"
    poison(logger)


def main():
    manufacturers = []
    counters = []

    read_file(manufacturers, counters)
    print_table(manufacturers, counters)

    producerMessages = Channel('Producers') * len(manufacturers)
    consumerMessages = Channel('Consumers') * len(counters)
    taken = Channel('Taken') * len(counters)
    logger = Channel('log')

    buffer = Buffer()

    Spawn(
        observer(logger.reader()),
        manager(
            logger.writer(),
            [produce.reader() for produce in producerMessages],
            [consume.reader() for consume in consumerMessages],
            [take.writer() for take in taken],
            buffer
        ),
        [producer(producerMessages[id].writer(), manufacturers[id]) for id in xrange(len(manufacturers))],
        [consumer(consumerMessages[id].writer(), taken[id].reader(), counters[id]) for id in xrange(len(counters))]
    )

    shutdown()

    buffer.print_list()


if __name__ == "__main__":
    sys.exit(main())
