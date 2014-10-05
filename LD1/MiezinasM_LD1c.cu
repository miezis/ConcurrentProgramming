/**
* 1. Kiek iteracijø ið eilës padaro vienas procesas?
* Ats.: vienà pilnai
*
* 2. Kokia tvarka atspausdinami to paties masyvo duomenys?
* Ats.: tokia, kokia suraðyti
*
* 3. Kokia tvarka vykdomi procesai?
* Ats.: tokia, kokia startuoja(tai yra 1 procesas atlieka viena iteracija, tuomet 2 procesas ir t.t.,
* kol atliktos visos iteracijos
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

using namespace std;

//Konstantos
const char DataFile[] = "MiezinasM.txt";
const int MAX_THREADS = 5;
const int MAX_ARRAY_SIZE = 10;

//Struktura gamintojo modeliams saugoti
struct model {
	char name[15];
	int quantity;
	double price;
};

//Gamintojo klase, tarnauja kaip modeliu konteineris
class Manufacturer {
public:
	string name;
	int quantity;
	model models[MAX_ARRAY_SIZE];
	//Konstruktorius
	Manufacturer(string name, int quantity, model models[]) {
		this->name = name;
		this->quantity = quantity;
		for (int i = 0; i < quantity; i++){
			this->models[i] = models[i];
		}
	}
	Manufacturer() {
		this->name = "";
		this->quantity = 0;
	}
};

//Funkciju prototipai
void ReadFile(string filename, Manufacturer(&AllModels)[MAX_THREADS]);
void PrintTable(Manufacturer(&printOut)[MAX_THREADS]);
__device__ void PrintManufacturerModels(Manufacturer printOut, int procNum);
__global__ void RunOnGPU(Manufacturer *printOut);

int main() {
	Manufacturer AllModels[MAX_THREADS];
	ReadFile(DataFile, AllModels);
	PrintTable(AllModels);
	//Paruosiame GPU, t.y. perkeliame duomenis is RAM i VRAM
	Manufacturer *manu;
	int i = MAX_THREADS * sizeof(Manufacturer);
	cudaMalloc((void**)&manu, i);
	cudaMemcpy(manu, AllModels, i, cudaMemcpyHostToDevice);
	//Iskvieciame GPU
	RunOnGPU << <1, MAX_THREADS >> >(manu);
	cudaDeviceSynchronize();
	//Atlaisviname VRAM atminti
	cudaFree(manu);

	system("Pause");
	return 0;
}

void ReadFile(string filename, Manufacturer(&AllModels)[MAX_THREADS]) {
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
			AllModels[j].name = title;
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
void PrintTable(Manufacturer(&printOut)[MAX_THREADS]){
	cout << "-----------------------------------------------------------------------------\n";
	for (Manufacturer & manu : printOut){
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
__device__ void PrintManufacturerModels(Manufacturer printOut, int procNum){
	for (int i = 0; i < printOut.quantity; i++) {
		model mod = printOut.models[i];
		printf("procesas_%d %d    %15s %2d %4.2f\n", procNum + 1, i + 1, mod.name, mod.quantity, mod.price);
	}
}
__global__ void RunOnGPU(Manufacturer *printOut) {
	PrintManufacturerModels(printOut[threadIdx.x], threadIdx.x);
}