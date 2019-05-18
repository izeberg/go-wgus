package wgus

import (
	"testing"
)

func TestReferenceRepository_Parse(t *testing.T) {
	data := []byte(`ILU.dll	155000	7F3CE4D1
Licenses.txt	96715	3676FC08
mods/1.5.0.2/readme.txt	74	E3301934
msvcp140.dll	440120	954DB51D`)

	repo := ReferenceRepository{}
	if err := repo.Parse(data); err != nil {
		t.Fatal(err)
	}
	if name := repo.GetModsFolder(); name != `1.5.0.2` {
		t.Fatal(`wrong mods folder`, name)
	}
	if repo[`Licenses.txt`].Name != `Licenses.txt` || repo[`Licenses.txt`].CRC != 0x3676FC08 || repo[`Licenses.txt`].Size != 96715 {
		t.Fatal(`wrong file info`, repo[`Licenses.txt`])
	}
}
