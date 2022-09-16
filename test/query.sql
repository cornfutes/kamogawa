psql -h localhost -p 5432 -U local-dev -d shimogawa_db -c "TRUNCATE gce_instance_dbs, project_dbs, gae_service_dbs, gae_version_dbs, gae_instance_dbs;";
psql -h localhost -p 5432 -U local-dev -d shimogawa_db -c "COPY gce_instance_dbs ($(head -1 ./test/instances.csv)) FROM STDIN delimiter ',' CSV HEADER;" < ./test/instances.csv ;
psql -h localhost -p 5432 -U local-dev -d shimogawa_db -c "COPY project_dbs ($(head -1 ./test/projects.csv)) FROM STDIN delimiter ',' CSV HEADER;" < ./test/projects.csv;
psql -h localhost -p 5432 -U local-dev -d shimogawa_db -c "COPY gae_service_dbs FROM STDIN delimiter ',';" < ./test/gaeservices.csv ;
psql -h localhost -p 5432 -U local-dev -d shimogawa_db -c "COPY gae_version_dbs FROM STDIN delimiter ',';" < ./test/gaeversions.csv ;
psql -h localhost -p 5432 -U local-dev -d shimogawa_db -c "COPY gae_instance_dbs FROM STDIN delimiter ',';" < ./test/gaeinstances.csv ;
