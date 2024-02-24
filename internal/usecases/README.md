## :bookmark_tabs: Use cases

The following use cases must be agnostic to communication approach.  It can be used in RESTful endpoints, gRPC and messaging.

They also represents one or more types of operations: *IO bound operations*, *CPU bound operations* and *Memory bound operations*

- `create.go`: create a product. *IO bound operation* because it uses database connection and query.
- `report.go`: generate a products report. *Memory bound operation* and *IO bound operation* because it uses database connection and query to get all products then create a file in the system (PDF) and it holds all products (from `products` table) in memory.
- `getbydiscount.go`: get products, apply a discount and sort. *Memory bound operation*, *IO bound operation* and *CPU bound operation* because it get products from database, holds into memory and then execute some logic such as calculate price with discount and order by discount applied.

:bulb: The `ìnterfaces.go` files holds the common interfaces.