## How to Test Postgres

To start the database and test the connection, follow these steps from the `blog-service` root directory.

1.  **Build and run the services:**
    This command will build your Go application, start the Postgres database, and run them together.

    ```bash
    docker-compose up --build
    ```

2.  **Connect to the database container:**
    Open a **second terminal window** and run the following command to get a shell inside the running Postgres container.

    ```bash
    docker exec -it gi25-group2-projectmonorepo-postgres-1 bash
    ```

3.  **Start the `psql` client:**
    Once inside the container, connect to your specific database using the credentials you defined.

    ```bash
    psql -U my_test_user -d my_test_db
    ```

4.  **Check your tables:**
    After your service has run its initialization logic, you can check if the `likes` and `comments` table was created correctly by describing it.
    ```sql
    \d comments
    \d likes
    ```

## How to Test MongoDB

To start your blog service and confirm that it successfully connects to MongoDB and creates the articles collection, follow these steps from the root of your project.

1.  **Build and run the services:**
    This command will build your Go application, start the Postgres database, and run them together.

    ```bash
    docker-compose up --build
    ```

2.  **Connect to the database container:**
    Open a **second terminal window** and run the following command to get a shell inside the running Postgres container.

    ```bash
    docker exec -it gi25-group2-projectmonorepo-mongodb-1 bash
    ```

3.  **Start the `mongo`:**

    Start the mongo shell

    ```bash
    mongosh
    ```

    Switch to the admin database to authenticate

    ```bash
      use admin
      Output: switched to db admin
    ```

    Authenticate with your root credentials

    ```bash
      db.auth('rootuser', 'rootpass')
      Output: 1
    ```

4.  **Check for the Articles Collection:**

    Switch to the blog_db

    ```bash
      use blog_db
      switched to db blog_db
    ```

    See the data you add

    ```bash
     db.articles.find().pretty()
      {

        "\_id" : ObjectId("..."),
        "content" : "This is a test article content.",
        "publisher_name" : "test_user",
        "category" : "test_category",
        "created_at" : ISODate("...")

    }
    ```

    Insert use to check connection

    ```bash
     db.articles.insertOne({ "test": "document" })
     {

            "acknowledged" : true,
            "insertedId" : ObjectId("68cbba8b2063ef2199f55e7c")

    }
    ```

    See db

    ```bash
      show dbs
    ```

    List the collections. You should see "articles"

    ```bash
      show collections
      articles
    ```

    Optional: Count the documents to confirm one was inserted (or more if you added)

    ```bash
      db.articles.countDocuments({})
      1 (or more if you added)
    ```
