package state

var (
	incomingUpdate bool = false
)

func scheduleUpdate() {
	incomingUpdate = true
}
