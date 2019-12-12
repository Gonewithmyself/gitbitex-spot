#!/bin/bash
host=""
user="john"
password="since1999"
conn_str="-u ${user} -p${password}"

pairs=("BTC-USDT" "HT-USDT" "ETH-USDT" "EOS-USDT")

create_db(){
    for pair in ${pairs[@]}
    do
        mysql ${conn_str} -e "create database IF NOT EXISTS \`${pair}\`"
        mysql $conn_str $pair -e "source market.sql"
    done
    mysql ${conn_str} -e "create database IF NOT EXISTS db_account"
    mysql $conn_str db_account -e "source account.sql"
}

clear_db(){
    for pair in ${pairs[@]}
    do
        mysql $conn_str -e "drop database IF EXISTS \`${pair}\`"
    done
    mysql $conn_str -e "drop database IF EXISTS db_account"
}

if [ $1 == "clean" ]; then
    echo "clean..."
    clear_db
else
    echo "create..."
    create_db
fi