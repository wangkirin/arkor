# Details of Arkor API

## Common _Request_ Headers

The following table describes headers that can be used by various types of Arkor REST requests.

|Name|Description|Required|
|-----------|-----------|--------|
|`Host`|Standard HTTP Host Header. Should be set to the registry host|Yes|
|`Content-Length`|Length of the message (without the headers) according to RFC 2616. This header is required for PUTs and operations that load XML|Conditional|
|`Content-Type`|The content type of the resource in case the request content in the body. Example: text/plain|No|
|`Date`|The current date and time according to the requester|No|

## Common _Response_ Headers

The following table describes response headers that are common to most Arkor REST responses.

|Name|Description|Type|Default|Valid Values|
|----|-----------|----|-------|------------|
|`Content-Length`|The length in bytes of the body in the response|String|None|-|
|`Content-Type`|The MIME type of the content. e.g, Content-Type: text/html|String|None|-|
|`Connection`|Specifies whether the connection to the server is open or closed|Enum|None|{open,close}|
|`Date`|The date and time Arkor responded, for example, Wed, 01 Mar 2006 12:00:00 GMT|String|None|-|
|`Server`|The name of the server that created the response.|String|Arkor|-|
|`x-arkor-request-id`|A value created by Arkor that uniquely identifies the request|String|None|-|

## Error Response Format

### Error Response Headers

When there is an error, the header information contains:

- Content-Type: application/json

- An appropriate 3xx, 4xx, or 5xx HTTP status code

### List of Error Codes

|Error Code|Description|HTTP Status Code|
|----------|-----------|----------------|
|`BucketAlreadyExists`|The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again|`409 Conflict`|
|`BucketNotEmpty`|The bucket you tried to delete is not empty|`409 Conflict`|
|`EntityTooSmall`|Your proposed upload is smaller than the minimum allowed object size|`400 Bad Request`|
|`EntityTooLarge`|Your proposed upload exceeds the maximum allowed object size|`400 Bad Request`|
|`IncompleteBody`|You did not provide the number of bytes specified by the Content-Length HTTP header|`400 Bad Request`|
|`InternalError`|We encountered an internal error. Please try again|`500 Internal Server Error`|
|`InvalidArgument`|Invalid Argument|`400 Bad Request`|
|`InvalidBucketName`|The specified bucket is not valid|`400 Bad Request`|
|`InvalidDigest`|The Content-MD5 you specified is not valid|`400 Bad Request`|
|`InvalidRequest`|Invalid Request|`400 Bad Request`|
|`InvalidSecurity`|The provided security credentials are not valid|`403 Forbidden`|
|`MaxMessageLengthExceeded`|Your request was too big|`400 Bad Request`|
|`MetadataTooLarge`|Your metadata headers exceed the maximum allowed metadata size|`400 Bad Request`|
|`MethodNotAllowed`|The specified method is not allowed against this resource|`405 Method Not Allowed`|
|`MissingContentLength`|You must provide the Content-Length HTTP header|`411 Length Required`|
|`MissingRequestBodyError`|This happens when the user sends an empty xml document as a request|`400 Bad Request`|
|`NoSuchBucket`|The specified bucket does not exist|`404 Not Found`|
|`NoSuchKey`|The specified key does not exist|`404 Not Found`|
|`OperationAborted`|A conflicting conditional operation is currently in progress against this resource. Try again|`409 Conflict`|
|`RequestTimeout`|Your socket connection to the server was not read from or written to within the timeout period|`400 Bad Request`|

### Error Response Body

The body or the response also contains information about the error. The following sample error response shows the structure of response elements common to all REST error responses.

```json
{
  "Type": "Error",
  "Code": "NoSuchKey",
  "Message": "The resource you requested does not exist",
  "Resource": "/mybucket/myfoto.jpg",
  "RequestId": "4442587FB7D0A2F9"
}
```
The following table explains the REST error response elements

|Name|Description|
|----|-----------|
|`Type`|Type of response, here we set it as "Error"|
|`Code`|The error code is a string that uniquely identifies an error condition|
|`Message`|The error message contains a generic description of the error condition in English. It is intended for a human audience. |
|`RequestId`|ID of the request associated with the error|
|`Resource`|The bucket or object that is involved in the error.|

## Bucket

Include create, list, delete and check bucket.

