#!/bin/bash
log="0b256a16111af3656cf6293498a46eb03e6b5f66 testing message. #patch ea9ff1c820a462e2ca081f5246277ca1543d7d49 Merge pull request #1 from anothrNick/anothrNick-patch-1 72d753711d5d2a1fd04c0a68a0ec0201550f894e Test commit dc8337c5f4990a6fa6de63e494dbee13db7b70c6 remove unused sql file. #patch 91071313f65ce477dd9f98588dfb9a804d257116 Initial version. #major"
# supports #major, #minor, #patch (anything else will be 'minor')
case "$log" in
    *#major* ) echo major in $log;;
    *#minor* ) echo minor;;
    *#patch* ) echo patch;;
    * ) echo other in $log;;
esac