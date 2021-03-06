/**
 * 1. Kiek iteracijų iš eilės padaro vienas procesas?
 * Ats.: atsitiktinį skaičių
 *
 * 2. Kokia tvarka atspausdinami to paties masyvo duomenys? 
 * Ats.: tokia, kokia surašyti
 * 
 * 3. Kokia tvarka vykdomi procesai?
 * Ats.: atsitiktine
*/

#include <iostream>
#include <string>
#include <fstream>
#include <vector>
#include <iomanip>
#include <omp.h>

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
	Manufacturer(string name, int quantity, vector<model> models) {
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

int main() {
	vector<Manufacturer> AllModels = ReadFile(DataFile);
	int i;	//proceso numeriui saugot


	PrintTable(AllModels);

	//nustatom kiek giju kursime
	omp_set_num_threads(AllModels.size());
	//lygiagrecioji dalis
	#pragma omp parallel private(i)
	{
		i = omp_get_thread_num();
		PrintManufacturerModels(AllModels.at(i), "procesas_" + to_string(i + 1));
	}

	return 0;
}

vector<Manufacturer> ReadFile(string filename){
	vector<Manufacturer> AllModels;
	string title;
	int count;
	ifstream fin(filename);
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
	cout << "-----------------------------------------------------------------------\n";
	for(Manufacturer & manu : printOut){
		cout << right << setw(35) << manu.getName() << "\n";
		cout << "-----------------------------------------------------------------------\n";
		cout 	<< left << setw(63) << "Modelio Pavadinimas"
				<< setw(8) << "Kiekis"
				<< setw(5) << "Kaina" << "\n";
		for(int i = 0; i < manu.getQuantity(); i++){
			model forPrinting = manu.getModels().at(i);
			cout 	<< left << setw(3) << to_string(i+1) + ")" 
					<< setw(60) << forPrinting.name 
					<< setw(8) << forPrinting.quantity 
					<< setw(5) << setprecision(4) << forPrinting.price 
					<< "\n";
		}
		cout << "-----------------------------------------------------------------------\n";
	}
}

void PrintManufacturerModels(Manufacturer printOut, string procNum){
	int i = 0;
	for(model &mod : printOut.getModels()){
		cout 	<< left << setw(20) << procNum 
				<< setw(5) << (i + 1) 
				<< setw(30) << mod.name 
				<< setw(5) << mod.quantity 
				<< setw(5) << setprecision(4) << mod.price << "\n";
		i++;
	}
}