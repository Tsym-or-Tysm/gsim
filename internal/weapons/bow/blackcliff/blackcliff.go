package blackcliff

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
)

func init() {
	combat.RegisterWeaponFunc("blackcliff warbow", weapon)
}

func weapon(c def.Character, s def.Sim, log def.Logger, r int, param map[string]int) {

	atk := 0.09 + float64(r)
	index := 0
	stacks := []int{-1, -1, -1}

	m := make([]float64, def.EndStatType)
	c.AddMod(def.CharStatMod{
		Key: "blackcliff",
		Amount: func(a def.AttackTag) ([]float64, bool) {
			count := 0
			for _, v := range stacks {
				if v > s.Frame() {
					count++
				}
			}
			m[def.ATKP] = atk * float64(count)
			return m, true
		},
		Expiry: -1,
	})

	s.AddOnTargetDefeated(func(t def.Target) {
		stacks[index] = s.Frame() + 1800
		index++
		if index == 3 {
			index = 0
		}
	}, fmt.Sprintf("blackcliff-warbow-%v", c.Name()))
}
