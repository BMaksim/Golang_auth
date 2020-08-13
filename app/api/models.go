package api

type Tokens struct {
	GUID          string `json:"_id,omitempty" bson:"_id,omitempty"`
	Access_token  string `json:"access_token,omitempty" bson:"access_token,omitempty"`
	Refresh_token string `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	IP            string `json:"ip,omitempty" bson:"ip,omitempty"`
}
