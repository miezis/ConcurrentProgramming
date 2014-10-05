/**
 * 1. Kiek iteracijų iš eilės padaro vienas procesas?
 * Ats.: atsitiktinį skaičių
 *
 * 2. Kokia tvarka atspausdinami to paties masyvo duomenys? 
 * Ats.: tokia, kokia surašyti
 * 
 * 3. Kokia tvarka vykdomi procesai?
 * Ats.: atsitiktine
 * 4. Kurioje programoje trumpiausias vieno proceso kodas? 
 * Ats.: Sukurti procesui trumpiausias kodas A (C++ 11) dalyje, taciau dar 
 *       reikia juos apjungti i pagrindini procesa, tad galima svarstyti, jog B (openMP) dalyje kodas trumpesnis.
 * 5. Kokiu kompiuteriu vykdėte savo programas? Nurodykite branduolių skaičių ir dažnius, OA apimtį, OS, 
 *    NVIDIA plokštės tipą.
 * Ats.: Asus K53SM,
 *       Procesorius:
 *       Intel® Core™ i3 2350M
 *		 2 fiziniai branduoliai (cores), 4 virtualios gijos (threads)
 *       2.3Ghz taktinis daznis
 *       http://ark.intel.com/products/53438/Intel-Core-i3-2350M-Processor-3M-Cache-2_30-GHz
 *       Operatyvioji atmintis:
 *		 4GB DDR3 1330Mhz (kylancio ir krintancio frontu kartu), 665Mhz (realus)
 *       Operacine sistema:
 *		 A ir B dalys kompiliuotos Ubuntu Linux 14.04 64bit OS aplinkoje su gcc 4.8.2
 *       C dalis kompiliuota Windows 8.1 64bit aplinkoje su Visual Studio 2013 Ultimate
 *       (nepavyko CUDA pasileisti Linux aplinkoje, nes CUDA toolkit iraso driver'ius, kurie matyt su mano GPU nedraugauja)
 *       D dalis parasyta Linux aplinkoje, kompiliuota mpilab.elen.ktu.lt aplinkoje
 *		 NVIDIA GPU:
 *       GeForce GT630M su 2GB GDDR5 VRAM
 *		 http://www.geforce.co.uk/hardware/notebook-gpus/geforce-gt-630M/specifications
 */

#include <stdio.h>
#include <stdlib.h>
#include <iostream>
#include <string>
#include <fstream>
#include <vector>
#include <iomanip>
#include <mpi.h>
#include <sstream>

using namespace std;

//konstanta duomenu failo pavadinimui
const char DataFile[] = "MiezinasM.txt";

//bazine struktura saugoti vienam modelio irasui
struct model {
	string name;
	int quantity;
	double price;
};

//konteinerine klase, skirta saugoti modeliams (vektoriu)
class Manufacturer {
private:
	string name;
	int quantity;
	vector<model> models;
public:
	Manufacturer(string name, int quantity, vector<model> &models) {
		this->name = name;
		this->quantity = quantity;
		this->models = models;
	}

	string getName() {
		return this->name;
	}

	int getQuantity() {
		return this->quantity;
	}

	vector<model> getModels(){
		return this->models;
	}
};

//Funkciju prototipai
vector<Manufacturer> ReadFile(string filename);
void PrintTable(vector<Manufacturer> printOut);
void PrintManufacturerModels(Manufacturer printOut, string procNum);

int main(int argc, char **argv) {
	vector<Manufacturer> AllModels = ReadFile(DataFile);
	Manufacturer T1 = AllModels.at(0);
	Manufacturer T2 = AllModels.at(1);
	Manufacturer T3 = AllModels.at(2);
	Manufacturer T4 = AllModels.at(3);
	Manufacturer T5 = AllModels.at(4);

	int procNr;
	//Lygiagrecioji dalis
	MPI_Init(&argc, &argv); // argumentu truksta
	MPI_Comm_rank(MPI_COMM_WORLD, &procNr);

	stringstream ss;
	ss << "procesas_" << procNr;

	switch (procNr) {
		case 0:
			PrintTable(AllModels);
			break;
		case 1:
			PrintManufacturerModels(T1, ss.str());
			break;
		case 2:
			PrintManufacturerModels(T2, ss.str());
			break;
		case 3:
			PrintManufacturerModels(T3, ss.str());
			break;
		case 4:
			PrintManufacturerModels(T4, ss.str());
			break;
		case 5:
			PrintManufacturerModels(T5, ss.str());
			break;
		default:
			cout << "Nera tokio proceso.";
			break;
	}
	/*
	if (procNr == 0)
	{
		PrintTable(AllModels);
	}

	stringstream ss;
	ss << "procesas_" << procNr;
	PrintManufacturerModels(AllModels.at(procNr), ss.str());
	*/
	//Uzbaigiame lygiagreciaja dali su finalize
	MPI_Finalize();

	return 0;
}

vector<Manufacturer> ReadFile(string filename){
	vector<Manufacturer> AllModels;
	string title;
	int count;
	ifstream fin(filename.c_str());
	if(!fin) {
		cerr << "Couldn't open file!\n";
	} else {
		while(!fin.eof()){
			fin >> title >> count;
			vector<model> models;
			for(int i = 0; i < count; i++){
				model modelis;
				fin >> modelis.name >> modelis.quantity >> modelis.price;
				models.push_back(modelis);
			}
			AllModels.push_back(Manufacturer(title, count, models));
		}
		fin.close();
	}	
	return AllModels;
}

void PrintTable(vector<Manufacturer> printOut){
	cout << "-----------------------------------------------------------------------------\n";
	for(int i = 0; i < printOut.size(); i++){
		cout << right << setw(35) << printOut.at(i).getName() << "\n";
		cout << "-----------------------------------------------------------------------------\n";
		cout 	<< left << setw(63) << "Modelio Pavadinimas"
				<< setw(8) << "Kiekis"
				<< setw(5) << "Kaina" << "\n";
		for(int j = 0; j < printOut.at(i).getQuantity(); j++){
			model forPrinting = printOut.at(i).getModels().at(j);
			stringstream ss;
			ss << j + 1 << ")";
			cout 	<< left << setw(3) << ss.str()
					<< setw(60) << forPrinting.name 
					<< setw(8) << forPrinting.quantity 
					<< setw(5) << setprecision(4) << forPrinting.price 
					<< "\n";
		}
		cout << "-----------------------------------------------------------------------------\n";
	}
}

void PrintManufacturerModels(Manufacturer printOut, string procNum){
	vector<model> models = printOut.getModels();
	for(int i = 0; i < models.size(); i++) {
		cout 	<< left << setw(20) << procNum 
				<< setw(5) << (i + 1) 
				<< setw(30) << models.at(i).name 
				<< setw(5) << models.at(i).quantity 
				<< setw(5) << setprecision(4) << models.at(i).price << "\n";
	}
}