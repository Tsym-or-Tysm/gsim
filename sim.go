package gcsim

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/genshinsim/gcsim/internal/tmpl/enemy"
	"github.com/genshinsim/gcsim/internal/tmpl/player"
	"github.com/genshinsim/gcsim/pkg/core"
)

type Simulation struct {
	// f    int
	skip int
	C    *core.Core
	cfg  *core.SimulationConfig
	// queue
	queue             []core.Command
	dropQueueIfFailed bool
	//hurt event
	lastHurt int
	//energy event
	lastEnergyDrop int
	//result
	stats Stats
}

func New(cfg *core.SimulationConfig, seed int64, cust ...func(*Simulation) error) (*Simulation, error) {
	var err error
	s := &Simulation{}
	s.cfg = cfg

	c, err := core.New(
		func(c *core.Core) error {
			c.Rand = rand.New(rand.NewSource(seed))
			// if seed > 0 {
			// 	c.Rand = rand.New(rand.NewSource(seed))
			// } else {
			// 	c.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
			// }
			c.F = -1
			c.Flags.DamageMode = cfg.DamageMode
			c.Flags.EnergyCalcMode = opts.ERCalcMode
			c.Log, err = core.NewDefaultLogger(c, opts.Debug, true, opts.DebugPaths)
			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	s.C = c

	err = s.initTargets(cfg)
	if err != nil {
		return nil, err
	}
	err = s.initChars(cfg)
	if err != nil {
		return nil, err
	}
	s.stats.IsDamageMode = cfg.DamageMode

	if s.opts.LogDetails {
		s.initDetailLog()
	}

	err = s.initQueuer(cfg)
	if err != nil {
		return nil, err
	}

	s.randomOnHitEnergy()

	for _, f := range cust {
		err := f(s)
		if err != nil {
			return nil, err
		}
	}

	c.Init()

	if s.opts.LogDetails {
		//grab a snapshot for each char
		for i, c := range s.C.Chars {
			stats := c.Snapshot(&core.AttackInfo{
				Abil:      "stats-check",
				AttackTag: core.AttackTagNone,
			})
			s.stats.CharDetails[i].SnapshotStats = stats.Stats[:]
			s.stats.CharDetails[i].Element = c.Ele().String()
			s.stats.CharDetails[i].Weapon.Name = c.WeaponKey()
		}
	}

	// log.Println(s.cfg.Energy)

	return s, nil
}

func (s *Simulation) randomOnHitEnergy() {
	/**
	WeaponClassSword
	WeaponClassClaymore
	WeaponClassSpear
	WeaponClassBow
	WeaponClassCatalyst
	**/
	current := make([]float64, core.EndWeaponClass)
	inc := []float64{
		0.05,
		0.05,
		0.04,
		0.01,
		0.01,
	}

	//TODO not sure if there's like a 0.2s icd on this. for now let's add it in to be safe
	icd := 0
	s.C.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.AttackTag != core.AttackTagNormal && atk.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		//check icd
		if icd > s.C.F {
			return false
		}
		//check chance
		char := s.C.Chars[atk.Info.ActorIndex]
		w := char.WeaponClass()
		if s.C.Rand.Float64() > current[w] {
			//increment chance
			current[w] += inc[w]
			return false
		}
		//add energy
		char.AddEnergy(1)
		s.C.Log.Debugw("random energy on normal", "frame", s.C.F, "event", core.LogEnergyEvent, "char", atk.Info.ActorIndex, "chance", current[w])
		//set icd
		icd = s.C.F + 12
		current[w] = 0
		return false
	}, "random-energy-restore-on-hit")
	s.C.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		//TODO: assuming we clear the probability on swap
		for i := range current {
			current[i] = 0
		}
		return false
	}, "random-energy-restore-on-hit-swap")
}

