package users

type user struct {
	name string
	room string
}

var mapMake = make(map[string]*user)

func UsersJoin(userID string, joinedRoom string, userName string) *user {

	mapMake[userID] = &user{name: userName, room: joinedRoom}

	return mapMake[userID]
}

func GetJoinedUserDetails(userID string) *user {

	return mapMake[userID]
}

func UserLeave(userID string) {

	delete(mapMake, userID)
}
