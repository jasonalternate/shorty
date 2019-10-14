#!/bin/bash

<<COMMENT1
    Note: this is the result of google search, not as much my own bash scripting knowledge.

     O. Tange (2018): GNU Parallel 2018, Mar 2018, ISBN 9781387509881,
     DOI https://doi.org/10.5281/zenodo.1146014
COMMENT1



mycurl() {
    START=$(date +%s)
    curl -d '{"destination":"http://www.hopefullynotawebsite.com"}' -H "Content-Type: application/json" -X POST http://localhost:8080/links
    END=$(date +%s)
    DIFF=$(( $END - $START ))
    echo "It took $DIFF seconds"
}
export -f mycurl

seq 100000 | parallel -j0 mycurl
