# Arkor HTTP APIs

## Introduction


The **Arkor HTTP APIs** is the interfaces between user and Arkor, user can easily manage objects in Arkor through standard RESTful style HTTP request. For better use, the Arkor APIs have similar concepts and operations with current Object storage service.

## Concepts

- **Object**

The fundamental entity type stored in Arkor. Objects consist of object key, object data and metadata. Object key defines the name of object, each object in a bucket must have a unique object name. Object date is the data content of object. Object metadata describes the information of Object data.

- **Bucket**

Bucket is a container for stored objects and its has a global unique name. Every object is contained in a bucket. Unlike file system, object storage is flat that all objects in bucket at the same logical level.


## API List

|Method|Path|Entity|Description|
|------|----|------|-----------|
|GET|[`/`](API_details.md/#getall)|bucket|Lists all buckets|
|PUT|[`/:bucket`](API_details.md/#putbucket)|bucket|Create a new bucket|
|HEAD|[`/:bucket`](API_details.md/#headbucket)|bucket|Checks if a bucket exists|
|DELETE|[`/:bucket`](API_details.md/#deletebucket)|bucket|Removes a bucket|
|GET|[`/:bucket`](API_details.md/#getbucket)|bucket|Lists objects in a bucket|
|GET|[`/:bucket/:object`](API_details.md/#getobject)|object|Downloads an object|
|PUT|[`/:bucket/:object`](API_details.md/#putobject)|object|Uploads an object via PUT|
|PUT|[`/:bucket/:destination`](API_details.md/#putdestination)|object|Copy a source object into a new one|
|HEAD|[`/:bucket/:object`](API_details.md/#headobject)|object|Gets metadata of an object|
|DELETE|[`/:bucket/:object`](API_details.md/#deleteobject)|object|Removes an object|
|POST|[`/objects`](API_details.md/#postobjects)|object|init the upload of an object from git-lfs API|
|GET|[`/objects/:oid`](API_details.md/#getoid)|object|gets the object's meta data from git-lfs API|
|POST|[`/objects/batch`](API_details.md/#postbatch)|object|retrieves the metadata for a batch objects from git-lfs API|

## Details

To get the syntax, parameters and the response of each API, see [HERE](API_details.md/) or click the `Path` of it. 
