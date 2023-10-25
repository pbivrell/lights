package data

import "time"

type Light struct {
	ID      int64     `firestore:"id"`
	IP      int64     `firestore:"ip"`
	Updated time.Time `firestore:"updated"`
}