<a name="getall"></a> 
### GET  / 
Fetch all buckets in Arkor


#### Request 

- **Syntax**
```http
GET /   HTTP/1.1
Host: <Arkor Host>
```

- **Parameters**

	_Noting Special_ , see [Common Request Headers]()

#### Response On Success

- **Syntax**
```
HTTP/1.1 200 OK
x-arkor-request-id: request id
Date: date
Content-Type: type
Content-Length: length
```

```json
{
  "Type": "ListAllMyBucketsResult",
  "Owner": {
    "ID": "id",
    "DisplayName": "arkor"
  },
  "Buckets": {
    "Bucket": {
      "Name": "quotes",
      "CreationDate": "date"
    }
  }
}
```

- **Elements:**

|Name|Description|Type|Children|Ancestor|
|----|-----------|----|--------|--------|
|`Type`|Type of response|-|-|None|
|`Bucket`|Container for bucket information|Container|Name, CreationDate|Buckets|
|`Buckets`|Container for one or more buckets|Container|Bucket|None|
|`CreationDate`|Date the bucket was created|Date|-|Bucket|
|`DisplayName`|Bucket owner's display name|String|-|Owner|
|`ID`|Bucket owner's user ID|String|-|Owner|
|`Name`|Bucket's name|String|-|Buckets.Bucket|
|`Owner`|Container for bucket owner information|Container|-|None|

#### Response On Failure

_Noting Special_ , see [Error Response Format]()


<a name="putbucket"></a> 
### PUT  /:bucket

#### Request 

- **Syntax:**
```http
HEAD /:bucket HTTP/1.1
Host: <Arkor Host>
Date: date
```

- **Parameters:**

  _Noting Special_ , see [Common Request Headers]()


#### Response On Success

- **Syntax:**
```
HTTP/1.1 200 OK
x-arkor-request-id: request id
Content-Type: type
Location: location
Date: date
Content-Length: 0
```

#### Response On Failure

  _Noting Special_ , see [Error Response Format]()

<a name="headbucket"></a> 
### HEAD  /:bucket

#### Request

- **Syntax:**
```http
HEAD /:bucket HTTP/1.1
Host: <Arkor Host>
Date: date
```
- **Parameters:**

  _Noting Special_ , see [Common Request Headers]()

#### Response On Success

- **Syntax:**
```
HTTP/1.1 200 OK
x-arkor-request-id: request id
Date: date
```

#### Response On Failure

  _Noting Special_ , see [Error Response Format]()

<a name="deletebucket"></a> 
### DELETE  /:bucket 

#### Request

- **Syntax:**
```http
DELETE /:bucket HTTP/1.1
Host: <Arkor Host>
Date: date
```
- **Parameters:**

  _Noting Special_ , see [Common Request Headers]()

#### Response On Success

- **Syntax:**
```
HTTP/1.1 200 OK
x-arkor-request-id: request id
Date: date
```

#### Response On Failure

  _Noting Special_ , see [Error Response Format]()


<a name="getbucket"></a> 
### GET  /:bucket 

#### Request

- **Syntax:**
```http
DELETE /:bucket HTTP/1.1
Host: <Arkor Host>
Date: date
```

- **Parameters:**

|Name|Location|Description|
|----|----|-----------|
|`prefix`|URL|Limits the response to keys that begin with the specified prefix|
|`max-keys`|URL|Sets the maximum number of keys returned in the response body|
|`delimiter`|URL|A delimiter is a character you use to group keys|


#### Response On Success

- **Syntax:**
```
HTTP/1.1 status_code
x-arkor-request-id: request id
Date: date
Content-Type: type
Content-Length: length

<Response Body>
```

```json
{
  "Type": "ListBucketResult",
  "Name": "oss-example",
  "MaxKeys": "100",
  "IsTruncated": "false",
  "Contents": [
    {
      "Key": "fun/movie/001.image",
      "LastModified": "2012-02-24T08:43:07.000Z",
      "ETag": "\"5B3C1A2E053D763E1B002CC607C5A0FE\"",
      "Type": "Normal",
      "Size": "344606",
      "StorageClass": "Standard",
      "Owner": {
        "ID": "00220120222",
        "DisplayName": "user-example"
      }
    },
    {
      "Key": "oss.jpg",
      "LastModified": "2012-02-24T06:07:48.000Z",
      "ETag": "\"5B3C1A2E053D763E1B002CC607C5A0FE\"",
      "Type": "Normal",
      "Size": "344606",
      "StorageClass": "Standard",
      "Owner": {
        "ID": "00220120222",
        "DisplayName": "user-example"
      }
    }
  ]
}
```

