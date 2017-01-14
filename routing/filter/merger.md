Merger Usage
===
### SimpleMerger
the simple merger provide filter to meger multi json result to single .

##### Client
flow example

* if i want to call thread api like blow
  * the get `http://localhost/api/a?val=abc` result is
  
  ```
  {
    "code": 0,
    "data": "abc"
  }
  ```
  * the get `http://localhost/api/b?val=123` result is
  
  ```
  {
    "code": 0,
    "data": {
        "val": "123"
    }
  }
  ```
  * the post `http://localhost/api/c` data and result is
  
  ```
  {
    "a": 0,
    "b": "xyz"
  }
  ```
  
  ```
  {
    "code": 0,
    "data": {
        "a": 0,
        "b": "xyz"
    }
  }
  ```
  * the post `http://localhost/api/d` data and result is
  
  ```
  {
    "c": 0,
    "d": "xyz"
  }
  ```
  
  ```
  {
    "code": 0,
    "data": {
        "c": 0,
        "d": "xyz"
    }
  }
  ```
* the server can merger to single api like `http://localhost/merger/x`, and the call case is
  * use `merger` query argument to call speical api, like `a,b,c` is call all api, `a,b` is ingore `c`
  * the `GET` argument will be parsed `http://localhost/merger/x?merger=a,b,c,d&a.val=abc&b.val=123`
  * the `POST` body will be parsed by 
  
   ```
  {
    "c": {
        "a": 0,
        "b": "xyz"
    },
    "d": {
        "c": 0,
        "d": "xyz"
    }
  }
  ```
  * the result will be parsed by
  
  ```
  {
    "code": 0,
    "data": {
        "a": "abc",
        "b": {
            "val": "123"
        },
        "c": {
            "a": 0,
            "b": "xyz"
        },
        "d": {
            "c": 0,
            "d": "xyz"
        }
    }
  }
```

##### Server
* for the merger configure, flowing `merger.properties`
* register merger by `filter.HandMerger`
  
  