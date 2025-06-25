export user="postgres"
export defaultDB='postgres'

export timestamp=$(date '+%Y_%m_%d_%H_%M_%S')

if [ "$host" == "" ] ;then
  echo "host info is required"
  return 1
fi

if [ "$password" == "" ] ;then
  echo "password is required"
  return 1
fi

echo "Host: $host"

if [ "$dbName" != "" ] ; then
  export dbName=$(echo $(PGPASSWORD=$password psql -h $host -U $user -d postgres -c "select array (select datname from pg_database where datname='$dbName')") | grep -o "{.*}" | grep -o "[^{}]*")
else
  export dbstr=$(echo $(PGPASSWORD=$password psql -h $host -U $user -d $defaultDB -c "select array(select datname from pg_database where datname not in ('postgres', 'template1', 'template0', 'cloudsqladmin'))") | grep -o "{.*}" | grep -o "[^{}]*")
fi
export dbName=${dbName:-"$dbstr"}

if [ "$dbName" == "" ] ;then
  echo "No db found"
  return 1
fi
echo "db: $dbName"
declare -a db_array
IFS=',' read -ra db_array <<< "$dbName"
mkdir $timestamp

for i in "${db_array[@]}"
do
    export fileName="$timestamp/$i.tar"
    export output=$((PGPASSWORD=$password pg_dump -h $host -U $user -F t $i > $fileName) 2>&1)
    export err=$(echo $output | grep -o "error")
    if [ "$err" != "" ]; then
          echo "$i - $output"
          rm $fileName
    else
        echo "$i - success - $fileName"
    fi
done

export files="$(ls $timestamp | wc -l)"

if [ $files -gt 0 ]; then
    echo "Total $files files"
    gsutil cp -r $timestamp gs://stock-x-342909-backup-db
fi

rm -rf $timestamp