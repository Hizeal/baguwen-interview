//懒汉模式
//单例的构造函数/析构函数应该总是私有，避免使用' new ' / ' delete '直接构造/销毁



#include<iostream>
#include<algorithm>
#include<string>
#include<vector>
#include<cstring>
#include<thread>
#include<mutex>


class Singleleton
{
private:
    static Singleleton* pininstance;
    static std::mutex mutex_;

protected:
    Singleleton(const std::string value):value_(value){}
    ~Singleleton(){}
    std::string value_;
public:
    Singleleton(Singleleton& other) = delete;
    void operator=(const Singleleton&) = delete;

    static Singleleton* Getinstance(const std::string& value);

    void SomeBusinessLogic(){

    }

    std::string value() const{
        return value_;
    }
};

//静态方法/成员在类外定义或初始化
Singleleton* Singleleton::pininstance{nullptr};
std::mutex Singleleton::mutex_;

Singleleton* Singleleton::Getinstance(const std::string& value){
    std::lock_guard<std::mutex> lock(mutex_); //构造时调用mutex_lock()，析构时调用mutex_unlock()
    if(pininstance==nullptr){
        pininstance = new Singleleton(value);
    }
    
    return pininstance;
}


void ThreadFoo(){
    // Following code emulates slow initialization.
    std::this_thread::sleep_for(std::chrono::milliseconds(1000));
    Singleleton* singleton = Singleleton::Getinstance("FOO");
    std::cout << singleton->value() << "\n";
}

void ThreadBar(){
    // Following code emulates slow initialization.
    std::this_thread::sleep_for(std::chrono::milliseconds(1000));
    Singleleton* singleton = Singleleton::Getinstance("BAR");
    std::cout << singleton->value() << "\n";
}

int main()
{   
    std::cout <<"If you see the same value, then singleton was reused (yay!\n" <<
                "If you see different values, then 2 singletons were created (booo!!)\n\n" <<
                "RESULT:\n";   
    std::thread t1(ThreadFoo);
    std::thread t2(ThreadBar);
    t1.join();
    t2.join();
    
    return 0;
}