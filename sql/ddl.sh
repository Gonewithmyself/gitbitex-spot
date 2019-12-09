#!/bin/bash
host=""
user="john"
password="since1999"
create_db_sql="create database IF NOT EXISTS ${db}"

conn_str="-u $user -p$password"

dbs=("eth" "eos" "btc" "ggc")

for db in ${dbs[@]}
do
    # mysql $conn_str -e "create database IF NOT EXISTS ${db}"
    # mysql $conn_str $db -e "source trade.sql"
    mysql $conn_str -e "drop database IF EXISTS ${db}"
done