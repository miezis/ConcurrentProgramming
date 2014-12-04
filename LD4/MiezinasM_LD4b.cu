/*
============================================================================
Mantas Miežinas, IFF-2
LD4b
============================================================================
*/

#include "cuda_runtime.h"
#include "device_launch_parameters.h"
#include <cuda.h>

#include <stdio.h>
#include <stdlib.h>
#include <iostream>
#include <string>
#include <fstream>
#include <iomanip>
#include <thrust/host_vector.h>
#include <thrust/device_vector.h>

using namespace std;

//Konstantos
const char DataFile[] = "MiezinasM_L4.txt";
const int MAX_THREADS = 5;
const int MAX_ARRAY_SIZE = 5;

//Struktura gamintojo modeliams saugoti
struct model {
	char name[75];
	int quantity;
	double price;
};

//Konteinerine struktura saugoti modeliams
struct manufacturer {
	char name[15];
	int quantity;
	model models[MAX_ARRAY_SIZE];
};

//Funkciju prototipai
void ReadFile(string filename, thrust::host_vector<manufacturer> &AllModels);
void PrintTable(thrust::host_vector<manufacturer> printOut);
void PrintResults(thrust::host_vector<model> printOut);
void Plus(thrust::device_vector<manufacturer> &manu, thrust::device_vector<model> &resultsArray);


int main() {
	thrust::host_vector<manufacturer> AllModels(MAX_THREADS);
	thrust::host_vector<model> results(MAX_ARRAY_SIZE);

	ReadFile(DataFile, AllModels);
	PrintTable(AllModels);
	//Nusinuliname rezultatu masyvo elementus
	for (int i = 0; i < MAX_ARRAY_SIZE; i++){
		model data;
		strcpy(data.name, "");
		data.price = 0.0;
		data.quantity = 0;
		results[i] = data;
	}

	//Perkeliame duomenis is RAM i VRAM
	thrust::device_vector<manufacturer> manu = AllModels;
	thrust::device_vector<model> resultsArray = results;

	//Iskvieciame sumuojancia funkcija
	Plus(manu, resultsArray);
	cudaDeviceSynchronize();

	//Persikeliame rezultatus is VRAM(GPU) i RAM(CPU)
	results = resultsArray;

	//Atspausdiname rezultatus
	PrintResults(results);
	system("Pause");
	return 0;
}

/*
============================================================================
ReadFile

Pradiniu duomenu nuskaitymo funkcija, per nuoroda grazina manufacturer
strukturos host vektoriu (CPU atmintyje)
============================================================================
*/
void ReadFile(string filename, thrust::host_vector<manufacturer> &AllModels) {
	string title;
	int count, j;
	j = 0;
	ifstream fin(filename);
	if (!fin) {
		cerr << "Couldn't open file!\n";
	}
	else {
		while (!fin.eof()){
			fin >> title >> count;
			strcpy(AllModels[j].name, title.c_str());
			//AllModels[j].name = title;
			AllModels[j].quantity = count;
			model models[MAX_ARRAY_SIZE];
			for (int i = 0; i < count; i++){
				model modelis;
				fin >> modelis.name >> modelis.quantity >> modelis.price;
				AllModels[j].models[i] = modelis;
			}
			j++;
		}
		fin.close();
	}
}

/*
============================================================================
PrintTable

Atspausdina pradinius duomenis lentelemis
============================================================================
*/
void PrintTable(thrust::host_vector<manufacturer> printOut){
	cout << "-----------------------------------------------------------------------------\n";
	for (manufacturer & manu : printOut){
		cout << right << setw(35) << manu.name << "\n";
		cout << "-----------------------------------------------------------------------------\n";
		cout << left << setw(63) << "Modelio Pavadinimas"
			<< setw(8) << "Kiekis"
			<< setw(5) << "Kaina" << "\n";
		cout << "-----------------------------------------------------------------------------\n";
		for (int i = 0; i < manu.quantity; i++){
			model forPrinting = manu.models[i];
			cout << left << setw(3) << to_string(i + 1) + ")"
				<< setw(60) << forPrinting.name
				<< setw(8) << forPrinting.quantity
				<< setw(5) << setprecision(4) << forPrinting.price
				<< "\n";
		}
		cout << "-----------------------------------------------------------------------------\n";
	}
}

/*
============================================================================
PrintResults

Atspausdina rezultatu vektoriu lenteleje
============================================================================
*/
void PrintResults(thrust::host_vector<model> printOut) {
	cout << "-----------------------------------------------------------------------------\n";
	cout << left << setw(63) << "Modelio Pavadinimas"
		<< setw(8) << "Kiekis"
		<< setw(5) << "Kaina" << "\n";
	cout << "-----------------------------------------------------------------------------\n";
	for (int i = 0; i < MAX_ARRAY_SIZE; i++) {
		model forPrinting = printOut[i];
		cout << left << setw(3) << to_string(i + 1) + ")"
			<< setw(60) << forPrinting.name
			<< setw(8) << forPrinting.quantity
			<< setw(5) << setprecision(4) << forPrinting.price
			<< "\n";
	}
}

/*
============================================================================
Plus

Susumuoja pradiniu duomenu masyvus i viena masyva ir grazina rezultata per
nuoroda. Naudojami device vektoriai, reiskias naudojama atmintis priklauso
GPU.
============================================================================
*/
void Plus(thrust::device_vector<manufacturer> &manu, thrust::device_vector<model> &resultsArray) {
	for (int i = 0; i < manu.size(); i++) {
		for (int j = 0; j < resultsArray.size(); j++) {
			model data = resultsArray[i];
			manufacturer addData = manu[j];
			strcat(data.name, addData.models[i].name);
			data.price += addData.models[i].price;
			data.quantity += addData.models[i].quantity;
			resultsArray[i] = data;
		}
	}
}