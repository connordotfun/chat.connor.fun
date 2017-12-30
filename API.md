# CHAT.CONNOR.FUN API

The API expects all data in JSON format. The generic response format is:
```json
{
  "error": {"type": "string", "message": "string"},
  "data": {}
}
```
All JSON data will be reported in this format. If a API call does not
return this format, it will return no content.

The one exception to this is API calls that get upgraded to a websocket.

## Create User
Creates a new user account if one does not exist for the given username. Currently
there are no special requirements for username or passwords. 

* **URL**

  */api/v1/users*
  
* **Method**

  `POST`

* **URL Params**

  None

* **Data Params**
  
  **Content**:
  ```json
  {
    "username": "someusername",
    "secret": "somepassword"
  }
  ```
* **Success Response**
  
  * **Code**: `201 CREATED`  <br />
  **Content**:
  ```json
  {
    "error": null,
    "data": {
      "id": 12,
      "username": "someusername"
    }
  }
  ```

* **Error Response**  
  TODO
  
  
## Login
  Checks the user credentials and returns a valid JWT for that user.

* **URL**

  `/api/v1/login`

* **Method:**
  
  `POST`
  
*  **URL Params**

   None

* **Data Params**

  * **Content**
  ```json
  {"username": "foobar", "secret": "123xyz"}
  ```
* **Success Response:**
  
  <_What should the status code be on success and is there any returned data? This is useful when people need to to know what their callbacks should expect!_>

  * **Code:** 200 <br />
    **Content:** 
    ```json
    {
        "error": null,
        "data": {
          "token": "jwt-token-string",
          "user": {
            "id": 12,
            "username": "foobar"
          }
        }
    }
    ```
 
* **Error Response:**

  TODO