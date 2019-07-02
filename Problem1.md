# Problem Statement

=> Load the data (in the accompanying file) into a datastore and analyse it
=> Build the following REST APIs using any Python API framework:
- Enable filter on category, brand, source, subcategory
- Search on title, sku
- List out all the products
- Update the product attributes like Brand, Category, Sub Category, Product
Type
=> Update the discount value for each product. Discount = ((mrp -
available_price)/(mrp)) * 100
=> Build an API to return count of products in each of the following discount buckets.
0%, 0-10%, 10-30%, 30-50%, >50%
*Bonus points if it’s optimized for speed
=> Feel free to try out extra things which you believe can improve the system in anyway
and most importantly enjoy the task!

Instructions
- Please use a relational datastore for storing data
- Please use a Python API framework to develop the APIs(optional)
- Please DO NOT USE an ORM
- Please make sure that the API responses are in standard and readable format
- Write a small doc that explains the data model and design decisions.
- Your code should be deployable on any standard Linux distro with minimal
intervention. Share a setup script and deployment instructions along with the rest of the
code.

# Solution

#### API's to be build

- GET /healthcheck
- GET /version
- GET /product -- List all products
- GET /product/{productId} -- Get info about a product
- GET /category/{categoryId}/ -- List all products in a category
- GET /category/{categoryId}/subcategory/{subcategoryId} -- List all products in a subcategory
- GET /category/{categoryId}/subcategory/{subcategoryId}/brand/{brandId} -- List all products in a brand

Note: All of the above API supports query params to enable filter on non hierarchical fields like sku, title, url...

- GET /products/count?discount_min=10&discount_max=20 -- Lists count of all the product having discount between 10 and 20
- PUT /products/discount/{discount_percentage} -- To update the discount

Note: Due to time constraint only implementing /product endpoint

#### Steps to run

- Install go >= 1.11
- Update config.json
- go build main.go
- copy the build/main and config to any linux machine
- Use ```./build/main config_path``` to run

#### Things which could be improved

- Currently all the data is in a single table. Ideally it should have been in different dimension and tables and then product info in one measures table.
So I would like to do ETL and split the tables accordingly during the load. 
- Add authentication/authorization
- Add unit tests
