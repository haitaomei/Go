

export MongoDB_HOST_LIST="mongodb1=192.168.9.1; mongodb2=10.10.2.3"
mongoDBs=(${MongoDB_HOST_LIST//;/ })


ARG=""
for (( i=0; i<${#mongoDBs[@]}; i++ ));
do
    echo ${mongoDBs[i]}
    ARG="${ARG} --set mongodbhost.${mongoDBs[i]}"
done

echo ${ARG}