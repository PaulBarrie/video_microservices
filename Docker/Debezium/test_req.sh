#!/bin/bash

curl 'http://localhost:9200/video/_search?pretty'
{
  "took" : 42,
  "timed_out" : false,
  "_shards" : {
    "total" : 5,
    "successful" : 5,
    "failed" : 0
  },
  "hits" : {
    "total" : 4,
    "max_score" : 1.0,
    "hits" : [
      {
        "_index" : "video",
        "_type" : "video",
        "_id" : "1",
        "_score" : 1.0,
        "_source" : {
          "first_name" : "Sally",
          "last_name" : "Thomas",
          "email" : "sally.thomas@acme.com"
        }
      }]
  }
}