#include <mutex>
#include <condition_variable>

using namespace std;

class semaphore {
private:
	int count;
	mutex mtx;
	condition_variable cv;
public:
	semaphore(int count_):count(count_){}
	void acquire(){
		unique_lock<mutex> lck(mtx);
		while(cont == 0){
			cv.wait(lck);
		}
		count--;
	}
	void release(){
		unique_lock<mutex> lck(mtx);
        ++count;
        cv.notify_one();
	}
}