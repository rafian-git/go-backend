export user="postgres"
export defaultDB='postgres'

if [ "$host" == "" ] ;then
  echo "host info is required"
  return 1
fi

if [ "$password" == "" ] ;then
  echo "password is required"
  return 1
fi

if [ "$bucket" == "" ] ;then
  echo "Gcloud bucket name is required"
  return 1
fi

if [ "$folder" == "" ] ;then
  echo "Gcloud bucket's folder name is required"
  return 1
fi

if [ "$file" == "" ] ;then
  gsutil cp -r gs://$bucket/$folder .
else
  gsutil cp gs://$bucket/$folder/$file $folder/$file
fi

cd $folder

for file in *; do
  export dbName="${file%.*}"
  export found=$(echo $(PGPASSWORD=$password psql -h $host -U $user -d postgres -c "select array (select datname from pg_database where datname='$dbName')") | grep -o "{.*}" | grep -o "[^{}]*")
  export restore_db=$dbName
  if [ "$found" != "" ] ;then
    export restore_db="restored_$dbName"
  fi
  echo "Restoring $dbName in $restore_db"


  PGPASSWORD=$password createdb -h $host  -U $user $restore_db

  PGPASSWORD=$password pg_restore -h $host -U $user --dbname $restore_db  --verbose $file
done
cd ..

rm -rf $folder