- **Elements:**

|Name|Description|Type|Ancestor|
|----|-----------|----|--------|
|`Contents`|Metadata about each object returned|Container|None|
|`CommonPrefixes`|All of the keys rolled up into a common prefix count as a single return when calculating the number of returns| String| None|
|`Delimiter`|Causes keys that contain the same string between the prefix and the first occurrence of the delimiter to be rolled up into a single result element in the CommonPrefixes collection|Date|None|
|`DisplayName`|Object owner's name|String|Contents.Owner|
|`ETag`|The entity tag is an MD5 hash of the object|String|Contents|
|`ID`|Object owner's ID|String|Contents.Owner|
|`IsTruncated`|Specifies whether (true) or not (false) all of the results were returned|Boolean|None|
|`Key`|The object's key|String|Contents|
|`LastModified`|Date and time the object was last modified|Date|Contents|
|`MaxKeys`|The maximum number of keys returned in the response body|String|None|
|`Name`|Name of the bucket|String|None|
|`Owner`|Bucket owner|String|Contents|
|`Prefix`|Keys that begin with the indicated prefix|String|None|
|`Size`|Size in bytes of the object|String|Contents|
|`StorageClass`|STANDARD,STANDARD_IA,REDUCED_REDUNDANCY,sGLACIER|String|Contents|
|`KeyCount`|Returns the number of keys included in the response|String|None|


#### Response On Failure

  _Noting Special_ , see [Error Response Format]()

## Object

Include create, list, delete and check object.


<a name="getobject"></a> 
### GET  /:bucket/:object
Downloads an object

#### Request

- **Syntax:**
```http
GET /:bucket/:object HTTP/1.1
Host: <Arkor Host>
Date: date
Range:bytes=byte_range

<Optional Additional Header>
```
- **Parameters:**

|Name|Location|Description|
|----|----|-----------|
|`Range`|Header|Downloads the specified range bytes of an object|
|`If-Modified-Since`|Header|Return the object only if it has been modified since the specified time, otherwise return 304|
|`If-Unmodified-Since`|Header|Return the object only if it has not been modified since the specified time, otherwise return 412|
|`If-Match`|Header|Return the object only if it has not been modified since the specified time, otherwise return 412|
|`If-None-Match`|Header|Return the object only if it has been modified since the specified time, otherwise return 304|
|`response-content-type`|URL|Sets the Content-Type header of the response|
|`response-content-language`|URL|Sets the Content-Language header of the response|
|`response-expires`|URL|Sets the Expires header of the response|
|`response-cache-control`|URL|Sets the Cache-Control header of the response|
|`response-content-disposition`|URL|Sets the Content-Disposition header of the response|
|`response-content-encoding`|URL|Sets the Content-Encoding header of the response|


#### Response On Success

- **Syntax:**
```
HTTP/1.1 200 OK
x-arkor-request-id: request id
Content-Type: type
Date: date
Content-Length: length
Etag： etag
Last-Modified： time

<Object Content>
```

#### Response On Failure

  _Noting Special_ , see [Error Response Format]()


<a name="putobject"></a> 
### PUT  /:bucket/:object
Uploads an object via PUT

#### Request

- **Syntax:**
```http
PUT /:bucket/:object HTTP/1.1
Host: Host Server
Content-Type: type
Content-Length: length
Authorization: authorization
Date: date

<Optional Additional Header>

<object Content>
```
- **Parameters:**

|Name|Location|Description|
|----|----|-----------|
|`Content-MD5`|Header|The base64-encoded 128-bit MD5 digest of the message (without the headers) according to RFC 1864|
|`Content-Disposition`|Header|Specifies presentational information for the object|
|`x-arkor-website-redirect-location`|Header|If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL.|


#### Response On Success

- **Syntax:**
```
HTTP/1.1 200 OK
x-arkor-request-id: request id
Date: date
Content-Length: length
Content-Type: type
```

#### Response On Failure

  _Noting Special_ , see [Error Response Format]()


