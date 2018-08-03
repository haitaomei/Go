# Go test if the client set the specific json field

For example, 

        type User struct {
            Name string `json:"name"`
            ID   int    `json:"id"`
        }
When ``ID`` is zero after decodin, we can not tell if the client sent zero, or omited it.
Therefore, we can use pointers:

        type UserForRead struct {
            Name *string `json:"name"`
            ID   *int    `json:"id"`
        }
Decoding use this struct, and then we can test if ID==nil.

``curl --header "Content-Type: application/json"   --request GET   --data '{"name":"james","id":123}' http://127.0.0.1:60443/``

Gets ``ID is set``

``curl --header "Content-Type: application/json"   --request GET   --data '{"name":"james"}' http://127.0.0.1:60443/``

Gets ``ID is omited in the request``



Reject Unknown Json fields
---------

        decoder := json.NewDecoder(httpReq.Body)
        decoder.DisallowUnknownFields()
        err = decoder.Decode(&cfg)
