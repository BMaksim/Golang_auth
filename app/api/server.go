package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	config *Config
	router *mux.Router
	db     *DB
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,
		router: mux.NewRouter(),
	}
}

func (s *Server) StartServer() error {
	if err := s.setDB(); err != nil {
		return err
	}
	s.setRoutes()
	port := os.Getenv("PORT")
	return http.ListenAndServe(":"+port, s.router)
}

func (s *Server) setRoutes() {
	s.router.HandleFunc("/getTokens", s.getTokens())
	s.router.HandleFunc("/refreshTokens", s.refreshTokens())
	s.router.HandleFunc("/deleteToken", s.deleteToken())
	s.router.HandleFunc("/deleteTokens", s.deleteTokens())
}

func (s *Server) setDB() error {
	db := NewDB(s.config)
	if err := db.Open(); err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Server) deleteTokens() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		t := Tokens{}
		err := decoder.Decode(&t)

		coll := s.db.conn.Collection("tokens")
		ctx := context.Background()

		//get data from db
		findResult := coll.FindOne(ctx, bson.M{"_id": t.GUID, "ip": r.RemoteAddr})
		if err := findResult.Err(); err != nil {
			fmt.Println(err)
		}

		data := Tokens{}
		err = findResult.Decode(&data)
		if err != nil {
			fmt.Println(err)
		}

		//decode refresh token from user
		decodetRt, err := base64.StdEncoding.DecodeString(t.Refresh_token)
		t.Refresh_token = string(decodetRt)
		if err != nil {
			fmt.Println(err)
		}

		// compare refresh tokens from user and db
		err = bcrypt.CompareHashAndPassword([]byte(data.Refresh_token), []byte(t.Refresh_token))
		if err == nil {
			_, err := coll.DeleteMany(ctx, bson.M{"_id": t.GUID})
			if err == nil {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Successfully deleted!"))
				return
			}

			fmt.Println(err)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Something bad happened!"))
		}
	}
}

func (s *Server) deleteToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		t := Tokens{}
		err := decoder.Decode(&t)

		coll := s.db.conn.Collection("tokens")
		ctx := context.Background()

		//get data from db
		findResult := coll.FindOne(ctx, bson.M{"_id": t.GUID, "ip": r.RemoteAddr})
		if err := findResult.Err(); err != nil {
			fmt.Println(err)
		}

		data := Tokens{}
		err = findResult.Decode(&data)
		if err != nil {
			fmt.Println(err)
		}

		//decode refresh token from user
		decodetRt, err := base64.StdEncoding.DecodeString(t.Refresh_token)
		t.Refresh_token = string(decodetRt)
		if err != nil {
			fmt.Println(err)
		}

		// compare refresh tokens from user and db
		err = bcrypt.CompareHashAndPassword([]byte(data.Refresh_token), []byte(t.Refresh_token))
		if err == nil {
			_, err := coll.DeleteOne(ctx, bson.M{"_id": t.GUID, "ip": r.RemoteAddr})
			if err == nil {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Successfully deleted!"))
				return
			}

			fmt.Println(err)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Something bad happened!"))
		}
	}
}

func (s *Server) refreshTokens() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		t := Tokens{}
		err := decoder.Decode(&t)
		if err != nil {
			fmt.Println(err)
		}
		coll := s.db.conn.Collection("tokens")
		ctx := context.Background()

		//get data from db
		findResult := coll.FindOne(ctx, bson.M{"_id": t.GUID, "ip": "[::1]:47124"})
		if err := findResult.Err(); err != nil {
			fmt.Println(err)
		}

		data := Tokens{}
		err = findResult.Decode(&data)
		if err != nil {
			fmt.Println(err)
		}

		//decode refresh token from user
		decodetRt, err := base64.StdEncoding.DecodeString(t.Refresh_token)
		t.Refresh_token = string(decodetRt)
		if err != nil {
			fmt.Println(err)
		}

		// compare refresh tokens from user and db
		err = bcrypt.CompareHashAndPassword([]byte(data.Refresh_token), []byte(t.Refresh_token))
		if err == nil {
			tNew, err := generateTokenPair(t.GUID)
			if err != nil {
				fmt.Println(err)
			}

			// update access and refresh tokens in db
			_, err = coll.UpdateOne(
				ctx,
				bson.M{"_id": t.GUID, "ip": r.RemoteAddr},
				bson.M{
					"$set": bson.M{
						"refresh_token": tNew["refresh_token"],
						"access_token":  tNew["access_token"],
					},
				},
			)
			if err != nil {
				fmt.Println(err)
			}

			//encode refresh token to base64
			tNew["refresh_token"] = base64.StdEncoding.EncodeToString([]byte(tNew["refresh_token"]))
			jsonResponse, err := json.Marshal(tNew)
			if err != nil {
				fmt.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResponse)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Something bad happened!"))
		}
	}
}

func (s *Server) getTokens() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		GUID := r.URL.Query().Get("GUID")

		// create access and refresh token for user with GUID
		t, err := generateTokenPair(GUID)
		if err != nil {
			fmt.Println(err)
		}
		// bcrypt for refresh token
		hashedRt, err := bcrypt.GenerateFromPassword([]byte(t["refresh_token"]), 10)
		if err != nil {
			fmt.Println(err)
		}
		data := Tokens{GUID, t["access_token"], string(hashedRt), r.RemoteAddr}

		coll := s.db.conn.Collection("tokens")
		ctx := context.Background()

		//insert created data into db
		_, err = coll.InsertOne(ctx, data)
		if err != nil {
			fmt.Println(err)
		}

		//encode refresh token to base64
		t["refresh_token"] = base64.StdEncoding.EncodeToString([]byte(t["refresh_token"]))
		jsonResponse, err := json.Marshal(t)
		if err != nil {
			fmt.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

//creat access and refresh tokens for user with GUID
func generateTokenPair(GUID string) (map[string]string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["expiresIn"] = time.Now().Add(time.Minute * 30).Unix()

	t, err := token.SignedString([]byte(GUID))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["expiresIn"] = time.Now().Add(time.Hour * 24 * 30).Unix()

	rt, err := refreshToken.SignedString([]byte(GUID))
	if err != nil {
		return nil, err
	}
	rt = base64.StdEncoding.EncodeToString([]byte(rt))

	return map[string]string{
		"access_token":  t,
		"refresh_token": rt,
	}, nil
}
