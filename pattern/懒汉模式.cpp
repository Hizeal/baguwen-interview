//方法调用前，创建好实例

#include<string>
#include<iostream>


class Singleton
{
private:
    static Singleton* instance;
    
    Singleton(std::string values):value_(values){}
    ~Singleton(){}
    std::string value_;

    

public:
    Singleton(Singleton& other) = delete;
    void operator=(Singleton& other) = delete;

    Singleton* Getinstance(std::string& value){}
};

Singleton* Singleton::instance(nullptr); 

Singleton* Singleton::Getinstance(std::string& value){
    instance = new Singleton(value);
    return instance;
}