func (s *Simulation) initTargets() error {
	s.C.Targets = make([]core.Target, len(s.cfg.Targets)+1)
	if s.opts.LogDetails {
		s.stats.ElementUptime = make([]map[core.EleType]int, len(s.C.Targets))
		s.stats.ElementUptime[0] = make(map[core.EleType]int)
	}
	s.C.Targets[0] = player.New(0, s.C)

	//first target is the player
	for i := 0; i < len(s.cfg.Targets); i++ {
		s.cfg.Targets[i].Size = 0.5
		if i > 0 {
			cfg.Targets[i].CoordX = 0.6
			cfg.Targets[i].CoordY = 0
		}
		s.C.Targets[i+1] = enemy.New(i+1, s.C, cfg.Targets[i])
		if s.opts.LogDetails {
			s.stats.ElementUptime[i+1] = make(map[core.EleType]int)
		}
	}
	return nil
}

func (s *Simulation) initChars(cfg core.SimulationConfig) error {
	dup := make(map[core.CharKey]bool)
	res := make(map[core.EleType]int)

	count := len(cfg.Characters.Profile)

	if count > 4 {
		return fmt.Errorf("more than 4 characters in a team detected")
	}

	if s.opts.LogDetails {
		s.stats.CharNames = make([]string, count)
		s.stats.CharDetails = make([]CharDetail, 0, count)
		s.stats.DamageByChar = make([]map[string]float64, count)
		s.stats.DamageInstancesByChar = make([]map[string]int, count)
		s.stats.DamageByCharByTargets = make([]map[int]float64, count)
		s.stats.DamageDetailByTime = make(map[DamageDetails]float64)
		s.stats.CharActiveTime = make([]int, count)
		s.stats.AbilUsageCountByChar = make([]map[string]int, count)
		s.stats.ParticleCount = make(map[string]int)
		s.stats.EnergyWhenBurst = make([][]float64, count)
	}

	s.C.ActiveChar = -1
	for i, v := range cfg.Characters.Profile {
		//call new char function
		char, err := s.C.AddChar(v)
		if err != nil {
			return err
		}

		if v.Base.Key == cfg.Characters.Initial {
			s.C.ActiveChar = i
		}

		if _, ok := dup[v.Base.Key]; ok {
			return fmt.Errorf("duplicated character %v", v.Base.Key)
		}
		dup[v.Base.Key] = true

		//track resonance
		res[char.Ele()]++

		//setup maps
		if s.opts.LogDetails {
			s.stats.DamageByChar[i] = make(map[string]float64)
			s.stats.DamageInstancesByChar[i] = make(map[string]int)
			s.stats.DamageByCharByTargets[i] = make(map[int]float64)
			s.stats.AbilUsageCountByChar[i] = make(map[string]int)
			s.stats.CharNames[i] = v.Base.Key.String()
			s.stats.EnergyWhenBurst[i] = make([]float64, 0, s.opts.Duration/12+2)

			//log the character data
			s.stats.CharDetails = append(s.stats.CharDetails, CharDetail{
				Name:     v.Base.Key.String(),
				Level:    v.Base.Level,
				MaxLevel: v.Base.Level,
				Cons:     v.Base.Cons,
				Weapon: WeaponDetail{
					Refine:   v.Weapon.Refine,
					Level:    v.Weapon.Level,
					MaxLevel: v.Weapon.MaxLevel,
				},
				Talents: TalentDetail{
					Attack: v.Talents.Attack,
					Skill:  v.Talents.Skill,
					Burst:  v.Talents.Burst,
				},
				Sets: v.Sets,
			})

		}

	}

	if s.C.ActiveChar == -1 {
		return errors.New("no active char set")
	}

	s.initResonance(res)

	return nil
}

func (s *Simulation) initQueuer(cfg core.SimulationConfig) error {
	s.queue = make([]core.Command, 0, 20)
	// cust := make(map[string]int)
	// for i, v := range cfg.Rotation {
	// 	if v.Label != "" {
	// 		cust[v.Name] = i
	// 	}
	// 	// log.Println(v.Conditions)
	// }
	for i, v := range cfg.Rotation {
		if _, ok := s.C.CharByName(v.SequenceChar); v.Type == core.ActionBlockTypeSequence && !ok {
			return fmt.Errorf("invalid char in rotation %v; %v", v.SequenceChar, v)
		}
		cfg.Rotation[i].LastQueued = -1
	}
	s.C.Log.Debugw(
		"setting queue",
		"frame", s.C.F,
		"event", core.LogSimEvent,
		"pq", cfg.Rotation,
	)

	err := s.C.Queue.SetActionList(cfg.Rotation)
	return err
}

