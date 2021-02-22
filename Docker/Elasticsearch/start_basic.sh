#!/bin/bash   

spawn curl -uelastic -XPOST 'http://127.0.0.1:9200/_xpack/license/start_basic'
expect "Enter host password for user 'elastic':"
send \n;
interact