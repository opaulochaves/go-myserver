#!/bin/bash
rungo () {
        if [ $# -eq 0 ]
                then nodemon --exec go run main.go --signal SIGTERM
        elif [ $# -eq 1 ]
                then nodemon -e go --exec go run $1 --signal SIGTERM
        fi
}

rungo $1

# https://medium.com/@thomaswsmith_63994/how-i-got-hot-reloading-to-work-with-golang-and-nodemon-unix-b6f78e5902fe
