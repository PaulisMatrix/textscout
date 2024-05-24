# circuitsearch
search API over in-memory movies dataset


# Approach

* Store the records in a database. 

* Query the db: `select * from movies where movie_name like "%'name'%" and movie_description like "%'desc'%"`
