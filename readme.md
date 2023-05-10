## Environment Variables

To run this project, you will need to add the following environment variables in your OS Env

`DB_HOST` : Your Database Host

`DB_PORT` : Your Database Port

`DB_USER` : Your Database Username

`DB_PASS` : Your Database Password

`DB_NAME` : Your Database Name

`APP_PORT` : Your APP Port


## Installation

1. Install makefile first
2. Install golang migrate for running database migration. Checkout <a href="https://github.com/golang-migrate/migrate">here</a> for details
3. Install MySQL Driver

## Running Migration

```bash
  migrate -path database/migrations -database "mysql://$DB_USER:$DB_PASS@tcp($DB_HOST:$DB_PORT)/$DB_NAME" up
```


## Running Test

```bash
  make test
```

## Running Application
```bash
  make run
```
    



## API Reference

#### Create Brand

```http
  POST /brand
```

| Body | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `title` | `string` | **Required**. Describe your brand title |

#### Create Product

```http
  POST /product
```

| Body | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `title`      | `string` | **Required**. title of the product |
| `description`      | `string` | **Optional**. Describe the detail of product |
| `brandId`      | `int` | **Required**. brandId of the product |
| `price`      | `decimal` | **Required**. Id of item to fetch |


#### Get Product By Id

```http
  GET /product?id=1
```

| Query Params | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Your Product Id |


#### Get Product By Brand

```http
  GET /product/brand?id=1
```

| Query Params | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Your Brand Id |


#### Create Order

```http
  POST /order
```
| Body | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `deliveryAddress`      | `string` | **Required**. title of the product |
| `details`      | `array` | **Required**. Your Detail Order. Check below for requirement |

#### Details
| Body | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `productId`      | `Int` | **Required**. Your Product |
| `qty`      | `Int` | **Required**. Qty you want to buy |


#### Get Order By Id

```http
  GET /order?id=1
```

| Query Params | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` | **Required**. Your Order Id |

I'm attached postman documentation in this repo too. You can check simple-ecommerce.postman_collection.json file for detail.
