package password

import "testing"

func Test_GeneratePassWorld(t *testing.T) {
	world := GeneratePassWorld("Qaz@1234")
	t.Log(world)
}