<a name="putdestination"></a> 
### PUT  /:bucket/:destination
Copy a source object into a new one

#### Request

- **Syntax:**
```http
PUT  /:bucket/:destination HTTP/1.1
Host: <Arkor Host>
x-arkor-copy-source: /sourceBucket/sourceObject
Date: date

```
- **Parameters:**

|Name|Location|Description|
|----|----|-----------|
|`x-arkor-copy-source`|Header|The name of the source bucket and key name of the source object, separated by a slash (/)|
|`x-arkor-metadata- directive`|Header|Specifies whether the metadata is copied from the source object or replaced with metadata provided in the request, COPY or REPLACE|
|`x-arkor-copy-source-if-match`|Header|Copies the object if its entity tag (ETag) matches the specified tag; otherwise, the request returns 412|
|`x-arkor-copy-source-if-none-match`|Header|Copies the object if its entity tag (ETag) is different than the specified ETag;otherwise, the request returns 412|
|`x-arkor-website-redirect-location`|Header|If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL.|
|`x-arkor-copy-source-if-unmodified- since`|Header|Copies the object if it hasn't been modified since the specified time; otherwise, the request returns 412|
|`x-arkor-copy-source-if-modified-since`|Header|Copies the object if it has been modified since the specified time;otherwise, the request returns 412|


#### Response On Success

- **Syntax:**
```
HTTP/1.1 200 OK
x-arkor-request-id: request id
x-reserved: Arkor
Content-Type: type 
Date: date
Content-Length: length
```
```json
{
  "Type": "CopyObjectResult",
  "LastModified": "Fri, 24 Feb 2012 07:18:48 GMT",
  "ETag": "\"5B3C1A2E053D763E1B002CC607C5A0FE\""
}
```
- **Parameters:**

|Name|Location|Description|
|----|----|-----------|
|`x-arkor-copy-source-version-id`|Header|Version of the source object that was copied|
|`x-arkor-version-id`|Header|Version of the copied object in the destination bucket|

- **Elements:**

|Name|Description|Type|Ancestor|
|----|-----------|----|--------|
|`ETag`|Returns the ETag of the new object.|String|None|
|`LastModified`|Returns the date the object was last modified|String|None|

#### Response On Failure

  _Noting Special_ , see [Error Response Format]()

<a name="deleteobject"></a> 
### DELETE  /:bucket/:object
Removes an object

#### Request

- **Syntax:**
```http
DELETE /:bucket/:object HTTP/1.1
Host: <Arkor Host>
Date: date
```
- **Parameters:**

  _Noting Special_ , see [Common Request Headers]()

#### Response On Success

- **Syntax:**
```
HTTP/1.1 200 OK
x-arkor-request-id: request id
Date: date
```

- **Parameters:**

|Name|Location|Description|
|----|----|-----------|
|`x-arkor-delete-marker`|Header|Specifies whether the versioned object that was permanently deleted was (true) or was not (false) a delete marker|
|`x-arkor-version-id`|Header|Version of the copied object in the destination bucket|



#### Response On Failure

  _Noting Special_ , see [Error Response Format]()


<a name="headobject"></a> 
### HEAD  /:bucket/:object
Gets metadata of an object

#### Request

- **Syntax:**
```http
HEAD /:bucket/:object HTTP/1.1
Host: <Arkor Host>
Date: date
```
- **Parameters:**

  _Noting Special_ , see [Common Request Headers]()

#### Response On Success

- **Syntax:**
```
HTTP/1.1 200 OK
x-arkor-request-id: request id
Content-Type: type
Date: date
Content-Length: length
Etag: etag
Last-Modified: time
```


#### Response On Failure

  _Noting Special_ , see [Error Response Format]()


<a name="postobjects"></a> 
### POST  /objects
init the upload of an object from git-lfs API

#### Request

- **Syntax:**
```http
POST  /objects HTTP/1.1
Accept: application/vnd.git-lfs+json
Content-Type: application/vnd.git-lfs+json
```

```json
{
   "oid": "1111111",
   "size": 123
}
```

#### Response On Success

- **Syntax:**
```
HTTP/1.1 202 Accepted
Content-Type: application/vnd.git-lfs+jso
```
```json
{
   "_links": {
     "upload": {
       "href": "https://some-upload.com",
       "header": {
         "Key": "value"
       }
     },
     "verify": {
       "href": "https://some-callback.com",
       "header": {
         "Key": "value"
       }
     }
   }
}
```

