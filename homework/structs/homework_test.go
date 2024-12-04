package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	if len(name) > 42 {
		panic("name is too long")
	}

	return func(person *GamePerson) {
		for i, c := range name {
			person.name[i] = encodeChar(c)
		}
	}
}

func encodeChar(c rune) uint8 {
	if c >= 'a' && c <= 'z' {
		return uint8(c-'a') + 0
	} else if c >= 'A' && c <= 'Z' {
		return uint8(c-'A') + 26
	} else if c >= '0' && c <= '9' {
		return uint8(c-'0') + 52
	} else if c == '_' {
		return 62
	}
	panic("invalid char")
}

func decodeChar(c uint8) rune {
	if c < 26 {
		return rune('a' + c)
	} else if c < 52 {
		return rune('A' + c)
	} else if c < 62 {
		return rune(0 + c)
	} else if c == 62 {
		return '_'
	}
	panic("invalid char")
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
		if gold < 0 || gold > 0x7FFFFFFF {
			panic("gold more than 31 bits can manage")
		}
		person.goldAndHome = (person.goldAndHome & 0x80000000) | uint32(gold&0x7FFFFFFF)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		if mana < 0 || mana > 1000 {
			panic("mana out of range")
		}
		person.data &= ^uint32(0x3FF)
		person.data |= uint32(mana & 0x3FF)

	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		if health < 0 || health > 1000 {
			panic("health out of range")
		}
		person.data &= ^uint32(0x3FF << 15)
		person.data |= uint32(health&0x3FF) << 15
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		if respect < 0 || respect > 10 {
			panic("respect out of range")
		}
		person.respect &= ^uint8(0xF)
		person.respect |= uint8(respect & 0xF)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		if strength < 0 || strength > 10 {
			panic("strength out of range")
		}
		person.data &= ^uint32(0xF << 10)
		person.data |= uint32(strength&0xF) << 10
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		if experience < 0 || experience > 10 {
			panic("experience out of range")
		}
		person.data &= ^uint32(0xF << 25)
		person.data |= uint32(experience&0xF) << 25
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		if level < 0 || level > 10 {
			panic("level out of range")
		}
		person.level &= ^uint8(0xF)
		person.level |= uint8(level & 0xF)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.goldAndHome |= 1 << 31
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.data |= 1 << 14
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.data |= 1 << 29
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		if personType < 0 || personType > 3 {
			panic("personType out of range")
		}
		person.data &= ^uint32(0x3 << 30)
		person.data |= uint32(personType&0x3) << 30
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x, y, z     int32     //12
	goldAndHome uint32    //4
	name        [42]uint8 //42
	level       uint8     //1
	respect     uint8     //1
	data        uint32    //4
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}
	for _, option := range options {
		option(&person)
	}
	return person
}

func (p *GamePerson) Name() string {
	var name [42]rune
	for i, c := range p.name {
		name[i] = decodeChar(c)
	}
	return string(name[:])
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
	return int(p.goldAndHome & 0x7FFFFFFF)
}

func (p *GamePerson) Mana() int {
	return int(p.data & 0x3FF)
}

func (p *GamePerson) Health() int {
	return int((p.data >> 15) & 0x3FF)
}

func (p *GamePerson) Respect() int {
	return int(p.respect & 0xF)
}

func (p *GamePerson) Strength() int {
	return int((p.data >> 10) & 0xF)
}

func (p *GamePerson) Experience() int {
	return int((p.data >> 25) & 0xF)
}

func (p *GamePerson) Level() int {
	return int(p.level & 0xF)
}

func (p *GamePerson) HasHouse() bool {
	return (p.goldAndHome>>31)&1 == 1
}

func (p *GamePerson) HasGun() bool {
	return (p.data>>14)&1 == 1
}

func (p *GamePerson) HasFamilty() bool {
	return (p.data>>29)&1 == 1
}

func (p *GamePerson) Type() int {
	return int((p.data >> 30) & 0x3)
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
