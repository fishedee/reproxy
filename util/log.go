package util

import (
	"log"
	"os"
	"time"
)

type LoggerType struct{
	log * log.Logger
}

var Logger *LoggerType

func InitLogger(fileAddress string)(error){
	logfile,err := os.OpenFile(fileAddress,os.O_RDWR | os.O_APPEND|os.O_CREATE,0660);
	if err != nil{
		return err
	}

	Logger = &LoggerType{}
	Logger.log = log.New(logfile,"",0)
	return nil
}

func (this *LoggerType)Info(args ...interface{}){
	this.Print("INFO",args...)
}

func (this *LoggerType)Warn(args ...interface{}){
	this.Print("WARN",args...)
}

func (this *LoggerType)Err(args ...interface{}){
	this.Print("ERROR",args...)
}

func (this *LoggerType)Print(level string,args ...interface{}){
	now := time.Now()
	logSomething := []interface{}{}
	logSomething = append(logSomething,level)
	logSomething = append(logSomething," ")
	logSomething = append(logSomething,now.Year())
	logSomething = append(logSomething,"/")
	logSomething = append(logSomething,int(now.Month()))
	logSomething = append(logSomething,"/")
	logSomething = append(logSomething,now.Day())
	logSomething = append(logSomething," ")
	logSomething = append(logSomething,now.Hour())
	logSomething = append(logSomething,":")
	logSomething = append(logSomething,now.Minute())
	logSomething = append(logSomething,":")
	logSomething = append(logSomething,now.Second())
	logSomething = append(logSomething," ")
	for _,value := range args{
		logSomething = append(logSomething,value)
	}

	this.log.Print(
		logSomething...
	)
}