A response can include one of multiple link relations, each with an `href`
property and an optional `header` property.

- **Elements:**

|Name|Description|
|----|-----------|
|`upload`|This relation describes how to upload the object. Expect this with a 202 status.|
|`verify`|Returns the ETag of the new object.|
|`download`|This relation describes how to download the object content. This only appears on a 200 status|

#### Response Code

|HTTP Status Code|Description|
|----------|-----------|
|200|The object already exists. Don't bother re-uploading|
|202|The object is ready to be uploaded. Follow the "upload" and optional "verify" links|
|401|The authentication credentials are incorrect|
|403|The user has read, but not write access|
|404|The repository does not exist for the user|

<a name="getoid"></a> 
### GET  /objects/:oid
gets the object's meta data from git-lfs API

#### Request

- **Syntax:**
```http
GET /objects/:oid HTTP/1.1
Accept: application/vnd.git-lfs+json
```


#### Response On Success

- **Syntax:**
```
HTTP/1.1 200 OK
Content-Type: application/vnd.git-lfs+jso
```
```json
 {
   "oid": "the-sha-256-signature",
   "size": 123456,
   "_links": {
     "self": {
       "href": "https://git-lfs-server.com/objects/OID",
     },
     "download": {
       "href": "https://some-download.com",
       "header": {
         "Key": "value"
       }
     }
   }
 }
```


The oid and size properties are required. A hypermedia _links section is included with a download link relation. Clients can follow this link to access the object content. See the "Hypermedia" section above for more.

#### Response Code

|HTTP Status Code|Description|
|----------|-----------|
|200|The object exists and the user has access to download it|
|401|The authentication credentials are incorrect|
|404|The user does not have access to the object, or it does not exist.|
|410|The object used to exist, but was deleted. The message should state why (user initiated, legal issues, etc).|

<a name="postbatch"></a> 
### POST  /objects/batch
retrieves the metadata for a batch objects from git-lfs API

#### Request

- **Syntax:**
```http
POST /objects/batch HTTP/1.1
Accept: application/vnd.git-lfs+json
Content-Type: application/vnd.git-lfs+json
```

```json
 {
  "$schema": "http://json-schema.org/draft-04/schema",
  "title": "Git LFS HTTPS Batch API v1.3 Request",
  "type": "object",
  "properties": {
    "transfers": {
      "type": "array",
      "items": {
        "type": "string"
      },
    },
    "operation": {
      "type": "string"
    },
    "objects": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "oid": {
            "type": "string"
          },
          "size": {
            "type": "number"
          }
        },
        "required": ["oid", "size"],
        "additionalProperties": false
      }
    }
  },
  "required": ["objects", "operation"]
}
```

#### Response On Success

- **Syntax:**
```
HTTP/1.1 200 OK
Content-Type: application/vnd.git-lfs+json
```
```json
{
  "$schema": "http://json-schema.org/draft-04/schema",
  "title": "Git LFS HTTPS Batch API v1.3 Response",
  "type": "object",

  "definitions": {
    "action": {
      "type": "object",
      "properties": {
        "href": {
          "type": "string"
        },
        "header": {
          "type": "object",
          "additionalProperties": true
        },
        "expires_at": {
          "type": "string"
        }
      },
      "required": ["href"],
      "additionalProperties": false
    }
  },

  "properties": {
    "transfer": {
      "type": "string"
    },
    "objects": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "oid": {
            "type": "string"
          },
          "size": {
            "type": "number"
          },
          "actions": {
            "type": "object",
            "properties": {
              "download": { "$ref": "#/definitions/action" },
              "upload": { "$ref": "#/definitions/action" },
              "verify": { "$ref": "#/definitions/action" }
            },
            "additionalProperties": false
          },
          "error": {
            "type": "object",
            "properties": {
              "code": {
                "type": "number"
              },
              "message": {
                "type": "string"
              }
            },
            "required": ["code", "message"],
            "additionalProperties": false
          }
        },
        "required": ["oid", "size"],
        "additionalProperties": false
      }
    },
    "message": {
      "type": "string"
    },
    "request_id": {
      "type": "string"
    },
    "documentation_url": {
      "type": "string"
    },
  },
  "required": ["objects"]
}
```


