#! /bin/bash
for i in `find ios`; do
	xxd -c 25600 -p $i | sed -e '1 s/320000001800000001/320000001800000002/; t' -e '1,// s//320000001800000002/' | xxd -c 256 -p -r > out_$i
done