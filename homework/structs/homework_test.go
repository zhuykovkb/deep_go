package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	// Маски и смещения для поля data
	DataManaMask        = 0x3FF // 10 бит для маны
	DataStrengthMask    = 0xF   // 4 бита для силы
	DataStrengthShift   = 10
	DataWeaponBit       = 14    // 1 бит для оружия
	DataHealthMask      = 0x3FF // 10 бит для здоровья
	DataHealthShift     = 15
	DataExperienceMask  = 0xF // 4 бита для опыта
	DataExperienceShift = 25
	DataFamilyBit       = 29  // 1 бит для семьи
	DataTypeMask        = 0x3 // 2 бита для типа игрока
	DataTypeShift       = 30

	// Маски и смещения для goldAndHome
	GoldMask = 0x7FFFFFFF // 31 бит для золота
	HomeBit  = 31         // 1 бит для дома

	// Маски для level и respect
	LevelMask   = 0xF // 4 бита для уровня
	RespectMask = 0xF // 4 бита для уважения
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	if len(name) > 42 {
		panic("name is too long")
	}

	return func(person *GamePerson) {
		for i, c := range name {
			person.name[i] = uint8(c)
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		if gold < 0 || gold > GoldMask {
			panic("gold out of range")
		}
		person.goldAndHome = (person.goldAndHome & ^uint32(GoldMask)) | uint32(gold&GoldMask)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		if mana < 0 || mana > 1000 {
			panic("mana out of range")
		}
		person.data &= ^uint32(DataManaMask)
		person.data |= uint32(mana & DataManaMask)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		if health < 0 || health > 1000 {
			panic("health out of range")
		}
		person.data &= ^uint32(DataHealthMask << DataHealthShift)
		person.data |= uint32(health&DataHealthMask) << DataHealthShift
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		if respect < 0 || respect > 10 {
			panic("respect out of range")
		}
		person.respect &= ^uint8(RespectMask)
		person.respect |= uint8(respect & RespectMask)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		if strength < 0 || strength > 10 {
			panic("strength out of range")
		}
		person.data &= ^uint32(DataStrengthMask << DataStrengthShift)
		person.data |= uint32(strength&DataStrengthMask) << DataStrengthShift
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		if experience < 0 || experience > 10 {
			panic("experience out of range")
		}
		person.data &= ^uint32(DataExperienceMask << DataExperienceShift)
		person.data |= uint32(experience&DataExperienceMask) << DataExperienceShift
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		if level < 0 || level > 10 {
			panic("level out of range")
		}
		person.level &= ^uint8(LevelMask)
		person.level |= uint8(level & LevelMask)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.goldAndHome |= 1 << HomeBit
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.data |= 1 << DataWeaponBit
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.data |= 1 << DataFamilyBit
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		if personType < 0 || personType > 3 {
			panic("personType out of range")
		}
		person.data &= ^uint32(DataTypeMask << DataTypeShift)
		person.data |= uint32(personType&DataTypeMask) << DataTypeShift
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x, y, z     int32     //12 байт
	goldAndHome uint32    //4 байта
	name        [42]uint8 //42 байт
	level       uint8     //1 байт
	respect     uint8     //1 байт
	data        uint32    //4 байта - характеристики персонажа:
	/*
		Биты:
		0–9:   Мана (10 бит, значения от 0 до 1000)
		10–13: Сила (4 бита, значения от 0 до 10)
		14:    Есть ли оружие (1 бит, true/false)
		15–24: Здоровье (10 бит, значения от 0 до 1000)
		25–28: Опыт (4 бита, значения от 0 до 10)
		29:    Есть ли семья (1 бит, true/false)
		30–31: Тип игрока (2 бита, значения: 0 - строитель, 1 - кузнец, 2 - воин)
	*/
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}
	for _, option := range options {
		option(&person)
	}
	return person
}

func (p *GamePerson) Name() string {
	return string(p.name[:])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.goldAndHome & GoldMask)
}

func (p *GamePerson) Mana() int {
	return int(p.data & DataManaMask)
}

func (p *GamePerson) Health() int {
	return int((p.data >> DataHealthShift) & DataHealthMask)
}

func (p *GamePerson) Respect() int {
	return int(p.respect & RespectMask)
}

func (p *GamePerson) Strength() int {
	return int((p.data >> DataStrengthShift) & DataStrengthMask)
}

func (p *GamePerson) Experience() int {
	return int((p.data >> DataExperienceShift) & DataExperienceMask)
}

func (p *GamePerson) Level() int {
	return int(p.level & LevelMask)
}

func (p *GamePerson) HasHouse() bool {
	return (p.goldAndHome>>HomeBit)&1 == 1
}

func (p *GamePerson) HasGun() bool {
	return (p.data>>DataWeaponBit)&1 == 1
}

func (p *GamePerson) HasFamilty() bool {
	return (p.data>>DataFamilyBit)&1 == 1
}

func (p *GamePerson) Type() int {
	return int((p.data >> DataTypeShift) & DataTypeMask)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
