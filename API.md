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

  `/api/v1/users`
  
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
  
  * **Code**: `400 Bad Request`  <br />
    **Content**:
    ```json
    {
      "error": {"Type": "BAD_BINDING"},
      "data": null
    }
    ```
    **Notes**: 
    This error indicates that the POST payload was bad.
    
  OR
  
  * **Code**: `400 Bad Request` <br />
    **Content**:
    ```json
    {
      "error": {"Type": "USER_NOT_FOUND"},
      "data": null
    }
    ```  
  
  OR
  
  * **Code**: `400 Bad Request` <br />
    **Content**:
    ```json
    {
      "error": {"Type": "BAD_PASSWORD"},
      "data": null
    }
    ```
    **Notes**: <br />
    This indicates the password was not suitable for hashing
  
  
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
  
  * **Code:** `200 Success` <br />
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

    * **Code:** `400 Bad Request`  <br />
      **Content:**
      ```json
      {
        "error": {"Type": "BAD_BINDING"},
        "data": null
      }
      ```
      **Notes:** 
      This error indicates that the POST payload was bad.
      
    * **Code:** `400 Bad Request`  <br />
      **Content:**
      ```json
      {
        "error": {"Type": "USER_NOT_FOUND"},
        "data": null
      }
      ```
      **Notes:**  <br/>
      This error will be removed in the future to protect user info
      
    * **Code:** `401 Unauthorized`  <br />
      **Content:**
      ```json
      {
        "error": {"Type": "PASSWORD_MATCH_FAILED"},
        "data": null
      }
      ```
  
## Join Chat Room

Requests a websocket connection into a chatroom. Use `new WebSocket()` to use this
api call.

To use authentication, pass the JWT token string as the first protocol parameter to
the `WebSocket` object.

Communication over the websocket must be in json format. < TODO >

This api call can be made with or with authentication. If without authentication,
the user will join the chat room in a read only state. Any attempt to send a message
into the chat room will cause the connection to be terminated.

* **URL**

    `/api/v1/rooms/<room:string>/messages/ws`
    
* **URL Params**

    `room: string` - name of the room
    
* **Data Params**

    None
    
* **Success Response**

    * **Code:** `101 Change of Protocol`
    
* **Error Response**

    TODO
    
