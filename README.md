# GoCrawler
Web crawler in GO

For a given set of urls perform a GET request and search for a specific field in that response ('Total' property).

An example of an expected url would be something like this: 

`https://www.farfetch.com/pt/shopping/women/fendi/items.aspx?format=json`

This is what we call a _*plp*_ (product list page). With this crawler I'm interested in retrieving the value of property _"Total"_ that is present inside the object _PagingOptions_.


```json
{
    ...
    "PagingOptions": {
        "Total": 1596,
        ...
    },
    ...
}
```

Run:

`go run crawler.go <input_file> <output_file>` 

