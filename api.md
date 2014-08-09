
# Router API


## <a name="contents"></a> Table Of Contents
* [General Info](#gen)
 - [Authentication](#auth)
 - [Response Status](#resstat)
* [Pools](#pools)
 - [Get List of Pool Names](#getpools)
 - [Get Pool](#getpool)
 - [Update/Create Pool](#setpool)
 - [Delete Pool](#delpool)
 - [Get List of Hosts](#listhosts)
 - [Add List of Hosts](#addhosts)
 - [Delete List of Hosts](#delhosts)
* [Tries](#tries)
 - [Get List of Trie Names](#gettries)
 - [Get Trie](#gettrie)
 - [Update/Create Trie](#settrie)
 - [Delete Trie](#deltrie)
* [Rules](#rules)
 - [Get List of Rule Names](#getrules)
 - [Get Rule](#getrule)
 - [Update/Create Rule](#setrule)
 - [Delete Rule](#delrule)
* [Ports](#ports)
 - [Get List of Port Names](#getports)
 - [Get Port](#getport)
 - [Update/Create Port](#setport)
 - [Delete Port](#delport)

## <a name="gen"></a> General

#### <a name="auth"></a> Authentication

This API uses HTTP simple auth which essentially means two headers must be passed in with EVERY request. Those headers should be named User and Secret, and their values should be the username as well as the associated secret with that user. For more information on HTTP Simple Auth. see:  http://en.wikipedia.org/wiki/Basic_access_authentication

#### <a name="resstat"></a> Response Status

All of the API Methods will return an HTTP Status code, a list of possible response codes and their in-context interpretations are below:


Some of the API Methods are shown to return a JSON field "Status". This field is used in conjunction with the HTTP Status codes to build an accurate representation of the outcome of associated API Request.

***
[top](#contents)

## <a name="pools"></a> Pools

#### <a name="getpools"></a> Get List of Pool Names

	GET /pools

###### Response

``` json
200

{
	"Pools": [
    	"pool1",
    	"pool2",
    	"pool3"
    ]
}
```
***

#### <a name="getpool"></a> Get Pool

	GET /pools/:poolname

###### Response

``` json
200

{
	"Name": "pool1",
    "Internal": false,
    "Hosts": {
    	"host1": {"Address": "host1Addr"},
        "host2": {"Address": "host2Addr"}
    },
    "Config": {
    	"HealthzEvery": "10",
        "HealthzTimeout": "5",
        "RequestTimeout": "2",
        "Status": "status string"
    }
}
```
***

#### <a name="setpool"></a> Update/Create Pool

	PUT /pools/:poolname

###### Request

``` json
{
	"Name": "pool1",
    "Internal": false,
    "Hosts": {
    	"host1": {"Address": "host1Addr"},
        "host2": {"Address": "host2Addr"}
    },
    "Config": {
    	"HealthzEvery": "10",
        "HealthzTimeout": "5",
        "RequestTimeout": "2",
        "Status": "status string"
    }
}

```

###### Response

``` json
200

{
	"Status": "Successfully added/updated :poolname"
}
```
***

#### <a name="delpool"></a> Delete Pool

	DELETE /pools/:poolname

###### Response

``` json
200

{
	"Status": "Successfully deleted :poolname"
}
```
[top](#contents)
***

#### <a name="gethosts"></a> Get A List of Hosts

	GET /pools/:poolname/hosts

###### Response

``` json
200

{
	"Hosts:"[
    	{"Address": "host1Address"},
        {"Address": "host2Address"},
        {"Address": "host3Address"}
    ]
}
```
***

#### <a name="addhosts"></a> Add A List of Hosts
	PUT /pools/:poolname/hosts

###### Request

``` json
{
	"Hosts":[
    	{"Address": "newHost1Addr"},
        {"Address": "newHost2Addr"},
        {"Address": "newHost3Addr"}
    ]
}
```

###### Response

``` json
200

{
	"Status": "Succesfully added hosts to :poolname"
}
```
***

#### <a name="delhosts"></a> Delete a List of Hosts

	DELETE /pools/:poolname/hosts

###### Request

``` json
{
	"Hosts":[
    	{"Address": "host1Addr"},
        {"Address": "host2Addr"},
        {"Address": "host3Addr"}
    ]
}
```

###### Response

``` json
200

{
	"Status": "Succesfully deleted hosts from :poolname"
}
```
[top](#contents)
***


## <a name="tries"></a> Tries


### <a name="gettries"></a> Get List of Trie Names

	GET /tries

###### Response

``` json
200

{
	"Tries": [
    	"trie1", "trie2", "trie3"
 	]
}
```
***

### <a name="gettrie"></a> Get Trie

	GET /tries/:triename

###### Response

``` json
200

{
	"Name": "tName",
    "Rules": [
    	"rule1", "rule2", "rule3"
    ],
    "Internal": false
}
```
***

### <a name="settrie"></a> Update/Create Trie

	PUT /tries/:triename

###### Request

``` json
{
	"Name": "newTrie",
    "Rules": [
    	"rule1", "rule2", "rule3"
    ],
    "Internal": false
}

```

###### Response

``` json
200

{
	"Status": "Successfully added/updated :triename"
}
```
***

### <a name="deltrie"></a> Delete Trie

	DELETE /tries/:triename

###### Response

``` json
200

{
	"Status": "Successfully deleted :triename"
}
```
[top](#contents)
***

## <a name="rules"></a>	Rules


### <a name="getrules"></a> Get List of Rule Names

	GET /rules

###### Response

``` json
200

{
	"Rules": [
    	"rule1", "rule2", "rule3"
    ]
}
```
***

### <a name="getrule"></a> Get Rule

	GET /rules/:rulename

###### Response

``` json
200

{
	"Name": "rName",
    "Type": "rType",
    "Value": "rVal",
    "Next": "next",
    "Pool": "pName",
    "Internal": false
}
```
***

### <a name="setrule"></a> Update/Create Rule

	PUT /rules/:rulename

###### Request

``` json
{
	"Name": "newRule",
    "Type": "rType",
    "Value": "rVal",
    "Next": "next",
    "Pool": "pName",
    "Internal": false
}
```

###### Response

``` json
200

{
	"Status": "Successfully added/updated :rulename"
}
```
***

### <a name="delrule"></a> Delete Rule

	DELETE /rules/:rulename

###### Response

``` json
200

{
	"Status": "Successfully deleted :rulename"
}
```
[top](#contents)
***
## <a name="ports"></a>	Ports

### <a name="getports"></a> Get List of Port Values

	GET /ports

###### Response

``` json
200

{
	"Ports": [
    	28080, 8080, 2033
    ]
}
```
***

### <a name="getport"></a> Get Port

	GET /ports/:portname

###### Response

``` json
200

{
	"Port": 28080,
    "Trie": "trie1",
    "Internal": false
}
```
***

### <a name="setport"></a> Update/Create Port

	PUT /ports/:portname

###### Request

``` json
{
	"Port": 28080,
    "Trie": "trie1",
    "Internal": false
}
```

###### Response

``` json
200

{
	"Status": "Successfully added/updated :portname"
}
```
***

### <a name="delport"></a> Delete Port

	DELETE /ports/:portname

###### Response

``` json
200

{
	"Status": "Successfully deleted :portname"
}
```
[top](#contents)
***

