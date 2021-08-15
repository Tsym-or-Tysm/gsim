package gladiator

import (
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("gladiator's finale", New)
}

func New(c core.Character, s *core.Core, count int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.ATKP] = 0.18
		c.AddMod(core.CharStatMod{
			Key: "glad-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		switch c.WeaponClass() {
		case core.WeaponClassSpear:
		case core.WeaponClassSword:
		case core.WeaponClassClaymore:
		default:
			//don't add this mod if wrong weapon class
			return
		}
		m := make([]float64, core.EndStatType)
		m[core.DmgP] = 0.35
		c.AddMod(core.CharStatMod{
			Key: "glad-4pc",
			Amount: func(ds core.AttackTag) ([]float64, bool) {
				if ds != core.AttackTagNormal && ds != core.AttackTagExtra {
					return nil, false
				}
				return m, true
			},
			Expiry: -1,
		})
	}
	//add flat stat to char
}
