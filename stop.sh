#!/bin/bash
app=main.go
kill -9 `ps -ef|grep $app|grep -v grep|awk '{print $2}'`
ps aux|grep $app
