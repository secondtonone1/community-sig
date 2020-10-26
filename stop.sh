sudo kill -9 $(netstat -nlp | grep :9699 | awk '{print $7}' | awk -F"/" '{ print $1 }')
