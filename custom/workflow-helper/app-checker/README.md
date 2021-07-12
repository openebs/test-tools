## API-Checker 
It present in the sock-shop namespace which sends API requests: GET and POST.
Api checker keeps sending API requests to sock-shop applications various services. 

Example logs:
```
[Status] : API Checker has been started
[Status]: 200  User successfully logged in
[Info]: Unique ID for user is: xP6RDO0yoPtZhaoHdgvjPJUsYSk8erOq
[Info]: Sending request to front-end
[Status]:  ResponseCode: 200  FrontEnd is accessible
[Info]: Adding new customer 
[Status]: ResponseCode: 200  Customer added successfully with user Name: test_user_895455cf-5832-4ce7-bfdb-bc95082db56c
[Status]: ResponseCode: 200  Catalogue get request successfully send
[Info]: Catalogue Item: {'id': '3395a43e-2d88-40de-b95f-e00e1502085b', 'name': 'Colourful', 'description': 'proident occaecat irure et excepteur labore minim nisi amet irure', 'imageUrl': ['/catalogue/images/colourful_socks.jpg', '/catalogue/images/colourful_socks.jpg'], 'price': 18, 'count': 438, 'tag': ['brown', 'blue']}
[Info]: Adding card details for purchase
[Status]: ResponseCode: 200 Card details has been successfully added
[Info]: Adding Address for user
[Status]: ResponseCode: 200 Address has been added successfully
[Status]: ResponseCode: 200 Card ID: {'longNum': '5544154011345918', 'expires': '08/19', 'ccv': '958', 'id': '57a98d98e4b00679b4a830b1', '_links': {'card': {'href': 'http://user/cards/57a98d98e4b00679b4a830b1'}, 'self': {'href': 'http://user/cards/57a98d98e4b00679b4a830b1'}}}
[Status]: ResponseCode: 200  Address has been retrieved successfully
[Info]: Getting item details
[Status]: ResponseCode: 200  Item details successfully received
```