func (s *Simulation) initDetailLog() {
	var sb strings.Builder
	s.stats.ReactionsTriggered = make(map[core.ReactionType]int)
	//add new targets
	s.C.Events.Subscribe(core.OnTargetAdded, func(args ...interface{}) bool {
		t := args[0].(core.Target)

		s.C.Log.Debugw("Target Added", "frame", s.C.F, "event", core.LogSimEvent, "target_type", t.Type())

		s.stats.ElementUptime = append(s.stats.ElementUptime, make(map[core.EleType]int))

		return false
	}, "sim-new-target-stats")
	//add call backs to track details
	s.C.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		t := args[0].(core.Target)

		// No need to pull damage stats for non-enemies
		if t.Type() != core.TargettableEnemy {
			return false
		}
		atk := args[1].(*core.AttackEvent)

		//skip if do not log
		if atk.Info.DoNotLog {
			return false
		}

		dmg := args[2].(float64)
		sb.Reset()
		sb.WriteString(atk.Info.Abil)
		if atk.Info.Amped {
			if atk.Info.AmpMult == 1.5 {
				sb.WriteString(" [amp: 1.5]")
			} else if atk.Info.AmpMult == 2 {
				sb.WriteString(" [amp: 2.0]")
			}
		}
		s.stats.DamageByChar[atk.Info.ActorIndex][sb.String()] += dmg
		if dmg > 0 {
			s.stats.DamageInstancesByChar[atk.Info.ActorIndex][sb.String()] += 1
		}
		s.stats.DamageByCharByTargets[atk.Info.ActorIndex][t.Index()] += dmg

		// Want to capture information in 0.25s intervals - allows more flexibility in bucketizing
		frameBucket := int(s.C.F/15) * 15
		details := DamageDetails{
			FrameBucket: frameBucket,
			Char:        atk.Info.ActorIndex,
			Target:      t.Index(),
		}
		// Go defaults to 0 for map values that don't exist
		s.stats.DamageDetailByTime[details] += dmg
		return false
	}, "dmg-log")

	eventSubFunc := func(t core.ReactionType) func(args ...interface{}) bool {
		return func(args ...interface{}) bool {
			s.stats.ReactionsTriggered[t]++
			return false
		}
	}

	var reactions = map[core.EventType]core.ReactionType{
		core.OnOverload:           core.Overload,
		core.OnSuperconduct:       core.Superconduct,
		core.OnMelt:               core.Melt,
		core.OnVaporize:           core.Vaporize,
		core.OnFrozen:             core.Freeze,
		core.OnElectroCharged:     core.ElectroCharged,
		core.OnSwirlHydro:         core.SwirlHydro,
		core.OnSwirlCryo:          core.SwirlCryo,
		core.OnSwirlElectro:       core.SwirlElectro,
		core.OnSwirlPyro:          core.SwirlPyro,
		core.OnCrystallizeCryo:    core.CrystallizeCryo,
		core.OnCrystallizeElectro: core.CrystallizeElectro,
		core.OnCrystallizeHydro:   core.CrystallizeHydro,
		core.OnCrystallizePyro:    core.CrystallizePyro,
	}

	for k, v := range reactions {
		s.C.Events.Subscribe(k, eventSubFunc(v), "reaction-log")
	}

	s.C.Events.Subscribe(core.OnParticleReceived, func(args ...interface{}) bool {
		p := args[0].(core.Particle)
		s.stats.ParticleCount[p.Source] += p.Num
		return false
	}, "particles-log")

	s.C.Events.Subscribe(core.PreBurst, func(args ...interface{}) bool {
		activeChar := s.C.Chars[s.C.ActiveChar]
		s.stats.EnergyWhenBurst[s.C.ActiveChar] = append(s.stats.EnergyWhenBurst[s.C.ActiveChar], activeChar.CurrentEnergy())
		return false
	}, "energy-calc-log")

}
