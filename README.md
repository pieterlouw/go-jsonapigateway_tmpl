## Go Template project for a JSON API gateway ##

The project mocks a payment gateway that interfaces into a 3rd party JSON API that will do the actual authorization of the transaction.
Why a gateway into an API you might ask? In this application there's some data translation (mobile number to API userID) where an external app (i.e mobile app) might not 
have all the details required to interface directly into the remote API, so the mobile app do have the user's mobile number. Also, it won't be a good idea to have details of the remote payment API stored on a mobile device.
 
When the mobile application request a payment authorization it sends all the transaction details as well as the user's mobile number, then the local application will do a lookup in a local database (in this case a local text file) to get the userid associated with the user to be used in the remote API.

Effectively, this project is a template for any system with that kind of architecture pattern.

The local API is also a JSON API with the following endpoints:

* `/api/v1/ping` - GET request to ping if system is up, will also ping remote system 
* `/api/v1/payments` - POST request that will accept a JSON Payment request, translate, and pass on to remote system
* `/api/v1/payments/{transactionID}` - GET request to query status of a specific transaction on remote system

Transaction responses are stored locally in an embedded key/value stored, in this case [BoltDB](https://github.com/boltdb/bolt).There's normal web app GET request function to query a transaction response that will output the transaction response received that is stored in the embedded database as HTML: `/report/payments/{transactionID}`

###Things not included, but should be:####
* https listener
* Local API Authentication
* Tests

###Resources used when developing the project:###
* [Standard Go Package Layout - Ben Johnson](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1)
* [Making a RESTful JSON API in Go - Cory Lanou](http://thenewstack.io/make-a-restful-json-api-go/)
* [Building Sourcegraph, a large-scale code search engine in Go](https://text.sourcegraph.com/google-i-o-talk-building-sourcegraph-a-large-scale-code-search-cross-reference-engine-in-go-1f911b78a82e#.n7wg94okz)
* [BoltDB standard docs](https://github.com/boltdb/bolt)
* [Bolt — an embedded key/value database for Go](https://www.progville.com/go/bolt-embedded-db-golang/)
* [Intro to BoltDB: Painless Performant Persistence](http://npf.io/2014/07/intro-to-boltdb-painless-performant-persistence/)



